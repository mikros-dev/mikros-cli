package config

import (
	"fmt"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

// CreateDefaultSettings creates a default settings file if it does not
// already exist, initializing it with default values.
func CreateDefaultSettings() error {
	if _, ok := settings.FileExists(); ok {
		fmt.Println("settings file already exists")
		return nil
	}

	cfg, err := settings.NewDefault()
	if err != nil {
		return err
	}

	return cfg.Write()
}
