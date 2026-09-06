package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	figure "github.com/common-nighthawk/go-figure"
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"github.com/sofelaisrael/gitective/internal/git"
)

func main() {
	rawFig := figure.NewFigure("Gitective", "slant", true).String()
	raw := strings.TrimSuffix(rawFig, "\n")
	styled := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7D56F4")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5E3FD0")).
		Padding(0, 2).
		MarginBottom(1).
		Bold(true).
		Render(raw)
	fmt.Println(styled)

	spinner, _ := pterm.DefaultSpinner.Start("Detecting git repository...")
	root, err := git.FindRepository(".")
	if err != nil {
		spinner.Fail("Repository check failed")
		if errors.Is(err, git.ErrNotRepository) {
			pterm.Error.Println("Not a git repository (or any parent). Run `git init` to initialize.")
			os.Exit(1)
		}
		pterm.Error.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	spinner.Success("Repository found")
	pterm.Success.Printf("Repository: %s\n", root)
	pterm.Info.Println("Investigating your Git history...")
}
