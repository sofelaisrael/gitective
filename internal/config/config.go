package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultStyle     string  `json:"default_style"`
	DefaultIntensity float64 `json:"default_intensity"`
	StyleEngineURL   string  `json:"style_engine_url"`
}

var defaults = Config{
	DefaultStyle:     "cyberpunk-commit",
	DefaultIntensity: 0.8,
	StyleEngineURL:   "http://localhost:8080",
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "gitective", "config.json"), nil
}

func Load() Config {
	path, err := configPath()
	if err != nil {
		return defaults
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return defaults
	}
	cfg := defaults
	json.Unmarshal(data, &cfg)
	return cfg
}

func Save(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
