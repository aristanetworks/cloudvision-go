// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package syslog

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

// ExportServiceClient is the client API for ExportService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExportServiceClient interface {
	GetOne(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error)
	GetAll(ctx context.Context, in *ExportStreamRequest, opts ...grpc.CallOption) (ExportService_GetAllClient, error)
	Subscribe(ctx context.Context, in *ExportStreamRequest, opts ...grpc.CallOption) (ExportService_SubscribeClient, error)
}

type exportServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExportServiceClient(cc grpc.ClientConnInterface) ExportServiceClient {
	return &exportServiceClient{cc}
}

func (c *exportServiceClient) GetOne(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error) {
	out := new(ExportResponse)
	err := c.cc.Invoke(ctx, "/arista.syslog.v1.ExportService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exportServiceClient) GetAll(ctx context.Context, in *ExportStreamRequest, opts ...grpc.CallOption) (ExportService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExportService_ServiceDesc.Streams[0], "/arista.syslog.v1.ExportService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &exportServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExportService_GetAllClient interface {
	Recv() (*ExportStreamResponse, error)
	grpc.ClientStream
}

type exportServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *exportServiceGetAllClient) Recv() (*ExportStreamResponse, error) {
	m := new(ExportStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exportServiceClient) Subscribe(ctx context.Context, in *ExportStreamRequest, opts ...grpc.CallOption) (ExportService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExportService_ServiceDesc.Streams[1], "/arista.syslog.v1.ExportService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &exportServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExportService_SubscribeClient interface {
	Recv() (*ExportStreamResponse, error)
	grpc.ClientStream
}

type exportServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *exportServiceSubscribeClient) Recv() (*ExportStreamResponse, error) {
	m := new(ExportStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExportServiceServer is the server API for ExportService service.
// All implementations must embed UnimplementedExportServiceServer
// for forward compatibility
type ExportServiceServer interface {
	GetOne(context.Context, *ExportRequest) (*ExportResponse, error)
	GetAll(*ExportStreamRequest, ExportService_GetAllServer) error
	Subscribe(*ExportStreamRequest, ExportService_SubscribeServer) error
	mustEmbedUnimplementedExportServiceServer()
}

// UnimplementedExportServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExportServiceServer struct {
}

func (UnimplementedExportServiceServer) GetOne(context.Context, *ExportRequest) (*ExportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedExportServiceServer) GetAll(*ExportStreamRequest, ExportService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedExportServiceServer) Subscribe(*ExportStreamRequest, ExportService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedExportServiceServer) mustEmbedUnimplementedExportServiceServer() {}

// UnsafeExportServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExportServiceServer will
// result in compilation errors.
type UnsafeExportServiceServer interface {
	mustEmbedUnimplementedExportServiceServer()
}

func RegisterExportServiceServer(s grpc.ServiceRegistrar, srv ExportServiceServer) {
	s.RegisterService(&ExportService_ServiceDesc, srv)
}

func _ExportService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExportServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.syslog.v1.ExportService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExportServiceServer).GetOne(ctx, req.(*ExportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExportService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExportStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExportServiceServer).GetAll(m, &exportServiceGetAllServer{stream})
}

type ExportService_GetAllServer interface {
	Send(*ExportStreamResponse) error
	grpc.ServerStream
}

type exportServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *exportServiceGetAllServer) Send(m *ExportStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ExportService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExportStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExportServiceServer).Subscribe(m, &exportServiceSubscribeServer{stream})
}

type ExportService_SubscribeServer interface {
	Send(*ExportStreamResponse) error
	grpc.ServerStream
}

type exportServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *exportServiceSubscribeServer) Send(m *ExportStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ExportService_ServiceDesc is the grpc.ServiceDesc for ExportService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExportService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.syslog.v1.ExportService",
	HandlerType: (*ExportServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _ExportService_GetOne_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _ExportService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _ExportService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/syslog.v1/services.gen.proto",
}

// ExportConfigServiceClient is the client API for ExportConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExportConfigServiceClient interface {
	GetOne(ctx context.Context, in *ExportConfigRequest, opts ...grpc.CallOption) (*ExportConfigResponse, error)
	GetAll(ctx context.Context, in *ExportConfigStreamRequest, opts ...grpc.CallOption) (ExportConfigService_GetAllClient, error)
	Subscribe(ctx context.Context, in *ExportConfigStreamRequest, opts ...grpc.CallOption) (ExportConfigService_SubscribeClient, error)
	Set(ctx context.Context, in *ExportConfigSetRequest, opts ...grpc.CallOption) (*ExportConfigSetResponse, error)
	Delete(ctx context.Context, in *ExportConfigDeleteRequest, opts ...grpc.CallOption) (*ExportConfigDeleteResponse, error)
	DeleteAll(ctx context.Context, in *ExportConfigDeleteAllRequest, opts ...grpc.CallOption) (ExportConfigService_DeleteAllClient, error)
}

type exportConfigServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExportConfigServiceClient(cc grpc.ClientConnInterface) ExportConfigServiceClient {
	return &exportConfigServiceClient{cc}
}

func (c *exportConfigServiceClient) GetOne(ctx context.Context, in *ExportConfigRequest, opts ...grpc.CallOption) (*ExportConfigResponse, error) {
	out := new(ExportConfigResponse)
	err := c.cc.Invoke(ctx, "/arista.syslog.v1.ExportConfigService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exportConfigServiceClient) GetAll(ctx context.Context, in *ExportConfigStreamRequest, opts ...grpc.CallOption) (ExportConfigService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExportConfigService_ServiceDesc.Streams[0], "/arista.syslog.v1.ExportConfigService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &exportConfigServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExportConfigService_GetAllClient interface {
	Recv() (*ExportConfigStreamResponse, error)
	grpc.ClientStream
}

type exportConfigServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *exportConfigServiceGetAllClient) Recv() (*ExportConfigStreamResponse, error) {
	m := new(ExportConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exportConfigServiceClient) Subscribe(ctx context.Context, in *ExportConfigStreamRequest, opts ...grpc.CallOption) (ExportConfigService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExportConfigService_ServiceDesc.Streams[1], "/arista.syslog.v1.ExportConfigService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &exportConfigServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExportConfigService_SubscribeClient interface {
	Recv() (*ExportConfigStreamResponse, error)
	grpc.ClientStream
}

type exportConfigServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *exportConfigServiceSubscribeClient) Recv() (*ExportConfigStreamResponse, error) {
	m := new(ExportConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exportConfigServiceClient) Set(ctx context.Context, in *ExportConfigSetRequest, opts ...grpc.CallOption) (*ExportConfigSetResponse, error) {
	out := new(ExportConfigSetResponse)
	err := c.cc.Invoke(ctx, "/arista.syslog.v1.ExportConfigService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exportConfigServiceClient) Delete(ctx context.Context, in *ExportConfigDeleteRequest, opts ...grpc.CallOption) (*ExportConfigDeleteResponse, error) {
	out := new(ExportConfigDeleteResponse)
	err := c.cc.Invoke(ctx, "/arista.syslog.v1.ExportConfigService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exportConfigServiceClient) DeleteAll(ctx context.Context, in *ExportConfigDeleteAllRequest, opts ...grpc.CallOption) (ExportConfigService_DeleteAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExportConfigService_ServiceDesc.Streams[2], "/arista.syslog.v1.ExportConfigService/DeleteAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &exportConfigServiceDeleteAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExportConfigService_DeleteAllClient interface {
	Recv() (*ExportConfigDeleteAllResponse, error)
	grpc.ClientStream
}

type exportConfigServiceDeleteAllClient struct {
	grpc.ClientStream
}

func (x *exportConfigServiceDeleteAllClient) Recv() (*ExportConfigDeleteAllResponse, error) {
	m := new(ExportConfigDeleteAllResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ExportConfigServiceServer is the server API for ExportConfigService service.
// All implementations must embed UnimplementedExportConfigServiceServer
// for forward compatibility
type ExportConfigServiceServer interface {
	GetOne(context.Context, *ExportConfigRequest) (*ExportConfigResponse, error)
	GetAll(*ExportConfigStreamRequest, ExportConfigService_GetAllServer) error
	Subscribe(*ExportConfigStreamRequest, ExportConfigService_SubscribeServer) error
	Set(context.Context, *ExportConfigSetRequest) (*ExportConfigSetResponse, error)
	Delete(context.Context, *ExportConfigDeleteRequest) (*ExportConfigDeleteResponse, error)
	DeleteAll(*ExportConfigDeleteAllRequest, ExportConfigService_DeleteAllServer) error
	mustEmbedUnimplementedExportConfigServiceServer()
}

// UnimplementedExportConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExportConfigServiceServer struct {
}

func (UnimplementedExportConfigServiceServer) GetOne(context.Context, *ExportConfigRequest) (*ExportConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedExportConfigServiceServer) GetAll(*ExportConfigStreamRequest, ExportConfigService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedExportConfigServiceServer) Subscribe(*ExportConfigStreamRequest, ExportConfigService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedExportConfigServiceServer) Set(context.Context, *ExportConfigSetRequest) (*ExportConfigSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedExportConfigServiceServer) Delete(context.Context, *ExportConfigDeleteRequest) (*ExportConfigDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedExportConfigServiceServer) DeleteAll(*ExportConfigDeleteAllRequest, ExportConfigService_DeleteAllServer) error {
	return status.Errorf(codes.Unimplemented, "method DeleteAll not implemented")
}
func (UnimplementedExportConfigServiceServer) mustEmbedUnimplementedExportConfigServiceServer() {}

// UnsafeExportConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExportConfigServiceServer will
// result in compilation errors.
type UnsafeExportConfigServiceServer interface {
	mustEmbedUnimplementedExportConfigServiceServer()
}

func RegisterExportConfigServiceServer(s grpc.ServiceRegistrar, srv ExportConfigServiceServer) {
	s.RegisterService(&ExportConfigService_ServiceDesc, srv)
}

func _ExportConfigService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExportConfigServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.syslog.v1.ExportConfigService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExportConfigServiceServer).GetOne(ctx, req.(*ExportConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExportConfigService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExportConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExportConfigServiceServer).GetAll(m, &exportConfigServiceGetAllServer{stream})
}

type ExportConfigService_GetAllServer interface {
	Send(*ExportConfigStreamResponse) error
	grpc.ServerStream
}

type exportConfigServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *exportConfigServiceGetAllServer) Send(m *ExportConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ExportConfigService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExportConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExportConfigServiceServer).Subscribe(m, &exportConfigServiceSubscribeServer{stream})
}

type ExportConfigService_SubscribeServer interface {
	Send(*ExportConfigStreamResponse) error
	grpc.ServerStream
}

type exportConfigServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *exportConfigServiceSubscribeServer) Send(m *ExportConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ExportConfigService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportConfigSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExportConfigServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.syslog.v1.ExportConfigService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExportConfigServiceServer).Set(ctx, req.(*ExportConfigSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExportConfigService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportConfigDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExportConfigServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.syslog.v1.ExportConfigService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExportConfigServiceServer).Delete(ctx, req.(*ExportConfigDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExportConfigService_DeleteAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExportConfigDeleteAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExportConfigServiceServer).DeleteAll(m, &exportConfigServiceDeleteAllServer{stream})
}

type ExportConfigService_DeleteAllServer interface {
	Send(*ExportConfigDeleteAllResponse) error
	grpc.ServerStream
}

type exportConfigServiceDeleteAllServer struct {
	grpc.ServerStream
}

func (x *exportConfigServiceDeleteAllServer) Send(m *ExportConfigDeleteAllResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ExportConfigService_ServiceDesc is the grpc.ServiceDesc for ExportConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExportConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.syslog.v1.ExportConfigService",
	HandlerType: (*ExportConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _ExportConfigService_GetOne_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _ExportConfigService_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ExportConfigService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _ExportConfigService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _ExportConfigService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DeleteAll",
			Handler:       _ExportConfigService_DeleteAll_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/syslog.v1/services.gen.proto",
}

// ExportFormatConfigServiceClient is the client API for ExportFormatConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExportFormatConfigServiceClient interface {
	GetOne(ctx context.Context, in *ExportFormatConfigRequest, opts ...grpc.CallOption) (*ExportFormatConfigResponse, error)
	Subscribe(ctx context.Context, in *ExportFormatConfigStreamRequest, opts ...grpc.CallOption) (ExportFormatConfigService_SubscribeClient, error)
	Set(ctx context.Context, in *ExportFormatConfigSetRequest, opts ...grpc.CallOption) (*ExportFormatConfigSetResponse, error)
}

type exportFormatConfigServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewExportFormatConfigServiceClient(cc grpc.ClientConnInterface) ExportFormatConfigServiceClient {
	return &exportFormatConfigServiceClient{cc}
}

func (c *exportFormatConfigServiceClient) GetOne(ctx context.Context, in *ExportFormatConfigRequest, opts ...grpc.CallOption) (*ExportFormatConfigResponse, error) {
	out := new(ExportFormatConfigResponse)
	err := c.cc.Invoke(ctx, "/arista.syslog.v1.ExportFormatConfigService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *exportFormatConfigServiceClient) Subscribe(ctx context.Context, in *ExportFormatConfigStreamRequest, opts ...grpc.CallOption) (ExportFormatConfigService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ExportFormatConfigService_ServiceDesc.Streams[0], "/arista.syslog.v1.ExportFormatConfigService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &exportFormatConfigServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ExportFormatConfigService_SubscribeClient interface {
	Recv() (*ExportFormatConfigStreamResponse, error)
	grpc.ClientStream
}

type exportFormatConfigServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *exportFormatConfigServiceSubscribeClient) Recv() (*ExportFormatConfigStreamResponse, error) {
	m := new(ExportFormatConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *exportFormatConfigServiceClient) Set(ctx context.Context, in *ExportFormatConfigSetRequest, opts ...grpc.CallOption) (*ExportFormatConfigSetResponse, error) {
	out := new(ExportFormatConfigSetResponse)
	err := c.cc.Invoke(ctx, "/arista.syslog.v1.ExportFormatConfigService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExportFormatConfigServiceServer is the server API for ExportFormatConfigService service.
// All implementations must embed UnimplementedExportFormatConfigServiceServer
// for forward compatibility
type ExportFormatConfigServiceServer interface {
	GetOne(context.Context, *ExportFormatConfigRequest) (*ExportFormatConfigResponse, error)
	Subscribe(*ExportFormatConfigStreamRequest, ExportFormatConfigService_SubscribeServer) error
	Set(context.Context, *ExportFormatConfigSetRequest) (*ExportFormatConfigSetResponse, error)
	mustEmbedUnimplementedExportFormatConfigServiceServer()
}

// UnimplementedExportFormatConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedExportFormatConfigServiceServer struct {
}

func (UnimplementedExportFormatConfigServiceServer) GetOne(context.Context, *ExportFormatConfigRequest) (*ExportFormatConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedExportFormatConfigServiceServer) Subscribe(*ExportFormatConfigStreamRequest, ExportFormatConfigService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedExportFormatConfigServiceServer) Set(context.Context, *ExportFormatConfigSetRequest) (*ExportFormatConfigSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedExportFormatConfigServiceServer) mustEmbedUnimplementedExportFormatConfigServiceServer() {
}

// UnsafeExportFormatConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExportFormatConfigServiceServer will
// result in compilation errors.
type UnsafeExportFormatConfigServiceServer interface {
	mustEmbedUnimplementedExportFormatConfigServiceServer()
}

func RegisterExportFormatConfigServiceServer(s grpc.ServiceRegistrar, srv ExportFormatConfigServiceServer) {
	s.RegisterService(&ExportFormatConfigService_ServiceDesc, srv)
}

func _ExportFormatConfigService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportFormatConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExportFormatConfigServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.syslog.v1.ExportFormatConfigService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExportFormatConfigServiceServer).GetOne(ctx, req.(*ExportFormatConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ExportFormatConfigService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ExportFormatConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ExportFormatConfigServiceServer).Subscribe(m, &exportFormatConfigServiceSubscribeServer{stream})
}

type ExportFormatConfigService_SubscribeServer interface {
	Send(*ExportFormatConfigStreamResponse) error
	grpc.ServerStream
}

type exportFormatConfigServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *exportFormatConfigServiceSubscribeServer) Send(m *ExportFormatConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ExportFormatConfigService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportFormatConfigSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExportFormatConfigServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.syslog.v1.ExportFormatConfigService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExportFormatConfigServiceServer).Set(ctx, req.(*ExportFormatConfigSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExportFormatConfigService_ServiceDesc is the grpc.ServiceDesc for ExportFormatConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExportFormatConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.syslog.v1.ExportFormatConfigService",
	HandlerType: (*ExportFormatConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _ExportFormatConfigService_GetOne_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _ExportFormatConfigService_Set_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _ExportFormatConfigService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/syslog.v1/services.gen.proto",
}
