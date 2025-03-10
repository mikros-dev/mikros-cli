package survey

// CLIFeature is a behavior that a client should have to let the mikros CLI
// know more information about it.
type CLIFeature interface {
	// IsCLISupported tells, by returning true or false, if the client is
	// available to be manipulated by the mikros CLI tool.
	IsCLISupported() bool
}

// FeatureSurvey is a behavior that a client should have to be used by the
// mikros CLI tool.
type FeatureSurvey interface {
	// GetSurvey must return a list of Question objects to tell mikros CLI how
	// the settings will be prompted to the user.
	GetSurvey() *Survey

	// Answers will receive the survey answers in a map where the key will be
	// the Name field of each Question. The value format directly depends on
	// how the survey was assembled. It should return the client settings
	// structure in case of success and a flag indicating if they should be
	// written in the service.toml file or not.
	Answers(answers map[string]interface{}) (interface{}, bool, error)
}

// FeatureSurveyUI when implemented by a mikros feature can override some of
// its information to be used by the mikros CLI UI.
type FeatureSurveyUI interface {
	UIName() string
}

// Survey is a structure that a client uses to tell mikros CLI how to present
// its survey for the user to answer questions.
type Survey struct {
	// AskOne when true sets that the survey will be executed each question
	// separately.
	AskOne bool

	// ConfirmQuestion is a question that will inform mikros CLI that the
	// following questions will be asked inside a loop, until the user
	// decides to stop.
	ConfirmQuestion *Question

	// Questions gathers a list of questions that will be presented to the
	// user.
	Questions []*Question
}

type Question struct {
	Required     bool
	ConfirmAfter bool
	Prompt       PromptKind `validate:"required"`
	Message      string
	Name         string `validate:"required"`
	Default      string
	Condition    *QuestionCondition
	Options      []string
	Validate     func(v interface{}) error
	Survey       *Survey `validate:"-"`
}

type QuestionCondition struct {
	Name  string
	Value interface{}
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
