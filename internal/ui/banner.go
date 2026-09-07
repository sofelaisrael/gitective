package ui

import (
	"strings"

	figure "github.com/common-nighthawk/go-figure"
)

func RenderBanner(text string) string {
	if text == "" {
		text = "Gitective"
	}
	raw := figure.NewFigure(text, "slant", true).String()
	raw = strings.TrimSuffix(raw, "\n")
	return BannerStyle.Render(raw)
}
