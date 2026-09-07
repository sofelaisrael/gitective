package git

import (
	"os/exec"
	"strings"
	"time"
)

type Commit struct {
	Hash      string
	Author    string
	Message   string
	Timestamp time.Time
}

func GetCommits(repoPath string) ([]Commit, error) {
	cmd := exec.Command("git", "-C", repoPath, "log", "--format=%H%x1f%an%x1f%aI%x1f%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var commits []Commit
	for _, line := range lines {
		parts := strings.Split(line, "\x1f")
		if len(parts) != 4 {
			continue
		}
		timestamp, err := time.Parse(time.RFC3339, parts[2])
		if err != nil {
			continue
		}
		commits = append(commits, Commit{Hash: parts[0], Author: parts[1], Timestamp: timestamp, Message: parts[3]})
	}
	return commits, nil
}
