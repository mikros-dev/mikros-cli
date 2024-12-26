package service

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/somatech1/mikros/components/definition"
	"github.com/somatech1/mikros/components/plugin"

	golang_templates "github.com/somatech1/mikros-cli/internal/assets/golang"
	"github.com/somatech1/mikros-cli/internal/golang"
	"github.com/somatech1/mikros-cli/internal/protobuf"
	"github.com/somatech1/mikros-cli/internal/templates"
	"github.com/somatech1/mikros-cli/pkg/definitions"
	"github.com/somatech1/mikros-cli/pkg/path"
	mtemplates "github.com/somatech1/mikros-cli/pkg/templates"
)

type InitOptions struct {
	Kind              Kind
	Path              string
	ProtoFilename     string
	FeatureNames      []string
	Features          *plugin.FeatureSet
	Services          *plugin.ServiceSet
	ExternalTemplates *TemplateFileOptions
}

type Kind int

const (
	KindGolang Kind = iota
	KindRust
)

type TemplateFileOptions struct {
	Files                   embed.FS
	Templates               []mtemplates.TemplateFile
	Api                     map[string]interface{}
	NewServiceArgs          map[string]string
	WithExternalFeaturesArg string
	WithExternalServicesArg string
}

// Init initializes a new service locally.
func Init(options *InitOptions) error {
	answers, err := runInitSurvey(options)
	if err != nil {
		return err
	}

	if err := generateTemplates(options, answers); err != nil {
		return err
	}

	return nil
}

func generateTemplates(options *InitOptions, answers *initSurveyAnswers) error {
	var (
		destinationPath = options.Path
	)

	// Set the project base path
	if destinationPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		destinationPath = filepath.Join(cwd, strings.ToLower(answers.Name))
	}

	if _, err := path.CreatePath(destinationPath); err != nil {
		return err
	}

	// Creates the service.toml file
	if err := writeServiceDefinitions(destinationPath, answers); err != nil {
		return err
	}

	// Switch to the destination path to create template sources
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
	if err := golang.ModInit(strcase.ToKebab(answers.Name)); err != nil {
		return err
	}

	// creates go source templates
	if err := generateSources(options, answers); err != nil {
		return err
	}

	return nil
}

func writeServiceDefinitions(path string, answers *initSurveyAnswers) error {
	defs := &definition.Definitions{
		Name:     answers.Name,
		Types:    []string{answers.Type},
		Version:  answers.Version,
		Language: answers.Language,
		Product:  strings.ToUpper(answers.Product),
	}

	if err := definitions.Write(path, defs); err != nil {
		return err
	}

	for name, d := range answers.FeatureDefinitions() {
		if d.ShouldBeSaved() {
			if err := definitions.AppendFeature(path, name, d.Definitions()); err != nil {
				return err
			}
		}
	}

	if svcDefs := answers.ServiceDefinitions(); svcDefs != nil && svcDefs.ShouldBeSaved() {
		if err := definitions.AppendService(path, answers.Type, svcDefs.Definitions()); err != nil {
			return err
		}
	}

	return nil
}

func generateSources(options *InitOptions, answers *initSurveyAnswers) error {
	context, err := generateTemplateContext(options, answers)
	if err != nil {
		return err
	}

	if err := createServiceTemplates(options, answers.TemplateNames(), context); err != nil {
		return err
	}

	return nil
}

func generateTemplateContext(options *InitOptions, answers *initSurveyAnswers) (TemplateContext, error) {
	var (
		svcDefs = answers.ServiceDefinitions()
		defs    interface{}
	)

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

	newServiceArgs, err := generateNewServiceArgs(options, answers)
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

	if options.ExternalTemplates != nil {
		context.ExternalServicesArg = options.ExternalTemplates.WithExternalServicesArg
		context.ExternalFeaturesArg = options.ExternalTemplates.WithExternalFeaturesArg
	}

	if filename := options.ProtoFilename; filename != "" {
		pbFile, err := protobuf.Parse(filename)
		if err != nil {
			return TemplateContext{}, err
		}
		context.GrpcMethods = pbFile.Methods
	}

	return context, nil
}

func generateNewServiceArgs(options *InitOptions, answers *initSurveyAnswers) (string, error) {
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
		if options.ExternalTemplates != nil {
			var (
				svcDefs = answers.ServiceDefinitions()
				defs    interface{}
			)

			if svcDefs != nil {
				defs = svcDefs.Definitions()
			}

			if tpl, ok := options.ExternalTemplates.NewServiceArgs[answers.Type]; ok {
				data := struct {
					ServiceName              string
					ServiceType              string
					ServiceTypeCustomAnswers interface{}
				}{
					ServiceName:              answers.Type,
					ServiceType:              answers.Name,
					ServiceTypeCustomAnswers: defs,
				}

				block, err := templates.ParseBlock(tpl, options.ExternalTemplates.Api, data)
				if err != nil {
					return "", err
				}

				return block, nil
			}
		}
	}

	return "", nil
}

func generateImports(answers *initSurveyAnswers) map[string][]ImportContext {
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

func createServiceTemplates(options *InitOptions, filenames []mtemplates.TemplateFile, context interface{}) error {
	if err := runTemplates(golang_templates.Files, filenames, context, nil); err != nil {
		return err
	}

	if options != nil && options.ExternalTemplates != nil {
		if err := runTemplates(options.ExternalTemplates.Files, options.ExternalTemplates.Templates, context, options.ExternalTemplates.Api); err != nil {
			return err
		}
	}

	return nil
}

func runTemplates(files embed.FS, filenames []mtemplates.TemplateFile, context interface{}, api map[string]interface{}) error {
	tpls, err := mtemplates.Load(&mtemplates.LoadOptions{
		TemplateNames: filenames,
		Files:         files,
		Api:           api,
	})
	if err != nil {
		return err
	}

	generated, err := tpls.Execute(context)
	if err != nil {
		return err
	}

	for _, gen := range generated {
		file, err := os.Create(gen.Filename())
		if err != nil {
			return err
		}

		if _, err := file.Write(gen.Content()); err != nil {
			_ = file.Close()
			return err
		}

		_ = file.Close()
	}

	return nil
}
