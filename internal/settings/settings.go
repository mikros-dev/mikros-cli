package settings

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"

	"github.com/mikros-dev/mikros-cli/internal/path"
)

const (
	settingsFilename = "$HOME/.mikros/config.toml"
)

type Settings struct {
	Paths Path `toml:"paths"`
}

type Path struct {
	Services string `toml:"services" default:"$HOME/.mikros/plugins/services"`
	Features string `toml:"features" default:"$HOME/.mikros/plugins/features"`
}

func New() (*Settings, error) {
	cfg := &Settings{}
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	if name, ok := settingsFileExists(); ok {
		if _, err := toml.DecodeFile(name, cfg); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func settingsFileExists() (string, bool) {
	name := os.ExpandEnv(settingsFilename)
	return name, path.FindPath(name)
}
