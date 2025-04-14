package protobuf_module

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

type Context struct {
	httpService      bool
	IsAuthenticated  bool
	ServiceName      string
	Version          string
	EntityName       string
	CustomAuthName   string
	RPCMethods       []*RPC
	CustomRPCs       []*RPC
	MainPackageName  string
	RepositoryName   string
	VCSProjectPrefix string
}

func generateTemplateContext(cfg *settings.Settings, answers *Answers, profileName string) *Context {
	var (
		isAuthenticated bool
		entityName      string
		rpcs            []*RPC
		customRPCs      []*RPC
		profile         = projectProfile(cfg, profileName)
	)

	if answers.Grpc != nil {
		entityName = answers.Grpc.EntityName
		customRPCs = generateRPCs(answers.Grpc.CustomRPCs)

		if answers.Grpc.UseDefaultRPCs {
			rpcs = generateCRUDRPCs(entityName)
		}
	}
	if answers.Http != nil {
		rpcs = answers.Http.RPCs
		isAuthenticated = answers.Http.IsAuthenticated
	}

	return &Context{
		httpService:      answers.Kind == "http",
		IsAuthenticated:  isAuthenticated,
		ServiceName:      answers.ServiceName,
		Version:          "v0.1.0",
		EntityName:       entityName,
		CustomAuthName:   profile.Project.Templates.Protobuf.CustomAuthName,
		RPCMethods:       rpcs,
		CustomRPCs:       customRPCs,
		MainPackageName:  profile.Project.ProtobufMonorepo.ProjectName,
		RepositoryName:   profile.Project.ProtobufMonorepo.RepositoryName,
		VCSProjectPrefix: profile.Project.ProtobufMonorepo.VcsPath,
	}
}

func projectProfile(cfg *settings.Settings, profileName string) *settings.Profile {
	profile := &cfg.App
	if profileName == "default" {
		return profile
	}

	d, ok := cfg.Profile[profileName]
	if !ok {
		return profile
	}

	return &d
}

func (c *Context) IsHTTPService() bool {
	return c.httpService
}

func (c *Context) Extension() string {
	return "proto"
}

type RPC struct {
	IsAuthenticated bool
	Name            string
	HTTPMethod      string
	HTTPEndpoint    string
	AuthArgMode     string
	RequestName     string
	ResponseName    string
	RequestBody     string
	ResponseBody    string
}

func generateCRUDRPCs(entityName string) []*RPC {
	var (
		messageName = strcase.ToCamel(entityName)
		fieldName   = strcase.ToSnake(entityName)
	)

	return []*RPC{
		{
			Name:         fmt.Sprintf("Get%sByID", messageName),
			RequestBody:  "string id = 1;",
			ResponseBody: fmt.Sprintf("%sWire %s = 1;", messageName, fieldName),
		},
		{
			Name:         fmt.Sprintf("Create%s", messageName),
			ResponseBody: fmt.Sprintf("%sWire %s = 1;", messageName, fieldName),
		},
		{
			Name:         fmt.Sprintf("Update%sByID", messageName),
			RequestBody:  "string id = 1;",
			ResponseBody: fmt.Sprintf("%sWire %s = 1;", messageName, fieldName),
		},
		{
			Name:         fmt.Sprintf("Delete%sByID", messageName),
			RequestBody:  "string id = 1;",
			ResponseBody: fmt.Sprintf("%sWire %s = 1;", messageName, fieldName),
		},
	}
}

func generateRPCs(names []string) []*RPC {
	var (
		rpcs []*RPC
	)

	for _, name := range names {
		messageName := strcase.ToCamel(name)
		rpcs = append(rpcs, &RPC{
			Name:         messageName,
			RequestName:  messageName + "Request",
			ResponseName: messageName + "Response",
		})
	}

	return rpcs
}

func (m *RPC) HasBody() bool {
	return m.HTTPMethod == "post" || m.HTTPMethod == "put"
}
