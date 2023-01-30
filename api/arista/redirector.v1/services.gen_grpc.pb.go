// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package redirector

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

// AssignmentServiceClient is the client API for AssignmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssignmentServiceClient interface {
	GetOne(ctx context.Context, in *AssignmentRequest, opts ...grpc.CallOption) (*AssignmentResponse, error)
	GetAll(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_GetAllClient, error)
	Subscribe(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_SubscribeClient, error)
}

type assignmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAssignmentServiceClient(cc grpc.ClientConnInterface) AssignmentServiceClient {
	return &assignmentServiceClient{cc}
}

func (c *assignmentServiceClient) GetOne(ctx context.Context, in *AssignmentRequest, opts ...grpc.CallOption) (*AssignmentResponse, error) {
	out := new(AssignmentResponse)
	err := c.cc.Invoke(ctx, "/arista.redirector.v1.AssignmentService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assignmentServiceClient) GetAll(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &AssignmentService_ServiceDesc.Streams[0], "/arista.redirector.v1.AssignmentService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &assignmentServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AssignmentService_GetAllClient interface {
	Recv() (*AssignmentStreamResponse, error)
	grpc.ClientStream
}

type assignmentServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *assignmentServiceGetAllClient) Recv() (*AssignmentStreamResponse, error) {
	m := new(AssignmentStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *assignmentServiceClient) Subscribe(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &AssignmentService_ServiceDesc.Streams[1], "/arista.redirector.v1.AssignmentService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &assignmentServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AssignmentService_SubscribeClient interface {
	Recv() (*AssignmentStreamResponse, error)
	grpc.ClientStream
}

type assignmentServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *assignmentServiceSubscribeClient) Recv() (*AssignmentStreamResponse, error) {
	m := new(AssignmentStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AssignmentServiceServer is the server API for AssignmentService service.
// All implementations must embed UnimplementedAssignmentServiceServer
// for forward compatibility
type AssignmentServiceServer interface {
	GetOne(context.Context, *AssignmentRequest) (*AssignmentResponse, error)
	GetAll(*AssignmentStreamRequest, AssignmentService_GetAllServer) error
	Subscribe(*AssignmentStreamRequest, AssignmentService_SubscribeServer) error
	mustEmbedUnimplementedAssignmentServiceServer()
}

// UnimplementedAssignmentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAssignmentServiceServer struct {
}

func (UnimplementedAssignmentServiceServer) GetOne(context.Context, *AssignmentRequest) (*AssignmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedAssignmentServiceServer) GetAll(*AssignmentStreamRequest, AssignmentService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedAssignmentServiceServer) Subscribe(*AssignmentStreamRequest, AssignmentService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedAssignmentServiceServer) mustEmbedUnimplementedAssignmentServiceServer() {}

// UnsafeAssignmentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AssignmentServiceServer will
// result in compilation errors.
type UnsafeAssignmentServiceServer interface {
	mustEmbedUnimplementedAssignmentServiceServer()
}

func RegisterAssignmentServiceServer(s grpc.ServiceRegistrar, srv AssignmentServiceServer) {
	s.RegisterService(&AssignmentService_ServiceDesc, srv)
}

func _AssignmentService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssignmentServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.redirector.v1.AssignmentService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssignmentServiceServer).GetOne(ctx, req.(*AssignmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssignmentService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AssignmentStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AssignmentServiceServer).GetAll(m, &assignmentServiceGetAllServer{stream})
}

type AssignmentService_GetAllServer interface {
	Send(*AssignmentStreamResponse) error
	grpc.ServerStream
}

type assignmentServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *assignmentServiceGetAllServer) Send(m *AssignmentStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _AssignmentService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AssignmentStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AssignmentServiceServer).Subscribe(m, &assignmentServiceSubscribeServer{stream})
}

type AssignmentService_SubscribeServer interface {
	Send(*AssignmentStreamResponse) error
	grpc.ServerStream
}

type assignmentServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *assignmentServiceSubscribeServer) Send(m *AssignmentStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// AssignmentService_ServiceDesc is the grpc.ServiceDesc for AssignmentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AssignmentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.redirector.v1.AssignmentService",
	HandlerType: (*AssignmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _AssignmentService_GetOne_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _AssignmentService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _AssignmentService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/redirector.v1/services.gen.proto",
}
