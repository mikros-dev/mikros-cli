package edit

import (
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func New() error {
	cfg, err := settings.Load()
	if err != nil {
		return err
	}

	if err := runForm(cfg); err != nil {
		return err
	}

	return nil
}
