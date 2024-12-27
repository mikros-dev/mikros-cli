package rust

import (
	"embed"

	"github.com/somatech1/mikros/components/definition"

	"github.com/somatech1/mikros-cli/internal/answers"
	"github.com/somatech1/mikros-cli/internal/assets/rust"
	"github.com/somatech1/mikros-cli/pkg/templates"
)

type Executioner struct {
	hasLifecycle bool
	serviceType  string
}

func (e *Executioner) PreExecution(serviceName, destinationPath string) error {
	if err := cargoInit(destinationPath, serviceName); err != nil {
		return err
	}

	return nil
}

func (e *Executioner) GenerateContext(answers *answers.InitSurveyAnswers, _ string, _ interface{}) (interface{}, error) {
	// Store values to be used later
	e.hasLifecycle = answers.RustLifecycle
	e.serviceType = answers.Type

	return nil, nil
}

func (e *Executioner) Templates() []templates.TemplateFile {
	return []templates.TemplateFile{}
}

func (e *Executioner) Files() embed.FS {
	return rust.Files
}

type dependency struct {
	Name    string
	Version string
	Git     string
	Path    string
	Feature []string
}

func (e *Executioner) PostExecution(destinationPath string) error {
	dependencies := []dependency{
		{
			Name: "mikros",
			// FIXME: change this when mikros-rs is published
			Path: "/Users/rodrigo/desenv/github/rsfreitas/mikros-rs",
		},
		{
			Name:    "tokio",
			Version: "1.41.1",
			Feature: []string{"full"},
		},
	}

	if e.serviceType == definition.ServiceType_Script.String() || e.serviceType == definition.ServiceType_Native.String() || e.hasLifecycle {
		dependencies = append(dependencies, dependency{
			Name:    "async-trait",
			Version: "0.1.83",
		})
	}

	if e.serviceType == definition.ServiceType_HTTP.String() {
		dependencies = append(dependencies, dependency{
			Name:    "axum",
			Version: "0.7.7",
		})
	}

	if e.serviceType == definition.ServiceType_gRPC.String() {
		dependencies = append(dependencies, []dependency{
			{
				Name:    "prost",
				Version: "0.13.3",
			},
			{
				Name:    "tonic",
				Version: "0.12.3",
			},
		}...)
	}

	for _, d := range dependencies {
		if err := cargoAdd(destinationPath, d.Name, d.Version, d.Git, d.Path, d.Feature); err != nil {
			return err
		}
	}

	return nil
}
