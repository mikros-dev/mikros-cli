package golang

import (
	"github.com/somatech1/mikros-cli/internal/process"
)

// ModInit executes a "go mod init" at the current working directory.
func ModInit(name string) error {
	_, err := process.Exec("go", "mod", "init", name)
	return err
}
