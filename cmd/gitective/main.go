package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/sofelaisrael/gitective/internal/analysis"
	commitPkg "github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/git"
	themes "github.com/sofelaisrael/gitective/internal/themes"
	"github.com/sofelaisrael/gitective/internal/ui"
)

func main() {
	theme := flag.String("theme", "cyberpunk", "theme to use")
	demo := flag.Bool("demo", false, "demo mode")
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
	if len(commits) > 0 {
		latest := commits[0]
		facts, err := commitPkg.ExtractFacts(root, latest.Hash)
		if err == nil {
			if *theme == "cyberpunk" {
				t := themes.CyberpunkTheme{}
				fmt.Println(t.Render(facts))
			} else {
				fmt.Println(facts.Message)
			}
		}
	}
	if *demo {
		return
	}
	timingData := make([]analysis.CommitData, 0, len(commits))
	for _, c := range commits {
		timingData = append(timingData, analysis.CommitData{Timestamp: c.Timestamp})
	}
	timing := analysis.AnalyzeTiming(timingData)
	activity := analysis.AnalyzeActivity(timingData)
	ui.Successf("Repository: %s\n", root)
	pterm.Info.Println(fmt.Sprintf("Commits found: %d", len(commits)))
	fmt.Println()
	fmt.Println("Coding hours:")
	fmt.Printf("Morning:   %d (%.1f%%)\n", timing.Morning, analysis.Percentage(timing.Morning, timing.Total))
	fmt.Printf("Afternoon: %d (%.1f%%)\n", timing.Afternoon, analysis.Percentage(timing.Afternoon, timing.Total))
	fmt.Printf("Evening:   %d (%.1f%%)\n", timing.Evening, analysis.Percentage(timing.Evening, timing.Total))
	fmt.Printf("Night:     %d (%.1f%%)\n", timing.Night, analysis.Percentage(timing.Night, timing.Total))
	fmt.Println()
	fmt.Println("Investigation:")
	fmt.Println("Most active period:", analysis.DominantPeriod(timing))
	fmt.Println()
	fmt.Println("Activity:")
	fmt.Println("Total commits:", activity.TotalCommits)
	fmt.Printf("Average commits/day: %.2f\n", activity.AveragePerDay)
	fmt.Println("Most active day:", activity.MostActiveDay)
	fmt.Println("Most active weekday:", activity.MostActiveWeekday)
	fmt.Println("Most commits in one day:", activity.MaxCommitsInDay)
	fmt.Println("Longest inactive gap:", activity.LongestInactiveGap)
	for i, commit := range commits {
		if i == 10 {
			break
		}
		fmt.Println(commit.Timestamp.Format("2006-01-02 15:04"))
		fmt.Println(commit.Message)
		fmt.Println()
	}
}
