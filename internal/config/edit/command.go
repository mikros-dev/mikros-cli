package edit

import (
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

// New manipulates the configuration file.
func New() error {
	cfg, err := settings.Load()
	if err != nil {
		return err
	}
	h1, err := cfg.Hash()
	if err != nil {
		return err
	}

	if err := runForm(cfg); err != nil {
		return err
	}

	h2, err := cfg.Hash()
	if err != nil {
		return err
	}

	if h1 != h2 {
		if err := confirmSave(cfg); err != nil {
			return err
		}
	}

	return nil
}
