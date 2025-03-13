package plugin

import (
	"os"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
)

func Error(err error) {
	encoder := plugin.NewEncoder()
	encoder.SetError(err)

	// Nothing to do here but ignore the error
	_ = encoder.Output()
	os.Exit(1)
}
