package themes

import (
	"fmt"

	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

type CyberpunkTheme struct{}

func (CyberpunkTheme) Name() string { return "cyberpunk" }

func (CyberpunkTheme) Render(facts commit.CommitFacts) string {
	switch personality.Classify(facts) {
	case personality.PanicPatch:
		return fmt.Sprintf("// EMERGENCY PATCH\nOperator: %s\nBreach: %s\nSectors: %d  +%d/-%d\nAction: hot-patched at %02d:00\nStatus: CONTAINED.", facts.Author, facts.Message, facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted, facts.Timestamp.Hour())
	case personality.NuclearJanitor:
		return fmt.Sprintf("// PURGE PROTOCOL\nOperator: %s\nErased: %d lines from %d sectors\nPayload: %s\nStatus: GRID CLEAN. No witnesses.", facts.Author, facts.LinesDeleted, facts.FilesChanged, facts.Message)
	case personality.MadScientist:
		return fmt.Sprintf("// NEO-FORGE\nOperator: %s\nForged: %d sectors + %d lines\nConstruct: %s\nStatus: UNSTABLE. Proceed.", facts.Author, facts.FilesChanged, facts.LinesAdded, facts.Message)
	case personality.SilentCommit:
		return fmt.Sprintf("// GHOST SIGNAL\nOperator: %s\nTransmission: \"%s\"\nSectors: %d\nStatus: NO PAYLOAD. Ghosted.", facts.Author, facts.Message, facts.FilesChanged)
	default:
		return fmt.Sprintf("// FREIGHT RUN\nOperator: %s\nPayload: %s\nSectors: %d  +%d/-%d\nStatus: DELIVERED.", facts.Author, facts.Message, facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted)
	}
}
