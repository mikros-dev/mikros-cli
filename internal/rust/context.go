package rust

import (
	"github.com/somatech1/mikros/components/definition"
)

type Context struct {
	HasLifecycle bool
	ServiceName  string
	ModuleName   string
	Methods      []*Method

	serviceType string
}

type Method struct {
	Name         string
	RequestName  string
	ResponseName string
}

func (t *Context) IsScriptService() bool {
	return t.serviceType == definition.ServiceType_Script.String()
}

func (t *Context) IsNativeService() bool {
	return t.serviceType == definition.ServiceType_Native.String()
}

func (t *Context) IsGrpcService() bool {
	return t.serviceType == definition.ServiceType_gRPC.String()
}

func (t *Context) IsHttpService() bool {
	return t.serviceType == definition.ServiceType_HTTP.String()
}

// ServiceTypeBuilderCall returns, in a string format, the function that
// initializes the service according its type.
func (t *Context) ServiceTypeBuilderCall() string {
	if t.IsScriptService() {
		return ".script(s)"
	}
	if t.IsNativeService() {
		return ".native(s)"
	}
	if t.IsGrpcService() {
		return ".grpc(s)"
	}
	if t.IsHttpService() {
		return ".http(s.routes())"
	}

	return ""
}
