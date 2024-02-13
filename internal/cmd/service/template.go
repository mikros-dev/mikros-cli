package service

import (
	"os"

	"github.com/somatech1/mikros/components/definition"

	assets "github.com/somatech1/mikros-cli/internal/assets/templates"
	"github.com/somatech1/mikros-cli/internal/protobuf"
	"github.com/somatech1/mikros-cli/internal/templates"
)

type TemplateContext struct {
	featuresExtensions bool
	servicesExtensions bool
	onStartLifecycle   bool
	onFinishLifecycle  bool
	serviceType        string

	NewServiceArgs string
	ServiceName    string
	GrpcMethods    []*protobuf.Method
	Imports        map[string][]ImportContext
}

type ImportContext struct {
	Alias string
	Path  string
}

func (t TemplateContext) IsScriptService() bool {
	return t.serviceType == definition.ServiceType_Script.String()
}

func (t TemplateContext) IsNativeService() bool {
	return t.serviceType == definition.ServiceType_Native.String()
}

func (t TemplateContext) IsGrpcService() bool {
	return t.serviceType == definition.ServiceType_gRPC.String()
}

func (t TemplateContext) IsHttpService() bool {
	return t.serviceType == definition.ServiceType_HTTP.String()
}

func (t TemplateContext) HasGrpcMethods() bool {
	return len(t.GrpcMethods) > 0
}

func (t TemplateContext) ServiceType() string {
	return t.serviceType
}

func (t TemplateContext) HasFeaturesExtensions() bool {
	return t.featuresExtensions
}

func (t TemplateContext) HasServicesExtensions() bool {
	return t.servicesExtensions
}

func (t TemplateContext) GetTemplateImports(templateName string) []ImportContext {
	return t.Imports[templateName]
}

func (t TemplateContext) HasOnStart() bool {
	return t.onStartLifecycle
}

func (t TemplateContext) HasOnFinish() bool {
	return t.onFinishLifecycle
}

func runTemplates(filenames []templates.TemplateName, context interface{}) error {
	tpls, err := templates.Load(&templates.LoadOptions{
		TemplateNames: filenames,
		Files:         assets.Files,
	})
	if err != nil {
		return err
	}

	generated, err := tpls.Execute(context)
	if err != nil {
		return err
	}

	for _, gen := range generated {
		file, err := os.Create(gen.Filename())
		if err != nil {
			return err
		}

		if _, err := file.Write(gen.Content()); err != nil {
			_ = file.Close()
			return err
		}

		_ = file.Close()
	}

	return nil
}
