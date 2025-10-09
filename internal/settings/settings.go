package settings

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/huh"
	"github.com/creasty/defaults"

	"github.com/mikros-dev/mikros-cli/internal/path"
)

const (
	settingsFilename = "$HOME/.mikros/config.toml"
)

// Settings represents the configuration structure file.
type Settings struct {
	Paths   Path               `toml:"paths"`
	UI      UI                 `toml:"ui"`
	App     Profile            `toml:"app"`
	Profile map[string]Profile `toml:"profile"`
}

// Path represents a configuration structure related to plugin directories.
type Path struct {
	Plugins Plugins `toml:"plugins"`
}

// Plugins represents the configuration structure for plugin directory paths.
type Plugins struct {
	// Services specifies the path to the services plugin directory.
	Services string `toml:"services" default:"$HOME/.mikros/plugins/services"`

	// Features specifies the path to the features plugin directory.
	Features string `toml:"features" default:"$HOME/.mikros/plugins/features"`
}

// Profile represents a configuration structure tied to a specific project.
type Profile struct {
	Project Project `toml:"project"`
}

// Project represents configuration details for a project, including protobuf
// monorepo and template definitions.
type Project struct {
	ProtobufMonorepo ProtobufMonorepo `toml:"protobuf_monorepo"`
	Templates        Templates        `toml:"templates"`
}

// ProtobufMonorepo represents the configuration for a protobuf monorepo.
type ProtobufMonorepo struct {
	RepositoryName string `toml:"repository_name" default:"protobuf-workspace"`
	ProjectName    string `toml:"project_name" default:"services"`
	VcsPath        string `toml:"vcs_path" default:"github.com/your-organization"`
}

// Templates represents configuration for defining template-specific settings
// for protobuf generation.
type Templates struct {
	Protobuf ProtobufTemplates `toml:"protobuf"`
}

// ProtobufTemplates represents configuration for Protobuf template-specific
// settings.
type ProtobufTemplates struct {
	CustomAuthName string `toml:"custom_auth_name" default:"scopes"`
}

// UI represents the configuration for user interface settings.
type UI struct {
	// Theme specifies the theme to be applied to the UI.
	Theme string `toml:"theme"`

	// Accessible indicates whether accessibility features are enabled.
	Accessible bool `toml:"accessible"`
}

// Load initializes and retrieves the application settings using default
// settings or loading them from a configuration file if available.
func Load() (*Settings, error) {
	cfg, err := NewDefault()
	if err != nil {
		return nil, err
	}

	if name, ok := FileExists(); ok {
		if _, err := toml.DecodeFile(name, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

// NewDefault initializes a Settings instance with default values and applies
// environmental variable expansion.
func NewDefault() (*Settings, error) {
	cfg := &Settings{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	cfg.Paths.Plugins.Services = os.ExpandEnv(cfg.Paths.Plugins.Services)
	cfg.Paths.Plugins.Features = os.ExpandEnv(cfg.Paths.Plugins.Features)

	return cfg, nil
}

// FileExists checks if the settings file exists and returns its expanded path.
func FileExists() (string, bool) {
	name := os.ExpandEnv(settingsFilename)
	return name, path.FindPath(name)
}

// Write saves the current Settings instance to a file at the predefined location
// using the TOML format.
func (s *Settings) Write() error {
	var basePath = os.ExpandEnv(settingsFilename)

	if _, err := path.CreatePath(filepath.Dir(basePath)); err != nil {
		return err
	}

	file, err := os.Create(os.ExpandEnv(settingsFilename))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	en := toml.NewEncoder(file)
	return en.Encode(s)
}

// GetTheme returns the appropriate theme based on the UI.Theme value,
// defaulting to a base theme if no match is found.
func (s *Settings) GetTheme() *huh.Theme {
	switch strings.ToLower(s.UI.Theme) {
	case "charm":
		return huh.ThemeCharm()
	case "dracula":
		return huh.ThemeDracula()
	case "catppuccin":
		return huh.ThemeCatppuccin()
	case "base16":
		return huh.ThemeBase16()
	}

	return huh.ThemeBase()
}

// Hash computes the SHA-256 hash of the Settings instance serialized in TOML
// format and returns it as a hex string.
func (s *Settings) Hash() (string, error) {
	b, err := toml.Marshal(s)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
