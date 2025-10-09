package service

import (
	"github.com/mikros-dev/mikros/components/definition"

	"github.com/mikros-dev/mikros-cli/internal/protobuf"
)

// TemplateContext represents the context required for template generation in a
// service.
type TemplateContext struct {
	featuresExtensions bool
	servicesExtensions bool
	onStartLifecycle   bool
	onFinishLifecycle  bool
	serviceType        string

	ExternalFeaturesArg      string
	ExternalServicesArg      string
	NewServiceArgs           string
	ServiceName              string
	GrpcMethods              []*protobuf.Method
	Imports                  map[string][]ImportContext
	ServiceTypeCustomAnswers interface{}
	PluginData               interface{}
}

// ImportContext represents an import statement with its path and an
// optional alias.
type ImportContext struct {
	Alias string
	Path  string
}

// IsScriptService checks if the service type in the TemplateContext is
// classified as a script-based service.
func (t TemplateContext) IsScriptService() bool {
	return t.serviceType == definition.ServiceType_Script.String()
}

// IsWorkerService checks if the service type in the TemplateContext is
// classified as a worker service.
func (t TemplateContext) IsWorkerService() bool {
	return t.serviceType == definition.ServiceType_Worker.String()
}

// IsGrpcService determines if the service type in the TemplateContext is
// classified as a gRPC service.
func (t TemplateContext) IsGrpcService() bool {
	return t.serviceType == definition.ServiceType_gRPC.String()
}

// IsHTTPSpecService checks if the service type in the TemplateContext is
// classified as an HTTP-based service.
func (t TemplateContext) IsHTTPSpecService() bool {
	return t.serviceType == definition.ServiceType_HTTPSpec.String()
}

// HasGrpcMethods checks if the TemplateContext contains any defined gRPC
// methods.
func (t TemplateContext) HasGrpcMethods() bool {
	return len(t.GrpcMethods) > 0
}

// ServiceType returns the type of service associated with the TemplateContext
// as a string.
func (t TemplateContext) ServiceType() string {
	return t.serviceType
}

// HasFeaturesExtensions checks if the TemplateContext includes feature
// extensions.
func (t TemplateContext) HasFeaturesExtensions() bool {
	return t.featuresExtensions
}

// HasServicesExtensions checks if the TemplateContext includes service
// extensions.
func (t TemplateContext) HasServicesExtensions() bool {
	return t.servicesExtensions
}

// GetTemplateImports retrieves the list of ImportContext entries
// associated with a specific template name.
func (t TemplateContext) GetTemplateImports(templateName string) []ImportContext {
	return t.Imports[templateName]
}

// HasOnStart checks whether the TemplateContext has the OnStart
// lifecycle flag set to true.
func (t TemplateContext) HasOnStart() bool {
	return t.onStartLifecycle
}

// HasOnFinish checks whether the TemplateContext has the OnFinish
// lifecycle flag set to true.
func (t TemplateContext) HasOnFinish() bool {
	return t.onFinishLifecycle
}
