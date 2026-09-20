package themes

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/sofelaisrael/gitective/internal/commit"
	"github.com/sofelaisrael/gitective/internal/personality"
)

var (
	accentColor = lipgloss.Color("#FF3366")
	accentStyle = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	labelStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#888888"))
	valueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	boxStyle    = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF3366")).
			PaddingLeft(1).PaddingRight(1).
			MarginTop(1)
)

type CyberpunkTheme struct{}

func (CyberpunkTheme) Name() string { return "cyberpunk" }

func (CyberpunkTheme) Render(facts commit.CommitFacts) string {
	var tag, detail string

	switch personality.Classify(facts) {
	case personality.PanicPatch:
		tag = "EMERGENCY PATCH"
		detail = fmt.Sprintf("hot-patched at %02d:00", facts.Timestamp.Hour())
	case personality.NuclearJanitor:
		tag = "PURGE PROTOCOL"
		detail = fmt.Sprintf("erased %d lines from %d sectors", facts.LinesDeleted, facts.FilesChanged)
	case personality.MadScientist:
		tag = "NEO-FORGE"
		detail = fmt.Sprintf("forged %d sectors + %d lines", facts.FilesChanged, facts.LinesAdded)
	case personality.SilentCommit:
		tag = "GHOST SIGNAL"
		detail = "no payload"
	default:
		tag = "FREIGHT RUN"
		detail = fmt.Sprintf("%d sectors  +%d/-%d", facts.FilesChanged, facts.LinesAdded, facts.LinesDeleted)
	}

	header := accentStyle.Render(fmt.Sprintf("// %s", tag))

	left := fmt.Sprintf("%s  %s\n%s  %s\n%s  %s",
		labelStyle.Render("Operator"), valueStyle.Render(facts.Author),
		labelStyle.Render("Target"), valueStyle.Render(truncate(facts.Message, 40)),
		labelStyle.Render("Detail"), dimStyle.Render(detail),
	)

	left = lipgloss.NewStyle().MarginLeft(0).Render(left)

	bottom := accentStyle.Render("■ ") + dimStyle.Render("STATUS: ") + valueStyle.Render(statusWord(personality.Classify(facts)))

	return boxStyle.Render(fmt.Sprintf("%s\n\n%s\n\n%s", header, left, bottom))
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func statusWord(p personality.Personality) string {
	switch p {
	case personality.PanicPatch:
		return "CONTAINED"
	case personality.NuclearJanitor:
		return "GRID CLEAN"
	case personality.MadScientist:
		return "UNSTABLE"
	case personality.SilentCommit:
		return "GHOSTED"
	default:
		return "DELIVERED"
	}
}
