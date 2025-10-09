package template

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
