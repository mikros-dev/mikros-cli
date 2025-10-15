package main

import (
	"github.com/mikros-dev/mikros-cli/pkg/plugin"
)

const (
	pluginName    = "complete-options"
	uiFeatureName = "complete options"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) UIName() string {
	return uiFeatureName
}

func (p *Plugin) Survey() *plugin.Survey {
	// We create a plugin.inside a loop here.
	return &plugin.Survey{
		ConfirmQuestion: &plugin.Question{
			ConfirmAfter: true,
			Message:      "Do you want to execute the form again?",
			Default:      "true",
		},
		Questions: []*plugin.Question{
			{
				Name:    "option-chosen",
				Prompt:  plugin.PromptSelect,
				Message: "Select your option:",
				Options: []string{
					"option1", "option2", "option3",
				},
				Default: "option2",
			},
		},
	}
}

func (p *Plugin) ValidateAnswers(_ map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
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
