package template

type Template struct {
	// NewServiceArgs Allows adding custom content for external service kind
	// when creating the main file of a new template service. It will be available
	// inside the {{.NewServiceArgs}} inside the template.
	//
	// It can be a template string that receives the internal default API and
	// the custom Api as functions available to be used as well as a short data
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

	Templates []*File `json:"templates,omitempty"`
}

type File struct {
	Content   string `json:"content,omitempty"`
	Name      string `json:"name,omitempty"`
	Output    string `json:"output,omitempty"`
	Extension string `json:"extension,omitempty"`
}
