package main

import (
	"github.com/mikros-dev/mikros-cli/pkg/plugin"
)

const (
	pluginName    = "follow-up"
	uiFeatureName = "follow-up"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) UIName() string {
	return uiFeatureName
}

func (p *Plugin) Survey() *plugin.Survey {
	// We create the plugin inside the loop here.
	return &plugin.Survey{
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
		FollowUp: []*plugin.FollowUpSurvey{
			// Only executed when 'option-chosen' is option3
			{
				Name: "name-to-choose",
				Condition: &plugin.QuestionCondition{
					Name:  "option-chosen",
					Value: "option3",
				},
				Survey: &plugin.Survey{
					Questions: []*plugin.Question{
						{
							Name:    "condition1-option3-chosen",
							Prompt:  plugin.PromptInput,
							Message: "Enter the name you want:",
							Default: "my name",
						},
					},
				},
			},
			// Only executed when 'option-chosen' is option1
			{
				Name: "age-to-choose",
				Condition: &plugin.QuestionCondition{
					Name:  "option-chosen",
					Value: "option1",
				},
				Survey: &plugin.Survey{
					Questions: []*plugin.Question{
						{
							Name:    "condition1-option1-chosen",
							Prompt:  plugin.PromptInput,
							Message: "Enter your age:",
							Default: "42",
						},
					},
				},
			},
			// Only executed when 'option-chosen' is option1 or option3
			{
				Name: "address-to-choose",
				Condition: &plugin.QuestionCondition{
					Name:  "option-chosen",
					Value: []string{"option1", "option3"},
				},
				Survey: &plugin.Survey{
					Questions: []*plugin.Question{
						{
							Name:    "condition1-option1-option3-chosen",
							Prompt:  plugin.PromptInput,
							Message: "Enter your address:",
							Default: "Nowhere",
						},
					},
				},
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
