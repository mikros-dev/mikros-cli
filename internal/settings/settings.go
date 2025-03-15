package settings

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"

	"github.com/mikros-dev/mikros-cli/internal/path"
)

const (
	settingsFilename = "$HOME/.mikros/config.toml"
)

type Settings struct {
	Paths   Path    `toml:"paths"`
	Project Project `toml:"project"`
}

type Path struct {
	Services string `toml:"services" default:"$HOME/.mikros/plugins/services"`
	Features string `toml:"features" default:"$HOME/.mikros/plugins/features"`
}

type Project struct {
	Template Template `toml:"template"`
}

type Template struct {
	VcsPath string `toml:"vcs_path"`
}

func New() (*Settings, error) {
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

func NewDefault() (*Settings, error) {
	cfg := &Settings{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	cfg.Paths.Services = os.ExpandEnv(cfg.Paths.Services)
	cfg.Paths.Features = os.ExpandEnv(cfg.Paths.Features)

	return cfg, nil
}

func FileExists() (string, bool) {
	name := os.ExpandEnv(settingsFilename)
	return name, path.FindPath(name)
}

func (s *Settings) Write() error {
	var (
		basePath = os.ExpandEnv(settingsFilename)
	)

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
	if err := en.Encode(s); err != nil {
		return err
	}

	return nil
}
