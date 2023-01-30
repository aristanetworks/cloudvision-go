// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package identityprovider

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

// OAuthConfigServiceClient is the client API for OAuthConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OAuthConfigServiceClient interface {
	GetOne(ctx context.Context, in *OAuthConfigRequest, opts ...grpc.CallOption) (*OAuthConfigResponse, error)
	GetAll(ctx context.Context, in *OAuthConfigStreamRequest, opts ...grpc.CallOption) (OAuthConfigService_GetAllClient, error)
	Subscribe(ctx context.Context, in *OAuthConfigStreamRequest, opts ...grpc.CallOption) (OAuthConfigService_SubscribeClient, error)
	Set(ctx context.Context, in *OAuthConfigSetRequest, opts ...grpc.CallOption) (*OAuthConfigSetResponse, error)
	Delete(ctx context.Context, in *OAuthConfigDeleteRequest, opts ...grpc.CallOption) (*OAuthConfigDeleteResponse, error)
	DeleteAll(ctx context.Context, in *OAuthConfigDeleteAllRequest, opts ...grpc.CallOption) (OAuthConfigService_DeleteAllClient, error)
}

type oAuthConfigServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOAuthConfigServiceClient(cc grpc.ClientConnInterface) OAuthConfigServiceClient {
	return &oAuthConfigServiceClient{cc}
}

