package golang

import (
	"embed"
	"fmt"
	"os"
	"slices"

	"github.com/iancoleman/strcase"
	"github.com/somatech1/mikros/components/definition"

	"github.com/somatech1/mikros-cli/internal/answers"
	"github.com/somatech1/mikros-cli/internal/assets/golang"
	"github.com/somatech1/mikros-cli/internal/protobuf"
	"github.com/somatech1/mikros-cli/internal/templates"
	"github.com/somatech1/mikros-cli/pkg/path"
	mtemplates "github.com/somatech1/mikros-cli/pkg/templates"
)

type ExternalTemplates struct {
	Files                   embed.FS
	Templates               []mtemplates.TemplateFile
	Api                     map[string]interface{}
	NewServiceArgs          map[string]string
	WithExternalFeaturesArg string
	WithExternalServicesArg string
}

type Executioner struct {
	hasLifecycle bool
}

func (e *Executioner) PreExecution(serviceName, destinationPath string) error {
	if _, err := path.CreatePath(destinationPath); err != nil {
		return err
	}

	// Switch to the destination path to create required language files
	cwd, err := path.ChangeDir(destinationPath)
	if err != nil {
		return err
	}

	defer func() {
		if e := os.Chdir(cwd); e != nil {
			err = e
		}
	}()

	// creates go.mod
	if err := ModInit(strcase.ToKebab(serviceName)); err != nil {
		return err
	}

	return nil
}

func (e *Executioner) GenerateContext(answers *answers.InitSurveyAnswers, protoFilename string, data interface{}) (interface{}, error) {
	var (
		externalTemplates *ExternalTemplates
		svcDefs           = answers.ServiceDefinitions()
		defs              interface{}
	)

	if e, ok := data.(*ExternalTemplates); ok {
		externalTemplates = e
	}

	if svcDefs != nil {
		defs = svcDefs.Definitions()
	}

	externalService := func() bool {
		switch answers.Type {
		case definition.ServiceType_gRPC.String(), definition.ServiceType_HTTP.String(), definition.ServiceType_Script.String(), definition.ServiceType_Native.String():
			return false
		}

		return true
	}

	newServiceArgs, err := generateNewServiceArgs(externalTemplates, answers)
	if err != nil {
		return TemplateContext{}, err
	}

	context := TemplateContext{
		featuresExtensions:       len(answers.Features) > 0,
		servicesExtensions:       externalService(),
		onStartLifecycle:         slices.Contains(answers.Lifecycle, "start"),
		onFinishLifecycle:        slices.Contains(answers.Lifecycle, "finish"),
		serviceType:              answers.Type,
		NewServiceArgs:           newServiceArgs,
		ServiceName:              answers.Name,
		Imports:                  generateImports(answers),
		ServiceTypeCustomAnswers: defs,
	}

	if externalTemplates != nil {
		context.ExternalServicesArg = externalTemplates.WithExternalServicesArg
		context.ExternalFeaturesArg = externalTemplates.WithExternalFeaturesArg
	}

	if filename := protoFilename; filename != "" {
		pbFile, err := protobuf.Parse(filename)
		if err != nil {
			return TemplateContext{}, err
		}
		context.GrpcMethods = pbFile.Methods
	}

	// Adjust some internal information to use inside other methods
	if context.onStartLifecycle || context.onFinishLifecycle {
		e.hasLifecycle = true
	}

	return context, nil
}

func generateNewServiceArgs(options *ExternalTemplates, answers *answers.InitSurveyAnswers) (string, error) {
	svcSnake := strcase.ToSnake(answers.Name)

	switch answers.Type {
	case definition.ServiceType_gRPC.String():
		return fmt.Sprintf(`Service: map[string]options.ServiceOptions{
			"grpc": &options.GrpcServiceOptions{
				ProtoServiceDescription: &%spb.%sService_ServiceDesc,
			},
		},`, svcSnake, strcase.ToCamel(answers.Name)), nil

	case definition.ServiceType_HTTP.String():
		return fmt.Sprintf(`Service: map[string]options.ServiceOptions{
			"http": &options.HttpServiceOptions{
				ProtoHttpServer: %spb.NewHttpServer(),
			},
		},`, svcSnake), nil

	case definition.ServiceType_Native.String():
		return `Service: map[string]options.ServiceOptions{
			"native": &options.NativeServiceOptions{},
		},`, nil

	case definition.ServiceType_Script.String():
		return `Service: map[string]options.ServiceOptions{
			"script": &options.ScriptServiceOptions{},
		},`, nil

	default:
		if options != nil {
			var (
				svcDefs = answers.ServiceDefinitions()
				defs    interface{}
			)

			if svcDefs != nil {
				defs = svcDefs.Definitions()
			}

			if tpl, ok := options.NewServiceArgs[answers.Type]; ok {
				data := struct {
					ServiceName              string
					ServiceType              string
					ServiceTypeCustomAnswers interface{}
				}{
					ServiceName:              answers.Type,
					ServiceType:              answers.Name,
					ServiceTypeCustomAnswers: defs,
				}

				block, err := templates.ParseBlock(tpl, options.Api, data)
				if err != nil {
					return "", err
				}

				return block, nil
			}
		}
	}

	return "", nil
}

func generateImports(answers *answers.InitSurveyAnswers) map[string][]ImportContext {
	imports := map[string][]ImportContext{
		"main": {
			{
				Path: "github.com/somatech1/mikros",
			},
			{
				Path: "github.com/somatech1/mikros/components/options",
			},
		},
		"service": {
			{
				Path: "github.com/somatech1/mikros",
			},
		},
	}

	if len(answers.Lifecycle) > 0 {
		imports["lifecycle"] = append(imports["lifecycle"], ImportContext{
			Path: "context",
		})
	}

	return imports
}

func (e *Executioner) Templates() []mtemplates.TemplateFile {
	names := []mtemplates.TemplateFile{
		{
			Name:      "main",
			Extension: "go",
		},
		{
			Name:      "service",
			Extension: "go",
		},
	}

	if e.hasLifecycle {
		names = append(names, mtemplates.TemplateFile{
			Name:      "lifecycle",
			Extension: "go",
		})
	}

	return names
}

func (e *Executioner) Files() embed.FS {
	return golang.Files
}

func (e *Executioner) PostExecution(_ string) error {
	return nil
}
