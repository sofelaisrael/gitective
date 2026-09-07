package ui

import (
	"errors"

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
