package templates

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

	"github.com/somatech1/mikros-cli/pkg/templates"
)

var defaultApi = template.FuncMap{
	"toCamel":      strcase.ToCamel,
	"toSnake":      strcase.ToSnake,
	"toUpperSnake": strcase.ToScreamingSnake,
	"basename":     path.Base,
}

type Templates struct {
	templates []*TemplateInfo
}

type TemplateInfo struct {
	tpl  *template.Template
	name templates.TemplateFile
}

type GeneratedTemplate struct {
	data *bytes.Buffer
	name string
}

func (g *GeneratedTemplate) Filename() string {
	return g.name
}

func (g *GeneratedTemplate) Content() []byte {
	return g.data.Bytes()
}

type LoadOptions struct {
	TemplateNames []templates.TemplateFile
	Api           map[string]interface{}
	Files         embed.FS
}

func Load(options *LoadOptions) (*Templates, error) {
	files, err := options.Files.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var tpls []*TemplateInfo
	for _, file := range files {
		data, err := options.Files.ReadFile(file.Name())
		if err != nil {
			return nil, err
		}

		var (
			name      = filenameWithoutExtension(file.Name())
			helperApi = defaultApi
		)

		helperApi["templateName"] = func() string {
			return name
		}
		for call, function := range options.Api {
			helperApi[call] = function
		}

		idx := slices.IndexFunc(options.TemplateNames, func(t templates.TemplateFile) bool {
			return t.Name == name
		})
		if idx == -1 {
			// The template is not being used at the moment.
			continue
		}

		tpl, err := parse(name, data, helperApi)
		if err != nil {
			return nil, err
		}

		tpls = append(tpls, &TemplateInfo{
			name: options.TemplateNames[idx],
			tpl:  tpl,
		})
	}

	return &Templates{
		templates: tpls,
	}, nil
}

func filenameWithoutExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func parse(key string, data []byte, helperApi template.FuncMap) (*template.Template, error) {
	t, err := template.New(key).Funcs(helperApi).Parse(string(data))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Templates) Execute(ctx interface{}) ([]*GeneratedTemplate, error) {
	var gen []*GeneratedTemplate

	for _, tpl := range t.templates {
		var buf bytes.Buffer
		w := bufio.NewWriter(&buf)

		if err := tpl.tpl.Execute(w, ctx); err != nil {
			return nil, err
		}

		_ = w.Flush()
		gen = append(gen, t.newGenerated(&buf, tpl.name))
	}

	return gen, nil
}

func (t *Templates) newGenerated(data *bytes.Buffer, name templates.TemplateFile) *GeneratedTemplate {
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
