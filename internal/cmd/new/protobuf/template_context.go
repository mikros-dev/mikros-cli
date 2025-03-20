package protobuf

type Context struct {
	httpService  bool
	ServiceName  string
	Version      string
	EntityPrefix string
	RPCMethods   []*RPC
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
}

func (m *RPC) HasBody() bool {
	return m.HTTPMethod == "post" || m.HTTPMethod == "put"
}
