package llm

import (
	"fmt"
	"strings"
	"time"

	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

func BuildContext(facts commit.CommitFacts) string {
	var b strings.Builder

	// Time context
	hour := facts.Timestamp.In(time.Local).Hour()
	var timeOfDay string
	switch {
	case hour >= 5 && hour < 12:
		timeOfDay = "morning"
	case hour >= 12 && hour < 17:
		timeOfDay = "afternoon"
	case hour >= 17 && hour < 22:
		timeOfDay = "evening"
	default:
		timeOfDay = "dead of night"
	}

	// Commit type from personality
	p := personality.Classify(facts)
	var commitType string
	switch p {
	case personality.PanicPatch:
		commitType = "bugfix (emergency)"
	case personality.NuclearJanitor:
		commitType = "massive cleanup"
	case personality.MadScientist:
		commitType = "feature (no docs/tests)"
	case personality.Architect:
		commitType = "feature (with docs and tests)"
	case personality.SilentCommit:
		commitType = "cryptic commit"
	case personality.Janitor:
		commitType = "maintenance"
	default:
		commitType = "standard commit"
	}

	// Build the context prompt
	b.WriteString(fmt.Sprintf("Git commit by %s at %s (%s).\n",
		facts.Author,
		facts.Timestamp.In(time.Local).Format("3:04 PM"),
		timeOfDay,
	))
	b.WriteString(fmt.Sprintf("Type: %s. Files changed: %d (+%d/-%d lines).\n",
		commitType,
		facts.FilesChanged,
		facts.LinesAdded,
		facts.LinesDeleted,
	))

	if len(facts.Files) > 0 {
		maxFiles := 8
		if len(facts.Files) < maxFiles {
			maxFiles = len(facts.Files)
		}
		b.WriteString(fmt.Sprintf("Files: %s", strings.Join(facts.Files[:maxFiles], ", ")))
		if len(facts.Files) > 8 {
			b.WriteString(fmt.Sprintf(" ... and %d more", len(facts.Files)-8))
		}
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("Message: '%s'", facts.Message))

	return b.String()
}
