package themes

import (
	"fmt"

	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

type RenaissanceTheme struct{}

func (RenaissanceTheme) Name() string { return "renaissance" }

func (RenaissanceTheme) Render(facts commit.CommitFacts) string {
	p := personality.Classify(facts)
	var verseName string
	switch p {
	case personality.PanicPatch:
		verseName = "THE HASTY MENDER"
	case personality.NuclearJanitor:
		verseName = "THE GRAND PURGATION"
	case personality.Janitor:
		verseName = "THE SCRIBE OF DUST"
	case personality.MadScientist:
		verseName = "THE MASTER ARTIFICER"
	case personality.Architect:
		verseName = "THE CATHEDRAL BUILDER"
	case personality.SilentCommit:
		verseName = "THE SILENT VOW"
	case personality.Shipper:
		verseName = "THE CARAVAN"
	}
	return fmt.Sprintf("⚜ THE CHRONICLES [%s]\n\nUpon this day, keeper %s hath wrought '%s'.\nScrolls: %d  +%d/-%d\nThus endureth.", verseName, facts.Author, facts.Message, facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted)
}
