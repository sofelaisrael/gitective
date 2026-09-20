package git

import (
	"os/exec"
	"strings"
)

func GetDiff(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "diff", "--staged")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	diff := strings.TrimSpace(string(out))
	if diff != "" {
		return diff, nil
	}

	cmd2 := exec.Command("git", "-C", repoPath, "diff")
	out2, err := cmd2.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out2)), nil
}

func GetDiffStat(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "diff", "--stat")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func GetStagedDiffStat(repoPath string) (string, error) {
	cmd := exec.Command("git", "-C", repoPath, "diff", "--staged", "--stat")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
