package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/pterm/pterm"
	"github.com/sofelaisrael/gitective/internal/analysis"
	commitPkg "github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/git"
	"github.com/sofelaisrael/gitective/internal/themes"
	"github.com/sofelaisrael/gitective/internal/ui"
)

func main() {
	themeFlag := flag.String("theme", "cyberpunk", "Theme to use: \"cyberpunk\" or \"renaissance\"")
	commitsFlag := flag.Int("commits", 1, "Number of recent commits to analyse (0 = all)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gitective [options]\n\n")
		fmt.Fprintf(os.Stderr, "Investigate your Git history.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  --theme string    Theme: \"cyberpunk\" or \"renaissance\" (default \"cyberpunk\")\n")
		fmt.Fprintf(os.Stderr, "  --commits int     Commits to analyse: 1 = latest, 0 = all (default 1)\n")
		fmt.Fprintf(os.Stderr, "  --help            Show this help\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  gitective                          # Latest commit, cyberpunk theme\n")
		fmt.Fprintf(os.Stderr, "  gitective --theme renaissance      # Latest commit, renaissance theme\n")
		fmt.Fprintf(os.Stderr, "  gitective --commits 10             # Last 10 commits\n")
		fmt.Fprintf(os.Stderr, "  gitective --theme renaissance --commits 0  # All commits, renaissance\n")
	}
	flag.Parse()

	fmt.Println(ui.RenderBanner("Gitective"))

	spinner, _ := ui.Spinner("Detecting git repository...")
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

	// Select theme
	var theme themes.Theme
	switch *themeFlag {
	case "cyberpunk":
		theme = themes.CyberpunkTheme{}
	case "renaissance":
		theme = themes.RenaissanceTheme{}
	default:
		pterm.Error.Printf("Unknown theme: %q (use \"cyberpunk\" or \"renaissance\")\n", *themeFlag)
		os.Exit(1)
	}

	// Determine how many commits to analyse
	numCommits := *commitsFlag
	if numCommits <= 0 || numCommits > len(commits) {
		numCommits = len(commits)
	}

	// Analyse individual commits
	for i := 0; i < numCommits; i++ {
		commit := commits[i]
		facts, err := commitPkg.ExtractFacts(root, commit.Hash)
		if err != nil {
			pterm.Error.Printf("Error analysing commit %s: %v\n", commit.Hash[:8], err)
			continue
		}
		fmt.Println(theme.Render(facts))
		fmt.Println()
	}

	// Activity & timing analysis
	var commitData []analysis.CommitData
	for _, c := range commits {
		commitData = append(commitData, analysis.CommitData{Timestamp: c.Timestamp})
	}
	activity := analysis.AnalyzeActivity(commitData)
	timing := analysis.AnalyzeTiming(commitData)

	fmt.Printf("── Activity Summary ──────────────────\n")
	fmt.Printf("  Total commits:  %d\n", activity.TotalCommits)
	fmt.Printf("  Active days:    %d\n", countActiveDays(commits))
	if activity.TotalCommits > 1 {
		fmt.Printf("  Avg/day:        %.1f\n", activity.AveragePerDay)
	}
	fmt.Printf("  Peak day:       %s (%d commits)\n", activity.MostActiveDay, activity.MaxCommitsInDay)
	fmt.Printf("  Dominant time:  %s (%.0f%%)\n", analysis.DominantPeriod(timing),
		analysis.Percentage(
			getPeriodCount(timing, analysis.DominantPeriod(timing)),
			timing.Total,
		))
	if activity.LongestInactiveGap > 0 {
		fmt.Printf("  Longest gap:    %s\n", formatDuration(activity.LongestInactiveGap))
	}
	fmt.Printf("─────────────────────────────────────\n")
}

func countActiveDays(commits []git.Commit) int {
	days := make(map[string]bool)
	for _, c := range commits {
		days[c.Timestamp.Format("2006-01-02")] = true
	}
	return len(days)
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if days > 0 {
		return fmt.Sprintf("%dd %dh", days, hours)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return "<1h"
}

func getPeriodCount(timing analysis.TimeAnalysis, period string) int {
	switch period {
	case "morning":
		return timing.Morning
	case "afternoon":
		return timing.Afternoon
	case "evening":
		return timing.Evening
	case "night":
		return timing.Night
	default:
		return 0
	}
}
