package service

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/mikros-dev/mikros-cli/internal/fs"
	"github.com/mikros-dev/mikros-cli/internal/git"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
)

// NewOptions represents the options for the New command.
type NewOptions struct {
	NoVCS bool
	Path  string
}

// New creates a new project based on the provided settings and options,
// running a survey and generating the project files.
func New(cfg *settings.Settings, options *NewOptions) error {
	answers, err := runSurvey(cfg)
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
	cwd, err := fs.ChangeDir(repositoryPath)
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

	if _, err := fs.CreatePath(p); err != nil {
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

func createProjectTemplates(answer *surveyAnswers, repositoryPath string) error {
	tplCtx := &TemplateContext{
		RepositoryName: answer.RepositoryName,
	}

	if err := createProjectRootTemplates(tplCtx); err != nil {
		return err
	}

	return createProjectScriptsTemplates(repositoryPath, tplCtx)
}

func createProjectRootTemplates(tplCtx *TemplateContext) error {
	templates := []template.File{
		{
			Name: "Makefile",
		},
		{
			Name: "README.md",
		},
		{
			Name: ".gitignore",
		},
	}

	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplatesToUse: templates,
		FilesBasePath:  "assets/root",
	}, rootFiles)
	if err != nil {
		return err
	}

	return runTemplates(session, tplCtx)
}

func createProjectScriptsTemplates(repositoryPath string, tplCtx *TemplateContext) error {
	templates := []template.File{
		{
			Name: "badges.sh",
		},
		{
			Name: "services.sh",
		},
		{
			Name: "tests.sh",
		},
		{
			Name: "utils.sh",
		},
		{
			Name: "check-service-toml.sh",
		},
	}

	// Create .scripts folder and dive into it
	scriptsPath := filepath.Join(repositoryPath, ".scripts")
	if _, err := fs.CreatePath(scriptsPath); err != nil {
		return err
	}
	cwd, err := fs.ChangeDir(scriptsPath)
	if err != nil {
		return err
	}
	defer func() {
		if e := os.Chdir(cwd); e != nil {
			err = e
		}
	}()

	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplatesToUse: templates,
		FilesBasePath:  "assets/scripts",
	}, scriptFiles)
	if err != nil {
		return err
	}

	if err := runTemplates(session, tplCtx); err != nil {
		return err
	}

	for _, file := range templates {
		if err := fs.SetExecutablePath(filepath.Join(scriptsPath, file.Name)); err != nil {
			return err
		}
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
