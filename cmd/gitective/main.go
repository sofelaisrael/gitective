package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	commitPkg "github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/git"
	"github.com/sofelaisrael/gitective/internal/llm"
	"github.com/sofelaisrael/gitective/internal/personality"
	"github.com/sofelaisrael/gitective/internal/themes"
	"github.com/sofelaisrael/gitective/internal/ui"
)

func main() {
	styleFlag := flag.String("style", "", "Style engine style (e.g. cyberpunk-commit, renaissance-commit, shakespeare, pirate)")
	themeFlag := flag.String("theme", "", "Legacy theme: \"cyberpunk\" or \"renaissance\" (fallback if style-engine is offline)")
	intensityFlag := flag.Float64("intensity", 0.8, "Style intensity (0.0-1.0)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gitective [options]\n\n")
		fmt.Fprintf(os.Stderr, "Investigate your Git history.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --style string     Style-engine style (e.g. cyberpunk-commit, shakespeare)\n")
		fmt.Fprintf(os.Stderr, "  --theme string     Legacy theme fallback (\"cyberpunk\" or \"renaissance\")\n")
		fmt.Fprintf(os.Stderr, "  --intensity float  Style intensity 0.0-1.0 (default 0.8)\n")
		fmt.Fprintf(os.Stderr, "  --help             Show this help\n")
	}
	flag.Parse()

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

	latest := commits[0]
	facts, err := commitPkg.ExtractFacts(root, latest.Hash)
	if err != nil {
		pterm.Error.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Try LLM path first if --style is set
	if *styleFlag != "" {
		pterm.Info.Printf("Connecting to style-engine (style: %s)...\n", *styleFlag)
		client := llm.NewClient()
		contextPrompt := llm.BuildContext(facts)

		result, err := client.Transform(contextPrompt, *styleFlag, *intensityFlag)
		if err != nil {
			pterm.Warning.Printf("Style-engine unavailable: %v\n", err)
			pterm.Info.Println("Falling back to local theme...")
			// Fall through to local theme
		} else {
			// Render LLM result with lipgloss box
			p := personality.Classify(facts)
			header := fmt.Sprintf("// %s", strings.ToUpper(string(p)))
			footer := fmt.Sprintf("STATUS: %s | Score: %.0f%%", strings.ToUpper(strings.ReplaceAll(string(p), "_", " ")), result.Score*100)
			fmt.Println(ui.RenderBox(header, result.Transformed, footer))
			return
		}
	}

	// Fallback: local theme rendering
	if *themeFlag == "" {
		*themeFlag = "cyberpunk"
	}
	var theme themes.Theme
	switch *themeFlag {
	case "cyberpunk":
		theme = themes.CyberpunkTheme{}
	case "renaissance":
		theme = themes.RenaissanceTheme{}
	default:
		pterm.Error.Printf("Unknown theme: %q\n", *themeFlag)
		os.Exit(1)
	}
	fmt.Println(theme.Render(facts))
}
