package template

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"path"
	"path/filepath"
	"slices"
	"text/template"

	"github.com/iancoleman/strcase"
)

var defaultAPI = template.FuncMap{
	"toCamel":      strcase.ToCamel,
	"toSnake":      strcase.ToSnake,
	"toUpperSnake": strcase.ToScreamingSnake,
	"basename":     path.Base,
	"toKebab":      strcase.ToKebab,
}

// Session is the template session.
type Session struct {
	loadedTemplates []*Info
}

// Info is the template information.
type Info struct {
	template *template.Template
	name     File
	context  interface{}
}

// File is the representation of a template file to be processed when
// creating new template services.
type File struct {
	// Name is the template name without its extensions (.tmpl). It is used
	// as the file final name if Output is empty.
	Name string

	// Output is an optional name that the template can have after it is
	// processed (it replaces the Name member).
	Output string

	// Extension is an optional field to set the file extension.
	Extension string
}

// LoadOptions is the template loading options.
type LoadOptions struct {
	TemplatesToUse []File
	API            map[string]interface{}
	FilesBasePath  string
}

// NewSessionFromFiles creates a new template session from a set of files.
func NewSessionFromFiles(options *LoadOptions, files embed.FS) (*Session, error) {
	dirFiles, err := files.ReadDir(options.FilesBasePath)
	if err != nil {
		return nil, fmt.Errorf("reading files: %w", err)
	}

	var templates []*Info
	for _, file := range dirFiles {
		data, err := files.ReadFile(filepath.Join(options.FilesBasePath, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading file: %w", err)
		}

		var name = filenameWithoutExtension(file.Name())

		idx := slices.IndexFunc(options.TemplatesToUse, func(t File) bool {
			return t.Name == name
		})
		if idx == -1 {
			// The template is not being used at the moment.
			continue
		}

		tpl, err := loadTemplate(name, data, options)
		if err != nil {
			return nil, err
		}

		templates = append(templates, &Info{
			name:     options.TemplatesToUse[idx],
			template: tpl,
		})
	}

	return &Session{
		loadedTemplates: templates,
	}, nil
}

// Data is the template data.
type Data struct {
	FileName string
	Content  []byte
	Context  interface{}
}

// NewSessionFromData creates a new template session from a set of data.
func NewSessionFromData(options *LoadOptions, files []*Data) (*Session, error) {
	var templates []*Info
	for _, file := range files {
		var name = filenameWithoutExtension(file.FileName)

		idx := slices.IndexFunc(options.TemplatesToUse, func(t File) bool {
			tplName := t.Name
			if tplName == "" {
				tplName = t.Output
			}

			return tplName == name
		})
		if idx == -1 {
			// The template is not being used at the moment.
			continue
		}

		tpl, err := loadTemplate(name, file.Content, options)
		if err != nil {
			return nil, err
		}

		templates = append(templates, &Info{
			name:     options.TemplatesToUse[idx],
			template: tpl,
			context:  file.Context,
		})
	}

	return &Session{
		loadedTemplates: templates,
	}, nil
}

func filenameWithoutExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func loadTemplate(name string, data []byte, options *LoadOptions) (*template.Template, error) {
	var helperAPI = defaultAPI

	helperAPI["templateName"] = func() string {
		return name
	}
	for call, function := range options.API {
		helperAPI[call] = function
	}

	tpl, err := parse(name, data, helperAPI)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func parse(key string, data []byte, helperAPI template.FuncMap) (*template.Template, error) {
	t, err := template.New(key).Funcs(helperAPI).Parse(string(data))
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	return t, nil
}

// ExecuteTemplates executes the templates in the session.
func (s *Session) ExecuteTemplates(ctx interface{}) ([]*GeneratedTemplate, error) {
	var gen []*GeneratedTemplate

	for _, t := range s.loadedTemplates {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		tplCtx := ctx
		if ctx == nil {
			tplCtx = t.context
		}

		if err := t.template.Execute(w, tplCtx); err != nil {
			return nil, err
		}

		_ = w.Flush()
		gen = append(gen, newGeneratedTemplate(&buf, t.name))
	}

	return gen, nil
}

// GeneratedTemplate is the generated template.
type GeneratedTemplate struct {
	data *bytes.Buffer
	name string
}

func newGeneratedTemplate(data *bytes.Buffer, name File) *GeneratedTemplate {
	filename := name.Name
	if name.Output != "" {
		filename = name.Output
	}
	if name.Extension != "" {
		filename += fmt.Sprintf(".%v", name.Extension)
	}

	return &GeneratedTemplate{
		data: data,
		name: filename,
	}
}

// Filename returns the generated template filename.
func (g *GeneratedTemplate) Filename() string {
	return g.name
}

// Content returns the generated template content.
func (g *GeneratedTemplate) Content() []byte {
	return g.data.Bytes()
}
