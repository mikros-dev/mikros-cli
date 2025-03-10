package service

import (
	"errors"
	"fmt"
	"slices"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-playground/validator/v10"
	"github.com/somatech1/mikros/components/definition"
	"github.com/somatech1/mikros/components/options"
	"github.com/somatech1/mikros/components/plugin"

	"github.com/somatech1/mikros-cli/internal/answers"
	"github.com/somatech1/mikros-cli/internal/templates"
	msurvey "github.com/somatech1/mikros-cli/pkg/survey"
)

func runInitSurvey(opt *InitOptions) (*answers.InitSurveyAnswers, error) {
	surveyAnswers := answers.NewInitSurveyAnswers(opt.Language)

	if err := survey.Ask(baseQuestions(opt), surveyAnswers); err != nil {
		return nil, err
	}

	if err := runServiceTypeSurvey(surveyAnswers, opt); err != nil {
		return nil, err
	}

	// Presents only questions from selected features
	for _, name := range surveyAnswers.Features {
		defs, save, err := runFeatureSurvey(name, opt)
		if err != nil {
			return nil, err
		}
		if defs != nil {
			surveyAnswers.AddFeatureDefinitions(name, defs, save)
		}
	}

	return surveyAnswers, nil
}

func baseQuestions(opt *InitOptions) []*survey.Question {
	supportedTypes := []string{
		definition.ServiceType_gRPC.String(),
		definition.ServiceType_HTTP.String(),
		definition.ServiceType_Native.String(),
		definition.ServiceType_Script.String(),
	}

	if opt.Services != nil {
		for name := range opt.Services.Services() {
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
	}

	if opt.Language == templates.LanguageGolang {
		questions = append(questions, &survey.Question{
			Name: "lifecycle",
			Prompt: &survey.MultiSelect{
				Message: "Select lifecycle events to handle in the service:",
				Options: []string{"start", "finish"},
			},
		})

		// Features is only supported for golang services by now.
		if opt.Features != nil {
			var (
				featureNames = opt.FeatureNames
				iter         = opt.Features.Iterator()
			)

			for f, next := iter.Next(); next; f, next = iter.Next() {
				if api, ok := f.(msurvey.CLIFeature); ok {
					if api.IsCLISupported() {
						featureNames = append(featureNames, getFeatureUIName(f))
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
	}

	if opt.Language == templates.LanguageRust {
		questions = append(questions, &survey.Question{
			Name: "lifecycle_rust",
			Prompt: &survey.Confirm{
				Message: "Enable lifecycle events to handle in the service?",
			},
		})
	}

	return questions
}

func getFeatureUIName(feature plugin.Feature) string {
	if api, ok := feature.(msurvey.FeatureSurveyUI); ok {
		return api.UIName()
	}

	return feature.Name()
}

// runServiceTypeSurvey executes the survey that a specific service type may
// have implemented.
func runServiceTypeSurvey(answers *answers.InitSurveyAnswers, opt *InitOptions) error {
	if opt.Services == nil {
		return nil
	}

	s, ok := opt.Services.Services()[answers.Type]
	if !ok {
		return nil
	}

	cli, ok := s.(msurvey.CLIFeature)
	if !ok || !cli.IsCLISupported() {
		return nil
	}

	api, ok := s.(msurvey.FeatureSurvey)
	if !ok {
		return nil
	}

	svcSurvey := api.GetSurvey()
	if s == nil {
		return nil
	}

	response, err := handleSurvey(answers.Type, svcSurvey)
	if err != nil {
		return err
	}

	d, save, err := api.Answers(response)
	if err != nil {
		return err
	}

	answers.SetServiceDefinitions(d, save)
	return nil
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

func buildSurveyPrompt(name string, question *msurvey.Question) survey.Prompt {
	switch question.Prompt {
	case msurvey.PromptInput:
		return &survey.Input{
			Message: fmt.Sprintf("[%s] %s", name, question.Message),
			Default: question.Default,
		}

	case msurvey.PromptSelect:
		return &survey.Select{
			Message:  fmt.Sprintf("[%s] %s", name, question.Message),
			Options:  question.Options,
			PageSize: len(question.Options),
			Default:  question.Default,
		}

	case msurvey.PromptMultiSelect:
		return &survey.MultiSelect{
			Message: fmt.Sprintf("[%s] %s", name, question.Message),
			Options: question.Options,
		}

	case msurvey.PromptMultiline:
		return &survey.Multiline{
			Message: fmt.Sprintf("[%s] %s", name, question.Message),
		}

	case msurvey.PromptConfirm:
		return &survey.Confirm{
			Message: fmt.Sprintf("[%s] %s", name, question.Message),
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

func runFeatureSurvey(name string, opt *InitOptions) (interface{}, bool, error) {
	f, err := opt.Features.Feature(name)
	if err != nil {
		// Search again using mikros feature prefix. Maybe it is an implementation
		// of a mikros feature.
		f, err = opt.Features.Feature(options.FeatureNamePrefix + name)
		if err != nil {
			return nil, false, err
		}
	}

	api, ok := f.(msurvey.FeatureSurvey)
	if !ok {
		return nil, false, nil
	}

	var response map[string]interface{}

	if s := api.GetSurvey(); s != nil {
		res, err := handleSurvey(name, s)
		if err != nil {
			return nil, false, err
		}
		response = res
	}

	defs, save, err := api.Answers(response)
	if err != nil {
		return nil, false, err
	}

	return defs, save, nil
}
