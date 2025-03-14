package service

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"github.com/somatech1/mikros/components/definition"

	assets "github.com/mikros-dev/mikros-cli/internal/assets/templates"
	"github.com/mikros-dev/mikros-cli/internal/definitions"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/protobuf"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
	"github.com/mikros-dev/mikros-cli/internal/ui"
	msurvey "github.com/mikros-dev/mikros-cli/pkg/survey"
	mtemplate "github.com/mikros-dev/mikros-cli/pkg/template"
)

type InitOptions struct {
	Path          string
	ProtoFilename string
}

// Init initializes a new service locally.
func Init(cfg *settings.Settings, options *InitOptions) error {
	questions, err := baseQuestions(cfg)
	if err != nil {
		return err
	}

	answers := &initSurveyAnswers{}
	if err := survey.Ask(questions, answers); err != nil {
		return err
	}

	svc, err := runServiceSurvey(cfg, answers)
	if err != nil {
		return err
	}

	// Presents only questions from selected features
	for _, name := range answers.Features {
		defs, save, err := runFeatureSurvey(cfg, name)
		if err != nil {
			return err
		}
		if defs != nil {
			answers.AddFeatureDefinitions(name, defs, save)
		}
	}

	if err := generateTemplates(options, answers, svc); err != nil {
		return err
	}

	return nil
}

func baseQuestions(cfg *settings.Settings) ([]*survey.Question, error) {
	supportedTypes := []string{
		definition.ServiceType_gRPC.String(),
		definition.ServiceType_HTTP.String(),
		definition.ServiceType_Native.String(),
		definition.ServiceType_Script.String(),
	}

	newTypes, err := plugin.GetNewServiceKinds(cfg)
	if err != nil {
		return nil, err
	}
	supportedTypes = append(supportedTypes, newTypes...)

	sort.Strings(supportedTypes)
	questions := []*survey.Question{
		// Service name
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "FileName. Can be a fully qualified service name (URL + name):",
			},
			Validate: survey.ComposeValidators(
				survey.Required,
				survey.MinLength(0),
				survey.MaxLength(512),
			),
		},
		// Service type
		{
			Name: "type",
			Prompt: &survey.Select{
				Message:  "Select the type of service:",
				Options:  supportedTypes,
				PageSize: len(supportedTypes),
			},
			Validate: survey.Required,
		},
		// Language
		{
			Name: "language",
			Prompt: &survey.Select{
				Message:  "Select the service main programming language:",
				Options:  definition.SupportedLanguages(),
				PageSize: len(definition.SupportedLanguages()),
			},
			Validate: survey.Required,
		},
		// Version
		{
			Name: "version",
			Prompt: &survey.Input{
				Message: "Version. A semver version string for the service, with 'v' as prefix (ex: v1.0.0):",
				Default: "v0.1.0",
			},
			Validate: func(val interface{}) error {
				if str, ok := val.(string); ok {
					if !definition.ValidateVersion(str) {
						return errors.New("invalid version format")
					}

					return nil
				}

				return errors.New("version has an invalid value type")
			},
		},
		// Product
		{
			Name: "product",
			Prompt: &survey.Input{
				Message: "Product name. Enter the product name that the service belongs to:",
			},
			Validate: survey.ComposeValidators(
				survey.Required,
				survey.MinLength(3),
				survey.MaxLength(512),
			),
		},
		// Lifecycle
		{
			Name: "lifecycle",
			Prompt: &survey.MultiSelect{
				Message: "Select lifecycle events to handle in the service:",
				Options: []string{"OnStart", "OnFinish"},
			},
		},
	}

	featureNames, err := plugin.GetFeaturesUINames(cfg)
	if err != nil {
		return nil, err
	}
	if len(featureNames) > 0 {
		// Features
		questions = append(questions, &survey.Question{
			Name: "features",
			Prompt: &survey.MultiSelect{
				Message:  "Select the features the service will have:",
				Options:  featureNames,
				PageSize: len(featureNames),
			},
		})
	}

	return questions, nil
}

// runServiceSurvey executes the survey that a service may have implemented.
func runServiceSurvey(cfg *settings.Settings, answers *initSurveyAnswers) (*client.Service, error) {
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

	response, err := handleSurvey(answers.Type, svcSurvey)
	if err != nil {
		return nil, err
	}

	d, save, err := svc.ValidateAnswers(response)
	if err != nil {
		return nil, err
	}

	answers.SetServiceDefinitions(d, save)
	return svc, nil
}

