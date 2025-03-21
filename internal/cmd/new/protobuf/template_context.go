package protobuf

import (
	"fmt"

	"github.com/iancoleman/strcase"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

type Context struct {
	httpService      bool
	ServiceName      string
	Version          string
	EntityName       string
	RPCMethods       []*RPC
	CustomRPCs       []*RPC
	MainPackageName  string
	RepositoryName   string
	VCSProjectPrefix string
}

func generateTemplateContext(cfg *settings.Settings, answers *Answers) *Context {
	var (
		entityName string
		rpcs       []*RPC
		customRPCs []*RPC
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
	}

	return &Context{
		httpService:      answers.Kind == "http",
		ServiceName:      answers.ServiceName,
		Version:          "v0.1.0",
		EntityName:       entityName,
		RPCMethods:       rpcs,
		CustomRPCs:       customRPCs,
		MainPackageName:  cfg.Project.ProtobufMonorepo.ProjectName,
		RepositoryName:   cfg.Project.ProtobufMonorepo.RepositoryName,
		VCSProjectPrefix: cfg.Project.ProtobufMonorepo.VcsPath,
	}
}

func (c *Context) IsHTTPService() bool {
	return c.httpService
}

func (c *Context) Extension() string {
	return "proto"
}

type RPC struct {
	Name         string
	HTTPMethod   string
	HTTPEndpoint string
	AuthArgMode  string
	RequestName  string
	ResponseName string
	RequestBody  string
	ResponseBody string
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
