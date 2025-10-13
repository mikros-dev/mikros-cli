package plugin

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

// Question defines the structure for a survey question, containing its
// attributes such as name, prompt, and messaging.
type Question struct {
	Required     bool       `json:"required"`
	ConfirmAfter bool       `json:"confirm_after"`
	Prompt       PromptKind `json:"prompt" validate:"required"`
	Message      string     `json:"message,omitempty"`
	Name         string     `json:"name" validate:"required"`
	Default      string     `json:"default,omitempty"`
	Options      []string   `json:"options,omitempty"`
}

// FollowUpSurvey defines a structure for secondary surveys triggered by
// specific conditions during a primary survey.
type FollowUpSurvey struct {
	Name      string             `json:"name,omitempty" validate:"required"`
	Condition *QuestionCondition `json:"condition,omitempty" validate:"required"`
	Survey    *Survey            `json:"survey,omitempty" validate:"-"`
}

// QuestionCondition defines a structure used to represent a condition for
// triggering specific survey actions.
type QuestionCondition struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
}

// PromptKind represents the type of prompt used in a survey or user interaction
// mechanism.
type PromptKind int

// Supported prompt types.
const (
	PromptInput PromptKind = iota + 1
	PromptSelect
	PromptMultiSelect
	PromptMultiline
	PromptConfirm
)

// Template represents a structure containing information related to creating
// and managing custom templates.
type Template struct {
	// NewServiceArgs Allows adding custom content for external service kind
	// when creating the main file of a new template service. It will be available
	// inside the {{.NewServiceArgs}} inside the template.
	//
	// It can be a template string that receives the internal default API and
	// the custom API as functions available to be used as well as a short data
	// context with the following fields:
	//
	// {{.ServiceName}}: holding the current service name.
	// {{.ServiceType}}: holding the current service type.
	// {{.ServiceTypeCustomAnswers}}: holding custom answers related to the current
	//								  service type.
	NewServiceArgs string `json:"new_service_args,omitempty"`

	// WithExternalFeaturesArg Sets the argument of method .WithExternalFeatures()
	// inside the main template.
	WithExternalFeaturesArg string `json:"with_external_features_arg,omitempty"`

	// WithExternalServicesArg Sets the argument of method .WithExternalServices()
	// inside the main template.
	WithExternalServicesArg string `json:"with_external_services_arg,omitempty"`

	// Templates contains a list of custom template files that will be generated
	// when the service is selected for a service.
	Templates []*File `json:"templates,omitempty"`
}

// File represents a file structure with customizable content, name, output path,
// extension, and additional context.
type File struct {
	Content   string `json:"content,omitempty"`
	Name      string `json:"name,omitempty"`
	Output    string `json:"output,omitempty"`
	Extension string `json:"extension,omitempty"`

	// Context is a custom context to be used inside custom templates exported
	// by the plugin.
	Context interface{} `json:"context,omitempty"`
}
