package plugin

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

// GetNewServiceKinds returns the list of new service kinds available in the
// plugins directory.
func GetNewServiceKinds(cfg *settings.Settings) ([]string, error) {
	var (
		basePath = cfg.Paths.Plugins.Services
		types    []string
	)

	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := listExecutableFiles(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p := client.NewService(basePath, file)

		newType, err := p.GetKind()
		if err != nil {
			return nil, err
		}
		types = append(types, newType)
	}

	return types, nil
}

// GetFeaturesUINames returns the list of feature names available in the
// plugins directory.
func GetFeaturesUINames(cfg *settings.Settings) ([]string, error) {
	var (
		basePath = cfg.Paths.Plugins.Features
		names    []string
	)

	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := listExecutableFiles(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p := client.NewFeature(basePath, file)

		newName, err := p.GetUIName()
		if err != nil {
			return nil, err
		}
		names = append(names, newName)
	}

	return names, nil
}

// GetServicePlugin returns the plugin for the given kind.
func GetServicePlugin(cfg *settings.Settings, kind string) (*client.Service, error) {
	var basePath = cfg.Paths.Plugins.Services
	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := listExecutableFiles(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p := client.NewService(basePath, file)

		pluginKind, err := p.GetKind()
		if err != nil {
			return nil, err
		}
		if pluginKind == kind {
			return p, nil
		}
	}

	return nil, nil
}

// GetFeaturePlugin returns the plugin for the given name.
func GetFeaturePlugin(cfg *settings.Settings, name string) (*client.Feature, error) {
	var basePath = cfg.Paths.Plugins.Features
	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := listExecutableFiles(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p := client.NewFeature(basePath, file)

		uiName, err := p.GetUIName()
		if err != nil {
			return nil, err
		}
		if uiName == name {
			return p, nil
		}
	}

	return nil, nil
}

func listExecutableFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if path.IsExecutable(filepath.Join(dir, entry.Name())) {
			files = append(files, entry.Name())
		}
	}

	sort.Strings(files)
	return files, nil
}
