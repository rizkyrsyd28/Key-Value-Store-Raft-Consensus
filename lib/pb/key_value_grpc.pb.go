// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.2
// source: protos/key_value.proto

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

// KeyValueServiceClient is the client API for KeyValueService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type KeyValueServiceClient interface {
	Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error)
	Set(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error)
	StrLn(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error)
	Del(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error)
	Append(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error)
	RequestLog(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error)
}

type keyValueServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewKeyValueServiceClient(cc grpc.ClientConnInterface) KeyValueServiceClient {
	return &keyValueServiceClient{cc}
}

func (c *keyValueServiceClient) Ping(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueServiceClient) Get(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueServiceClient) Set(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueServiceClient) StrLn(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/StrLn", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueServiceClient) Del(ctx context.Context, in *KeyRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/Del", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueServiceClient) Append(ctx context.Context, in *KeyValueRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/Append", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *keyValueServiceClient) RequestLog(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/dark_syster.KeyValueService/RequestLog", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// KeyValueServiceServer is the server API for KeyValueService service.
// All implementations must embed UnimplementedKeyValueServiceServer
// for forward compatibility
type KeyValueServiceServer interface {
	Ping(context.Context, *Empty) (*Response, error)
	Get(context.Context, *KeyRequest) (*Response, error)
	Set(context.Context, *KeyValueRequest) (*Response, error)
	StrLn(context.Context, *KeyRequest) (*Response, error)
	Del(context.Context, *KeyRequest) (*Response, error)
	Append(context.Context, *KeyValueRequest) (*Response, error)
	RequestLog(context.Context, *Empty) (*Response, error)
	mustEmbedUnimplementedKeyValueServiceServer()
}

// UnimplementedKeyValueServiceServer must be embedded to have forward compatible implementations.
type UnimplementedKeyValueServiceServer struct {
}

func (UnimplementedKeyValueServiceServer) Ping(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedKeyValueServiceServer) Get(context.Context, *KeyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedKeyValueServiceServer) Set(context.Context, *KeyValueRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedKeyValueServiceServer) StrLn(context.Context, *KeyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StrLn not implemented")
}
func (UnimplementedKeyValueServiceServer) Del(context.Context, *KeyRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Del not implemented")
}
func (UnimplementedKeyValueServiceServer) Append(context.Context, *KeyValueRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Append not implemented")
}
func (UnimplementedKeyValueServiceServer) RequestLog(context.Context, *Empty) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestLog not implemented")
}
func (UnimplementedKeyValueServiceServer) mustEmbedUnimplementedKeyValueServiceServer() {}

// UnsafeKeyValueServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to KeyValueServiceServer will
// result in compilation errors.
type UnsafeKeyValueServiceServer interface {
	mustEmbedUnimplementedKeyValueServiceServer()
}

func RegisterKeyValueServiceServer(s grpc.ServiceRegistrar, srv KeyValueServiceServer) {
	s.RegisterService(&KeyValueService_ServiceDesc, srv)
}

func _KeyValueService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).Ping(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).Get(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).Set(ctx, req.(*KeyValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueService_StrLn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).StrLn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/StrLn",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).StrLn(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueService_Del_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).Del(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/Del",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).Del(ctx, req.(*KeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueService_Append_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).Append(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/Append",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).Append(ctx, req.(*KeyValueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _KeyValueService_RequestLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(KeyValueServiceServer).RequestLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/dark_syster.KeyValueService/RequestLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KeyValueServiceServer).RequestLog(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// KeyValueService_ServiceDesc is the grpc.ServiceDesc for KeyValueService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var KeyValueService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dark_syster.KeyValueService",
	HandlerType: (*KeyValueServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _KeyValueService_Ping_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _KeyValueService_Get_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _KeyValueService_Set_Handler,
		},
		{
			MethodName: "StrLn",
			Handler:    _KeyValueService_StrLn_Handler,
		},
		{
			MethodName: "Del",
			Handler:    _KeyValueService_Del_Handler,
		},
		{
			MethodName: "Append",
			Handler:    _KeyValueService_Append_Handler,
		},
		{
			MethodName: "RequestLog",
			Handler:    _KeyValueService_RequestLog_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/key_value.proto",
}
