package ui

import (
	"errors"
	"fmt"
	"slices"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/plugin/survey"
)

// FormOptions defines options for configuring a form, including theme settings
// and accessibility preferences.
type FormOptions struct {
	Theme      *huh.Theme
	Accessible bool
}

// RunFormFromSurvey executes a survey form and returns the collected data as a map.
func RunFormFromSurvey(name string, s *survey.Survey, options *FormOptions) (map[string]interface{}, error) {
	if SurveyNeedsConfirmation(s) {
		return runFormWithConfirmation(name, s, options)
	}

	return runFormSurvey(name, s, options)
}

func runFormWithConfirmation(name string, s *survey.Survey, options *FormOptions) (map[string]interface{}, error) {
	var (
		results    = make([]map[string]interface{}, 0, 1)
		askConfirm = func(enabled bool) (bool, error) {
			if !enabled {
				return true, nil
			}

			return yesNo(s.ConfirmQuestion.Message, s.ConfirmQuestion.Default)
		}
	)

	for {
		ok, err := askConfirm(SurveyConfirmBefore(s))
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}

		response, err := runFormSurvey(name, s, options)
		if err != nil {
			return nil, err
		}
		results = append(results, response)

		ok, err = askConfirm(SurveyConfirmAfter(s))
		if err != nil {
			return nil, err
		}
		if !ok {
			break
		}
	}

	return map[string]interface{}{
		name: results,
	}, nil
}

func runFormSurvey(name string, s *survey.Survey, options *FormOptions) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	elements, err := buildFormSurveyElements(name, s, values)
	if err != nil {
		return nil, err
	}

	form := huh.NewForm(huh.NewGroup(elements...)).
		WithTheme(options.Theme).
		WithAccessible(options.Accessible)

	if err := form.Run(); err != nil {
		return nil, err
	}

	results := extractResults(values)

	// Check if we have a follow-up survey to execute
	if len(s.FollowUp) != 0 {
		followUpResults, err := executeFollowUpSurvey(s.FollowUp, results, options)
		if err != nil {
			return nil, err
		}
		results["follow-up"] = followUpResults
	}

	return results, nil
}

func buildFormSurveyElements(name string, s *survey.Survey, values map[string]interface{}) ([]huh.Field, error) {
	elements := make([]huh.Field, 0, len(s.Questions))

	for _, q := range s.Questions {
		title := fmt.Sprintf("[%s] %s", name, q.Message)
		field, err := buildFormElementQuestion(q, title, values)
		if err != nil {
			return nil, err
		}
		elements = append(elements, field)
	}

	return elements, nil
}

func buildFormElementQuestion(q *survey.Question, title string, values map[string]interface{}) (huh.Field, error) {
	switch q.Prompt {
	case survey.PromptInput:
		return buildPromptInputQuestion(q, title, values), nil

	case survey.PromptSelect:
		return buildPromptSelectQuestion(q, title, values), nil

	case survey.PromptMultiSelect:
		return buildPromptMultiSelectQuestion(q, title, values)

	case survey.PromptMultiline:
		return buildPromptMultilineQuestion(q, title, values), nil

	case survey.PromptConfirm:
		return buildPromptConfirmQuestion(q, title, values), nil
	}

	return nil, errors.New("unsupported prompt type")
}

func buildPromptInputQuestion(q *survey.Question, title string, values map[string]interface{}) huh.Field {
	defaultValue := q.Default
	values[q.Name] = &defaultValue

	input := huh.NewInput().Title(title).Value(values[q.Name].(*string))
	if q.Required {
		input = input.Validate(IsEmpty("cannot be empty"))
	}

	return input
}

func buildPromptSelectQuestion(q *survey.Question, title string, values map[string]interface{}) huh.Field {
	options := make([]huh.Option[string], len(q.Options))
	for i, option := range q.Options {
		opt := huh.NewOption(option, option)
		if option == q.Default {
			opt = opt.Selected(true)
		}
		options[i] = opt
	}

	values[q.Name] = new(string)
	return huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(values[q.Name].(*string))
}

func buildPromptMultiSelectQuestion(
	q *survey.Question,
	title string,
	values map[string]interface{},
) (huh.Field, error) {
	options := make([]huh.Option[string], len(q.Options))
	for i, option := range q.Options {
		options[i] = huh.NewOption(option, option)
	}

	values[q.Name] = new([]string)
	prompt := huh.NewMultiSelect[string]().
		Title(title).Options(options...).
		Value(values[q.Name].(*[]string))

	if q.Required {
		prompt = prompt.Validate(func(strings []string) error {
			if len(strings) == 0 {
				return errors.New("must choose at least one option")
			}

			return nil
		})
	}

	return prompt, nil
}

func buildPromptMultilineQuestion(q *survey.Question, title string, values map[string]interface{}) huh.Field {
	values[q.Name] = new(string)
	text := huh.NewText().Title(title).Value(values[q.Name].(*string))
	if q.Required {
		text = text.Validate(IsEmpty("cannot be empty"))
	}

	return text
}

func buildPromptConfirmQuestion(q *survey.Question, title string, values map[string]interface{}) huh.Field {
	defaultValue := false
	if q.Default != "" {
		if b, err := strconv.ParseBool(q.Default); err == nil {
			defaultValue = b
		}
	}

	values[q.Name] = &defaultValue
	return huh.NewConfirm().
		Title(title).
		Value(values[q.Name].(*bool))
}

func extractResults(values map[string]interface{}) map[string]interface{} {
	results := make(map[string]interface{})

	// Dereferences collected values into a plain map.
	for k, v := range values {
		switch vv := v.(type) {
		case *string:
			results[k] = *vv
		case *bool:
			results[k] = *vv
		case *[]string:
			results[k] = *vv
		}
	}

	return results
}

func executeFollowUpSurvey(
	surveys []*survey.FollowUpSurvey,
	previousResults map[string]interface{},
	options *FormOptions,
) (map[string]map[string]interface{}, error) {
	results := make(map[string]map[string]interface{})

	for _, s := range surveys {
		// Check if the condition is met
		ok, err := checkFollowUpSurveyCondition(s, previousResults)
		if err != nil {
			return nil, err
		}
		if ok {
			r, err := RunFormFromSurvey(s.Name, s.Survey, options)
			if err != nil {
				return nil, err
			}

			results[s.Name] = r
			continue
		}
	}

	return results, nil
}

func checkFollowUpSurveyCondition(s *survey.FollowUpSurvey, previousResults map[string]interface{}) (bool, error) {
	switch v := s.Condition.Value.(type) {
	case string:
		result, ok := previousResults[s.Condition.Name]
		if !ok {
			return false, nil
		}

		resultValue, ok := result.(string)
		if !ok {
			return false, errors.New("invalid result value type found")
		}

		return resultValue == v, nil

	case []interface{}:
		result, ok := previousResults[s.Condition.Name]
		if !ok {
			return false, nil
		}

		resultValue, ok := result.(string)
		if !ok {
			return false, errors.New("invalid result value type found")
		}

		values := make([]string, len(v))
		for i, item := range v {
			values[i] = item.(string)
		}

		return slices.Contains(values, resultValue), nil
	}

	return false, nil
}

func yesNo(message, defaultValue string) (bool, error) {
	confirm := false
	if defaultValue != "" {
		if b, err := strconv.ParseBool(defaultValue); err == nil {
			confirm = b
		}
	}

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
