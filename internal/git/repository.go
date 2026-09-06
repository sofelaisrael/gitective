package git

import (
	"os"
	"path/filepath"
)

func FindRepository(path string) (string, error) {
	current, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	for {
		gitPath := filepath.Join(current, ".git")

		if _, err := os.Stat(gitPath); err == nil {
			return current, nil
		}

		parent := filepath.Dir(current)

		if parent == current {
			break
		}

		current = parent
	}

	return "", os.ErrNotExist
}
