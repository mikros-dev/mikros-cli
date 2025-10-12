package golang

import (
	"fmt"

	"github.com/mikros-dev/mikros-cli/internal/process"
)

// ModInit executes a "go mod init" at the current working directory.
func ModInit(name string) error {
	if out, err := process.Exec("go", "mod", "init", name); err != nil {
		return fmt.Errorf("go mod init: %w\n%s", err, string(out))
	}

	return nil
}
