package ui

import (
	"errors"

	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"github.com/sofelaisrael/gitective/internal/git"
)

func Spinner(msg string) (*pterm.SpinnerPrinter, error) { return pterm.DefaultSpinner.Start(msg) }

func Success(msg string) { pterm.Success.Println(msg) }

func Successf(format string, a ...any) { pterm.Success.Printf(format, a...) }

func Info(msg string) { pterm.Info.Println(msg) }

func HandleRepoError(err error) bool {
	if errors.Is(err, git.ErrNotRepository) {
		pterm.Error.Println("Not a git repository (or any parent). Run `git init` to initialize.")
		return true
	}
	return false
}

func RenderBox(header, body, footer string) string {
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF3366")).
		PaddingLeft(1).PaddingRight(1).
		MarginTop(1)

	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FF3366")).Bold(true)
	dimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))

	content := accentStyle.Render(header) + "\n\n" + body + "\n\n" + dimStyle.Render(footer)
	return box.Render(content)
}
