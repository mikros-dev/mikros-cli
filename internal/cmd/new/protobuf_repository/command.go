package protobuf_repository

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	proto_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/protobuf_repository/proto"
	root_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/protobuf_repository/root"
	scripts_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/protobuf_repository/scripts"
	"github.com/mikros-dev/mikros-cli/internal/git"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
)

// NewOptions represents the options for the New command.
type NewOptions struct {
	NoVCS   bool
	Path    string
	Profile string
}

// New creates a new protobuf repository based on the provided settings.
func New(cfg *settings.Settings, options *NewOptions) error {
	answers, err := runSurvey(cfg, options.Profile)
	if err != nil {
		return err
	}

	return generateProject(options, answers)
}

func generateProject(options *NewOptions, answers *surveyAnswers) error {
	repositoryPath, err := createProjectDirectory(options, answers.RepositoryName)
	if err != nil {
		return err
	}

	// Switch to the destination path so we can work inside
	cwd, err := path.ChangeDir(repositoryPath)
	if err != nil {
		return err
	}

	defer func() {
		if e := os.Chdir(cwd); e != nil {
			err = e
		}
	}()

	// Notice that, starting from here, we're inside the project directory.
	if err := createProjectTemplates(answers, repositoryPath); err != nil {
		return err
	}

	// Initialize go module for the new repository
	if err := golang.ModInit(projectModuleName(answers)); err != nil {
		return err
	}

	if !options.NoVCS {
		if _, err := git.Init(); err != nil {
			return err
		}
	}

	return nil
}

func createProjectDirectory(options *NewOptions, repositoryName string) (string, error) {
	p, err := projectBasePath(options, repositoryName)
	if err != nil {
		return "", err
	}

	if _, err := path.CreatePath(p); err != nil {
		return "", err
	}

	return p, nil
}

func projectBasePath(options *NewOptions, repositoryName string) (string, error) {
	var name = strings.ToLower(strcase.ToKebab(repositoryName))

	if options.Path == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		return filepath.Join(cwd, name), nil
	}

	return filepath.Join(options.Path, name), nil
}

func projectModuleName(answers *surveyAnswers) string {
	return fmt.Sprintf("%s/%s", answers.VcsPath, strings.ToLower(strcase.ToKebab(answers.RepositoryName)))
}

func createProjectTemplates(answer *surveyAnswers, repositoryPath string) error {
	tplCtx := &TemplateContext{
		MainPackageName:  answer.ProjectName,
		RepositoryName:   answer.RepositoryName,
		VCSProjectPrefix: answer.VcsPath,
	}

	if err := createProjectRootTemplates(tplCtx); err != nil {
		return err
	}

	if err := createProjectScriptsTemplates(repositoryPath, tplCtx); err != nil {
		return err
	}

	return createProjectProtoTemplates(repositoryPath, tplCtx)
}

func createProjectRootTemplates(tplCtx *TemplateContext) error {
	templates := []template.File{
		{
			Name: "buf.gen.yaml",
		},
		{
			Name: "buf.yaml",
		},
		{
			Name: "Makefile",
		},
		{
			Name: "README.md",
		},
	}

	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplateNames: templates,
	}, root_tpl.Files)
	if err != nil {
		return err
	}

	return runTemplates(session, tplCtx)
}

func createProjectScriptsTemplates(repositoryPath string, tplCtx *TemplateContext) error {
	templates := []template.File{
		{
			Name: "generate.sh",
		},
		{
			Name: "go.sh",
		},
		{
			Name: "setup.sh",
		},
	}

	// Create .scripts folder and dive into it
	scriptsPath := filepath.Join(repositoryPath, ".scripts")
	if _, err := path.CreatePath(scriptsPath); err != nil {
		return err
	}
	cwd, err := path.ChangeDir(scriptsPath)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.Chdir(cwd); e != nil {
			err = e
		}
	}()

	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplateNames: templates,
	}, scripts_tpl.Files)
	if err != nil {
		return err
	}

	if err := runTemplates(session, tplCtx); err != nil {
		return err
	}

	for _, file := range templates {
		if err := path.SetExecutablePath(filepath.Join(scriptsPath, file.Name)); err != nil {
			return err
		}
	}

	return nil
}

func createProjectProtoTemplates(repositoryPath string, tplCtx *TemplateContext) error {
	templates := []template.File{
		{
			Name: "example.proto",
		},
	}

	// Create proto folder and dive into it
	scriptsPath := filepath.Join(repositoryPath, "proto", tplCtx.MainPackageName, "example")
	if _, err := path.CreatePath(scriptsPath); err != nil {
		return err
	}
	cwd, err := path.ChangeDir(scriptsPath)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.Chdir(cwd); e != nil {
			err = e
		}
	}()

	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplateNames: templates,
	}, proto_tpl.Files)
	if err != nil {
		return err
	}

	return runTemplates(session, tplCtx)
}

func runTemplates(session *template.Session, context interface{}) error {
	generated, err := session.ExecuteTemplates(context)
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
