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
	setupFlag := flag.Bool("setup", false, "Run the interactive setup wizard")

	styleFlag := flag.String("style", "", "Override default style")
	themeFlag := flag.String("theme", "", "Legacy theme fallback if style-engine is offline")
	intensityFlag := flag.Float64("intensity", -1, "Override default intensity (0.0-1.0)")
	commitsFlag := flag.Int("commits", 1, "Number of recent commits to process")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gitective [options]\n\n")
		fmt.Fprintf(os.Stderr, "Investigate your Git history.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --setup          Run interactive setup wizard\n")
		fmt.Fprintf(os.Stderr, "  --style string   Override default style\n")
		fmt.Fprintf(os.Stderr, "  --theme string   Legacy theme fallback\n")
		fmt.Fprintf(os.Stderr, "  --intensity float Override default intensity\n")
		fmt.Fprintf(os.Stderr, "  --commits int    Number of recent commits (default 1)\n")
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

	n := *commitsFlag
	if n > len(commits) {
		n = len(commits)
	}

	client := llm.NewClient(cfg.StyleEngineURL)
	color := ui.StyleColor(style)

	for i := 0; i < n; i++ {
		latest := commits[i]
		facts, err := commitPkg.ExtractFacts(root, latest.Hash)
		if err != nil {
			pterm.Error.Printf("Error extracting facts for %s: %v\n", latest.Hash[:8], err)
			continue
		}

		if i == 0 || *commitsFlag == 1 {
			pterm.Info.Printf("Style: %s | Intensity: %.0f%%\n", style, intensity*100)
		}

		contextPrompt := llm.BuildContext(facts)
		result, err := client.Transform(contextPrompt, style, intensity)
		if err != nil {
			pterm.Warning.Printf("Style-engine unavailable: %v\n", err)
			pterm.Info.Println("Falling back to local theme...")
			fallbackRender(facts, *themeFlag)
			continue
		}

		header := fmt.Sprintf("// %s", strings.ToUpper(style))
		footer := fmt.Sprintf("Score: %.0f%% | %d retries", result.Score*100, result.Retries)
		fmt.Println(ui.RenderBox(header, result.Transformed, footer, color))
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
