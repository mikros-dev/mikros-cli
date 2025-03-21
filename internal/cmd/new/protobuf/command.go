package protobuf

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/mikros-dev/mikros-cli/internal/assets/templates/protobuf"
	"github.com/mikros-dev/mikros-cli/internal/git"
	"github.com/mikros-dev/mikros-cli/internal/path"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/template"
)

func New(cfg *settings.Settings) error {
	name, kind, err := chooseService(cfg)
	if err != nil {
		return err
	}

	answers := &Answers{
		ServiceName: name,
		Kind:        kind,
	}

	switch kind {
	case "grpc":
		entity, defaultRPCs, customRPCs, err := runGrpcForm(cfg)
		if err != nil {
			return err
		}

		answers.Grpc = &GrpcAnswers{
			EntityName:     entity,
			UseDefaultRPCs: defaultRPCs,
			CustomRPCs:     customRPCs,
		}

	case "http":
		isAuthenticated, rpcs, err := runHttpForm(cfg)
		if err != nil {
			return err
		}

		answers.Http = &HttpAnswers{
			IsAuthenticated: isAuthenticated,
			RPCs:            rpcs,
		}
	}

	if err := generateTemplates(cfg, answers); err != nil {
		return err
	}

	return nil
}

func generateTemplates(cfg *settings.Settings, answers *Answers) error {
	templateBasePath, err := getTemplatesBasePath(answers.ServiceName)
	if err != nil {
		return err
	}

	if _, err := path.CreatePath(templateBasePath); err != nil {
		return err
	}

	// Switch to the destination path to create template sources
	cwd, err := path.ChangeDir(templateBasePath)
	if err != nil {
		return err
	}

	defer func() {
		if e := os.Chdir(cwd); e != nil {
			err = e
		}
	}()

	if err := generateProtobufFiles(cfg, templateBasePath, answers); err != nil {
		return err
	}

	return nil
}

func getTemplatesBasePath(serviceName string) (string, error) {
	repo, err := git.LoadFromCwd()
	if err != nil {
		return "", err
	}

	// If we're inside proto repo, create service path inside proto/project
	// and generate there.
	if repo.IsValidRepository() {
		files, err := os.ReadDir(filepath.Join(repo.RootPath, "proto"))
		if err != nil {
			// If it's not "our" protobuf repository, save the templates in
			// the current directory itself.
			return os.Getwd()
		}

		projectPath, err := findProtoMainProjectPath(repo.RootPath, serviceName, files)
		if err != nil {
			return os.Getwd()
		}

		return projectPath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Is there any proto/ folder from where we are now? If so, use the same
	// approach.
	if path.FindPath(filepath.Join(cwd, "proto")) {
		files, err := os.ReadDir(filepath.Join(cwd, "proto"))
		if err != nil {
			return cwd, nil
		}

		projectPath, err := findProtoMainProjectPath(cwd, serviceName, files)
		if err != nil {
			return cwd, nil
		}

		return projectPath, nil
	}

	return cwd, nil
}

func findProtoMainProjectPath(basePath, serviceName string, files []os.DirEntry) (string, error) {
	var projectPath string
	for _, file := range files {
		if file.IsDir() {
			projectPath = filepath.Join(basePath, "proto", file.Name())
			break
		}
	}
	if projectPath == "" {
		return "", errors.New("could not find protobuf main project folder")
	}

	return filepath.Join(projectPath, strings.ToLower(strcase.ToSnake(serviceName))), nil
}

func generateProtobufFiles(cfg *settings.Settings, basePath string, answers *Answers) error {
	var (
		filename = strings.ToLower(strcase.ToSnake(answers.ServiceName))
		tplFiles = []template.File{
			{
				Name:      "protobuf_api",
				Output:    filename + "_api",
				Extension: "proto",
			},
		}
	)

	if answers.Grpc != nil {
		tplFiles = append(tplFiles, template.File{
			Name:      "protobuf",
			Output:    filename,
			Extension: "proto",
		})
	}

	session, err := template.NewSessionFromFiles(&template.LoadOptions{
		TemplateNames: tplFiles,
	}, protobuf.Files)
	if err != nil {
		return err
	}

	ctx := generateTemplateContext(cfg, answers)
	if err := runTemplates(basePath, session, ctx); err != nil {
		return err
	}

	return nil
}

func runTemplates(basePath string, session *template.Session, context interface{}) error {
	generated, err := session.ExecuteTemplates(context)
	if err != nil {
		return err
	}

	for _, gen := range generated {
		file, err := os.Create(filepath.Join(basePath, gen.Filename()))
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