func handleSurvey(name string, featureSurvey *msurvey.Survey) (map[string]interface{}, error) {
	if featureSurvey.ConfirmQuestion != nil {
		var responses []map[string]interface{}

	loop:
		for {
			if !featureSurvey.ConfirmQuestion.ConfirmAfter {
				res := ui.YesNo(featureSurvey.ConfirmQuestion.Message)
				if !res {
					break loop
				}
			}

			response, err := surveyFromQuestion(name, featureSurvey)
			if err != nil {
				return nil, err
			}
			responses = append(responses, response)

			if featureSurvey.ConfirmQuestion.ConfirmAfter {
				res := ui.YesNo(featureSurvey.ConfirmQuestion.Message)
				if !res {
					break loop
				}
			}
		}

		return map[string]interface{}{
			name: responses,
		}, nil
	}

	response, err := surveyFromQuestion(name, featureSurvey)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func surveyFromQuestion(name string, entrySurvey *msurvey.Survey) (map[string]interface{}, error) {
	var (
		s        []*survey.Question
		response = make(map[string]interface{})
		validate = validator.New()
	)

	for _, q := range entrySurvey.Questions {
		if err := validate.Struct(q); err != nil {
			return nil, err
		}

		question := &survey.Question{
			Name: q.Name,
			Validate: func() func(v interface{}) error {
				//if q.Validate != nil {
				//	return q.Validate
				//}

				if q.Required {
					return survey.Required
				}

				return nil
			}(),
		}

		switch q.Prompt {
		case msurvey.PromptSurvey:
			if !validateInnerSurveyCondition(response, q.Condition) {
				continue
			}

			r, err := handleSurvey(q.Name, q.Survey)
			if err != nil {
				return nil, err
			}
			if r != nil {
				r = sanitizeResponse(r)
				response[q.Name] = r[q.Name]
			}

			continue

		default:
			question.Prompt = buildSurveyPrompt(name, q)
		}

		if q.Prompt != msurvey.PromptSurvey {
			s = append(s, question)
		}

		if entrySurvey.AskOne {
			if validateInnerSurveyCondition(response, q.Condition) {
				r, err := askOne(question.Prompt, q)
				if err != nil {
					return nil, err
				}

				response[question.Name] = r
				response = sanitizeResponse(response)
			}
		}
	}

	// If we don't have response we need to execute the survey entirely.
	if len(response) == 0 {
		if err := survey.Ask(s, &response); err != nil {
			return nil, err
		}
	}

	return sanitizeResponse(response), nil
}

func buildSurveyPrompt(name string, q *msurvey.Question) survey.Prompt {
	switch q.Prompt {
	case msurvey.PromptInput:
		return &survey.Input{
			Message: fmt.Sprintf("[%s] %s", name, q.Message),
			Default: q.Default,
		}

	case msurvey.PromptSelect:
		return &survey.Select{
			Message:  fmt.Sprintf("[%s] %s", name, q.Message),
			Options:  q.Options,
			PageSize: len(q.Options),
			Default:  q.Default,
		}

	case msurvey.PromptMultiSelect:
		return &survey.MultiSelect{
			Message: fmt.Sprintf("[%s] %s", name, q.Message),
			Options: q.Options,
		}

	case msurvey.PromptMultiline:
		return &survey.Multiline{
			Message: fmt.Sprintf("[%s] %s", name, q.Message),
		}

	case msurvey.PromptConfirm:
		return &survey.Confirm{
			Message: fmt.Sprintf("[%s] %s", name, q.Message),
		}

	default:
	}

	return nil
}

func validateInnerSurveyCondition(response map[string]interface{}, condition *msurvey.QuestionCondition) bool {
	if condition != nil {
		if r, ok := response[condition.Name]; ok {
			switch value := condition.Value.(type) {
			case []string:
				if slices.Contains(value, r.(string)) {
					return true
				}

			case string:
				if v, ok := r.(string); ok && v == value {
					return true
				}
			}
		}

		return false
	}

	return true
}

func askOne(prompt survey.Prompt, question *msurvey.Question) (interface{}, error) {
	getOptions := func() survey.AskOpt {
		//if question.Validate != nil {
		//	return survey.WithValidator(question.Validate)
		//}
		if question.Required {
			return survey.WithValidator(survey.Required)
		}

		return nil
	}

	if question.Prompt == msurvey.PromptSelect {
		index := 0
		if err := survey.AskOne(prompt, &index, getOptions()); err != nil {
			return nil, err
		}

		return question.Options[index], nil
	}

	if question.Prompt == msurvey.PromptMultiSelect {
		var r []string
		if err := survey.AskOne(prompt, &r, getOptions()); err != nil {
			return nil, err
		}

		return r, nil
	}

	var r string
	if err := survey.AskOne(prompt, &r, getOptions()); err != nil {
		return nil, err
	}

	return r, nil
}

// sanitizeResponse sanitizes the response to avoid sending internal formats to
// the client.
func sanitizeResponse(response map[string]interface{}) map[string]interface{} {
	for k, v := range response {
		if s, ok := v.(survey.OptionAnswer); ok {
			response[k] = s.Value
		}

		if opts, ok := v.([]survey.OptionAnswer); ok {
			var s []string
			for _, o := range opts {
				s = append(s, o.Value)
			}

			response[k] = s
		}
	}

	return response
}

func runFeatureSurvey(cfg *settings.Settings, name string) (interface{}, bool, error) {
	f, err := plugin.GetFeaturePlugin(cfg, name)
	if err != nil {
		return nil, false, err
	}
	if f == nil {
		return nil, false, nil
	}

	s, err := f.GetSurvey()
	if err != nil {
		return nil, false, err
	}
	if s == nil {
		return nil, false, nil
	}

	res, err := handleSurvey(name, s)
	if err != nil {
		return nil, false, err
	}

	defs, save, err := f.ValidateAnswers(res)
	if err != nil {
		return nil, false, err
	}

	return defs, save, nil
}

func generateTemplates(options *InitOptions, answers *initSurveyAnswers, svc *client.Service) error {
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
	if err := generateSources(options, answers, svc); err != nil {
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

func generateSources(options *InitOptions, answers *initSurveyAnswers, svc *client.Service) error {
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

func generateTemplateContext(options *InitOptions, answers *initSurveyAnswers, externalTemplate *mtemplate.Template) (TemplateContext, error) {
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

func generateNewServiceArgs(answers *initSurveyAnswers, externalTemplate *mtemplate.Template) (string, error) {
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

func generateImports(answers *initSurveyAnswers) map[string][]ImportContext {
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
	}, assets.Files)
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
