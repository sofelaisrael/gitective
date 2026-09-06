package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var ErrNotRepository = errors.New("not a git repository")

func FindRepository(path string) (string, error) {
	if path == "" {
		path = "."
	}

	current, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve path %q: %w", path, err)
	}

	for {
		gitPath := filepath.Join(current, ".git")

		if _, err := os.Stat(gitPath); err == nil {
			return current, nil
		} else if err != nil && !os.IsNotExist(err) {
			return "", fmt.Errorf("stat %s: %w", gitPath, err)
		}

		parent := filepath.Dir(current)

		if parent == current {
			break
		}

		current = parent
	}

	return "", fmt.Errorf("not a git repository (searched from %s): %w", path, ErrNotRepository)
}
