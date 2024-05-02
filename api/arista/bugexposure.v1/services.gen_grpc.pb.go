// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

//
// Code generated by boomtown. DO NOT EDIT.
//

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.24.4
// source: arista/bugexposure.v1/services.gen.proto

package bugexposure

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
	BugExposureService_GetOne_FullMethodName        = "/arista.bugexposure.v1.BugExposureService/GetOne"
	BugExposureService_GetAll_FullMethodName        = "/arista.bugexposure.v1.BugExposureService/GetAll"
	BugExposureService_Subscribe_FullMethodName     = "/arista.bugexposure.v1.BugExposureService/Subscribe"
	BugExposureService_GetMeta_FullMethodName       = "/arista.bugexposure.v1.BugExposureService/GetMeta"
	BugExposureService_SubscribeMeta_FullMethodName = "/arista.bugexposure.v1.BugExposureService/SubscribeMeta"
)

// BugExposureServiceClient is the client API for BugExposureService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BugExposureServiceClient interface {
	GetOne(ctx context.Context, in *BugExposureRequest, opts ...grpc.CallOption) (*BugExposureResponse, error)
	GetAll(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (BugExposureService_GetAllClient, error)
	Subscribe(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (BugExposureService_SubscribeClient, error)
	GetMeta(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error)
	SubscribeMeta(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (BugExposureService_SubscribeMetaClient, error)
}

type bugExposureServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBugExposureServiceClient(cc grpc.ClientConnInterface) BugExposureServiceClient {
	return &bugExposureServiceClient{cc}
}

func (c *bugExposureServiceClient) GetOne(ctx context.Context, in *BugExposureRequest, opts ...grpc.CallOption) (*BugExposureResponse, error) {
	out := new(BugExposureResponse)
	err := c.cc.Invoke(ctx, BugExposureService_GetOne_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bugExposureServiceClient) GetAll(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (BugExposureService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &BugExposureService_ServiceDesc.Streams[0], BugExposureService_GetAll_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &bugExposureServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BugExposureService_GetAllClient interface {
	Recv() (*BugExposureStreamResponse, error)
	grpc.ClientStream
}

type bugExposureServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *bugExposureServiceGetAllClient) Recv() (*BugExposureStreamResponse, error) {
	m := new(BugExposureStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bugExposureServiceClient) Subscribe(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (BugExposureService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &BugExposureService_ServiceDesc.Streams[1], BugExposureService_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &bugExposureServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BugExposureService_SubscribeClient interface {
	Recv() (*BugExposureStreamResponse, error)
	grpc.ClientStream
}

type bugExposureServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *bugExposureServiceSubscribeClient) Recv() (*BugExposureStreamResponse, error) {
	m := new(BugExposureStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bugExposureServiceClient) GetMeta(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error) {
	out := new(MetaResponse)
	err := c.cc.Invoke(ctx, BugExposureService_GetMeta_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bugExposureServiceClient) SubscribeMeta(ctx context.Context, in *BugExposureStreamRequest, opts ...grpc.CallOption) (BugExposureService_SubscribeMetaClient, error) {
	stream, err := c.cc.NewStream(ctx, &BugExposureService_ServiceDesc.Streams[2], BugExposureService_SubscribeMeta_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &bugExposureServiceSubscribeMetaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BugExposureService_SubscribeMetaClient interface {
	Recv() (*MetaResponse, error)
	grpc.ClientStream
}

type bugExposureServiceSubscribeMetaClient struct {
	grpc.ClientStream
}

func (x *bugExposureServiceSubscribeMetaClient) Recv() (*MetaResponse, error) {
	m := new(MetaResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BugExposureServiceServer is the server API for BugExposureService service.
// All implementations must embed UnimplementedBugExposureServiceServer
// for forward compatibility
type BugExposureServiceServer interface {
	GetOne(context.Context, *BugExposureRequest) (*BugExposureResponse, error)
	GetAll(*BugExposureStreamRequest, BugExposureService_GetAllServer) error
	Subscribe(*BugExposureStreamRequest, BugExposureService_SubscribeServer) error
	GetMeta(context.Context, *BugExposureStreamRequest) (*MetaResponse, error)
	SubscribeMeta(*BugExposureStreamRequest, BugExposureService_SubscribeMetaServer) error
	mustEmbedUnimplementedBugExposureServiceServer()
}

// UnimplementedBugExposureServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBugExposureServiceServer struct {
}

func (UnimplementedBugExposureServiceServer) GetOne(context.Context, *BugExposureRequest) (*BugExposureResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedBugExposureServiceServer) GetAll(*BugExposureStreamRequest, BugExposureService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedBugExposureServiceServer) Subscribe(*BugExposureStreamRequest, BugExposureService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedBugExposureServiceServer) GetMeta(context.Context, *BugExposureStreamRequest) (*MetaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeta not implemented")
}
func (UnimplementedBugExposureServiceServer) SubscribeMeta(*BugExposureStreamRequest, BugExposureService_SubscribeMetaServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeMeta not implemented")
}
func (UnimplementedBugExposureServiceServer) mustEmbedUnimplementedBugExposureServiceServer() {}

// UnsafeBugExposureServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BugExposureServiceServer will
// result in compilation errors.
type UnsafeBugExposureServiceServer interface {
	mustEmbedUnimplementedBugExposureServiceServer()
}

func RegisterBugExposureServiceServer(s grpc.ServiceRegistrar, srv BugExposureServiceServer) {
	s.RegisterService(&BugExposureService_ServiceDesc, srv)
}

func _BugExposureService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BugExposureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BugExposureServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BugExposureService_GetOne_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BugExposureServiceServer).GetOne(ctx, req.(*BugExposureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BugExposureService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BugExposureStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BugExposureServiceServer).GetAll(m, &bugExposureServiceGetAllServer{stream})
}

type BugExposureService_GetAllServer interface {
	Send(*BugExposureStreamResponse) error
	grpc.ServerStream
}

type bugExposureServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *bugExposureServiceGetAllServer) Send(m *BugExposureStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BugExposureService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BugExposureStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BugExposureServiceServer).Subscribe(m, &bugExposureServiceSubscribeServer{stream})
}

type BugExposureService_SubscribeServer interface {
	Send(*BugExposureStreamResponse) error
	grpc.ServerStream
}

type bugExposureServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *bugExposureServiceSubscribeServer) Send(m *BugExposureStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _BugExposureService_GetMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BugExposureStreamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BugExposureServiceServer).GetMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BugExposureService_GetMeta_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BugExposureServiceServer).GetMeta(ctx, req.(*BugExposureStreamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BugExposureService_SubscribeMeta_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(BugExposureStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BugExposureServiceServer).SubscribeMeta(m, &bugExposureServiceSubscribeMetaServer{stream})
}

type BugExposureService_SubscribeMetaServer interface {
	Send(*MetaResponse) error
	grpc.ServerStream
}

type bugExposureServiceSubscribeMetaServer struct {
	grpc.ServerStream
}

func (x *bugExposureServiceSubscribeMetaServer) Send(m *MetaResponse) error {
	return x.ServerStream.SendMsg(m)
}

// BugExposureService_ServiceDesc is the grpc.ServiceDesc for BugExposureService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BugExposureService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.bugexposure.v1.BugExposureService",
	HandlerType: (*BugExposureServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _BugExposureService_GetOne_Handler,
		},
		{
			MethodName: "GetMeta",
			Handler:    _BugExposureService_GetMeta_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _BugExposureService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _BugExposureService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeMeta",
			Handler:       _BugExposureService_SubscribeMeta_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/bugexposure.v1/services.gen.proto",
}
