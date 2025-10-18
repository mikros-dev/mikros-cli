package lint

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"text/template"

	"github.com/charmbracelet/log"
	"github.com/go-playground/validator/v10"

	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/process"
)

const (
	lintErrorExitCode = 42
)

// Options represents configurable parameters used for code analysis and
// linting execution.
type Options struct {
	Debug   bool
	Format  string `validate:"omitempty,oneof=default friendly json ndjson stylish checkstyle"`
	Config  string
	Path    string `validate:"required"`
	Exclude []string
}

// Run executes a linter for analyzing code style and issues.
func Run(ctx context.Context, opts Options) error {
	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return err
	}
	if opts.Debug {
		log.SetLevel(log.DebugLevel)
	}

	if err := locateDependencies(); err != nil {
		return err
	}

	return executeLint(ctx, opts)
}

func locateDependencies() error {
	if _, err := path.FindBinary("revive"); err != nil {
		return err
	}

	return nil
}

func executeLint(ctx context.Context, opts Options) error {
	tool, err := path.FindBinary("revive")
	if err != nil {
		return err
	}

	// Set up the revive command call
	configPath, err := saveReviveConfig(opts.Config, "default")
	if err != nil {
		return err
	}
	defer func() {
		_ = os.Remove(configPath)
	}()

	args := []string{
		tool,
		"-config",
		configPath,
	}
	if opts.Format != "" {
		args = append(args, "-formatter", opts.Format)
	}
	if len(opts.Exclude) > 0 {
		for _, exclude := range opts.Exclude {
			args = append(args, "-exclude", exclude)
		}
	}
	args = append(args, opts.Path)

	if opts.Debug {
		log.Debug(args)
	}

	out, code, err := process.ExecWithPTY(ctx, args...)
	if err != nil {
		// We show lint errors for the user.
		if code == lintErrorExitCode {
			_, _ = os.Stdout.Write(out)
			return nil
		}

		return fmt.Errorf("revive: %w - %s", err, string(out))
	}

	return nil
}

func saveReviveConfig(config, profile string) (string, error) {
	if config != "" {
		// User custom config path
		return config, nil
	}

	data, err := loadReviveConfig(profile)
	if err != nil {
		return "", err
	}

	f, err := os.CreateTemp("", "revive-config-*.toml")
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()

	if _, err := f.Write(data); err != nil {
		return "", err
	}

	return f.Name(), nil
}

func loadReviveConfig(profile string) ([]byte, error) {
	data, err := reviveConfig.ReadFile("config/revive.toml.tmpl")
	if err != nil {
		return nil, err
	}

	t, err := template.New("config").Parse(string(data))
	if err != nil {
		return nil, err
	}

	var (
		buf bytes.Buffer
		w   = bufio.NewWriter(&buf)
		ctx = availableProfiles[profile]
	)

	if err := t.Execute(w, ctx); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
