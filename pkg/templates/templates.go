package templates

// TemplateFile is representation of a template file to be processed when
// creating new template services.
type TemplateFile struct {
	// Name is the template name without its extensions (.tmpl). It is used
	// as the file final name if Output is empty.
	Name string

	// Output is an optional name that the template can have after it is
	// processed (it replaces the Name member).
	Output string

	// Extension is an optional field to set the file extension.
	Extension string
}
