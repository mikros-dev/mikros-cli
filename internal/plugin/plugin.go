package plugin

import (
	"os"
	"path/filepath"

	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func GetNewServiceKinds(cfg *settings.Settings) ([]string, error) {
	var (
		basePath = cfg.Paths.Plugins.Services
		types    []string
	)

	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !path.IsExecutable(file.Name()) {
			continue
		}

		p := client.NewService(basePath, file.Name())

		newType, err := p.GetKind()
		if err != nil {
			return nil, err
		}
		types = append(types, newType)
	}

	return types, nil
}

func GetFeaturesUINames(cfg *settings.Settings) ([]string, error) {
	var (
		basePath = cfg.Paths.Plugins.Features
		names    []string
	)

	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !path.IsExecutable(filepath.Join(basePath, file.Name())) {
			continue
		}

		p := client.NewFeature(basePath, file.Name())

		newName, err := p.GetUIName()
		if err != nil {
			return nil, err
		}
		names = append(names, newName)
	}

	return names, nil
}

func GetServicePlugin(cfg *settings.Settings, kind string) (*client.Service, error) {
	var (
		basePath = cfg.Paths.Plugins.Services
	)

	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !path.IsExecutable(file.Name()) {
			continue
		}

		p := client.NewService(basePath, file.Name())

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

func GetFeaturePlugin(cfg *settings.Settings, name string) (*client.Feature, error) {
	var (
		basePath = cfg.Paths.Plugins.Features
	)

	if !path.FindPath(basePath) {
		return nil, nil
	}

	files, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !path.IsExecutable(filepath.Join(basePath, file.Name())) {
			continue
		}

		p := client.NewFeature(basePath, file.Name())

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
