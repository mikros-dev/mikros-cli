package plugin

import (
	"context"
	"os"

	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func GetNewServiceKinds(cfg *settings.Settings) ([]string, error) {
	var (
		ctx      = context.Background()
		basePath = os.ExpandEnv(cfg.Paths.Services)
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
		if err := p.Start(); err != nil {
			return nil, err
		}

		newType, err := p.GetKind(ctx)
		if err != nil {
			return nil, err
		}
		types = append(types, newType)

		if err := p.Stop(ctx); err != nil {
			return nil, err
		}
	}

	return types, nil
}

func GetFeaturesUINames(cfg *settings.Settings) ([]string, error) {
	var (
		ctx      = context.Background()
		basePath = os.ExpandEnv(cfg.Paths.Features)
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
		if !path.IsExecutable(file.Name()) {
			continue
		}

		p := client.NewFeature(basePath, file.Name())
		if err := p.Start(); err != nil {
			return nil, err
		}

		newName, err := p.GetUIName(ctx)
		if err != nil {
			return nil, err
		}
		names = append(names, newName)

		if err := p.Stop(ctx); err != nil {
			return nil, err
		}
	}

	return names, nil
}

func GetServicePlugin(ctx context.Context, cfg *settings.Settings, name string) (*client.Service, error) {
	var (
		basePath = os.ExpandEnv(cfg.Paths.Services)
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
		if err := p.Start(); err != nil {
			return nil, err
		}

		pluginName, err := p.GetName(ctx)
		if err != nil {
			return nil, err
		}
		if pluginName == name {
			// Now it's the caller responsibility to Stop the plugin
			return p, nil
		}

		if err := p.Stop(ctx); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func GetFeaturePlugin(ctx context.Context, cfg *settings.Settings, name string) (*client.Feature, error) {
	var (
		basePath = os.ExpandEnv(cfg.Paths.Features)
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

		p := client.NewFeature(basePath, file.Name())
		if err := p.Start(); err != nil {
			return nil, err
		}

		uiName, err := p.GetUIName(ctx)
		if err != nil {
			return nil, err
		}
		if uiName == name {
			// Now it's the caller responsibility to Stop the plugin
			return p, nil
		}

		if err := p.Stop(ctx); err != nil {
			return nil, err
		}
	}

	return nil, nil
}
