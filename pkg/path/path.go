package path

import (
	"os"
	"path"
	"path/filepath"
)

// CreatePath creates a local directory. If name is the destination name, it
// creates the directory inside the current path. Otherwise, it expects a
// complete path.
//
// If the directory already exists, it does not return an error.
func CreatePath(name string) (string, error) {
	if !path.IsAbs(name) {
		p, err := getFullPath(name)
		if err != nil {
			return "", err
		}

		name = p
	}

	if !FindPath(name) {
		if err := os.MkdirAll(name, os.ModeDir|os.ModePerm); err != nil {
			return "", err
		}
	}

	return name, nil
}

// getFullPath returns a full path string with name at its end.
func getFullPath(name string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, name), nil
}

// FindPath checks if p is a valid local path.
func FindPath(p string) bool {
	_, err := os.Stat(p)
	return !os.IsNotExist(err)
}

// ChangeDir executes a chdir into path returning the old current working
// directory.
func ChangeDir(path string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if err := os.Chdir(path); err != nil {
		return "", err
	}

	return cwd, nil
}
