package main

import (
	"fmt"
	"os"

	"github.com/pterm/pterm"
	commitPkg "github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/git"
	"github.com/sofelaisrael/gitective/internal/themes"
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
	t := themes.CyberpunkTheme{}
	fmt.Println(t.Render(facts))
}
