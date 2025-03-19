package main

import (
	"github.com/iancoleman/strcase"

	"github.com/mikros-dev/mikros-cli/pkg/plugin"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	mtemplate "github.com/mikros-dev/mikros-cli/pkg/template"

	"worker/assets"
)

type Context struct {
	EventName string
}

type Plugin struct{}

func (p *Plugin) Kind() string {
	return "worker"
}

func (p *Plugin) Survey() *survey.Survey {
	return &survey.Survey{
		ConfirmQuestion: &survey.Question{
			Message:      "Do you want to add another event?",
			Default:      "true",
			ConfirmAfter: true,
		},
		Questions: []*survey.Question{
			{
				Name:     "topic_name",
				Prompt:   survey.PromptInput,
				Message:  "Topic name. The subscription topic name to subscribe into:",
				Required: true,
			},
			{
				Name:     "topic_service_name",
				Prompt:   survey.PromptInput,
				Message:  "The service that emits the event. Enter the service name:",
				Required: true,
			},
		},
	}
}

func (p *Plugin) ValidateAnswers(_ map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"stream_kind": "kinesis",
		"stream_name": "some-random-stream",
	}, nil
}

func (p *Plugin) Template(in map[string]interface{}) *mtemplate.Template {
	files := make(map[string]string)
	templateFiles, err := assets.Files.ReadDir(".")
	if err != nil {
		return nil
	}

	for _, file := range templateFiles {
		data, err := assets.Files.ReadFile(file.Name())
		if err != nil {
			continue
		}

		files[file.Name()] = string(data)
	}

	return &mtemplate.Template{
		NewServiceArgs:          `"worker": &mikros_extensions.WorkerService{},`,
		WithExternalFeaturesArg: "",
		WithExternalServicesArg: "",
		Templates:               createTemplateFiles(in, files),
	}
}

func createTemplateFiles(in map[string]interface{}, files map[string]string) []*mtemplate.File {
	data, ok := in["worker"].([]interface{})
	if !ok {
		return nil
	}

	tplFiles := make([]*mtemplate.File, len(data))
	for i, d := range data {
		entry, ok := d.(map[string]interface{})
		if !ok {
			continue
		}

		tplFiles[i] = &mtemplate.File{
			Content:   files["event.go.tmpl"],
			Output:    strcase.ToSnake(entry["topic_name"].(string)),
			Extension: "go",
			Context: Context{
				EventName: strcase.ToCamel(entry["topic_name"].(string)),
			},
		}
	}

	return tplFiles
}

func main() {
	p, err := plugin.NewService(&Plugin{})
	if err != nil {
		plugin.Error(err)
	}

	if err := p.Run(); err != nil {
		plugin.Error(err)
	}
}
