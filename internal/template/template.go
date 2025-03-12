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

var defaultApi = template.FuncMap{
	"toCamel":      strcase.ToCamel,
	"toSnake":      strcase.ToSnake,
	"toUpperSnake": strcase.ToScreamingSnake,
	"basename":     path.Base,
	"toKebab":      strcase.ToKebab,
}

type Session struct {
	loadedTemplates []*Info
}

type Info struct {
	template *template.Template
	name     File
}

// File is representation of a template file to be processed when
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

type LoadOptions struct {
	TemplateNames []File
	Api           map[string]interface{}
}

func NewSessionFromFiles(options *LoadOptions, files embed.FS) (*Session, error) {
	dirFiles, err := files.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var templates []*Info
	for _, file := range dirFiles {
		data, err := files.ReadFile(file.Name())
		if err != nil {
			return nil, err
		}

		var (
			name = filenameWithoutExtension(file.Name())
		)

		idx := slices.IndexFunc(options.TemplateNames, func(t File) bool {
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
			name:     options.TemplateNames[idx],
			template: tpl,
		})
	}

	return &Session{
		loadedTemplates: templates,
	}, nil
}

type Data struct {
	FileName string
	Content  []byte
}

func NewSessionFromData(options *LoadOptions, files []*Data) (*Session, error) {
	var templates []*Info
	for _, file := range files {
		var (
			name = filenameWithoutExtension(file.FileName)
		)

		idx := slices.IndexFunc(options.TemplateNames, func(t File) bool {
			return t.Name == name
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
			name:     options.TemplateNames[idx],
			template: tpl,
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
	var (
		helperApi = defaultApi
	)

	helperApi["templateName"] = func() string {
		return name
	}
	for call, function := range options.Api {
		helperApi[call] = function
	}

	tpl, err := parse(name, data, helperApi)
	if err != nil {
		return nil, err
	}

	return tpl, nil
}

func parse(key string, data []byte, helperApi template.FuncMap) (*template.Template, error) {
	t, err := template.New(key).Funcs(helperApi).Parse(string(data))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Session) ExecuteTemplates(ctx interface{}) ([]*GeneratedTemplate, error) {
	var gen []*GeneratedTemplate

	for _, t := range s.loadedTemplates {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		if err := t.template.Execute(w, ctx); err != nil {
			return nil, err
		}

		_ = w.Flush()
		gen = append(gen, newGeneratedTemplate(&buf, t.name))
	}

	return gen, nil
}

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

func (g *GeneratedTemplate) Filename() string {
	return g.name
}

func (g *GeneratedTemplate) Content() []byte {
	return g.data.Bytes()
}
