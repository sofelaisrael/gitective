package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sofelaisrael/gitective/internal/git"
)

func main() {
	root, err := git.FindRepository(".")
	if err != nil {
		if errors.Is(err, git.ErrNotRepository) {
			fmt.Fprintln(os.Stderr, "Error: not a git repository (or any parent). Run `git init` to initialize.")
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Repository: %s\n", root)
	fmt.Println("Gitective")
	fmt.Println("Investigating your Git history...")
}
