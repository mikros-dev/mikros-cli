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
	"github.com/somatech1/mikros/components/plugin"

	"github.com/somatech1/mikros-cli/internal/definitions"
	"github.com/somatech1/mikros-cli/internal/golang"
	"github.com/somatech1/mikros-cli/internal/path"
	"github.com/somatech1/mikros-cli/internal/protobuf"
	msurvey "github.com/somatech1/mikros-cli/pkg/survey"
)

type InitOptions struct {
	Path          string
	ProtoFilename string
	FeatureNames  []string
	Features      *plugin.FeatureSet
	Services      *plugin.ServiceSet
}

// Init initializes a new service locally.
func Init(options *InitOptions) error {
	var (
		featuresDefinitions = make(map[string]interface{})
		serviceDefinitions  = make(map[string]interface{})
		answers             initSurveyAnswers
	)

	if err := survey.Ask(baseQuestions(options), &answers); err != nil {
		return err
	}

	defs, err := runServiceAPISurvey(&answers, options)
	if err != nil {
		return err
	}
	if defs != nil {
		serviceDefinitions[answers.Type] = defs
	}

	// Presents only questions from selected features
	for _, name := range answers.Features {
		defs, err := runFeatureSurvey(name, options)
		if err != nil {
			return err
		}
		if defs != nil {
			featuresDefinitions[name] = defs
		}
	}

	if err := generateTemplates(options, &answers, featuresDefinitions, serviceDefinitions); err != nil {
		return err
	}

	return nil
}

func baseQuestions(options *InitOptions) []*survey.Question {
	supportedTypes := []string{
		definition.ServiceType_gRPC.String(),
		definition.ServiceType_HTTP.String(),
		definition.ServiceType_Native.String(),
		definition.ServiceType_Script.String(),
	}

	if options.Services != nil {
		for name := range options.Services.Services() {
			supportedTypes = append(supportedTypes, name)
		}
	}

	sort.Strings(supportedTypes)
	questions := []*survey.Question{
		// Service name
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Name. Can be a fully qualified service name (URL + name):",
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

	if options.Features != nil {
		var (
			featureNames = options.FeatureNames
			iter         = options.Features.Iterator()
		)

		for f, next := iter.Next(); next; f, next = iter.Next() {
			if api, ok := f.(msurvey.CLIFeature); ok {
				if api.IsCLISupported() {
					featureNames = append(featureNames, f.Name())
				}
			}
		}

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

	return questions
}

// runServiceAPISurvey executes the survey that a service may have implemented.
func runServiceAPISurvey(answers *initSurveyAnswers, options *InitOptions) (interface{}, error) {
	if options.Services == nil {
		return nil, nil
	}

	s, ok := options.Services.Services()[answers.Type]
	if !ok {
		return nil, nil
	}

	cli, ok := s.(msurvey.CLIFeature)
	if !ok || !cli.IsCLISupported() {
		return nil, nil
	}

	api, ok := s.(msurvey.FeatureSurvey)
	if !ok {
		return nil, nil
	}

	svcSurvey := api.GetSurvey()
	if s == nil {
		return nil, nil
	}

	response, err := handleSurvey(answers.Type, svcSurvey)
	if err != nil {
		return nil, err
	}

	d, err := api.Answers(response)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func handleSurvey(name string, featureSurvey *msurvey.Survey) (map[string]interface{}, error) {
	if featureSurvey.ConfirmQuestion != nil {
		var responses []map[string]interface{}

	loop:
		for {
			if !featureSurvey.ConfirmQuestion.ConfirmAfter {
				res := msurvey.YesNo(featureSurvey.ConfirmQuestion.Message)
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
				res := msurvey.YesNo(featureSurvey.ConfirmQuestion.Message)
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
				if q.Validate != nil {
					return q.Validate
				}

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
			if v, ok := r.(string); ok && v == condition.Value {
				return true
			}
		}

		return false
	}

	return true
}

func askOne(prompt survey.Prompt, question *msurvey.Question) (interface{}, error) {
	getOptions := func() survey.AskOpt {
		if question.Validate != nil {
			return survey.WithValidator(question.Validate)
		}
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

func runFeatureSurvey(name string, options *InitOptions) (interface{}, error) {
	f, err := options.Features.Feature(name)
	if err != nil {
		return nil, err
	}

	api, ok := f.(msurvey.FeatureSurvey)
	if !ok {
		return nil, nil
	}

	var response map[string]interface{}

	if s := api.GetSurvey(); s != nil {
		res, err := handleSurvey(name, s)
		if err != nil {
			return nil, err
		}
		response = res
	}

	defs, err := api.Answers(response)
	if err != nil {
		return nil, err
	}

	return defs, nil
}

func generateTemplates(options *InitOptions, answers *initSurveyAnswers, featuresDefinitions, serviceDefinitions map[string]interface{}) error {
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
	if err := writeServiceDefinitions(destinationPath, answers, featuresDefinitions, serviceDefinitions); err != nil {
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

func writeServiceDefinitions(path string, answers *initSurveyAnswers, featureDefinitions, serviceDefinitions map[string]interface{}) error {
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

	for name, d := range featureDefinitions {
		if err := definitions.AppendFeature(path, name, d); err != nil {
			return err
		}
	}

	for n, d := range serviceDefinitions {
		if err := definitions.AppendService(path, n, d); err != nil {
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

	if err := runTemplates(answers.TemplateNames(), context); err != nil {
		return err
	}

	return nil
}

func generateTemplateContext(options *InitOptions, answers *initSurveyAnswers) (TemplateContext, error) {
	externalService := func() bool {
		switch answers.Type {
		case definition.ServiceType_gRPC.String(), definition.ServiceType_HTTP.String(), definition.ServiceType_Script.String(), definition.ServiceType_Native.String():
			return false
		}

		return true
	}

	context := TemplateContext{
		featuresExtensions: len(answers.Features) > 0,
		servicesExtensions: externalService(),
		onStartLifecycle:   slices.Contains(answers.Lifecycle, "OnStart"),
		onFinishLifecycle:  slices.Contains(answers.Lifecycle, "OnFinish"),
		serviceType:        answers.Type,
		NewServiceArgs:     generateNewServiceArgs(answers),
		ServiceName:        answers.Name,
		Imports:            generateImports(answers),
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

func generateNewServiceArgs(answers *initSurveyAnswers) string {
	svcSnake := strcase.ToSnake(answers.Name)
	switch answers.Type {
	case definition.ServiceType_gRPC.String():
		return fmt.Sprintf(`Service: map[string]options.ServiceOptions{
			"grpc": &options.GrpcServiceOptions{
				ProtoServiceDescription: &%spb.GrpcService_ServiceDesc,
			},
		},`, svcSnake)

	case definition.ServiceType_HTTP.String():
		return fmt.Sprintf(`Service: map[string]options.ServiceOptions{
			"http": &options.HttpServiceOptions{
				ProtoHttpServer: %spb.NewHttpServer(),
			},
		},`, svcSnake)

	case definition.ServiceType_Native.String():
		return `Service: map[string]options.ServiceOptions{
			"native": &options.NativeServiceOptions{},
		},`

	case definition.ServiceType_Script.String():
		return `Service: map[string]options.ServiceOptions{
			"script": &options.ScriptServiceOptions{},
		},`
	}

	return ""
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
