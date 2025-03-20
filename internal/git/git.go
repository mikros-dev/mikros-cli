package git

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mikros-dev/mikros-cli/internal/process"
)

type Git struct {
	Name         string
	RootPath     string
	isRepository bool
}

func LoadFromCwd() (*Git, error) {
	tmp, err := process.Exec("git", "rev-parse", "--git-dir")
	if err != nil {
		return &Git{isRepository: false}, nil
	}

	var (
		gitDir   = strings.TrimSuffix(string(tmp), "\n")
		rootPath = filepath.Dir(gitDir)
	)

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if gitDir == ".git" {
		rootPath = cwd
	}

	return &Git{
		RootPath:     rootPath,
		Name:         filepath.Base(rootPath),
		isRepository: true,
	}, nil
}

func (g *Git) IsValidRepository() bool {
	return g.isRepository
}
