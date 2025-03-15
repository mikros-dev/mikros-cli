package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/iancoleman/strcase"

	proto_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/project/proto"
	root_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/project/root"
	scripts_tpl "github.com/mikros-dev/mikros-cli/internal/assets/templates/project/scripts"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/template"
)

type surveyAnswers struct {
	RepositoryName string `survey:"repository_name"`
	ProjectName    string `survey:"project_name"`
}

func New() error {
	answers := &surveyAnswers{}
	if err := survey.Ask(baseQuestions(), answers); err != nil {
		return err
	}

	if err := generateProject(answers); err != nil {
		return err
	}

	return nil
}

func baseQuestions() []*survey.Question {
	return []*survey.Question{
		// Repository name
		{
			Name: "repository_name",
			Prompt: &survey.Input{
				Message: "Repository name. Enter the name of the repository to create:",
				Default: "protobuf-workspace",
			},
			Validate: survey.Required,
		},
		// Project name
		{
			Name: "project_name",
			Prompt: &survey.Input{
				Message: "Project name. Enter your protobuf project name:",
				Default: "services",
			},
			Validate: survey.Required,
		},
	}
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
	if err := golang.ModInit(projectModuleName(answers.RepositoryName)); err != nil {
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

func projectModuleName(repositoryName string) string {
	return fmt.Sprintf("github.com/your-org/%s", strings.ToLower(strcase.ToKebab(repositoryName)))
}

func createProjectTemplates(answer *surveyAnswers, repositoryPath string) error {
	tplCtx := &TemplateContext{
		MainPackageName: answer.ProjectName,
		RepositoryName:  answer.RepositoryName,
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
