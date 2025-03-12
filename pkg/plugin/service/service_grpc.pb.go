// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: plugin/service/service.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Plugin_GetName_FullMethodName         = "/plugins.service.Plugin/GetName"
	Plugin_GetKind_FullMethodName         = "/plugins.service.Plugin/GetKind"
	Plugin_GetSurvey_FullMethodName       = "/plugins.service.Plugin/GetSurvey"
	Plugin_ValidateAnswers_FullMethodName = "/plugins.service.Plugin/ValidateAnswers"
	Plugin_GetTemplate_FullMethodName     = "/plugins.service.Plugin/GetTemplate"
	Plugin_Stop_FullMethodName            = "/plugins.service.Plugin/Stop"
)

// PluginClient is the client API for Plugin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginClient interface {
	GetName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetNameResponse, error)
	GetKind(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetKindResponse, error)
	GetSurvey(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetSurveyResponse, error)
	ValidateAnswers(ctx context.Context, in *ValidateAnswersRequest, opts ...grpc.CallOption) (*ValidateAnswersResponse, error)
	GetTemplate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetTemplateResponse, error)
	Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type pluginClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginClient(cc grpc.ClientConnInterface) PluginClient {
	return &pluginClient{cc}
}

func (c *pluginClient) GetName(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetNameResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetNameResponse)
	err := c.cc.Invoke(ctx, Plugin_GetName_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetKind(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetKindResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetKindResponse)
	err := c.cc.Invoke(ctx, Plugin_GetKind_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetSurvey(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetSurveyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetSurveyResponse)
	err := c.cc.Invoke(ctx, Plugin_GetSurvey_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) ValidateAnswers(ctx context.Context, in *ValidateAnswersRequest, opts ...grpc.CallOption) (*ValidateAnswersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ValidateAnswersResponse)
	err := c.cc.Invoke(ctx, Plugin_ValidateAnswers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) GetTemplate(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetTemplateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTemplateResponse)
	err := c.cc.Invoke(ctx, Plugin_GetTemplate_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pluginClient) Stop(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, Plugin_Stop_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginServer is the server API for Plugin service.
// All implementations should embed UnimplementedPluginServer
// for forward compatibility.
type PluginServer interface {
	GetName(context.Context, *Empty) (*GetNameResponse, error)
	GetKind(context.Context, *Empty) (*GetKindResponse, error)
	GetSurvey(context.Context, *Empty) (*GetSurveyResponse, error)
	ValidateAnswers(context.Context, *ValidateAnswersRequest) (*ValidateAnswersResponse, error)
	GetTemplate(context.Context, *Empty) (*GetTemplateResponse, error)
	Stop(context.Context, *Empty) (*Empty, error)
}

// UnimplementedPluginServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPluginServer struct{}

func (UnimplementedPluginServer) GetName(context.Context, *Empty) (*GetNameResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetName not implemented")
}
func (UnimplementedPluginServer) GetKind(context.Context, *Empty) (*GetKindResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKind not implemented")
}
func (UnimplementedPluginServer) GetSurvey(context.Context, *Empty) (*GetSurveyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSurvey not implemented")
}
func (UnimplementedPluginServer) ValidateAnswers(context.Context, *ValidateAnswersRequest) (*ValidateAnswersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateAnswers not implemented")
}
func (UnimplementedPluginServer) GetTemplate(context.Context, *Empty) (*GetTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTemplate not implemented")
}
func (UnimplementedPluginServer) Stop(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stop not implemented")
}
func (UnimplementedPluginServer) testEmbeddedByValue() {}

// UnsafePluginServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginServer will
// result in compilation errors.
type UnsafePluginServer interface {
	mustEmbedUnimplementedPluginServer()
}

func RegisterPluginServer(s grpc.ServiceRegistrar, srv PluginServer) {
	// If the following call pancis, it indicates UnimplementedPluginServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Plugin_ServiceDesc, srv)
}

func _Plugin_GetName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetName(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetKind_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetKind(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetKind_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetKind(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetSurvey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetSurvey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetSurvey_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetSurvey(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_ValidateAnswers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ValidateAnswersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).ValidateAnswers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_ValidateAnswers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).ValidateAnswers(ctx, req.(*ValidateAnswersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_GetTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).GetTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_GetTemplate_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).GetTemplate(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Plugin_Stop_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginServer).Stop(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Plugin_Stop_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginServer).Stop(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Plugin_ServiceDesc is the grpc.ServiceDesc for Plugin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Plugin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "plugins.service.Plugin",
	HandlerType: (*PluginServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetName",
			Handler:    _Plugin_GetName_Handler,
		},
		{
			MethodName: "GetKind",
			Handler:    _Plugin_GetKind_Handler,
		},
		{
			MethodName: "GetSurvey",
			Handler:    _Plugin_GetSurvey_Handler,
		},
		{
			MethodName: "ValidateAnswers",
			Handler:    _Plugin_ValidateAnswers_Handler,
		},
		{
			MethodName: "GetTemplate",
			Handler:    _Plugin_GetTemplate_Handler,
		},
		{
			MethodName: "Stop",
			Handler:    _Plugin_Stop_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin/service/service.proto",
}
