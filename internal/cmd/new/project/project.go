package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/creasty/defaults"
	"github.com/iancoleman/strcase"

	proto_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/project/proto"
	root_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/project/root"
	scripts_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/project/scripts"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

type surveyAnswers struct {
	RepositoryName string `survey:"repository_name" default:"protobuf-workspace"`
	ProjectName    string `survey:"project_name" default:"services"`
	VcsPath        string `survey:"vcs_path"`
}

func newSurveyAnswers(cfg *settings.Settings) *surveyAnswers {
	a := &surveyAnswers{}
	if err := defaults.Set(a); err != nil {
		// Without default values
		return a
	}

	a.VcsPath = cfg.Project.Template.VcsPath
	return a
}

type NewOptions struct {
	Path          string
	ProtoFilename string
}

func New(cfg *settings.Settings, options *NewOptions) error {
	answers, err := runSurvey(cfg)
	if err != nil {
		return err
	}

	if err := generateProject(answers); err != nil {
		return err
	}

	return nil
}

func runSurvey(cfg *settings.Settings) (*surveyAnswers, error) {
	answers := newSurveyAnswers(cfg)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Repository name. Enter the name of the repository to create:").
				Value(&answers.RepositoryName).
				Validate(ui.IsEmpty("repository name cannot be empty")),

			huh.NewInput().
				Title("Project name. Enter your protobuf project name:").
				Value(&answers.ProjectName).
				Validate(ui.IsEmpty("project name cannot be empty")),

			huh.NewInput().
				Title("VCS path prefix. Enter your VCS path prefix to use for the project:").
				Value(&answers.VcsPath).
				Validate(ui.IsEmpty("VCS path prefix cannot be empty")),
		),
	)

	if err := form.WithTheme(cfg.GetTheme()).Run(); err != nil {
		return nil, err
	}

	return answers, nil
}

func generateProject(answers *surveyAnswers) error {
	repositoryPath, err := createProjectDirectory(answers.RepositoryName)
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

	return nil
}

func createProjectDirectory(repositoryName string) (string, error) {
	p, err := projectBasePath(repositoryName)
	if err != nil {
		return "", err
	}

	if _, err := path.CreatePath(p); err != nil {
		return "", err
	}

	return p, nil
}

func projectBasePath(repositoryName string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, strings.ToLower(strcase.ToKebab(repositoryName))), nil
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

	if err := createProjectProtoTemplates(repositoryPath, tplCtx); err != nil {
		return err
	}

	return nil
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

	if err := runTemplates(session, tplCtx); err != nil {
		return err
	}

	return nil
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

	if err := runTemplates(session, tplCtx); err != nil {
		return err
	}

	return nil
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
