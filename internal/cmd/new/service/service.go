package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/iancoleman/strcase"
	"github.com/somatech1/mikros/components/definition"

	service_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/service"
	"github.com/mikros-dev/mikros-cli/internal/definitions"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/protobuf"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
	"github.com/mikros-dev/mikros-cli/internal/ui"
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

func runSurvey(cfg *settings.Settings) (*surveyAnswers, error) {
	var (
		answers        = newSurveyAnswers()
		supportedTypes = []huh.Option[string]{
			huh.NewOption(definition.ServiceType_gRPC.String(), definition.ServiceType_gRPC.String()),
			huh.NewOption(definition.ServiceType_HTTP.String(), definition.ServiceType_HTTP.String()),
			huh.NewOption(definition.ServiceType_Native.String(), definition.ServiceType_Native.String()),
			huh.NewOption(definition.ServiceType_Script.String(), definition.ServiceType_Script.String()),
		}
	)

	newTypes, err := plugin.GetNewServiceKinds(cfg)
	if err != nil {
		return nil, err
	}
	for _, t := range newTypes {
		supportedTypes = append(supportedTypes, huh.NewOption(t, t))
	}
	sort.Slice(supportedTypes, func(i, j int) bool {
		return supportedTypes[i].String() < supportedTypes[j].String()
	})

	var languages []huh.Option[string]
	for _, t := range definition.SupportedLanguages() {
		languages = append(languages, huh.NewOption(t, t))
	}

	questions := []huh.Field{
		huh.NewInput().
			Title("Service name. Can be a fully qualified name (URL + name):").
			Value(&answers.Name).
			Validate(ui.IsEmpty("service name cannot be empty")),

		huh.NewSelect[string]().
			Title("Select the type of service:").
			Options(supportedTypes...).
			Value(&answers.Type).
			Validate(ui.IsEmpty("service type cannot be empty")),

		huh.NewSelect[string]().
			Title("Select the service programming language:").
			Options(languages...).
			Value(&answers.Language).
			Validate(ui.IsEmpty("service programming language cannot be empty")),

		huh.NewInput().
			Title("Version. A semver version string for the service, with 'v' as prefix (ex: v1.0.0):").
			Value(&answers.Version).
			Validate(func(s string) error {
				if !definition.ValidateVersion(s) {
					return errors.New("invalid version format")
				}

				return nil
			}),

		huh.NewInput().
			Title("Product name. Enter the product name that the service belongs to:").
			Value(&answers.Product).
			Validate(ui.IsEmpty("product name cannot be empty")),

		huh.NewMultiSelect[string]().
			Title("Select lifecycle events to handle in the service:").
			Options(
				huh.NewOption("OnStart", "OnStart"),
				huh.NewOption("OnFinish", "OnFinish"),
			).
			Value(&answers.Lifecycle),
	}

	featureNames, err := plugin.GetFeaturesUINames(cfg)
	if err != nil {
		return nil, err
	}
	if len(featureNames) > 0 {
		features := make([]huh.Option[string], len(featureNames))
		for i, f := range featureNames {
			features[i] = huh.NewOption(f, f)
		}

		questions = append(questions, huh.NewMultiSelect[string]().
			Title("Select the features the service will have").
			Options(features...).
			Value(&answers.Features),
		)
	}

	form := huh.NewForm(huh.NewGroup(questions...)).
		WithTheme(cfg.GetTheme()).
		WithAccessible(cfg.UI.Accessible)

	if err := form.Run(); err != nil {
		return nil, err
	}

	return answers, nil
}

// runServiceSurvey executes the survey that a service may have implemented.
func runServiceSurvey(cfg *settings.Settings, answers *surveyAnswers) (*client.Service, error) {
	svc, err := plugin.GetServicePlugin(cfg, answers.Type)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// No plugin for the chosen service type.
		return nil, nil
	}

	svcSurvey, err := svc.GetSurvey()
	if err != nil {
		return nil, err
	}

	response, err := ui.RunFormFromSurvey(answers.Type, svcSurvey, &ui.FormOptions{
		Theme:      cfg.GetTheme(),
		Accessible: cfg.UI.Accessible,
	})
	if err != nil {
		return nil, err
	}

	d, err := svc.ValidateAnswers(response)
	if err != nil {
		return nil, err
	}

	answers.SetServiceDefinitions(d)
	return svc, nil
}

func runFeatureSurvey(cfg *settings.Settings, name string) (string, interface{}, error) {
	f, err := plugin.GetFeaturePlugin(cfg, name)
	if err != nil {
		return "", nil, err
	}
	if f == nil {
		return "", nil, nil
	}

	s, err := f.GetSurvey()
	if err != nil {
		return "", nil, err
	}
	if s == nil {
		return "", nil, nil
	}

	res, err := ui.RunFormFromSurvey(name, s, &ui.FormOptions{
		Theme:      cfg.GetTheme(),
		Accessible: cfg.UI.Accessible,
	})
	if err != nil {
		return "", nil, err
	}

	defs, err := f.ValidateAnswers(res)
	if err != nil {
		return "", nil, err
	}

	featureName, err := f.GetName()
	if err != nil {
		return "", nil, err
	}

	return featureName, defs, nil
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
		res, err := svc.GetTemplates()
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

				return block, nil
			}
		}
	}

	return "", nil
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

func createServiceTemplates(filenames []template.File, tplContext interface{}, externalTemplate *mtemplate.Template) error {
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
			files[i] = &template.Data{
				FileName: templateNames[i].Name,
				Content:  []byte(t.Content),
			}
		}

		session, err := template.NewSessionFromData(&template.LoadOptions{
			TemplateNames: templateNames,
		}, files)
		if err != nil {
			return err
		}

		if err := runTemplates(session, tplContext); err != nil {
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
