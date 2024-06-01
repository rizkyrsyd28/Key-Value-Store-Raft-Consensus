// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.2
// source: protos/execute.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ExecuteServiceClient is the client API for ExecuteService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExecuteServiceClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error)
	Set(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error)
	StrLn(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error)
	Del(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error)
	Append(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error)
}

type executeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExecuteServiceClient(cc grpc.ClientConnInterface) ExecuteServiceClient {
	return &executeServiceClient{cc}
}

func (c *executeServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.ExecuteService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executeServiceClient) Get(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.ExecuteService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executeServiceClient) Set(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.ExecuteService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executeServiceClient) StrLn(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.ExecuteService/StrLn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executeServiceClient) Del(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.ExecuteService/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executeServiceClient) Append(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.ExecuteService/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExecuteServiceServer is the server API for ExecuteService service.
// All implementations must embed UnimplementedExecuteServiceServer
// for forward compatibility
type ExecuteServiceServer interface {
	Ping(context.Context, *Empty) (*Response, error)
	Get(context.Context, *KeyRequest) (*Response, error)
	Set(context.Context, *KeyValueRequest) (*Response, error)
	StrLn(context.Context, *KeyRequest) (*Response, error)
	Del(context.Context, *KeyRequest) (*Response, error)
	Append(context.Context, *KeyValueRequest) (*Response, error)
	mustEmbedUnimplementedExecuteServiceServer()
}

// UnimplementedExecuteServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExecuteServiceServer struct {
}

func (UnimplementedExecuteServiceServer) Ping(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedExecuteServiceServer) Get(context.Context, *KeyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedExecuteServiceServer) Set(context.Context, *KeyValueRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedExecuteServiceServer) StrLn(context.Context, *KeyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StrLn not implemented")
}
func (UnimplementedExecuteServiceServer) Del(context.Context, *KeyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedExecuteServiceServer) Append(context.Context, *KeyValueRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedExecuteServiceServer) mustEmbedUnimplementedExecuteServiceServer() {}

// UnsafeExecuteServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExecuteServiceServer will
// result in compilation errors.
type UnsafeExecuteServiceServer interface {
	mustEmbedUnimplementedExecuteServiceServer()
}

func RegisterExecuteServiceServer(s grpc.ServiceRegistrar, srv ExecuteServiceServer) {
	s.RegisterService(&ExecuteService_ServiceDesc, srv)
}

func _ExecuteService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.ExecuteService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServiceServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExecuteService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.ExecuteService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServiceServer).Get(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExecuteService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.ExecuteService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServiceServer).Set(ctx, req.(*KeyValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExecuteService_StrLn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServiceServer).StrLn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.ExecuteService/StrLn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServiceServer).StrLn(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExecuteService_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServiceServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.ExecuteService/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServiceServer).Del(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExecuteService_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecuteServiceServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.ExecuteService/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecuteServiceServer).Append(ctx, req.(*KeyValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExecuteService_ServiceDesc is the grpc.ServiceDesc for ExecuteService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExecuteService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dark_syster.ExecuteService",
	HandlerType: (*ExecuteServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _ExecuteService_Ping_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ExecuteService_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _ExecuteService_Set_Handler,
		},
		{
			MethodName: "StrLn",
			Handler:    _ExecuteService_StrLn_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _ExecuteService_Del_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _ExecuteService_Append_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/execute.proto",
}
