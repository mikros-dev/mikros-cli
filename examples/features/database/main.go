package main

import (
	"github.com/mikros-dev/mikros-cli/pkg/plugin"
)

const (
	featureName   = "database"
	uiFeatureName = "nosql database"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return featureName
}

func (p *Plugin) UIName() string {
	return uiFeatureName
}

func (p *Plugin) Survey() *plugin.Survey {
	return &plugin.Survey{
		Questions: []*plugin.Question{
			{
				Name:    "database_cache",
				Message: "Use cache to optimize the queries?",
				Prompt:  plugin.PromptConfirm,
			},
			{
				Name:    "database_kind",
				Message: "Select the database kind:",
				Default: "mongo",
				Options: []string{"mongo", "postgres", "mysql", "sqlserver", "sqlite"},
				Prompt:  plugin.PromptSelect,
			},
			{
				Name:    "database_ttl",
				Message: "Enter the TTL of the entity, if it needs to be cooled:",
				Default: "0",
				Prompt:  plugin.PromptInput,
			},
			{
				Name:    "database_collections",
				Message: "Enter the name of additional collections (one by line):",
				Prompt:  plugin.PromptMultiline,
			},
		},
	}
}

func (p *Plugin) ValidateAnswers(in map[string]interface{}) (map[string]interface{}, error) {
	values := map[string]interface{}{
		"enabled":     true,
		"collections": []string{"name1", "name2"},
		"ttl":         0,
	}

	return values, nil
}

func main() {
	p, err := plugin.NewFeature(&Plugin{})
	if err != nil {
		plugin.Error(err)
	}

	if err := p.Run(); err != nil {
		plugin.Error(err)
	}
}
