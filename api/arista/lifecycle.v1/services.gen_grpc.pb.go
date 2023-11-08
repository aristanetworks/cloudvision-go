// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package lifecycle

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

// DeviceLifecycleSummaryServiceClient is the client API for DeviceLifecycleSummaryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeviceLifecycleSummaryServiceClient interface {
	GetOne(ctx context.Context, in *DeviceLifecycleSummaryRequest, opts ...grpc.CallOption) (*DeviceLifecycleSummaryResponse, error)
	GetAll(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (DeviceLifecycleSummaryService_GetAllClient, error)
	Subscribe(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (DeviceLifecycleSummaryService_SubscribeClient, error)
	GetMeta(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error)
	SubscribeMeta(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (DeviceLifecycleSummaryService_SubscribeMetaClient, error)
}

type deviceLifecycleSummaryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeviceLifecycleSummaryServiceClient(cc grpc.ClientConnInterface) DeviceLifecycleSummaryServiceClient {
	return &deviceLifecycleSummaryServiceClient{cc}
}

func (c *deviceLifecycleSummaryServiceClient) GetOne(ctx context.Context, in *DeviceLifecycleSummaryRequest, opts ...grpc.CallOption) (*DeviceLifecycleSummaryResponse, error) {
	out := new(DeviceLifecycleSummaryResponse)
	err := c.cc.Invoke(ctx, "/arista.lifecycle.v1.DeviceLifecycleSummaryService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceLifecycleSummaryServiceClient) GetAll(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (DeviceLifecycleSummaryService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &DeviceLifecycleSummaryService_ServiceDesc.Streams[0], "/arista.lifecycle.v1.DeviceLifecycleSummaryService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &deviceLifecycleSummaryServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DeviceLifecycleSummaryService_GetAllClient interface {
	Recv() (*DeviceLifecycleSummaryStreamResponse, error)
	grpc.ClientStream
}

type deviceLifecycleSummaryServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *deviceLifecycleSummaryServiceGetAllClient) Recv() (*DeviceLifecycleSummaryStreamResponse, error) {
	m := new(DeviceLifecycleSummaryStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deviceLifecycleSummaryServiceClient) Subscribe(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (DeviceLifecycleSummaryService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &DeviceLifecycleSummaryService_ServiceDesc.Streams[1], "/arista.lifecycle.v1.DeviceLifecycleSummaryService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &deviceLifecycleSummaryServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DeviceLifecycleSummaryService_SubscribeClient interface {
	Recv() (*DeviceLifecycleSummaryStreamResponse, error)
	grpc.ClientStream
}

type deviceLifecycleSummaryServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *deviceLifecycleSummaryServiceSubscribeClient) Recv() (*DeviceLifecycleSummaryStreamResponse, error) {
	m := new(DeviceLifecycleSummaryStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *deviceLifecycleSummaryServiceClient) GetMeta(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (*MetaResponse, error) {
	out := new(MetaResponse)
	err := c.cc.Invoke(ctx, "/arista.lifecycle.v1.DeviceLifecycleSummaryService/GetMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceLifecycleSummaryServiceClient) SubscribeMeta(ctx context.Context, in *DeviceLifecycleSummaryStreamRequest, opts ...grpc.CallOption) (DeviceLifecycleSummaryService_SubscribeMetaClient, error) {
	stream, err := c.cc.NewStream(ctx, &DeviceLifecycleSummaryService_ServiceDesc.Streams[2], "/arista.lifecycle.v1.DeviceLifecycleSummaryService/SubscribeMeta", opts...)
	if err != nil {
		return nil, err
	}
	x := &deviceLifecycleSummaryServiceSubscribeMetaClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type DeviceLifecycleSummaryService_SubscribeMetaClient interface {
	Recv() (*MetaResponse, error)
	grpc.ClientStream
}

type deviceLifecycleSummaryServiceSubscribeMetaClient struct {
	grpc.ClientStream
}

func (x *deviceLifecycleSummaryServiceSubscribeMetaClient) Recv() (*MetaResponse, error) {
	m := new(MetaResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// DeviceLifecycleSummaryServiceServer is the server API for DeviceLifecycleSummaryService service.
// All implementations must embed UnimplementedDeviceLifecycleSummaryServiceServer
// for forward compatibility
type DeviceLifecycleSummaryServiceServer interface {
	GetOne(context.Context, *DeviceLifecycleSummaryRequest) (*DeviceLifecycleSummaryResponse, error)
	GetAll(*DeviceLifecycleSummaryStreamRequest, DeviceLifecycleSummaryService_GetAllServer) error
	Subscribe(*DeviceLifecycleSummaryStreamRequest, DeviceLifecycleSummaryService_SubscribeServer) error
	GetMeta(context.Context, *DeviceLifecycleSummaryStreamRequest) (*MetaResponse, error)
	SubscribeMeta(*DeviceLifecycleSummaryStreamRequest, DeviceLifecycleSummaryService_SubscribeMetaServer) error
	mustEmbedUnimplementedDeviceLifecycleSummaryServiceServer()
}

// UnimplementedDeviceLifecycleSummaryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDeviceLifecycleSummaryServiceServer struct {
}

func (UnimplementedDeviceLifecycleSummaryServiceServer) GetOne(context.Context, *DeviceLifecycleSummaryRequest) (*DeviceLifecycleSummaryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedDeviceLifecycleSummaryServiceServer) GetAll(*DeviceLifecycleSummaryStreamRequest, DeviceLifecycleSummaryService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedDeviceLifecycleSummaryServiceServer) Subscribe(*DeviceLifecycleSummaryStreamRequest, DeviceLifecycleSummaryService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedDeviceLifecycleSummaryServiceServer) GetMeta(context.Context, *DeviceLifecycleSummaryStreamRequest) (*MetaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeta not implemented")
}
func (UnimplementedDeviceLifecycleSummaryServiceServer) SubscribeMeta(*DeviceLifecycleSummaryStreamRequest, DeviceLifecycleSummaryService_SubscribeMetaServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeMeta not implemented")
}
func (UnimplementedDeviceLifecycleSummaryServiceServer) mustEmbedUnimplementedDeviceLifecycleSummaryServiceServer() {
}

// UnsafeDeviceLifecycleSummaryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeviceLifecycleSummaryServiceServer will
// result in compilation errors.
type UnsafeDeviceLifecycleSummaryServiceServer interface {
	mustEmbedUnimplementedDeviceLifecycleSummaryServiceServer()
}

func RegisterDeviceLifecycleSummaryServiceServer(s grpc.ServiceRegistrar, srv DeviceLifecycleSummaryServiceServer) {
	s.RegisterService(&DeviceLifecycleSummaryService_ServiceDesc, srv)
}

func _DeviceLifecycleSummaryService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceLifecycleSummaryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceLifecycleSummaryServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.lifecycle.v1.DeviceLifecycleSummaryService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceLifecycleSummaryServiceServer).GetOne(ctx, req.(*DeviceLifecycleSummaryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceLifecycleSummaryService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeviceLifecycleSummaryStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeviceLifecycleSummaryServiceServer).GetAll(m, &deviceLifecycleSummaryServiceGetAllServer{stream})
}

type DeviceLifecycleSummaryService_GetAllServer interface {
	Send(*DeviceLifecycleSummaryStreamResponse) error
	grpc.ServerStream
}

type deviceLifecycleSummaryServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *deviceLifecycleSummaryServiceGetAllServer) Send(m *DeviceLifecycleSummaryStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _DeviceLifecycleSummaryService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeviceLifecycleSummaryStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeviceLifecycleSummaryServiceServer).Subscribe(m, &deviceLifecycleSummaryServiceSubscribeServer{stream})
}

type DeviceLifecycleSummaryService_SubscribeServer interface {
	Send(*DeviceLifecycleSummaryStreamResponse) error
	grpc.ServerStream
}

type deviceLifecycleSummaryServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *deviceLifecycleSummaryServiceSubscribeServer) Send(m *DeviceLifecycleSummaryStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _DeviceLifecycleSummaryService_GetMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeviceLifecycleSummaryStreamRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceLifecycleSummaryServiceServer).GetMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.lifecycle.v1.DeviceLifecycleSummaryService/GetMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceLifecycleSummaryServiceServer).GetMeta(ctx, req.(*DeviceLifecycleSummaryStreamRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceLifecycleSummaryService_SubscribeMeta_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DeviceLifecycleSummaryStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DeviceLifecycleSummaryServiceServer).SubscribeMeta(m, &deviceLifecycleSummaryServiceSubscribeMetaServer{stream})
}

type DeviceLifecycleSummaryService_SubscribeMetaServer interface {
	Send(*MetaResponse) error
	grpc.ServerStream
}

type deviceLifecycleSummaryServiceSubscribeMetaServer struct {
	grpc.ServerStream
}

func (x *deviceLifecycleSummaryServiceSubscribeMetaServer) Send(m *MetaResponse) error {
	return x.ServerStream.SendMsg(m)
}

// DeviceLifecycleSummaryService_ServiceDesc is the grpc.ServiceDesc for DeviceLifecycleSummaryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeviceLifecycleSummaryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.lifecycle.v1.DeviceLifecycleSummaryService",
	HandlerType: (*DeviceLifecycleSummaryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _DeviceLifecycleSummaryService_GetOne_Handler,
		},
		{
			MethodName: "GetMeta",
			Handler:    _DeviceLifecycleSummaryService_GetMeta_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _DeviceLifecycleSummaryService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _DeviceLifecycleSummaryService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeMeta",
			Handler:       _DeviceLifecycleSummaryService_SubscribeMeta_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/lifecycle.v1/services.gen.proto",
}