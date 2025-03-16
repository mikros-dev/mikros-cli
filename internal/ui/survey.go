package ui

import (
	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

func SurveyConfirmBefore(s *survey.Survey) bool {
	return s.ConfirmQuestion != nil && !s.ConfirmQuestion.ConfirmAfter
}

func SurveyConfirmAfter(s *survey.Survey) bool {
	return s.ConfirmQuestion != nil && s.ConfirmQuestion.ConfirmAfter
}

func SurveyNeedsConfirmation(s *survey.Survey) bool {
	return s.ConfirmQuestion != nil
}
