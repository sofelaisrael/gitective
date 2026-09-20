package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type CommitStyle struct {
	Style     string
	Intensity float64
}

func ReadTrailers(repoPath, hash string) CommitStyle {
	cmd := exec.Command("git", "-C", repoPath, "show", "-s", "--format=%b", hash)
	out, err := cmd.Output()
	if err != nil {
		return CommitStyle{}
	}

	cs := CommitStyle{Style: "cyberpunk-commit", Intensity: 0.8}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Gitective-Style:") {
			cs.Style = strings.TrimSpace(strings.TrimPrefix(line, "Gitective-Style:"))
		} else if strings.HasPrefix(line, "Gitective-Intensity:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "Gitective-Intensity:"))
			val = strings.TrimSuffix(val, "%")
			var f float64
			for _, ch := range val {
				if ch >= '0' && ch <= '9' {
					f = f*10 + float64(ch-'0')
				} else if ch == '.' {
				}
			}
			if f > 0 && f <= 1 {
				cs.Intensity = f
			} else if f > 1 && f <= 100 {
				cs.Intensity = f / 100
			}
		}
	}
	return cs
}

func AppendTrailers(message, style string, intensity float64) string {
	pct := int(intensity * 100)
	intensityStr := fmt.Sprintf("%d.%02d", pct/100, pct%100)
	intensityStr = strings.TrimRight(intensityStr, "0")
	intensityStr = strings.TrimRight(intensityStr, ".")
	if intensityStr == "" {
		intensityStr = "0.8"
	}

	trailer := "\n\nGitective-Style: " + style + "\nGitective-Intensity: " + intensityStr
	return message + trailer
}
