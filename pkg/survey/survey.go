package survey

// Survey is a structure that a client uses to tell mikros CLI how to present
// its survey for the user to answer questions.
type Survey struct {
	// ConfirmQuestion is a question that will inform mikros CLI that the
	// following questions will be asked inside a loop, until the user
	// decides to stop.
	ConfirmQuestion *Question `json:"confirm_question,omitempty"`

	// Questions gathers a list of questions that will be presented to the
	// user.
	Questions []*Question `json:"questions,omitempty"`

	// FollowUp holds any other survey that must run according to an
	// internal condition (these must have Condition adjusted or won't
	// be validated and executed).
	FollowUp []*FollowUpSurvey `json:"sub_survey,omitempty"`
}

type Question struct {
	Required     bool       `json:"required"`
	ConfirmAfter bool       `json:"confirm_after"`
	Prompt       PromptKind `json:"prompt" validate:"required"`
	Message      string     `json:"message,omitempty"`
	Name         string     `json:"name" validate:"required"`
	Default      string     `json:"default,omitempty"`
	Options      []string   `json:"options,omitempty"`
}

type FollowUpSurvey struct {
	Name      string             `json:"name,omitempty" validate:"required"`
	Condition *QuestionCondition `json:"condition,omitempty" validate:"required"`
	Survey    *Survey            `json:"survey,omitempty" validate:"-"`
}

type QuestionCondition struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

type PromptKind int

const (
	PromptInput PromptKind = iota + 1
	PromptSelect
	PromptMultiSelect
	PromptMultiline
	PromptConfirm
)
