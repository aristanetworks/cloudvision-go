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
// source: arista/endpointlocation.v1/services.gen.proto

package endpointlocation

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
	EndpointLocationService_GetOne_FullMethodName        = "/arista.endpointlocation.v1.EndpointLocationService/GetOne"
	EndpointLocationService_GetSome_FullMethodName       = "/arista.endpointlocation.v1.EndpointLocationService/GetSome"
	EndpointLocationService_GetAll_FullMethodName        = "/arista.endpointlocation.v1.EndpointLocationService/GetAll"
	EndpointLocationService_Subscribe_FullMethodName     = "/arista.endpointlocation.v1.EndpointLocationService/Subscribe"
	EndpointLocationService_GetMeta_FullMethodName       = "/arista.endpointlocation.v1.EndpointLocationService/GetMeta"
	EndpointLocationService_SubscribeMeta_FullMethodName = "/arista.endpointlocation.v1.EndpointLocationService/SubscribeMeta"
)

// EndpointLocationServiceClient is the client API for EndpointLocationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EndpointLocationServiceClient interface {
	GetOne(ctx context.Context, in *EndpointLocationRequest, opts ...grpc.CallOption) (*EndpointLocationResponse, error)
	GetSome(ctx context.Context, in *EndpointLocationSomeRequest, opts ...grpc.CallOption) (EndpointLocationService_GetSomeClient, error)
	GetAll(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (EndpointLocationService_GetAllClient, error)
	Subscribe(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (EndpointLocationService_SubscribeClient, error)
	GetMeta(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error)
	SubscribeMeta(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (EndpointLocationService_SubscribeMetaClient, error)
}

type endpointLocationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEndpointLocationServiceClient(cc grpc.ClientConnInterface) EndpointLocationServiceClient {
	return &endpointLocationServiceClient{cc}
}

func (c *endpointLocationServiceClient) GetOne(ctx context.Context, in *EndpointLocationRequest, opts ...grpc.CallOption) (*EndpointLocationResponse, error) {
	out := new(EndpointLocationResponse)
	err := c.cc.Invoke(ctx, EndpointLocationService_GetOne_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointLocationServiceClient) GetSome(ctx context.Context, in *EndpointLocationSomeRequest, opts ...grpc.CallOption) (EndpointLocationService_GetSomeClient, error) {
	stream, err := c.cc.NewStream(ctx, &EndpointLocationService_ServiceDesc.Streams[0], EndpointLocationService_GetSome_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &endpointLocationServiceGetSomeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EndpointLocationService_GetSomeClient interface {
	Recv() (*EndpointLocationSomeResponse, error)
	grpc.ClientStream
}

type endpointLocationServiceGetSomeClient struct {
	grpc.ClientStream
}

func (x *endpointLocationServiceGetSomeClient) Recv() (*EndpointLocationSomeResponse, error) {
	m := new(EndpointLocationSomeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *endpointLocationServiceClient) GetAll(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (EndpointLocationService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &EndpointLocationService_ServiceDesc.Streams[1], EndpointLocationService_GetAll_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &endpointLocationServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EndpointLocationService_GetAllClient interface {
	Recv() (*EndpointLocationStreamResponse, error)
	grpc.ClientStream
}

type endpointLocationServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *endpointLocationServiceGetAllClient) Recv() (*EndpointLocationStreamResponse, error) {
	m := new(EndpointLocationStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *endpointLocationServiceClient) Subscribe(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (EndpointLocationService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &EndpointLocationService_ServiceDesc.Streams[2], EndpointLocationService_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &endpointLocationServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EndpointLocationService_SubscribeClient interface {
	Recv() (*EndpointLocationStreamResponse, error)
	grpc.ClientStream
}

type endpointLocationServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *endpointLocationServiceSubscribeClient) Recv() (*EndpointLocationStreamResponse, error) {
	m := new(EndpointLocationStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *endpointLocationServiceClient) GetMeta(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error) {
	out := new(MetaResponse)
	err := c.cc.Invoke(ctx, EndpointLocationService_GetMeta_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *endpointLocationServiceClient) SubscribeMeta(ctx context.Context, in *EndpointLocationStreamRequest, opts ...grpc.CallOption) (EndpointLocationService_SubscribeMetaClient, error) {
	stream, err := c.cc.NewStream(ctx, &EndpointLocationService_ServiceDesc.Streams[3], EndpointLocationService_SubscribeMeta_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &endpointLocationServiceSubscribeMetaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type EndpointLocationService_SubscribeMetaClient interface {
	Recv() (*MetaResponse, error)
	grpc.ClientStream
}

type endpointLocationServiceSubscribeMetaClient struct {
	grpc.ClientStream
}

func (x *endpointLocationServiceSubscribeMetaClient) Recv() (*MetaResponse, error) {
	m := new(MetaResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// EndpointLocationServiceServer is the server API for EndpointLocationService service.
// All implementations must embed UnimplementedEndpointLocationServiceServer
// for forward compatibility
type EndpointLocationServiceServer interface {
	GetOne(context.Context, *EndpointLocationRequest) (*EndpointLocationResponse, error)
	GetSome(*EndpointLocationSomeRequest, EndpointLocationService_GetSomeServer) error
	GetAll(*EndpointLocationStreamRequest, EndpointLocationService_GetAllServer) error
	Subscribe(*EndpointLocationStreamRequest, EndpointLocationService_SubscribeServer) error
	GetMeta(context.Context, *EndpointLocationStreamRequest) (*MetaResponse, error)
	SubscribeMeta(*EndpointLocationStreamRequest, EndpointLocationService_SubscribeMetaServer) error
	mustEmbedUnimplementedEndpointLocationServiceServer()
}

// UnimplementedEndpointLocationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEndpointLocationServiceServer struct {
}

func (UnimplementedEndpointLocationServiceServer) GetOne(context.Context, *EndpointLocationRequest) (*EndpointLocationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedEndpointLocationServiceServer) GetSome(*EndpointLocationSomeRequest, EndpointLocationService_GetSomeServer) error {
	return status.Errorf(codes.Unimplemented, "method GetSome not implemented")
}
func (UnimplementedEndpointLocationServiceServer) GetAll(*EndpointLocationStreamRequest, EndpointLocationService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedEndpointLocationServiceServer) Subscribe(*EndpointLocationStreamRequest, EndpointLocationService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedEndpointLocationServiceServer) GetMeta(context.Context, *EndpointLocationStreamRequest) (*MetaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeta not implemented")
}
func (UnimplementedEndpointLocationServiceServer) SubscribeMeta(*EndpointLocationStreamRequest, EndpointLocationService_SubscribeMetaServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeMeta not implemented")
}
func (UnimplementedEndpointLocationServiceServer) mustEmbedUnimplementedEndpointLocationServiceServer() {
}

// UnsafeEndpointLocationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EndpointLocationServiceServer will
// result in compilation errors.
type UnsafeEndpointLocationServiceServer interface {
	mustEmbedUnimplementedEndpointLocationServiceServer()
}

func RegisterEndpointLocationServiceServer(s grpc.ServiceRegistrar, srv EndpointLocationServiceServer) {
	s.RegisterService(&EndpointLocationService_ServiceDesc, srv)
}

func _EndpointLocationService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndpointLocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointLocationServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointLocationService_GetOne_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointLocationServiceServer).GetOne(ctx, req.(*EndpointLocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointLocationService_GetSome_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EndpointLocationSomeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EndpointLocationServiceServer).GetSome(m, &endpointLocationServiceGetSomeServer{stream})
}

type EndpointLocationService_GetSomeServer interface {
	Send(*EndpointLocationSomeResponse) error
	grpc.ServerStream
}

type endpointLocationServiceGetSomeServer struct {
	grpc.ServerStream
}

func (x *endpointLocationServiceGetSomeServer) Send(m *EndpointLocationSomeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _EndpointLocationService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EndpointLocationStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EndpointLocationServiceServer).GetAll(m, &endpointLocationServiceGetAllServer{stream})
}

type EndpointLocationService_GetAllServer interface {
	Send(*EndpointLocationStreamResponse) error
	grpc.ServerStream
}

type endpointLocationServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *endpointLocationServiceGetAllServer) Send(m *EndpointLocationStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _EndpointLocationService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EndpointLocationStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EndpointLocationServiceServer).Subscribe(m, &endpointLocationServiceSubscribeServer{stream})
}

type EndpointLocationService_SubscribeServer interface {
	Send(*EndpointLocationStreamResponse) error
	grpc.ServerStream
}

type endpointLocationServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *endpointLocationServiceSubscribeServer) Send(m *EndpointLocationStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _EndpointLocationService_GetMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EndpointLocationStreamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EndpointLocationServiceServer).GetMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EndpointLocationService_GetMeta_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EndpointLocationServiceServer).GetMeta(ctx, req.(*EndpointLocationStreamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EndpointLocationService_SubscribeMeta_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EndpointLocationStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EndpointLocationServiceServer).SubscribeMeta(m, &endpointLocationServiceSubscribeMetaServer{stream})
}

type EndpointLocationService_SubscribeMetaServer interface {
	Send(*MetaResponse) error
	grpc.ServerStream
}

type endpointLocationServiceSubscribeMetaServer struct {
	grpc.ServerStream
}

func (x *endpointLocationServiceSubscribeMetaServer) Send(m *MetaResponse) error {
	return x.ServerStream.SendMsg(m)
}

// EndpointLocationService_ServiceDesc is the grpc.ServiceDesc for EndpointLocationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EndpointLocationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.endpointlocation.v1.EndpointLocationService",
	HandlerType: (*EndpointLocationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _EndpointLocationService_GetOne_Handler,
		},
		{
			MethodName: "GetMeta",
			Handler:    _EndpointLocationService_GetMeta_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetSome",
			Handler:       _EndpointLocationService_GetSome_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetAll",
			Handler:       _EndpointLocationService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _EndpointLocationService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeMeta",
			Handler:       _EndpointLocationService_SubscribeMeta_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/endpointlocation.v1/services.gen.proto",
}
