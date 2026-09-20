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
	cfg := config.Load()

	styleFlag := flag.String("style", cfg.DefaultStyle, "Style engine style")
	themeFlag := flag.String("theme", "", "Legacy theme fallback if style-engine is offline")
	intensityFlag := flag.Float64("intensity", cfg.DefaultIntensity, "Style intensity (0.0-1.0)")
	commitsFlag := flag.Int("commits", 1, "Number of recent commits to process")
	setStyleFlag := flag.String("set-style", "", "Save a default style to config and exit")
	setIntensityFlag := flag.Float64("set-intensity", -1, "Save a default intensity to config and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gitective [options]\n\n")
		fmt.Fprintf(os.Stderr, "Investigate your Git history.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --style string        Style-engine style (default: from config, or cyberpunk-commit)\n")
		fmt.Fprintf(os.Stderr, "  --theme string        Legacy theme fallback if style-engine is offline\n")
		fmt.Fprintf(os.Stderr, "  --intensity float     Style intensity 0.0-1.0 (default: from config, or 0.8)\n")
		fmt.Fprintf(os.Stderr, "  --commits int         Number of recent commits to process (default 1)\n")
		fmt.Fprintf(os.Stderr, "  --set-style string    Save default style to config and exit\n")
		fmt.Fprintf(os.Stderr, "  --set-intensity float Save default intensity to config and exit\n")
		fmt.Fprintf(os.Stderr, "  --help                Show this help\n")
	}
	flag.Parse()

	if *setStyleFlag != "" {
		cfg.DefaultStyle = *setStyleFlag
		if err := config.Save(cfg); err != nil {
			pterm.Error.Printf("Could not save config: %v\n", err)
			os.Exit(1)
		}
		pterm.Success.Printf("Default style set to: %s\n", *setStyleFlag)
		return
	}

	if *setIntensityFlag >= 0 {
		cfg.DefaultIntensity = *setIntensityFlag
		if err := config.Save(cfg); err != nil {
			pterm.Error.Printf("Could not save config: %v\n", err)
			os.Exit(1)
		}
		pterm.Success.Printf("Default intensity set to: %.0f%%\n", *setIntensityFlag*100)
		return
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
	color := ui.StyleColor(*styleFlag)

	for i := 0; i < n; i++ {
		latest := commits[i]
		facts, err := commitPkg.ExtractFacts(root, latest.Hash)
		if err != nil {
			pterm.Error.Printf("Error extracting facts for %s: %v\n", latest.Hash[:8], err)
			continue
		}

		if i == 0 || *commitsFlag == 1 {
			pterm.Info.Printf("Style: %s | Intensity: %.0f%%\n", *styleFlag, *intensityFlag*100)
		}

		contextPrompt := llm.BuildContext(facts)
		result, err := client.Transform(contextPrompt, *styleFlag, *intensityFlag)
		if err != nil {
			pterm.Warning.Printf("Style-engine unavailable: %v\n", err)
			pterm.Info.Println("Falling back to local theme...")
			fallbackRender(facts, *themeFlag)
			continue
		}

		header := fmt.Sprintf("// %s", strings.ToUpper(*styleFlag))
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
