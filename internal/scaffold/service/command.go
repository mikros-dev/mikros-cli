package service

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/mikros-dev/mikros/components/definition"

	"github.com/mikros-dev/mikros-cli/internal/definitions"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/protobuf"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
	mtemplate "github.com/mikros-dev/mikros-cli/pkg/template"
)

// NewOptions holds configuration options for creating a service template.
type NewOptions struct {
	// Path specifies the target directory for the service files.
	Path string

	// ProtoFilename defines the location of the protobuf file used for the
	// service.
	ProtoFilename string
}

// New creates a new service template directory with initial source files.
func New(cfg *settings.Settings, options *NewOptions) error {
	// Execute the base survey
	answers, err := runSurvey(cfg, options.ProtoFilename)
	if err != nil {
		return err
	}

	// Then execute everything specific for the selected service type.
	svc, err := runServiceTypeSurvey(cfg, answers)
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

	return generateTemplates(options, answers, svc)
}

func generateTemplates(options *NewOptions, answers *surveyAnswers, svc *client.Service) error {
	var destinationPath = filepath.Join(options.Path, strings.ToLower(answers.Name))

	// Set the project base path
	if destinationPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		destinationPath = filepath.Join(cwd, strings.ToLower(answers.Name))
	}

	if _, err := path.CreatePath(destinationPath); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	// Creates the service.toml file
	if err := writeServiceDefinitions(destinationPath, answers); err != nil {
		return err
	}

	// Switch to the destination path to create template sources
	cwd, err := path.ChangeDir(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	defer func() {
		_ = os.Chdir(cwd)
	}()

	// creates go.mod
	if err := golang.ModInit(strcase.ToKebab(answers.Name)); err != nil {
		return fmt.Errorf("failed to create go.mod: %w", err)
	}

	// creates go source templates
	return generateSources(options, answers, svc)
}

func writeServiceDefinitions(path string, answers *surveyAnswers) error {
	defs := &definition.Definitions{
		Name:     answers.Name,
		Types:    []string{answers.ServiceType()},
		Version:  answers.Version,
		Language: answers.Language,
		Product:  strings.ToUpper(answers.Product),
	}

	if err := definitions.Write(path, defs); err != nil {
		return fmt.Errorf("failed to write service definitions file: %w", err)
	}

	for name, d := range answers.FeatureDefinitions() {
		if d.ShouldBeSaved() {
			if err := definitions.AppendFeature(path, name, d.Definitions()); err != nil {
				return fmt.Errorf("failed to write feature definitions: %w", err)
			}
		}
	}

	if svcDefs := answers.ServiceDefinitions(); svcDefs != nil && svcDefs.ShouldBeSaved() {
		if err := definitions.AppendService(path, answers.Type, svcDefs.Definitions()); err != nil {
			return fmt.Errorf("failed to write service definitions: %w", err)
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

	return createServiceTemplates(answers.TemplateNames(), tplCtx, externalTemplate)
}

func generateTemplateContext(
	options *NewOptions,
	answers *surveyAnswers,
	externalTemplate *mtemplate.Template,
) (TemplateContext, error) {
	var (
		svcDefs = answers.ServiceDefinitions()
		defs    interface{}
	)

	if svcDefs != nil {
		defs = svcDefs.Definitions()
	}

	externalService := func() bool {
		switch answers.Type {
		case definition.ServiceTypeGRPC.String(),
			definition.ServiceTypeHTTP.String(),
			definition.ServiceTypeHTTPSpec.String(),
			definition.ServiceTypeScript.String(),
			definition.ServiceTypeWorker.String():
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
		serviceType:              answers.ServiceType(),
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

	switch answers.ServiceType() {
	case definition.ServiceTypeGRPC.String():
		svcInitBlock = fmt.Sprintf(`"grpc": &options.GrpcServiceOptions{
				ProtoServiceDescription: &%spb.%sService_ServiceDesc,
			},`, svcSnake, strcase.ToCamel(answers.Name))

	case definition.ServiceTypeHTTPSpec.String():
		svcInitBlock = fmt.Sprintf(`"http-spec": &options.HTTPSpecServiceOptions{
				ProtoHttpServer: %spb.NewHttpServer(),
			},`, svcSnake)

	case definition.ServiceTypeHTTP.String():
		svcInitBlock = `"http": &options.HTTPServiceOptions{},`

	case definition.ServiceTypeWorker.String():
		svcInitBlock = `"worker": &options.WorkerServiceOptions{},`

	case definition.ServiceTypeScript.String():
		svcInitBlock = `"script": &options.ScriptServiceOptions{},`

	default:
		b, err := externalTemplateInitBlock(answers, externalTemplate)
		if err != nil {
			return "", err
		}
		svcInitBlock = b
	}

	return fmt.Sprintf(`Service: map[string]options.ServiceOptions{
			%s
		},`, svcInitBlock), nil
}

func externalTemplateInitBlock(answers *surveyAnswers, externalTemplate *mtemplate.Template) (string, error) {
	if externalTemplate == nil {
		return "", nil
	}

	var (
		svcDefs   = answers.ServiceDefinitions()
		defs      interface{}
		initBlock string
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
			return "", fmt.Errorf("failed to parse external template: %w", err)
		}
		initBlock = block
	}

	return initBlock, nil
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
				Path:  "github.com/mikros-dev/mikros/apis/features/logger",
				Alias: "logger_api",
			},
			{
				Path:  "github.com/mikros-dev/mikros/apis/features/errors",
				Alias: "errors_api",
			},
		},
	}

	if len(answers.Lifecycle) > 0 {
		imports["lifecycle"] = append(imports["lifecycle"], ImportContext{
			Path: "context",
		})
	}

	if answers.HTTPType == definition.ServiceTypeHTTP.String() {
		imports["http"] = append(imports["http"], []ImportContext{
			{
				Path: "net/http",
			},
			{
				Path: "context",
			},
		}...)
	}

	return imports
}

func createServiceTemplates(
	filenames []template.File,
	tplContext TemplateContext,
	externalTemplate *mtemplate.Template,
) error {
	// Execute our templates
	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplatesToUse: filenames,
		FilesBasePath:  "assets",
	}, templateFiles)
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
			TemplatesToUse: templateNames,
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
