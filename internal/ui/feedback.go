package ui

import (
	"errors"

	"github.com/pterm/pterm"
	"github.com/sofelaisrael/gitective/internal/git"
)

// Spinner starts a pterm spinner with the given message.
func Spinner(msg string) (*pterm.SpinnerPrinter, error) { return pterm.DefaultSpinner.Start(msg) }

// Success prints a success message via pterm.
func Success(msg string) { pterm.Success.Println(msg) }

// Successf prints a formatted success message via pterm.
func Successf(format string, a ...any) { pterm.Success.Printf(format, a...) }

// Info prints an info message via pterm.
func Info(msg string) { pterm.Info.Println(msg) }

// HandleRepoError checks if err wraps git.ErrNotRepository.
// If it does, it prints a user-friendly error and returns true.
// Otherwise it returns false. It is intended for Brutalist C error handling.
func HandleRepoError(err error) bool {
	if errors.Is(err, git.ErrNotRepository) {
		pterm.Error.Println("Not a git repository (or any parent). Run `git init` to initialize.")
		return true
	}
	return false
}
