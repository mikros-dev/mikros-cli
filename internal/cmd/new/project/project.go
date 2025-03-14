package project

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/iancoleman/strcase"
	"github.com/mikros-dev/mikros-cli/internal/golang"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"os"
	"path/filepath"
	"strings"
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
	if err := createProjectTemplates(); err != nil {
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

func createProjectTemplates() error {
	//	- create root templates
	//  - create scripts templates
	//	- create proto templates
	return nil
}
