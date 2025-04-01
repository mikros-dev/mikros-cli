//package service
//
//import (
//	"embed"
//	"errors"
//	"os"
//	"path/filepath"
//	"strings"
//
//	"github.com/somatech1/mikros/components/definition"
//	"github.com/somatech1/mikros/components/plugin"
//
//	"github.com/somatech1/mikros-cli/internal/answers"
//	"github.com/somatech1/mikros-cli/internal/golang"
//	"github.com/somatech1/mikros-cli/internal/rust"
//	"github.com/somatech1/mikros-cli/internal/templates"
//	"github.com/somatech1/mikros-cli/pkg/definitions"
//	"github.com/somatech1/mikros-cli/pkg/path"
//	mtemplates "github.com/somatech1/mikros-cli/pkg/templates"
//)
//
//type TemplateExecutioner interface {
//	// PreExecution must be responsible for initializing the service base for
//	// the new service that will be created.
//	PreExecution(serviceName, destinationPath string) error
//
//	// GenerateContext is responsible for generating the context that will be
//	// used inside the templates when executed.
//	GenerateContext(answers *answers.InitSurveyAnswers, protoFilename string, externalTemplates interface{}) (interface{}, error)
//
//	// Templates must return all templates that the executioner will use.
//	Templates() []mtemplates.TemplateFile
//
//	// Files must return all .tmpl files that can be used by the executioner.
//	Files() embed.FS
//
//	// PostExecution is another handler where custom modifications can be made
//	// into the service that is being created. At this point, both initial source
//	// code and the service.toml file are already created.
//	PostExecution(destinationPath string) error
//}

//type InitOptions struct {
//	Language          templates.Language
//	Path              string
//	ProtoFilename     string
//	FeatureNames      []string
//	Features          *plugin.FeatureSet
//	Services          *plugin.ServiceSet
//	ExternalTemplates *golang.ExternalTemplates
//}
//
//// Init initializes a new service locally.
//func Init(options *InitOptions) error {
//	surveyAnswers, err := runInitSurvey(options)
//	if err != nil {
//		return err
//	}
//
//	if err := generateTemplates(options, surveyAnswers); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func generateTemplates(options *InitOptions, answers *answers.InitSurveyAnswers) error {
//	var (
//		destinationPath = options.Path
//	)
//
//	// Set the project base path
//	if destinationPath == "" {
//		cwd, err := os.Getwd()
//		if err != nil {
//			return err
//		}
//
//		destinationPath = filepath.Join(cwd, strings.ToLower(answers.Name))
//	}
//
//	executioner, err := getTemplateExecutioner(options.Language)
//	if err != nil {
//		return err
//	}
//
//	if err := executioner.PreExecution(answers.Name, destinationPath); err != nil {
//		return err
//	}
//
//	if err := generateSources(destinationPath, options, answers, executioner); err != nil {
//		return err
//	}
//
//	if err := writeServiceDefinitions(destinationPath, answers); err != nil {
//		return err
//	}
//
//	if err := executioner.PostExecution(destinationPath); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func getTemplateExecutioner(language templates.Language) (TemplateExecutioner, error) {
//	if language == templates.LanguageGolang {
//		return &golang.Executioner{}, nil
//	}
//	if language == templates.LanguageRust {
//		return &rust.Executioner{}, nil
//	}
//
//	return nil, errors.New("unsupported language")
//}
//
//func writeServiceDefinitions(path string, answers *answers.InitSurveyAnswers) error {
//	defs := &definition.Definitions{
//		Name:     answers.Name,
//		Types:    []string{answers.Type},
//		Version:  answers.Version,
//		Language: answers.Language(),
//		Product:  strings.ToUpper(answers.Product),
//	}
//
//	if err := definitions.Write(path, defs); err != nil {
//		return err
//	}
//
//	for name, d := range answers.FeatureDefinitions() {
//		if d.ShouldBeSaved() {
//			if err := definitions.AppendFeature(path, name, d.Definitions()); err != nil {
//				return err
//			}
//		}
//	}
//
//	if svcDefs := answers.ServiceDefinitions(); svcDefs != nil && svcDefs.ShouldBeSaved() {
//		if err := definitions.AppendService(path, answers.Type, svcDefs.Definitions()); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func generateSources(destinationPath string, options *InitOptions, answers *answers.InitSurveyAnswers, executioner TemplateExecutioner) error {
//	// Switch to the service folder so the initial source code can be created.
//	cwd, err := path.ChangeDir(destinationPath)
//	if err != nil {
//		return err
//	}
//
//	defer func() {
//		if e := os.Chdir(cwd); e != nil {
//			err = e
//		}
//	}()
//
//	context, err := executioner.GenerateContext(answers, options.ProtoFilename, options.ExternalTemplates)
//	if err != nil {
//		return err
//	}
//
//	if err := createServiceTemplates(options, executioner, context); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func createServiceTemplates(options *InitOptions, executioner TemplateExecutioner, context interface{}) error {
//	if err := runTemplates(executioner.Files(), executioner.Templates(), context, nil); err != nil {
//		return err
//	}
//
//	if options != nil && options.ExternalTemplates != nil {
//		if err := runTemplates(options.ExternalTemplates.Files, options.ExternalTemplates.Templates, context, options.ExternalTemplates.Api); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//func runTemplates(files embed.FS, filenames []mtemplates.TemplateFile, context interface{}, api map[string]interface{}) error {
//	tpls, err := mtemplates.Load(&mtemplates.LoadOptions{
//		TemplateNames: filenames,
//		Files:         files,
//		Api:           api,
//	})
//	if err != nil {
//		return err
//	}
//
//	generated, err := tpls.Execute(context)
//	if err != nil {
//		return err
//	}
//
//	for _, gen := range generated {
//		file, err := os.Create(gen.Filename())
//		if err != nil {
//			return err
//		}
//
//		if _, err := file.Write(gen.Content()); err != nil {
//			_ = file.Close()
//			return err
//		}
//
//		_ = file.Close()
//	}
//
//	return nil
//}
