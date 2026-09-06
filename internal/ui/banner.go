// Package ui provides Brutalist C theme styles for gitective.
package ui

import (
	"strings"

	figure "github.com/common-nighthawk/go-figure"
)

// RenderBanner renders text as a styled ASCII art banner.
//
// It uses the "slant" figure font via go-figure to generate ASCII art.
// If text is empty, it defaults to "Gitective". The raw figure output
// is trimmed of its trailing newline and then styled with [BannerStyle]
// before being returned as a string.
//
// RenderBanner is pure and has no side effects: it does not write to
// stdout/stderr, does not depend on os or pterm, and is safe to call
// from tests.
func RenderBanner(text string) string {
	if text == "" {
		text = "Gitective"
	}
	raw := figure.NewFigure(text, "slant", true).String()
	raw = strings.TrimSuffix(raw, "\n")
	return BannerStyle.Render(raw)
}
