package themes

import (
	"fmt"

	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

type RenaissanceTheme struct{}

func (RenaissanceTheme) Name() string { return "renaissance" }

func (RenaissanceTheme) Render(facts commit.CommitFacts) string {
	return fmt.Sprintf("⚜ THE CHRONICLES\n\nUpon this day, the keeper %s hath wrought '%s'.\nScrolls altered: %d\nLines summoned: %d\nLines cast away: %d\nCharacter: %s\nThus the realm endureth.", facts.Author, facts.Message, facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted, personality.Classify(facts))
}
