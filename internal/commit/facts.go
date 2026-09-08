package commit

import (
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type CommitFacts struct {
	Hash         string
	Message      string
	Author       string
	FilesChanged int
	LinesAdded   int
	LinesDeleted int
	Files        []string
	Timestamp    time.Time
	IsFix        bool
}

func ExtractFacts(repoPath string, hash string) (CommitFacts, error) {
	facts := CommitFacts{Hash: hash}
	cmd := exec.Command("git", "-C", repoPath, "show", "-s", "--format=%s%x1f%an%x1f%aI", hash)
	out, err := cmd.Output()
	if err != nil {
		return facts, err
	}
	parts := strings.Split(strings.TrimSpace(string(out)), "\x1f")
	if len(parts) >= 1 {
		facts.Message = parts[0]
	}
	if len(parts) >= 2 {
		facts.Author = parts[1]
	}
	if len(parts) >= 3 {
		if t, err := time.Parse(time.RFC3339, parts[2]); err == nil {
			facts.Timestamp = t
		}
	}
	if len(parts) >= 1 {
		facts.IsFix = strings.Contains(strings.ToLower(parts[0]), "fix")
	}
	cmd2 := exec.Command("git", "-C", repoPath, "show", "--numstat", "--format=", hash)
	out2, err := cmd2.Output()
	if err != nil {
		return facts, err
	}
	lines := strings.Split(strings.TrimSpace(string(out2)), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		added, _ := strconv.Atoi(fields[0])
		deleted, _ := strconv.Atoi(fields[1])
		if fields[0] != "-" {
			facts.LinesAdded += added
		}
		if fields[1] != "-" {
			facts.LinesDeleted += deleted
		}
		facts.FilesChanged++
		facts.Files = append(facts.Files, fields[2])
	}
	return facts, nil
}
