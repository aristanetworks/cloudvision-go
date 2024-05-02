// Copyright (c) 2022 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

//
// Code generated by boomtown. DO NOT EDIT.
//

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: arista/redirector.v1/services.gen.proto

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

const (
	AssignmentService_GetOne_FullMethodName        = "/arista.redirector.v1.AssignmentService/GetOne"
	AssignmentService_GetSome_FullMethodName       = "/arista.redirector.v1.AssignmentService/GetSome"
	AssignmentService_GetAll_FullMethodName        = "/arista.redirector.v1.AssignmentService/GetAll"
	AssignmentService_Subscribe_FullMethodName     = "/arista.redirector.v1.AssignmentService/Subscribe"
	AssignmentService_GetMeta_FullMethodName       = "/arista.redirector.v1.AssignmentService/GetMeta"
	AssignmentService_SubscribeMeta_FullMethodName = "/arista.redirector.v1.AssignmentService/SubscribeMeta"
)

// AssignmentServiceClient is the client API for AssignmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssignmentServiceClient interface {
	GetOne(ctx context.Context, in *AssignmentRequest, opts ...grpc.CallOption) (*AssignmentResponse, error)
	GetSome(ctx context.Context, in *AssignmentSomeRequest, opts ...grpc.CallOption) (AssignmentService_GetSomeClient, error)
	GetAll(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_GetAllClient, error)
	Subscribe(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_SubscribeClient, error)
	GetMeta(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error)
	SubscribeMeta(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_SubscribeMetaClient, error)
}

type assignmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAssignmentServiceClient(cc grpc.ClientConnInterface) AssignmentServiceClient {
	return &assignmentServiceClient{cc}
}

func (c *assignmentServiceClient) GetOne(ctx context.Context, in *AssignmentRequest, opts ...grpc.CallOption) (*AssignmentResponse, error) {
	out := new(AssignmentResponse)
	err := c.cc.Invoke(ctx, AssignmentService_GetOne_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assignmentServiceClient) GetSome(ctx context.Context, in *AssignmentSomeRequest, opts ...grpc.CallOption) (AssignmentService_GetSomeClient, error) {
	stream, err := c.cc.NewStream(ctx, &AssignmentService_ServiceDesc.Streams[0], AssignmentService_GetSome_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &assignmentServiceGetSomeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AssignmentService_GetSomeClient interface {
	Recv() (*AssignmentSomeResponse, error)
	grpc.ClientStream
}

type assignmentServiceGetSomeClient struct {
	grpc.ClientStream
}

func (x *assignmentServiceGetSomeClient) Recv() (*AssignmentSomeResponse, error) {
	m := new(AssignmentSomeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *assignmentServiceClient) GetAll(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &AssignmentService_ServiceDesc.Streams[1], AssignmentService_GetAll_FullMethodName, opts...)
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
	stream, err := c.cc.NewStream(ctx, &AssignmentService_ServiceDesc.Streams[2], AssignmentService_Subscribe_FullMethodName, opts...)
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

func (c *assignmentServiceClient) GetMeta(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error) {
	out := new(MetaResponse)
	err := c.cc.Invoke(ctx, AssignmentService_GetMeta_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assignmentServiceClient) SubscribeMeta(ctx context.Context, in *AssignmentStreamRequest, opts ...grpc.CallOption) (AssignmentService_SubscribeMetaClient, error) {
	stream, err := c.cc.NewStream(ctx, &AssignmentService_ServiceDesc.Streams[3], AssignmentService_SubscribeMeta_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &assignmentServiceSubscribeMetaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type AssignmentService_SubscribeMetaClient interface {
	Recv() (*MetaResponse, error)
	grpc.ClientStream
}

type assignmentServiceSubscribeMetaClient struct {
	grpc.ClientStream
}

func (x *assignmentServiceSubscribeMetaClient) Recv() (*MetaResponse, error) {
	m := new(MetaResponse)
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
	GetSome(*AssignmentSomeRequest, AssignmentService_GetSomeServer) error
	GetAll(*AssignmentStreamRequest, AssignmentService_GetAllServer) error
	Subscribe(*AssignmentStreamRequest, AssignmentService_SubscribeServer) error
	GetMeta(context.Context, *AssignmentStreamRequest) (*MetaResponse, error)
	SubscribeMeta(*AssignmentStreamRequest, AssignmentService_SubscribeMetaServer) error
	mustEmbedUnimplementedAssignmentServiceServer()
}

// UnimplementedAssignmentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAssignmentServiceServer struct {
}

func (UnimplementedAssignmentServiceServer) GetOne(context.Context, *AssignmentRequest) (*AssignmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedAssignmentServiceServer) GetSome(*AssignmentSomeRequest, AssignmentService_GetSomeServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSome not implemented")
}
func (UnimplementedAssignmentServiceServer) GetAll(*AssignmentStreamRequest, AssignmentService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedAssignmentServiceServer) Subscribe(*AssignmentStreamRequest, AssignmentService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedAssignmentServiceServer) GetMeta(context.Context, *AssignmentStreamRequest) (*MetaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeta not implemented")
}
func (UnimplementedAssignmentServiceServer) SubscribeMeta(*AssignmentStreamRequest, AssignmentService_SubscribeMetaServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeMeta not implemented")
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
		FullMethod: AssignmentService_GetOne_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssignmentServiceServer).GetOne(ctx, req.(*AssignmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssignmentService_GetSome_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AssignmentSomeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AssignmentServiceServer).GetSome(m, &assignmentServiceGetSomeServer{stream})
}

type AssignmentService_GetSomeServer interface {
	Send(*AssignmentSomeResponse) error
	grpc.ServerStream
}

type assignmentServiceGetSomeServer struct {
	grpc.ServerStream
}

func (x *assignmentServiceGetSomeServer) Send(m *AssignmentSomeResponse) error {
	return x.ServerStream.SendMsg(m)
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

func _AssignmentService_GetMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AssignmentStreamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssignmentServiceServer).GetMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AssignmentService_GetMeta_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssignmentServiceServer).GetMeta(ctx, req.(*AssignmentStreamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssignmentService_SubscribeMeta_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(AssignmentStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AssignmentServiceServer).SubscribeMeta(m, &assignmentServiceSubscribeMetaServer{stream})
}

type AssignmentService_SubscribeMetaServer interface {
	Send(*MetaResponse) error
	grpc.ServerStream
}

type assignmentServiceSubscribeMetaServer struct {
	grpc.ServerStream
}

func (x *assignmentServiceSubscribeMetaServer) Send(m *MetaResponse) error {
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
		{
			MethodName: "GetMeta",
			Handler:    _AssignmentService_GetMeta_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSome",
			Handler:       _AssignmentService_GetSome_Handler,
			ServerStreams: true,
		},
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
		{
			StreamName:    "SubscribeMeta",
			Handler:       _AssignmentService_SubscribeMeta_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/redirector.v1/services.gen.proto",
}
