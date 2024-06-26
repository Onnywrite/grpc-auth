// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: app.proto

package gen

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AppClient is the client API for App service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AppClient interface {
	Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Login(ctx context.Context, in *AppRequest, opts ...grpc.CallOption) (*AppResponse, error)
	Logout(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Relogin(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*AppResponse, error)
	Check(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
	SetProfile(ctx context.Context, in *ProfileChangeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetProfile(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*Profile, error)
	GetSessions(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*Sessions, error)
	Delete(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Recover(ctx context.Context, in *AppRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type appClient struct {
	cc grpc.ClientConnInterface
}

func NewAppClient(cc grpc.ClientConnInterface) AppClient {
	return &appClient{cc}
}

func (c *appClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/app.App/Ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Login(ctx context.Context, in *AppRequest, opts ...grpc.CallOption) (*AppResponse, error) {
	out := new(AppResponse)
	err := c.cc.Invoke(ctx, "/app.App/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Logout(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/app.App/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Relogin(ctx context.Context, in *RefreshToken, opts ...grpc.CallOption) (*AppResponse, error) {
	out := new(AppResponse)
	err := c.cc.Invoke(ctx, "/app.App/Relogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Check(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/app.App/Check", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) SetProfile(ctx context.Context, in *ProfileChangeRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/app.App/SetProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) GetProfile(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*Profile, error) {
	out := new(Profile)
	err := c.cc.Invoke(ctx, "/app.App/GetProfile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) GetSessions(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*Sessions, error) {
	out := new(Sessions)
	err := c.cc.Invoke(ctx, "/app.App/GetSessions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Delete(ctx context.Context, in *AccessToken, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/app.App/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *appClient) Recover(ctx context.Context, in *AppRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/app.App/Recover", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AppServer is the server API for App service.
// All implementations must embed UnimplementedAppServer
// for forward compatibility
type AppServer interface {
	Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error)
	Login(context.Context, *AppRequest) (*AppResponse, error)
	Logout(context.Context, *RefreshToken) (*emptypb.Empty, error)
	Relogin(context.Context, *RefreshToken) (*AppResponse, error)
	Check(context.Context, *AccessToken) (*emptypb.Empty, error)
	SetProfile(context.Context, *ProfileChangeRequest) (*emptypb.Empty, error)
	GetProfile(context.Context, *AccessToken) (*Profile, error)
	GetSessions(context.Context, *AccessToken) (*Sessions, error)
	Delete(context.Context, *AccessToken) (*emptypb.Empty, error)
	Recover(context.Context, *AppRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedAppServer()
}

// UnimplementedAppServer must be embedded to have forward compatible implementations.
type UnimplementedAppServer struct {
}

func (UnimplementedAppServer) Ping(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedAppServer) Login(context.Context, *AppRequest) (*AppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAppServer) Logout(context.Context, *RefreshToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedAppServer) Relogin(context.Context, *RefreshToken) (*AppResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Relogin not implemented")
}
func (UnimplementedAppServer) Check(context.Context, *AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}
func (UnimplementedAppServer) SetProfile(context.Context, *ProfileChangeRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetProfile not implemented")
}
func (UnimplementedAppServer) GetProfile(context.Context, *AccessToken) (*Profile, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProfile not implemented")
}
func (UnimplementedAppServer) GetSessions(context.Context, *AccessToken) (*Sessions, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSessions not implemented")
}
func (UnimplementedAppServer) Delete(context.Context, *AccessToken) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedAppServer) Recover(context.Context, *AppRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Recover not implemented")
}
func (UnimplementedAppServer) mustEmbedUnimplementedAppServer() {}

// UnsafeAppServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AppServer will
// result in compilation errors.
type UnsafeAppServer interface {
	mustEmbedUnimplementedAppServer()
}

func RegisterAppServer(s grpc.ServiceRegistrar, srv AppServer) {
	s.RegisterService(&App_ServiceDesc, srv)
}

func _App_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Ping(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Login(ctx, req.(*AppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Logout(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Relogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RefreshToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Relogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Relogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Relogin(ctx, req.(*RefreshToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Check_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Check(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Check",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Check(ctx, req.(*AccessToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_SetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProfileChangeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).SetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/SetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).SetProfile(ctx, req.(*ProfileChangeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_GetProfile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).GetProfile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/GetProfile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).GetProfile(ctx, req.(*AccessToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_GetSessions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).GetSessions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/GetSessions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).GetSessions(ctx, req.(*AccessToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccessToken)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Delete(ctx, req.(*AccessToken))
	}
	return interceptor(ctx, in, info, handler)
}

func _App_Recover_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AppServer).Recover(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/app.App/Recover",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AppServer).Recover(ctx, req.(*AppRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// App_ServiceDesc is the grpc.ServiceDesc for App service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var App_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "app.App",
	HandlerType: (*AppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _App_Ping_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _App_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _App_Logout_Handler,
		},
		{
			MethodName: "Relogin",
			Handler:    _App_Relogin_Handler,
		},
		{
			MethodName: "Check",
			Handler:    _App_Check_Handler,
		},
		{
			MethodName: "SetProfile",
			Handler:    _App_SetProfile_Handler,
		},
		{
			MethodName: "GetProfile",
			Handler:    _App_GetProfile_Handler,
		},
		{
			MethodName: "GetSessions",
			Handler:    _App_GetSessions_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _App_Delete_Handler,
		},
		{
			MethodName: "Recover",
			Handler:    _App_Recover_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "app.proto",
}
