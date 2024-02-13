package process

import (
	"errors"
	"os/exec"
)

// Exec executes a known command locally.
func Exec(args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("can't execute a nil command")
	}

	cmd := exec.Command(args[0], args[1:]...) //nolint
	return cmd.CombinedOutput()
}