func (c *oAuthConfigServiceClient) GetOne(ctx context.Context, in *OAuthConfigRequest, opts ...grpc.CallOption) (*OAuthConfigResponse, error) {
	out := new(OAuthConfigResponse)
	err := c.cc.Invoke(ctx, "/arista.identityprovider.v1.OAuthConfigService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oAuthConfigServiceClient) GetAll(ctx context.Context, in *OAuthConfigStreamRequest, opts ...grpc.CallOption) (OAuthConfigService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &OAuthConfigService_ServiceDesc.Streams[0], "/arista.identityprovider.v1.OAuthConfigService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &oAuthConfigServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OAuthConfigService_GetAllClient interface {
	Recv() (*OAuthConfigStreamResponse, error)
	grpc.ClientStream
}

type oAuthConfigServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *oAuthConfigServiceGetAllClient) Recv() (*OAuthConfigStreamResponse, error) {
	m := new(OAuthConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *oAuthConfigServiceClient) Subscribe(ctx context.Context, in *OAuthConfigStreamRequest, opts ...grpc.CallOption) (OAuthConfigService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &OAuthConfigService_ServiceDesc.Streams[1], "/arista.identityprovider.v1.OAuthConfigService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &oAuthConfigServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OAuthConfigService_SubscribeClient interface {
	Recv() (*OAuthConfigStreamResponse, error)
	grpc.ClientStream
}

type oAuthConfigServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *oAuthConfigServiceSubscribeClient) Recv() (*OAuthConfigStreamResponse, error) {
	m := new(OAuthConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *oAuthConfigServiceClient) Set(ctx context.Context, in *OAuthConfigSetRequest, opts ...grpc.CallOption) (*OAuthConfigSetResponse, error) {
	out := new(OAuthConfigSetResponse)
	err := c.cc.Invoke(ctx, "/arista.identityprovider.v1.OAuthConfigService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oAuthConfigServiceClient) Delete(ctx context.Context, in *OAuthConfigDeleteRequest, opts ...grpc.CallOption) (*OAuthConfigDeleteResponse, error) {
	out := new(OAuthConfigDeleteResponse)
	err := c.cc.Invoke(ctx, "/arista.identityprovider.v1.OAuthConfigService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oAuthConfigServiceClient) DeleteAll(ctx context.Context, in *OAuthConfigDeleteAllRequest, opts ...grpc.CallOption) (OAuthConfigService_DeleteAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &OAuthConfigService_ServiceDesc.Streams[2], "/arista.identityprovider.v1.OAuthConfigService/DeleteAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &oAuthConfigServiceDeleteAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type OAuthConfigService_DeleteAllClient interface {
	Recv() (*OAuthConfigDeleteAllResponse, error)
	grpc.ClientStream
}

type oAuthConfigServiceDeleteAllClient struct {
	grpc.ClientStream
}

func (x *oAuthConfigServiceDeleteAllClient) Recv() (*OAuthConfigDeleteAllResponse, error) {
	m := new(OAuthConfigDeleteAllResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// OAuthConfigServiceServer is the server API for OAuthConfigService service.
// All implementations must embed UnimplementedOAuthConfigServiceServer
// for forward compatibility
type OAuthConfigServiceServer interface {
	GetOne(context.Context, *OAuthConfigRequest) (*OAuthConfigResponse, error)
	GetAll(*OAuthConfigStreamRequest, OAuthConfigService_GetAllServer) error
	Subscribe(*OAuthConfigStreamRequest, OAuthConfigService_SubscribeServer) error
	Set(context.Context, *OAuthConfigSetRequest) (*OAuthConfigSetResponse, error)
	Delete(context.Context, *OAuthConfigDeleteRequest) (*OAuthConfigDeleteResponse, error)
	DeleteAll(*OAuthConfigDeleteAllRequest, OAuthConfigService_DeleteAllServer) error
	mustEmbedUnimplementedOAuthConfigServiceServer()
}

// UnimplementedOAuthConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOAuthConfigServiceServer struct {
}

func (UnimplementedOAuthConfigServiceServer) GetOne(context.Context, *OAuthConfigRequest) (*OAuthConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedOAuthConfigServiceServer) GetAll(*OAuthConfigStreamRequest, OAuthConfigService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedOAuthConfigServiceServer) Subscribe(*OAuthConfigStreamRequest, OAuthConfigService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedOAuthConfigServiceServer) Set(context.Context, *OAuthConfigSetRequest) (*OAuthConfigSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedOAuthConfigServiceServer) Delete(context.Context, *OAuthConfigDeleteRequest) (*OAuthConfigDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedOAuthConfigServiceServer) DeleteAll(*OAuthConfigDeleteAllRequest, OAuthConfigService_DeleteAllServer) error {
	return status.Errorf(codes.Unimplemented, "method DeleteAll not implemented")
}
func (UnimplementedOAuthConfigServiceServer) mustEmbedUnimplementedOAuthConfigServiceServer() {}

// UnsafeOAuthConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OAuthConfigServiceServer will
// result in compilation errors.
type UnsafeOAuthConfigServiceServer interface {
	mustEmbedUnimplementedOAuthConfigServiceServer()
}

func RegisterOAuthConfigServiceServer(s grpc.ServiceRegistrar, srv OAuthConfigServiceServer) {
	s.RegisterService(&OAuthConfigService_ServiceDesc, srv)
}

func _OAuthConfigService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OAuthConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OAuthConfigServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.identityprovider.v1.OAuthConfigService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OAuthConfigServiceServer).GetOne(ctx, req.(*OAuthConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OAuthConfigService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OAuthConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OAuthConfigServiceServer).GetAll(m, &oAuthConfigServiceGetAllServer{stream})
}

type OAuthConfigService_GetAllServer interface {
	Send(*OAuthConfigStreamResponse) error
	grpc.ServerStream
}

type oAuthConfigServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *oAuthConfigServiceGetAllServer) Send(m *OAuthConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _OAuthConfigService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OAuthConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OAuthConfigServiceServer).Subscribe(m, &oAuthConfigServiceSubscribeServer{stream})
}

type OAuthConfigService_SubscribeServer interface {
	Send(*OAuthConfigStreamResponse) error
	grpc.ServerStream
}

type oAuthConfigServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *oAuthConfigServiceSubscribeServer) Send(m *OAuthConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _OAuthConfigService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OAuthConfigSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OAuthConfigServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.identityprovider.v1.OAuthConfigService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OAuthConfigServiceServer).Set(ctx, req.(*OAuthConfigSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OAuthConfigService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OAuthConfigDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OAuthConfigServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.identityprovider.v1.OAuthConfigService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OAuthConfigServiceServer).Delete(ctx, req.(*OAuthConfigDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OAuthConfigService_DeleteAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(OAuthConfigDeleteAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(OAuthConfigServiceServer).DeleteAll(m, &oAuthConfigServiceDeleteAllServer{stream})
}

type OAuthConfigService_DeleteAllServer interface {
	Send(*OAuthConfigDeleteAllResponse) error
	grpc.ServerStream
}

type oAuthConfigServiceDeleteAllServer struct {
	grpc.ServerStream
}

func (x *oAuthConfigServiceDeleteAllServer) Send(m *OAuthConfigDeleteAllResponse) error {
	return x.ServerStream.SendMsg(m)
}

// OAuthConfigService_ServiceDesc is the grpc.ServiceDesc for OAuthConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OAuthConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.identityprovider.v1.OAuthConfigService",
	HandlerType: (*OAuthConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _OAuthConfigService_GetOne_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _OAuthConfigService_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _OAuthConfigService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _OAuthConfigService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _OAuthConfigService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DeleteAll",
			Handler:       _OAuthConfigService_DeleteAll_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/identityprovider.v1/services.gen.proto",
}

// SAMLConfigServiceClient is the client API for SAMLConfigService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SAMLConfigServiceClient interface {
	GetOne(ctx context.Context, in *SAMLConfigRequest, opts ...grpc.CallOption) (*SAMLConfigResponse, error)
	GetAll(ctx context.Context, in *SAMLConfigStreamRequest, opts ...grpc.CallOption) (SAMLConfigService_GetAllClient, error)
	Subscribe(ctx context.Context, in *SAMLConfigStreamRequest, opts ...grpc.CallOption) (SAMLConfigService_SubscribeClient, error)
	Set(ctx context.Context, in *SAMLConfigSetRequest, opts ...grpc.CallOption) (*SAMLConfigSetResponse, error)
	Delete(ctx context.Context, in *SAMLConfigDeleteRequest, opts ...grpc.CallOption) (*SAMLConfigDeleteResponse, error)
	DeleteAll(ctx context.Context, in *SAMLConfigDeleteAllRequest, opts ...grpc.CallOption) (SAMLConfigService_DeleteAllClient, error)
}

type sAMLConfigServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSAMLConfigServiceClient(cc grpc.ClientConnInterface) SAMLConfigServiceClient {
	return &sAMLConfigServiceClient{cc}
}

func (c *sAMLConfigServiceClient) GetOne(ctx context.Context, in *SAMLConfigRequest, opts ...grpc.CallOption) (*SAMLConfigResponse, error) {
	out := new(SAMLConfigResponse)
	err := c.cc.Invoke(ctx, "/arista.identityprovider.v1.SAMLConfigService/GetOne", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sAMLConfigServiceClient) GetAll(ctx context.Context, in *SAMLConfigStreamRequest, opts ...grpc.CallOption) (SAMLConfigService_GetAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &SAMLConfigService_ServiceDesc.Streams[0], "/arista.identityprovider.v1.SAMLConfigService/GetAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &sAMLConfigServiceGetAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SAMLConfigService_GetAllClient interface {
	Recv() (*SAMLConfigStreamResponse, error)
	grpc.ClientStream
}

type sAMLConfigServiceGetAllClient struct {
	grpc.ClientStream
}

func (x *sAMLConfigServiceGetAllClient) Recv() (*SAMLConfigStreamResponse, error) {
	m := new(SAMLConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sAMLConfigServiceClient) Subscribe(ctx context.Context, in *SAMLConfigStreamRequest, opts ...grpc.CallOption) (SAMLConfigService_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &SAMLConfigService_ServiceDesc.Streams[1], "/arista.identityprovider.v1.SAMLConfigService/Subscribe", opts...)
	if err != nil {
		return nil, err
	}
	x := &sAMLConfigServiceSubscribeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SAMLConfigService_SubscribeClient interface {
	Recv() (*SAMLConfigStreamResponse, error)
	grpc.ClientStream
}

type sAMLConfigServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *sAMLConfigServiceSubscribeClient) Recv() (*SAMLConfigStreamResponse, error) {
	m := new(SAMLConfigStreamResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *sAMLConfigServiceClient) Set(ctx context.Context, in *SAMLConfigSetRequest, opts ...grpc.CallOption) (*SAMLConfigSetResponse, error) {
	out := new(SAMLConfigSetResponse)
	err := c.cc.Invoke(ctx, "/arista.identityprovider.v1.SAMLConfigService/Set", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sAMLConfigServiceClient) Delete(ctx context.Context, in *SAMLConfigDeleteRequest, opts ...grpc.CallOption) (*SAMLConfigDeleteResponse, error) {
	out := new(SAMLConfigDeleteResponse)
	err := c.cc.Invoke(ctx, "/arista.identityprovider.v1.SAMLConfigService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sAMLConfigServiceClient) DeleteAll(ctx context.Context, in *SAMLConfigDeleteAllRequest, opts ...grpc.CallOption) (SAMLConfigService_DeleteAllClient, error) {
	stream, err := c.cc.NewStream(ctx, &SAMLConfigService_ServiceDesc.Streams[2], "/arista.identityprovider.v1.SAMLConfigService/DeleteAll", opts...)
	if err != nil {
		return nil, err
	}
	x := &sAMLConfigServiceDeleteAllClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type SAMLConfigService_DeleteAllClient interface {
	Recv() (*SAMLConfigDeleteAllResponse, error)
	grpc.ClientStream
}

type sAMLConfigServiceDeleteAllClient struct {
	grpc.ClientStream
}

func (x *sAMLConfigServiceDeleteAllClient) Recv() (*SAMLConfigDeleteAllResponse, error) {
	m := new(SAMLConfigDeleteAllResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SAMLConfigServiceServer is the server API for SAMLConfigService service.
// All implementations must embed UnimplementedSAMLConfigServiceServer
// for forward compatibility
type SAMLConfigServiceServer interface {
	GetOne(context.Context, *SAMLConfigRequest) (*SAMLConfigResponse, error)
	GetAll(*SAMLConfigStreamRequest, SAMLConfigService_GetAllServer) error
	Subscribe(*SAMLConfigStreamRequest, SAMLConfigService_SubscribeServer) error
	Set(context.Context, *SAMLConfigSetRequest) (*SAMLConfigSetResponse, error)
	Delete(context.Context, *SAMLConfigDeleteRequest) (*SAMLConfigDeleteResponse, error)
	DeleteAll(*SAMLConfigDeleteAllRequest, SAMLConfigService_DeleteAllServer) error
	mustEmbedUnimplementedSAMLConfigServiceServer()
}

// UnimplementedSAMLConfigServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSAMLConfigServiceServer struct {
}

func (UnimplementedSAMLConfigServiceServer) GetOne(context.Context, *SAMLConfigRequest) (*SAMLConfigResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOne not implemented")
}
func (UnimplementedSAMLConfigServiceServer) GetAll(*SAMLConfigStreamRequest, SAMLConfigService_GetAllServer) error {
	return status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedSAMLConfigServiceServer) Subscribe(*SAMLConfigStreamRequest, SAMLConfigService_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedSAMLConfigServiceServer) Set(context.Context, *SAMLConfigSetRequest) (*SAMLConfigSetResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Set not implemented")
}
func (UnimplementedSAMLConfigServiceServer) Delete(context.Context, *SAMLConfigDeleteRequest) (*SAMLConfigDeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedSAMLConfigServiceServer) DeleteAll(*SAMLConfigDeleteAllRequest, SAMLConfigService_DeleteAllServer) error {
	return status.Errorf(codes.Unimplemented, "method DeleteAll not implemented")
}
func (UnimplementedSAMLConfigServiceServer) mustEmbedUnimplementedSAMLConfigServiceServer() {}

// UnsafeSAMLConfigServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SAMLConfigServiceServer will
// result in compilation errors.
type UnsafeSAMLConfigServiceServer interface {
	mustEmbedUnimplementedSAMLConfigServiceServer()
}

func RegisterSAMLConfigServiceServer(s grpc.ServiceRegistrar, srv SAMLConfigServiceServer) {
	s.RegisterService(&SAMLConfigService_ServiceDesc, srv)
}

func _SAMLConfigService_GetOne_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SAMLConfigRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SAMLConfigServiceServer).GetOne(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.identityprovider.v1.SAMLConfigService/GetOne",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SAMLConfigServiceServer).GetOne(ctx, req.(*SAMLConfigRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SAMLConfigService_GetAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SAMLConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SAMLConfigServiceServer).GetAll(m, &sAMLConfigServiceGetAllServer{stream})
}

type SAMLConfigService_GetAllServer interface {
	Send(*SAMLConfigStreamResponse) error
	grpc.ServerStream
}

type sAMLConfigServiceGetAllServer struct {
	grpc.ServerStream
}

func (x *sAMLConfigServiceGetAllServer) Send(m *SAMLConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SAMLConfigService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SAMLConfigStreamRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SAMLConfigServiceServer).Subscribe(m, &sAMLConfigServiceSubscribeServer{stream})
}

type SAMLConfigService_SubscribeServer interface {
	Send(*SAMLConfigStreamResponse) error
	grpc.ServerStream
}

type sAMLConfigServiceSubscribeServer struct {
	grpc.ServerStream
}

func (x *sAMLConfigServiceSubscribeServer) Send(m *SAMLConfigStreamResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _SAMLConfigService_Set_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SAMLConfigSetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SAMLConfigServiceServer).Set(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.identityprovider.v1.SAMLConfigService/Set",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SAMLConfigServiceServer).Set(ctx, req.(*SAMLConfigSetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SAMLConfigService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SAMLConfigDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SAMLConfigServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/arista.identityprovider.v1.SAMLConfigService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SAMLConfigServiceServer).Delete(ctx, req.(*SAMLConfigDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SAMLConfigService_DeleteAll_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(SAMLConfigDeleteAllRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SAMLConfigServiceServer).DeleteAll(m, &sAMLConfigServiceDeleteAllServer{stream})
}

type SAMLConfigService_DeleteAllServer interface {
	Send(*SAMLConfigDeleteAllResponse) error
	grpc.ServerStream
}

type sAMLConfigServiceDeleteAllServer struct {
	grpc.ServerStream
}

func (x *sAMLConfigServiceDeleteAllServer) Send(m *SAMLConfigDeleteAllResponse) error {
	return x.ServerStream.SendMsg(m)
}

// SAMLConfigService_ServiceDesc is the grpc.ServiceDesc for SAMLConfigService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SAMLConfigService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "arista.identityprovider.v1.SAMLConfigService",
	HandlerType: (*SAMLConfigServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetOne",
			Handler:    _SAMLConfigService_GetOne_Handler,
		},
		{
			MethodName: "Set",
			Handler:    _SAMLConfigService_Set_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _SAMLConfigService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetAll",
			Handler:       _SAMLConfigService_GetAll_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Subscribe",
			Handler:       _SAMLConfigService_Subscribe_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "DeleteAll",
			Handler:       _SAMLConfigService_DeleteAll_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "arista/identityprovider.v1/services.gen.proto",
}
