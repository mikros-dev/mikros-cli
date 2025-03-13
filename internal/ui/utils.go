package ui

import (
	"github.com/AlecAivazis/survey/v2"
)

func YesNo(message string) bool {
	res := false
	prompt := &survey.Confirm{
		Message: message,
	}

	if err := survey.AskOne(prompt, &res); err != nil {
		return false
	}

	return res
}
