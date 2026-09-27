package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Settings are the application options, the same for all host lists.
// Everything extra is off by default, so the app stays simple for those who do not need it.
type Settings struct {
	ShowMin    bool
	ShowJitter bool
	ShowSince  bool

	AlwaysOnTop bool

	// Language of the interface as a tag like "ru"; "" - the system's
	Language string `json:",omitempty"`

	// Color theme: "light"; "" - the dark one
	Theme string `json:",omitempty"`
}

// Limits of the per-host options, see ConfigHost
const (
	DefaultIntervalMs = 1000
	MinIntervalMs     = 100
	MaxIntervalMs     = 60000
	MaxSlowMs         = 10000

	DefaultTimeoutMs = 1000
	MinTimeoutMs     = 100
	MaxTimeoutMs     = 10000
)

var (
	settingsMtx sync.Mutex
	settings    = defaultSettings()
)

func defaultSettings() Settings {
	return Settings{}
}

func settingsPath() string {
	return filepath.Join(ConfigDirectory(), "settings.json")
}

// LoadSettings reads the settings file; missing values get the defaults
func LoadSettings() {
	s := defaultSettings()
	if bs, err := os.ReadFile(settingsPath()); err == nil {
		_ = json.Unmarshal(bs, &s)
	}
	settingsMtx.Lock()
	settings = s
	settingsMtx.Unlock()
}

// GetSettings returns the current settings; safe to call from any goroutine
func GetSettings() Settings {
	settingsMtx.Lock()
	defer settingsMtx.Unlock()
	return settings
}

// SetSettings applies and saves the settings
func SetSettings(s Settings) error {
	settingsMtx.Lock()
	settings = s
	settingsMtx.Unlock()

	bs, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(ConfigDirectory(), 0755); err != nil {
		return err
	}
	return os.WriteFile(settingsPath(), bs, 0644)
}
