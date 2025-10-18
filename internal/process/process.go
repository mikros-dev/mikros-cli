package process

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
)

// Exec executes a known command locally.
func Exec(args ...string) ([]byte, error) {
	if len(args) == 0 {
		return nil, errors.New("can't execute a nil command")
	}

	cmd := exec.Command(args[0], args[1:]...)
	return cmd.CombinedOutput()
}

// ExecWithPTY executes a command in a pseudo-terminal (PTY) and captures its
// output as a byte slice.
func ExecWithPTY(ctx context.Context, args ...string) ([]byte, int, error) {
	if len(args) == 0 {
		return nil, -1, errors.New("não é possível executar um comando vazio")
	}

	cmd := exec.CommandContext(ctx, args[0], args[1:]...)

	// Start with a pseudo-terminal so we can have pretty outputs
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return nil, -1, err
	}
	defer func(ptmx *os.File) {
		_ = ptmx.Close()
	}(ptmx)

	// Capture output from PTY
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, ptmx)

	// Wait for the command to complete
	err = cmd.Wait()
	code := -1

	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}

	return buf.Bytes(), code, err
}
