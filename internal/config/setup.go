package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

var styleOptions = []struct {
	Name string
	Desc string
}{
	{"cyberpunk-commit", "Neon operator reports — Blade Runner / CP2077 vibes"},
	{"renaissance-commit", "Medieval chronicle prose — gilded manuscripts"},
	{"shakespeare", "Iambic pentameter — theatrical and grand"},
	{"pirate", "Nautical swashbuckler — sea shanties and treasure"},
	{"victorian", "Proper English — leather-bound and gaslit"},
	{"corporate", "Synergy speak — jargon-filled and professional"},
	{"genz", "No cap fam — TikTok energy and vibes"},
	{"noir-detective", "Hardboiled detective noir — rain, smoke, and cold cases"},
}

func RunSetup() error {
	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Println("  GITECTIVE — FIRST TIME SETUP  ")
	pterm.Println()

	styleNames := make([]string, len(styleOptions))
	for i, s := range styleOptions {
		styleNames[i] = fmt.Sprintf("%s — %s", s.Name, s.Desc)
	}

	pterm.Info.Println("Pick your default style:")
	pterm.Println()
	chosenDisplay, err := pterm.DefaultInteractiveSelect.
		WithOptions(styleNames).
		WithDefaultOption(styleNames[0]).
		Show()
	if err != nil {
		return fmt.Errorf("style select: %w", err)
	}

	chosenStyle := ""
	for _, s := range styleOptions {
		if fmt.Sprintf("%s — %s", s.Name, s.Desc) == chosenDisplay {
			chosenStyle = s.Name
			break
		}
	}
	if chosenStyle == "" {
		chosenStyle = "cyberpunk-commit"
	}
	pterm.Success.Printf("Style: %s\n\n", chosenStyle)

	pterm.Info.Println("Default intensity? (0.1 = subtle, 1.0 = maximum):")
	intensityStr, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("0.8").
		Show()
	if err != nil {
		return fmt.Errorf("intensity input: %w", err)
	}
	intensityStr = strings.TrimSpace(intensityStr)
	if intensityStr == "" {
		intensityStr = "0.8"
	}
	intensity, err := strconv.ParseFloat(intensityStr, 64)
	if err != nil || intensity < 0.1 || intensity > 1.0 {
		intensity = 0.8
	}
	pterm.Success.Printf("Intensity: %.0f%%\n\n", intensity*100)

	pterm.Info.Println("Style-engine URL?")
	urlStr, err := pterm.DefaultInteractiveTextInput.
		WithDefaultText("http://localhost:8080").
		Show()
	if err != nil {
		return fmt.Errorf("url input: %w", err)
	}
	urlStr = strings.TrimSpace(urlStr)
	if urlStr == "" {
		urlStr = "http://localhost:8080"
	}
	pterm.Success.Printf("URL: %s\n\n", urlStr)

	cfg := Config{
		DefaultStyle:     chosenStyle,
		DefaultIntensity: intensity,
		StyleEngineURL:   urlStr,
	}

	if err := Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}

	pterm.Success.Println("Config saved! Run 'gitective' to start.")
	pterm.Println()
	return nil
}

func IsConfigured() bool {
	path, err := configPath()
	if err != nil {
		return false
	}
	_, err = os.ReadFile(path)
	return err == nil
}
