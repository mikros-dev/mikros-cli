package service

import (
	"errors"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/mikros-dev/mikros/components/definition"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/internal/plugin/client"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func runSurvey(cfg *settings.Settings, protoFilename string) (*surveyAnswers, error) {
	answers, err := newSurveyAnswers(protoFilename)
	if err != nil {
		return nil, err
	}

	questions, err := getBaseQuestions(answers, cfg)
	if err != nil {
		return nil, err
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

func getBaseQuestions(answers *surveyAnswers, cfg *settings.Settings) ([]huh.Field, error) {
	supportedTypes, err := getSupportedServiceTypes(cfg)
	if err != nil {
		return nil, err
	}

	var languages []huh.Option[string]
	for _, t := range definition.SupportedLanguages() {
		languages = append(languages, huh.NewOption(t, t))
	}

	return []huh.Field{
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
	}, nil
}

func getSupportedServiceTypes(cfg *settings.Settings) ([]huh.Option[string], error) {
	types := []huh.Option[string]{
		huh.NewOption(definition.ServiceTypeGRPC.String(), definition.ServiceTypeGRPC.String()),
		huh.NewOption(definition.ServiceTypeHTTP.String(), definition.ServiceTypeHTTP.String()),
		huh.NewOption(definition.ServiceTypeWorker.String(), definition.ServiceTypeWorker.String()),
		huh.NewOption(definition.ServiceTypeScript.String(), definition.ServiceTypeScript.String()),
	}

	newTypes, err := plugin.GetNewServiceKinds(cfg)
	if err != nil {
		return nil, err
	}
	for _, t := range newTypes {
		types = append(types, huh.NewOption(t, t))
	}
	sort.Slice(types, func(i, j int) bool {
		return types[i].String() < types[j].String()
	})

	return types, nil
}

// runServiceTypeSurvey executes the survey that a service may have implemented.
func runServiceTypeSurvey(cfg *settings.Settings, answers *surveyAnswers) (*client.Service, error) {
	svc, err := plugin.GetServicePlugin(cfg, answers.Type)
	if err != nil {
		return nil, err
	}
	if svc == nil {
		// No plugin for the chosen service type. But do we have specific settings
		// for the service type?
		if err := runCoreServiceTypeSurvey(cfg, answers); err != nil {
			return nil, err
		}

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

	answers.SetServiceAnswers(response)
	answers.SetServiceDefinitions(d)

	return svc, nil
}

func runCoreServiceTypeSurvey(cfg *settings.Settings, answers *surveyAnswers) error {
	// We only have specific questions for HTTP services by now
	if answers.Type != definition.ServiceTypeHTTP.String() {
		return nil
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select the HTTP service type").
				Options(
					huh.NewOption("Standard library HTTP", "http"),
					huh.NewOption("Protobuf spec HTTP", "http-spec"),
				).
				Value(&answers.HTTPType),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	return form.Run()
}

// runFeatureSurvey executes the survey that a feature may have implemented.
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
