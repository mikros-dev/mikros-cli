package setup

import (
	"fmt"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func CreateDefaultSettings() error {
	if _, ok := settings.FileExists(); ok {
		fmt.Println("settings file already exists")
		return nil
	}

	cfg, err := settings.NewDefault()
	if err != nil {
		return err
	}

	if err := cfg.Write(); err != nil {
		return err
	}

	return nil
}
