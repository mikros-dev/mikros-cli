package protobuf

// TemplateContext holds contextual information used during template execution.
type TemplateContext struct {
	MainPackageName  string
	RepositoryName   string
	VCSProjectPrefix string
}
