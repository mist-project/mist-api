// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: appuser.proto

package protos

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
	AppuserService_CreateAppuser_FullMethodName = "/v1.appuser.AppuserService/CreateAppuser"
)

// AppuserServiceClient is the client API for AppuserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppuserServiceClient interface {
	// ----- APPUSER ----
	CreateAppuser(ctx context.Context, in *CreateAppuserRequest, opts ...grpc.CallOption) (*CreateAppuserResponse, error)
}

type appuserServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAppuserServiceClient(cc grpc.ClientConnInterface) AppuserServiceClient {
	return &appuserServiceClient{cc}
}

func (c *appuserServiceClient) CreateAppuser(ctx context.Context, in *CreateAppuserRequest, opts ...grpc.CallOption) (*CreateAppuserResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateAppuserResponse)
	err := c.cc.Invoke(ctx, AppuserService_CreateAppuser_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppuserServiceServer is the server API for AppuserService service.
// All implementations must embed UnimplementedAppuserServiceServer
// for forward compatibility.
type AppuserServiceServer interface {
	// ----- APPUSER ----
	CreateAppuser(context.Context, *CreateAppuserRequest) (*CreateAppuserResponse, error)
	mustEmbedUnimplementedAppuserServiceServer()
}

// UnimplementedAppuserServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAppuserServiceServer struct{}

func (UnimplementedAppuserServiceServer) CreateAppuser(context.Context, *CreateAppuserRequest) (*CreateAppuserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAppuser not implemented")
}
func (UnimplementedAppuserServiceServer) mustEmbedUnimplementedAppuserServiceServer() {}
func (UnimplementedAppuserServiceServer) testEmbeddedByValue()                        {}

// UnsafeAppuserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppuserServiceServer will
// result in compilation errors.
type UnsafeAppuserServiceServer interface {
	mustEmbedUnimplementedAppuserServiceServer()
}

func RegisterAppuserServiceServer(s grpc.ServiceRegistrar, srv AppuserServiceServer) {
	// If the following call pancis, it indicates UnimplementedAppuserServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&AppuserService_ServiceDesc, srv)
}

func _AppuserService_CreateAppuser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAppuserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppuserServiceServer).CreateAppuser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AppuserService_CreateAppuser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppuserServiceServer).CreateAppuser(ctx, req.(*CreateAppuserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AppuserService_ServiceDesc is the grpc.ServiceDesc for AppuserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AppuserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "v1.appuser.AppuserService",
	HandlerType: (*AppuserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAppuser",
			Handler:    _AppuserService_CreateAppuser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "appuser.proto",
}
