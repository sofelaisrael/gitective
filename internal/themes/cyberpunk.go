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
	return fmt.Sprintf("⚡ NIGHT//PROTOCOL\n\nTarget: %s\nAuthor: %s\nFiles: %d  +%d/-%d\nFiles list: %v\nPersonality: %s\nStatus: SYSTEM SECURED.", facts.Message, facts.Author, facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted, facts.Files, p)
}
