package survey

// Survey is a structure that a client uses to tell mikros CLI how to present
// its survey for the user to answer questions.
type Survey struct {
	// AskOne when true sets that the survey will be executed each question
	// separately.
	AskOne bool `json:"ask_one"`

	// ConfirmQuestion is a question that will inform mikros CLI that the
	// following questions will be asked inside a loop, until the user
	// decides to stop.
	ConfirmQuestion *Question `json:"confirm_question"`

	// Questions gathers a list of questions that will be presented to the
	// user.
	Questions []*Question `json:"questions"`
}

type Question struct {
	Required     bool               `json:"required"`
	ConfirmAfter bool               `json:"confirm_after"`
	Prompt       PromptKind         `json:"prompt" validate:"required"`
	Message      string             `json:"message,omitempty"`
	Name         string             `json:"name" validate:"required"`
	Default      string             `json:"default,omitempty"`
	Condition    *QuestionCondition `json:"condition,omitempty"`
	Options      []string           `json:"options,omitempty"`
	Survey       *Survey            `json:"survey,omitempty" validate:"-"`
}

type QuestionCondition struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type PromptKind int

const (
	PromptInput PromptKind = iota + 1
	PromptSelect
	PromptMultiSelect
	PromptMultiline
	PromptConfirm
	PromptSurvey
)
