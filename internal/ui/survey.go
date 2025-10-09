package ui

import (
	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

// SurveyConfirmBefore checks if the survey's confirmation question should be
// asked before the main survey execution.
func SurveyConfirmBefore(s *survey.Survey) bool {
	return s.ConfirmQuestion != nil && !s.ConfirmQuestion.ConfirmAfter
}

// SurveyConfirmAfter checks if the survey's confirmation question should be
// asked after the main survey execution.
func SurveyConfirmAfter(s *survey.Survey) bool {
	return s.ConfirmQuestion != nil && s.ConfirmQuestion.ConfirmAfter
}

// SurveyNeedsConfirmation determines if the survey includes a confirmation
// question.
func SurveyNeedsConfirmation(s *survey.Survey) bool {
	return s.ConfirmQuestion != nil
}
