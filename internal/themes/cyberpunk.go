package themes

import (
	"fmt"

	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

type CyberpunkTheme struct{}

func (CyberpunkTheme) Name() string { return "cyberpunk" }

func (CyberpunkTheme) Render(facts commit.CommitFacts) string {
	p := personality.Classify(facts)
	var verseName string
	switch p {
	case personality.PanicPatch:
		verseName = "CORPORATE TRIAGE"
	case personality.NuclearJanitor:
		verseName = "SECTOR PURGE"
	case personality.Janitor:
		verseName = "SANITATION CREW"
	case personality.MadScientist:
		verseName = "NEO-FORGE"
	case personality.Architect:
		verseName = "NEO-ARCHITECT"
	case personality.SilentCommit:
		verseName = "GHOST SIGNAL"
	case personality.Shipper:
		verseName = "FREIGHT RUNNER"
	}
	return fmt.Sprintf("⚡ NIGHT//PROTOCOL [%s]\n\nTarget: %s\nAuthor: %s\nFiles: %d  +%d/-%d\nStatus: SYSTEM SECURED.", verseName, facts.Message, facts.Author, facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted)
}
