package main

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	"github.com/sofelaisrael/gitective/internal/git"
	"github.com/sofelaisrael/gitective/internal/ui"
)

func main() {
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
	pterm.Info.Println(fmt.Sprintf("Commits found: %d", len(commits)))
	ui.Successf("Repository: %s\n", root)
	for i, commit := range commits {
		if i == 10 {
			break
		}
		fmt.Println(commit.Timestamp.Format("2006-01-02 15:04"))
		fmt.Println(commit.Message)
		fmt.Println()
	}
}
