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
	ui.Successf("Repository: %s\n", root)
	ui.Info("Investigating your Git history...")
}
