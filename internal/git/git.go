package git

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mikros-dev/mikros-cli/internal/process"
)

// Git represents a Git repository with basic metadata.
type Git struct {
	Name         string
	RootPath     string
	isRepository bool
}

// LoadFromCwd identifies if the current directory is part of a Git repository
// and retrieves its metadata. If the directory is not part of a repository,
// it returns a valid object with a proper flag indicating this information.
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

// Init initializes a Git repository in the current working directory.
func Init() (*Git, error) {
	if _, err := process.Exec("git", "init"); err != nil {
		return nil, err
	}

	return LoadFromCwd()
}

// IsValidRepository returns true if the Git instance represents a valid
// Git repository, otherwise false.
func (g *Git) IsValidRepository() bool {
	return g.isRepository
}
