package ui

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

func RunFormFromSurvey(name string, s *survey.Survey, theme *huh.Theme) (map[string]interface{}, error) {
	if SurveyNeedsConfirmation(s) {
		return runFormWithConfirmation(name, s, theme)
	}

	return runFormSurvey(name, s, theme)
}

func runFormWithConfirmation(name string, s *survey.Survey, theme *huh.Theme) (map[string]interface{}, error) {
	var results []map[string]interface{}

loop:
	for {
		if SurveyConfirmBefore(s) {
			confirm, err := yesNo(s.ConfirmQuestion.Message)
			if err != nil {
				return nil, err
			}
			if !confirm {
				break loop
			}
		}

		response, err := runFormSurvey(name, s, theme)
		if err != nil {
			return nil, err
		}
		results = append(results, response)

		if SurveyConfirmAfter(s) {
			confirm, err := yesNo(s.ConfirmQuestion.Message)
			if err != nil {
				return nil, err
			}
			if !confirm {
				break loop
			}
		}
	}

	return map[string]interface{}{
		name: results,
	}, nil
}

func runFormSurvey(name string, s *survey.Survey, theme *huh.Theme) (map[string]interface{}, error) {
	var (
		values   = make(map[string]interface{})
		results  = make(map[string]interface{})
		elements []huh.Field
	)

	for _, q := range s.Questions {
		var (
			title = fmt.Sprintf("[%s] %s", name, q.Message)
		)

		switch q.Prompt {
		case survey.PromptSurvey:
			res, err := RunFormFromSurvey(name, q.Survey, theme)
			if err != nil {
				return nil, err
			}
			if res != nil {
				results[q.Name] = res[q.Name]
			}

			continue

		case survey.PromptInput:
			defaultValue := q.Default
			values[q.Name] = &defaultValue

			input := huh.NewInput().Title(title).Value(values[q.Name].(*string))
			if q.Required {
				input = input.Validate(IsEmpty("cannot be empty"))
			}

			elements = append(elements, input)

		case survey.PromptSelect:
			options := make([]huh.Option[string], len(q.Options))
			for i, option := range q.Options {
				opt := huh.NewOption(option, option)
				if option == q.Default {
					opt.Selected(true)
				}
				options[i] = opt
			}

			values[q.Name] = new(string)
			elements = append(elements, huh.NewSelect[string]().
				Title(title).
				Options(options...).
				Value(values[q.Name].(*string)))

		case survey.PromptMultiSelect:
			options := make([]huh.Option[string], len(q.Options))
			for i, option := range q.Options {
				options[i] = huh.NewOption(option, option)
			}

			values[q.Name] = new([]string)
			prompt := huh.NewMultiSelect[string]().Title(title).Options(options...).Value(values[q.Name].(*[]string))
			if q.Required {
				prompt = prompt.Validate(func(strings []string) error {
					if len(strings) == 0 {
						return errors.New("must choose at least one option")
					}

					return nil
				})
			}

			elements = append(elements, prompt)

		case survey.PromptMultiline:
			values[q.Name] = new(string)
			text := huh.NewText().Title(title).Value(values[q.Name].(*string))
			if q.Required {
				text = text.Validate(IsEmpty("cannot be empty"))
			}

			elements = append(elements, text)

		case survey.PromptConfirm:
			defaultValue := false
			if q.Default != "" {
				if b, err := strconv.ParseBool(q.Default); err == nil {
					defaultValue = b
				}
			}

			values[q.Name] = &defaultValue
			elements = append(elements, huh.NewConfirm().
				Title(title).
				Value(values[q.Name].(*bool)))
		}
	}

	form := huh.NewForm(huh.NewGroup(elements...))
	if err := form.WithTheme(theme).Run(); err != nil {
		return nil, err
	}

	for k, v := range values {
		switch v := v.(type) {
		case *string:
			results[k] = *v
		case *bool:
			results[k] = *v
		case *[]string:
			results[k] = *v
		}
	}

	return results, nil
}

func yesNo(message string) (bool, error) {
	var confirm bool
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&confirm),
		),
	)

	if err := f.Run(); err != nil {
		return false, err
	}

	return confirm, nil
}
