package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	commitPkg "github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/config"
	"github.com/sofelaisrael/gitective/internal/git"
	"github.com/sofelaisrael/gitective/internal/llm"
	"github.com/sofelaisrael/gitective/internal/themes"
	"github.com/sofelaisrael/gitective/internal/ui"
)

func main() {
	commitFlag := flag.Bool("commit", false, "Suggest a commit message from uncommitted changes")
	shortCommitFlag := flag.Bool("c", false, "Same as --commit")
	setupFlag := flag.Bool("setup", false, "Run the interactive setup wizard")
	styleFlag := flag.String("style", "", "Override default style")
	themeFlag := flag.String("theme", "", "Legacy theme fallback if style-engine is offline")
	intensityFlag := flag.Float64("intensity", -1, "Override default intensity (0.0-1.0)")
	commitsFlag := flag.Int("commits", 1, "Number of recent commits to visualize")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gitective [command] [options]\n\n")
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  (none)       Visualize recent commits\n")
		fmt.Fprintf(os.Stderr, "  -c           Suggest a commit message from uncommitted changes\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --setup          Run interactive setup wizard\n")
		fmt.Fprintf(os.Stderr, "  --style string   Override default style\n")
		fmt.Fprintf(os.Stderr, "  --theme string   Legacy theme fallback\n")
		fmt.Fprintf(os.Stderr, "  --intensity float Override default intensity\n")
		fmt.Fprintf(os.Stderr, "  --commits int    Number of recent commits to visualize (default 1)\n")
		fmt.Fprintf(os.Stderr, "  --help           Show this help\n")
	}
	flag.Parse()

	if *setupFlag {
		if err := config.RunSetup(); err != nil {
			pterm.Error.Printf("Setup failed: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if !config.IsConfigured() {
		pterm.Warning.Println("No config found. Let's set you up.")
		if err := config.RunSetup(); err != nil {
			pterm.Error.Printf("Setup failed: %v\n", err)
			os.Exit(1)
		}
		pterm.Println()
	}

	cfg := config.Load()

	style := cfg.DefaultStyle
	if *styleFlag != "" {
		style = *styleFlag
	}
	intensity := cfg.DefaultIntensity
	if *intensityFlag >= 0 {
		intensity = *intensityFlag
	}

	if *commitFlag || *shortCommitFlag {
		runCommitSuggest(style, intensity, cfg.StyleEngineURL)
		return
	}

	runVisualizer(style, intensity, cfg.StyleEngineURL, *commitsFlag, *themeFlag)
}

func runCommitSuggest(style string, intensity float64, engineURL string) {
	spinner, _ := ui.Spinner("Scanning working directory...")
	root, err := git.FindRepository(".")
	if err != nil {
		spinner.Fail("Not a git repository")
		os.Exit(1)
	}

	diff, err := git.GetDiff(root)
	if err != nil {
		spinner.Fail("Could not read diff")
		pterm.Error.Printf("%v\n", err)
		os.Exit(1)
	}

	if diff == "" {
		spinner.Fail("No changes to commit")
		pterm.Info.Println("Stage some changes first: git add <files>")
		return
	}

	spinner.Success(fmt.Sprintf("Found %d bytes of changes", len(diff)))

	pterm.Println()
	pterm.Info.Printf("Generating commit message (style: %s)...\n", style)

	client := llm.NewClient(engineURL)
	userPrompt := llm.BuildStyledCommitPrompt(diff)
	result, err := client.Transform(userPrompt, style, intensity, 25)
	if err != nil {
		pterm.Error.Printf("Style-engine error: %v\n", err)
		return
	}

	styled := strings.TrimSpace(result.Transformed)
	color := ui.StyleColor(style)
	header := "// SUGGESTED COMMIT"
	footer := fmt.Sprintf("Style: %s | Score: %.0f%%", style, result.Score*100)
	fmt.Println(ui.RenderBox(header, styled, footer, color))

	fullMsg := git.AppendTrailers(styled, style, intensity)

	fmt.Println()
	pterm.Info.Println("To use this message:")
	fmt.Println()
	fmt.Printf("  git commit -m \"%s\"\n", strings.ReplaceAll(fullMsg, "\n", "\\n\\n"))
	fmt.Println()
}

func runVisualizer(style string, intensity float64, engineURL string, numCommits int, theme string) {
	fmt.Println(ui.RenderBanner("Gitective"))

	spinner, _ := ui.Spinner("Scanning repository...")
	root, err := git.FindRepository(".")
	if err != nil {
		spinner.Fail("Repository check failed")
		if ui.HandleRepoError(err) {
			os.Exit(1)
		}
		pterm.Error.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	spinner.Success("Repository found")

	commits, err := git.GetCommits(root)
	if err != nil {
		pterm.Error.Printf("Could not read Git history: %v\n", err)
		os.Exit(1)
	}
	if len(commits) == 0 {
		pterm.Warning.Println("No commits found.")
		return
	}

	n := numCommits
	if n > len(commits) {
		n = len(commits)
	}

	client := llm.NewClient(engineURL)

	commitStyles := make([]git.CommitStyle, n)
	for i := 0; i < n; i++ {
		commitStyles[i] = git.ReadTrailers(root, commits[i].Hash)
		if commitStyles[i].Style == "" {
			commitStyles[i].Style = style
			commitStyles[i].Intensity = intensity
		}
	}

	fmt.Println(ui.RenderCommitTree(commits[:n], commitStyles))

	for i := 0; i < n; i++ {
		latest := commits[i]
		cs := commitStyles[i]

		facts, err := commitPkg.ExtractFacts(root, latest.Hash)
		if err != nil {
			pterm.Error.Printf("Error extracting facts for %s: %v\n", latest.Hash[:8], err)
			continue
		}

		contextPrompt := llm.BuildContext(facts)
		result, err := client.Transform(contextPrompt, cs.Style, cs.Intensity, 0)
		if err != nil {
			pterm.Warning.Printf("Style-engine unavailable for %s: %v\n", latest.Hash[:8], err)
			continue
		}

		header := fmt.Sprintf("// %s", strings.ToUpper(cs.Style))
		footer := fmt.Sprintf("%s · %s · Score: %.0f%%", latest.Hash[:7], latest.Author, result.Score*100)
		fmt.Println(ui.RenderBox(header, result.Transformed, footer, ui.StyleColor(cs.Style)))
	}
}

func fallbackRender(facts commitPkg.CommitFacts, theme string) {
	if theme == "" {
		theme = "cyberpunk"
	}
	var t themes.Theme
	switch theme {
	case "cyberpunk":
		t = themes.CyberpunkTheme{}
	case "renaissance":
		t = themes.RenaissanceTheme{}
	default:
		pterm.Error.Printf("Unknown theme: %q\n", theme)
		return
	}
	fmt.Println(t.Render(facts))
}
