// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package configstatus

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

// ConfigDiffServiceClient is the client API for ConfigDiffService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigDiffServiceClient interface {
	GetOne(ctx context.Context, in *ConfigDiffRequest, opts ...grpc.CallOption) (*ConfigDiffResponse, error)
	GetAll(ctx context.Context, in *ConfigDiffStreamRequest, opts ...grpc.CallOption) (ConfigDiffService_GetAllClient, error)
	Subscribe(ctx context.Context, in *ConfigDiffStreamRequest, opts ...grpc.CallOption) (ConfigDiffService_SubscribeClient, error)
}

type configDiffServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigDiffServiceClient(cc grpc.ClientConnInterface) ConfigDiffServiceClient {
	return &configDiffServiceClient{cc}
}

func (c *configDiffServiceClient) GetOne(ctx context.Context, in *ConfigDiffRequest, opts ...grpc.CallOption) (*ConfigDiffResponse, error) {
	out := new(ConfigDiffResponse)
	err := c.cc.Invoke(ctx, "/arista.configstatus.v1.ConfigDiffService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configDiffServiceClient) GetAll(ctx context.Context, in *ConfigDiffStreamRequest, opts ...grpc.CallOption) (ConfigDiffService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConfigDiffService_ServiceDesc.Streams[0], "/arista.configstatus.v1.ConfigDiffService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &configDiffServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConfigDiffService_GetAllClient interface {
	Recv() (*ConfigDiffStreamResponse, error)
	grpc.ClientStream
}

type configDiffServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *configDiffServiceGetAllClient) Recv() (*ConfigDiffStreamResponse, error) {
	m := new(ConfigDiffStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *configDiffServiceClient) Subscribe(ctx context.Context, in *ConfigDiffStreamRequest, opts ...grpc.CallOption) (ConfigDiffService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConfigDiffService_ServiceDesc.Streams[1], "/arista.configstatus.v1.ConfigDiffService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &configDiffServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConfigDiffService_SubscribeClient interface {
	Recv() (*ConfigDiffStreamResponse, error)
	grpc.ClientStream
}

type configDiffServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *configDiffServiceSubscribeClient) Recv() (*ConfigDiffStreamResponse, error) {
	m := new(ConfigDiffStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConfigDiffServiceServer is the server API for ConfigDiffService service.
// All implementations must embed UnimplementedConfigDiffServiceServer
// for forward compatibility
type ConfigDiffServiceServer interface {
	GetOne(context.Context, *ConfigDiffRequest) (*ConfigDiffResponse, error)
	GetAll(*ConfigDiffStreamRequest, ConfigDiffService_GetAllServer) error
	Subscribe(*ConfigDiffStreamRequest, ConfigDiffService_SubscribeServer) error
	mustEmbedUnimplementedConfigDiffServiceServer()
}

// UnimplementedConfigDiffServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConfigDiffServiceServer struct {
}

func (UnimplementedConfigDiffServiceServer) GetOne(context.Context, *ConfigDiffRequest) (*ConfigDiffResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedConfigDiffServiceServer) GetAll(*ConfigDiffStreamRequest, ConfigDiffService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedConfigDiffServiceServer) Subscribe(*ConfigDiffStreamRequest, ConfigDiffService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedConfigDiffServiceServer) mustEmbedUnimplementedConfigDiffServiceServer() {}

// UnsafeConfigDiffServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigDiffServiceServer will
// result in compilation errors.
type UnsafeConfigDiffServiceServer interface {
	mustEmbedUnimplementedConfigDiffServiceServer()
}

func RegisterConfigDiffServiceServer(s grpc.ServiceRegistrar, srv ConfigDiffServiceServer) {
	s.RegisterService(&ConfigDiffService_ServiceDesc, srv)
}

func _ConfigDiffService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigDiffRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigDiffServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.configstatus.v1.ConfigDiffService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigDiffServiceServer).GetOne(ctx, req.(*ConfigDiffRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigDiffService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConfigDiffStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigDiffServiceServer).GetAll(m, &configDiffServiceGetAllServer{stream})
}

type ConfigDiffService_GetAllServer interface {
	Send(*ConfigDiffStreamResponse) error
	grpc.ServerStream
}

type configDiffServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *configDiffServiceGetAllServer) Send(m *ConfigDiffStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ConfigDiffService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConfigDiffStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigDiffServiceServer).Subscribe(m, &configDiffServiceSubscribeServer{stream})
}

type ConfigDiffService_SubscribeServer interface {
	Send(*ConfigDiffStreamResponse) error
	grpc.ServerStream
}

type configDiffServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *configDiffServiceSubscribeServer) Send(m *ConfigDiffStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ConfigDiffService_ServiceDesc is the grpc.ServiceDesc for ConfigDiffService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigDiffService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.configstatus.v1.ConfigDiffService",
	HandlerType: (*ConfigDiffServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _ConfigDiffService_GetOne_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _ConfigDiffService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _ConfigDiffService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/configstatus.v1/services.gen.proto",
}

// ConfigurationServiceClient is the client API for ConfigurationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ConfigurationServiceClient interface {
	GetOne(ctx context.Context, in *ConfigurationRequest, opts ...grpc.CallOption) (*ConfigurationResponse, error)
	GetAll(ctx context.Context, in *ConfigurationStreamRequest, opts ...grpc.CallOption) (ConfigurationService_GetAllClient, error)
	Subscribe(ctx context.Context, in *ConfigurationStreamRequest, opts ...grpc.CallOption) (ConfigurationService_SubscribeClient, error)
}

type configurationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewConfigurationServiceClient(cc grpc.ClientConnInterface) ConfigurationServiceClient {
	return &configurationServiceClient{cc}
}

func (c *configurationServiceClient) GetOne(ctx context.Context, in *ConfigurationRequest, opts ...grpc.CallOption) (*ConfigurationResponse, error) {
	out := new(ConfigurationResponse)
	err := c.cc.Invoke(ctx, "/arista.configstatus.v1.ConfigurationService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *configurationServiceClient) GetAll(ctx context.Context, in *ConfigurationStreamRequest, opts ...grpc.CallOption) (ConfigurationService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConfigurationService_ServiceDesc.Streams[0], "/arista.configstatus.v1.ConfigurationService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &configurationServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConfigurationService_GetAllClient interface {
	Recv() (*ConfigurationStreamResponse, error)
	grpc.ClientStream
}

type configurationServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *configurationServiceGetAllClient) Recv() (*ConfigurationStreamResponse, error) {
	m := new(ConfigurationStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *configurationServiceClient) Subscribe(ctx context.Context, in *ConfigurationStreamRequest, opts ...grpc.CallOption) (ConfigurationService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &ConfigurationService_ServiceDesc.Streams[1], "/arista.configstatus.v1.ConfigurationService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &configurationServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ConfigurationService_SubscribeClient interface {
	Recv() (*ConfigurationStreamResponse, error)
	grpc.ClientStream
}

type configurationServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *configurationServiceSubscribeClient) Recv() (*ConfigurationStreamResponse, error) {
	m := new(ConfigurationStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ConfigurationServiceServer is the server API for ConfigurationService service.
// All implementations must embed UnimplementedConfigurationServiceServer
// for forward compatibility
type ConfigurationServiceServer interface {
	GetOne(context.Context, *ConfigurationRequest) (*ConfigurationResponse, error)
	GetAll(*ConfigurationStreamRequest, ConfigurationService_GetAllServer) error
	Subscribe(*ConfigurationStreamRequest, ConfigurationService_SubscribeServer) error
	mustEmbedUnimplementedConfigurationServiceServer()
}

// UnimplementedConfigurationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedConfigurationServiceServer struct {
}

func (UnimplementedConfigurationServiceServer) GetOne(context.Context, *ConfigurationRequest) (*ConfigurationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedConfigurationServiceServer) GetAll(*ConfigurationStreamRequest, ConfigurationService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedConfigurationServiceServer) Subscribe(*ConfigurationStreamRequest, ConfigurationService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedConfigurationServiceServer) mustEmbedUnimplementedConfigurationServiceServer() {}

// UnsafeConfigurationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ConfigurationServiceServer will
// result in compilation errors.
type UnsafeConfigurationServiceServer interface {
	mustEmbedUnimplementedConfigurationServiceServer()
}

func RegisterConfigurationServiceServer(s grpc.ServiceRegistrar, srv ConfigurationServiceServer) {
	s.RegisterService(&ConfigurationService_ServiceDesc, srv)
}

func _ConfigurationService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConfigurationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ConfigurationServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.configstatus.v1.ConfigurationService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ConfigurationServiceServer).GetOne(ctx, req.(*ConfigurationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ConfigurationService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConfigurationStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigurationServiceServer).GetAll(m, &configurationServiceGetAllServer{stream})
}

type ConfigurationService_GetAllServer interface {
	Send(*ConfigurationStreamResponse) error
	grpc.ServerStream
}

type configurationServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *configurationServiceGetAllServer) Send(m *ConfigurationStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _ConfigurationService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConfigurationStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ConfigurationServiceServer).Subscribe(m, &configurationServiceSubscribeServer{stream})
}

type ConfigurationService_SubscribeServer interface {
	Send(*ConfigurationStreamResponse) error
	grpc.ServerStream
}

type configurationServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *configurationServiceSubscribeServer) Send(m *ConfigurationStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// ConfigurationService_ServiceDesc is the grpc.ServiceDesc for ConfigurationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ConfigurationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.configstatus.v1.ConfigurationService",
	HandlerType: (*ConfigurationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _ConfigurationService_GetOne_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _ConfigurationService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _ConfigurationService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/configstatus.v1/services.gen.proto",
}

// SummaryServiceClient is the client API for SummaryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SummaryServiceClient interface {
	GetOne(ctx context.Context, in *SummaryRequest, opts ...grpc.CallOption) (*SummaryResponse, error)
	GetAll(ctx context.Context, in *SummaryStreamRequest, opts ...grpc.CallOption) (SummaryService_GetAllClient, error)
	Subscribe(ctx context.Context, in *SummaryStreamRequest, opts ...grpc.CallOption) (SummaryService_SubscribeClient, error)
}

type summaryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSummaryServiceClient(cc grpc.ClientConnInterface) SummaryServiceClient {
	return &summaryServiceClient{cc}
}

func (c *summaryServiceClient) GetOne(ctx context.Context, in *SummaryRequest, opts ...grpc.CallOption) (*SummaryResponse, error) {
	out := new(SummaryResponse)
	err := c.cc.Invoke(ctx, "/arista.configstatus.v1.SummaryService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *summaryServiceClient) GetAll(ctx context.Context, in *SummaryStreamRequest, opts ...grpc.CallOption) (SummaryService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &SummaryService_ServiceDesc.Streams[0], "/arista.configstatus.v1.SummaryService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &summaryServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SummaryService_GetAllClient interface {
	Recv() (*SummaryStreamResponse, error)
	grpc.ClientStream
}

type summaryServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *summaryServiceGetAllClient) Recv() (*SummaryStreamResponse, error) {
	m := new(SummaryStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *summaryServiceClient) Subscribe(ctx context.Context, in *SummaryStreamRequest, opts ...grpc.CallOption) (SummaryService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &SummaryService_ServiceDesc.Streams[1], "/arista.configstatus.v1.SummaryService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &summaryServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SummaryService_SubscribeClient interface {
	Recv() (*SummaryStreamResponse, error)
	grpc.ClientStream
}

type summaryServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *summaryServiceSubscribeClient) Recv() (*SummaryStreamResponse, error) {
	m := new(SummaryStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SummaryServiceServer is the server API for SummaryService service.
// All implementations must embed UnimplementedSummaryServiceServer
// for forward compatibility
type SummaryServiceServer interface {
	GetOne(context.Context, *SummaryRequest) (*SummaryResponse, error)
	GetAll(*SummaryStreamRequest, SummaryService_GetAllServer) error
	Subscribe(*SummaryStreamRequest, SummaryService_SubscribeServer) error
	mustEmbedUnimplementedSummaryServiceServer()
}

// UnimplementedSummaryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSummaryServiceServer struct {
}

func (UnimplementedSummaryServiceServer) GetOne(context.Context, *SummaryRequest) (*SummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedSummaryServiceServer) GetAll(*SummaryStreamRequest, SummaryService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedSummaryServiceServer) Subscribe(*SummaryStreamRequest, SummaryService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedSummaryServiceServer) mustEmbedUnimplementedSummaryServiceServer() {}

// UnsafeSummaryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SummaryServiceServer will
// result in compilation errors.
type UnsafeSummaryServiceServer interface {
	mustEmbedUnimplementedSummaryServiceServer()
}

func RegisterSummaryServiceServer(s grpc.ServiceRegistrar, srv SummaryServiceServer) {
	s.RegisterService(&SummaryService_ServiceDesc, srv)
}

func _SummaryService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SummaryServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.configstatus.v1.SummaryService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SummaryServiceServer).GetOne(ctx, req.(*SummaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SummaryService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SummaryStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SummaryServiceServer).GetAll(m, &summaryServiceGetAllServer{stream})
}

type SummaryService_GetAllServer interface {
	Send(*SummaryStreamResponse) error
	grpc.ServerStream
}

type summaryServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *summaryServiceGetAllServer) Send(m *SummaryStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SummaryService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SummaryStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SummaryServiceServer).Subscribe(m, &summaryServiceSubscribeServer{stream})
}

type SummaryService_SubscribeServer interface {
	Send(*SummaryStreamResponse) error
	grpc.ServerStream
}

type summaryServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *summaryServiceSubscribeServer) Send(m *SummaryStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

// SummaryService_ServiceDesc is the grpc.ServiceDesc for SummaryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SummaryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.configstatus.v1.SummaryService",
	HandlerType: (*SummaryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _SummaryService_GetOne_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _SummaryService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _SummaryService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/configstatus.v1/services.gen.proto",
}
