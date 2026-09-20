package themes

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

var (
	renaColor  = lipgloss.Color("#C9A96E")
	renaAccent = lipgloss.NewStyle().Foreground(renaColor).Bold(true)
	renaDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	renaLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA"))
	renaValue  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	renaBox    = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(renaColor).
			PaddingLeft(1).PaddingRight(1).
			MarginTop(1)
)

type RenaissanceTheme struct{}

func (RenaissanceTheme) Name() string { return "renaissance" }

func (RenaissanceTheme) Render(facts commit.CommitFacts) string {
	var tag string

	switch personality.Classify(facts) {
	case personality.PanicPatch:
		tag = "THE HASTY MENDER"
	case personality.NuclearJanitor:
		tag = "THE GRAND PURGATION"
	case personality.Janitor:
		tag = "THE SCRIBE OF DUST"
	case personality.MadScientist:
		tag = "THE MASTER ARTIFICER"
	case personality.Architect:
		tag = "THE CATHEDRAL BUILDER"
	case personality.SilentCommit:
		tag = "THE SILENT VOW"
	default:
		tag = "THE CARAVAN"
	}

	header := renaAccent.Render(fmt.Sprintf("⚜ %s", tag))

	body := fmt.Sprintf("%s  %s\n%s  '%s'\n%s  %d scrolls  +%d/-%d",
		renaLabel.Render("Keeper"), renaValue.Render(facts.Author),
		renaLabel.Render("Wrought"), renaValue.Render(facts.Message),
		renaLabel.Render("Scribe"), facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted,
	)

	bottom := renaAccent.Render("■ ") + renaDim.Render("Thus endureth.")

	return renaBox.Render(fmt.Sprintf("%s\n\n%s\n\n%s", header, body, bottom))
}
