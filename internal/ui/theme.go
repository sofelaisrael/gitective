package ui

import "github.com/charmbracelet/lipgloss"

var (
	ColorFG     = lipgloss.Color("#FAFAFA")
	ColorBG     = lipgloss.Color("#1A1A1E")
	BannerStyle = lipgloss.NewStyle().
			Foreground(ColorFG).
			Background(ColorBG).
			Border(lipgloss.ThickBorder()).
			BorderForeground(ColorFG).
			Padding(0, 2).
			MarginBottom(1).
			Bold(true)
)
