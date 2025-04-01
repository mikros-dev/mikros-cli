package service

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/somatech1/mikros/components/definition"

	service_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/service"
	"github.com/mikros-dev/mikros-cli/internal/definitions"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/protobuf"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
	mtemplate "github.com/mikros-dev/mikros-cli/pkg/template"
)

type NewOptions struct {
	Path          string
	ProtoFilename string
}

// New creates a new service template directory with initial source files.
func New(cfg *settings.Settings, options *NewOptions) error {
	answers, err := runSurvey(cfg)
	if err != nil {
		return err
	}

	svc, err := runServiceSurvey(cfg, answers)
	if err != nil {
		return err
	}

	// Presents only questions from selected features
	for _, name := range answers.Features {
		featureName, defs, err := runFeatureSurvey(cfg, name)
		if err != nil {
			return err
		}
		if d, ok := defs.(map[string]interface{}); ok && len(d) != 0 {
			answers.AddFeatureDefinitions(featureName, defs)
		}
	}

	if err := generateTemplates(options, answers, svc); err != nil {
		return err
	}

	return nil
}

func generateTemplates(options *NewOptions, answers *surveyAnswers, svc *client.Service) error {
	var (
		destinationPath = filepath.Join(options.Path, strings.ToLower(answers.Name))
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
	if err := generateSources(options, answers, svc); err != nil {
		return err
	}

	return nil
}

func writeServiceDefinitions(path string, answers *surveyAnswers) error {
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

func generateSources(options *NewOptions, answers *surveyAnswers, svc *client.Service) error {
	var externalTemplate *mtemplate.Template
	if svc != nil {
		res, err := svc.GetTemplates(answers.ServiceAnswers())
		if err != nil {
			return err
		}
		externalTemplate = res
	}

	tplCtx, err := generateTemplateContext(options, answers, externalTemplate)
	if err != nil {
		return err
	}

	if err := createServiceTemplates(answers.TemplateNames(), tplCtx, externalTemplate); err != nil {
		return err
	}

	return nil
}

func generateTemplateContext(options *NewOptions, answers *surveyAnswers, externalTemplate *mtemplate.Template) (TemplateContext, error) {
	var (
		svcDefs = answers.ServiceDefinitions()
		defs    interface{}
	)

	if svcDefs != nil {
		defs = svcDefs.Definitions()
	}

	externalService := func() bool {
		switch answers.Type {
		case definition.ServiceType_gRPC.String(),
			definition.ServiceType_HTTP.String(),
			definition.ServiceType_Script.String(),
			definition.ServiceType_Native.String():
			return false
		}

		return true
	}

	newServiceArgs, err := generateNewServiceArgs(answers, externalTemplate)
	if err != nil {
		return TemplateContext{}, err
	}

	tplCtx := TemplateContext{
		featuresExtensions:       len(answers.Features) > 0,
		servicesExtensions:       externalService(),
		onStartLifecycle:         slices.Contains(answers.Lifecycle, "OnStart"),
		onFinishLifecycle:        slices.Contains(answers.Lifecycle, "OnFinish"),
		serviceType:              answers.Type,
		NewServiceArgs:           newServiceArgs,
		ServiceName:              answers.Name,
		Imports:                  generateImports(answers),
		ServiceTypeCustomAnswers: defs,
	}

	if externalTemplate != nil {
		tplCtx.ExternalServicesArg = externalTemplate.WithExternalServicesArg
		tplCtx.ExternalFeaturesArg = externalTemplate.WithExternalFeaturesArg
	}

	if filename := options.ProtoFilename; filename != "" {
		pbFile, err := protobuf.Parse(filename)
		if err != nil {
			return TemplateContext{}, err
		}
		tplCtx.GrpcMethods = pbFile.Methods
	}

	return tplCtx, nil
}

func generateNewServiceArgs(answers *surveyAnswers, externalTemplate *mtemplate.Template) (string, error) {
	var (
		svcSnake     = strcase.ToSnake(answers.Name)
		svcInitBlock string
	)

	switch answers.Type {
	case definition.ServiceType_gRPC.String():
		svcInitBlock = fmt.Sprintf(`
			"grpc": &options.GrpcServiceOptions{
				ProtoServiceDescription: &%spb.%sService_ServiceDesc,
			},`, svcSnake)

	case definition.ServiceType_HTTP.String():
		svcInitBlock = fmt.Sprintf(`
			"http": &options.HttpServiceOptions{
				ProtoHttpServer: %spb.NewHttpServer(),
			},`, svcSnake)

	case definition.ServiceType_Native.String():
		svcInitBlock = `
			"native": &options.NativeServiceOptions{},
`

	case definition.ServiceType_Script.String():
		svcInitBlock = `
			"script": &options.ScriptServiceOptions{},
`

	default:
		if externalTemplate != nil {
			var (
				svcDefs = answers.ServiceDefinitions()
				defs    interface{}
			)

			if svcDefs != nil {
				defs = svcDefs.Definitions()
			}

			if args := externalTemplate.NewServiceArgs; args != "" {
				data := struct {
					ServiceName              string
					ServiceType              string
					ServiceTypeCustomAnswers interface{}
				}{
					ServiceName:              answers.Type,
					ServiceType:              answers.Name,
					ServiceTypeCustomAnswers: defs,
				}

				block, err := template.ParseBlock(args, nil, data)
				if err != nil {
					return "", err
				}
				svcInitBlock = block
			}
		}
	}

	return fmt.Sprintf(`Service: map[string]options.ServiceOptions{
			%s
		},`, svcInitBlock), nil
}

func generateImports(answers *surveyAnswers) map[string][]ImportContext {
	imports := map[string][]ImportContext{
		"main": {
			{
				Path: "github.com/mikros-dev/mikros",
			},
			{
				Path: "github.com/mikros-dev/mikros/components/options",
			},
		},
		"service": {
			{
				Path: "github.com/mikros-dev/mikros",
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

func createServiceTemplates(filenames []template.File, tplContext TemplateContext, externalTemplate *mtemplate.Template) error {
	// Execute our templates
	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplateNames: filenames,
	}, service_tpl.Files)
	if err != nil {
		return err
	}

	if err := runTemplates(session, tplContext); err != nil {
		return err
	}

	// Then execute templates from the selected plugin (if any).
	if externalTemplate != nil {
		templateNames := make([]template.File, len(externalTemplate.Templates))
		for i, t := range externalTemplate.Templates {
			templateNames[i] = template.File{
				Name:      t.Name,
				Output:    t.Output,
				Extension: t.Extension,
			}
		}

		files := make([]*template.Data, len(templateNames))
		for i, t := range externalTemplate.Templates {
			name := templateNames[i].Name
			if name == "" {
				name = templateNames[i].Output
			}

			// Set the context PluginData with custom context from the plugin
			tplContext.PluginData = t.Context
			files[i] = &template.Data{
				FileName: name,
				Content:  []byte(t.Content),
				Context:  tplContext,
			}
		}

		session, err := template.NewSessionFromData(&template.LoadOptions{
			TemplateNames: templateNames,
		}, files)
		if err != nil {
			return err
		}

		if err := runTemplates(session, nil); err != nil {
			return err
		}
	}

	return nil
}

func runTemplates(session *template.Session, context interface{}) error {
	generated, err := session.ExecuteTemplates(context)
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
