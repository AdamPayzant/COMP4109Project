// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package smvshost

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

// ClientHostClient is the client API for ClientHost service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientHostClient interface {
	ReKey(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Status, error)
	DeleteMessage(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*Status, error)
	// Stuff for a handshake, basically just establishes a secret between users
	InitializeConvo(ctx context.Context, in *InitMessage, opts ...grpc.CallOption) (*Status, error)
	ConfirmConvo(ctx context.Context, in *InitMessage, opts ...grpc.CallOption) (*Status, error)
	// Messaging calls
	SendText(ctx context.Context, in *ClientText, opts ...grpc.CallOption) (*Status, error)
	RecieveText(ctx context.Context, in *H2HText, opts ...grpc.CallOption) (*Status, error)
	GetConversation(ctx context.Context, in *Username, opts ...grpc.CallOption) (*Conversation, error)
}

type clientHostClient struct {
	cc grpc.ClientConnInterface
}

func NewClientHostClient(cc grpc.ClientConnInterface) ClientHostClient {
	return &clientHostClient{cc}
}

func (c *clientHostClient) ReKey(ctx context.Context, in *Token, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/ReKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientHostClient) DeleteMessage(ctx context.Context, in *DeleteReq, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/DeleteMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientHostClient) InitializeConvo(ctx context.Context, in *InitMessage, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/InitializeConvo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientHostClient) ConfirmConvo(ctx context.Context, in *InitMessage, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/ConfirmConvo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientHostClient) SendText(ctx context.Context, in *ClientText, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/SendText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientHostClient) RecieveText(ctx context.Context, in *H2HText, opts ...grpc.CallOption) (*Status, error) {
	out := new(Status)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/RecieveText", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clientHostClient) GetConversation(ctx context.Context, in *Username, opts ...grpc.CallOption) (*Conversation, error) {
	out := new(Conversation)
	err := c.cc.Invoke(ctx, "/smvs.clientHost/GetConversation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientHostServer is the server API for ClientHost service.
// All implementations must embed UnimplementedClientHostServer
// for forward compatibility
type ClientHostServer interface {
	ReKey(context.Context, *Token) (*Status, error)
	DeleteMessage(context.Context, *DeleteReq) (*Status, error)
	// Stuff for a handshake, basically just establishes a secret between users
	InitializeConvo(context.Context, *InitMessage) (*Status, error)
	ConfirmConvo(context.Context, *InitMessage) (*Status, error)
	// Messaging calls
	SendText(context.Context, *ClientText) (*Status, error)
	RecieveText(context.Context, *H2HText) (*Status, error)
	GetConversation(context.Context, *Username) (*Conversation, error)
	mustEmbedUnimplementedClientHostServer()
}

// UnimplementedClientHostServer must be embedded to have forward compatible implementations.
type UnimplementedClientHostServer struct {
}

func (UnimplementedClientHostServer) ReKey(context.Context, *Token) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReKey not implemented")
}
func (UnimplementedClientHostServer) DeleteMessage(context.Context, *DeleteReq) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteMessage not implemented")
}
func (UnimplementedClientHostServer) InitializeConvo(context.Context, *InitMessage) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InitializeConvo not implemented")
}
func (UnimplementedClientHostServer) ConfirmConvo(context.Context, *InitMessage) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfirmConvo not implemented")
}
func (UnimplementedClientHostServer) SendText(context.Context, *ClientText) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendText not implemented")
}
func (UnimplementedClientHostServer) RecieveText(context.Context, *H2HText) (*Status, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RecieveText not implemented")
}
func (UnimplementedClientHostServer) GetConversation(context.Context, *Username) (*Conversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConversation not implemented")
}
func (UnimplementedClientHostServer) mustEmbedUnimplementedClientHostServer() {}

// UnsafeClientHostServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientHostServer will
// result in compilation errors.
type UnsafeClientHostServer interface {
	mustEmbedUnimplementedClientHostServer()
}

func RegisterClientHostServer(s grpc.ServiceRegistrar, srv ClientHostServer) {
	s.RegisterService(&ClientHost_ServiceDesc, srv)
}

func _ClientHost_ReKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).ReKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/ReKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).ReKey(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientHost_DeleteMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).DeleteMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/DeleteMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).DeleteMessage(ctx, req.(*DeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientHost_InitializeConvo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).InitializeConvo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/InitializeConvo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).InitializeConvo(ctx, req.(*InitMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientHost_ConfirmConvo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).ConfirmConvo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/ConfirmConvo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).ConfirmConvo(ctx, req.(*InitMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientHost_SendText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientText)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).SendText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/SendText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).SendText(ctx, req.(*ClientText))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientHost_RecieveText_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(H2HText)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).RecieveText(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/RecieveText",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).RecieveText(ctx, req.(*H2HText))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClientHost_GetConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Username)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientHostServer).GetConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/smvs.clientHost/GetConversation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientHostServer).GetConversation(ctx, req.(*Username))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientHost_ServiceDesc is the grpc.ServiceDesc for ClientHost service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientHost_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "smvs.clientHost",
	HandlerType: (*ClientHostServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReKey",
			Handler:    _ClientHost_ReKey_Handler,
		},
		{
			MethodName: "DeleteMessage",
			Handler:    _ClientHost_DeleteMessage_Handler,
		},
		{
			MethodName: "InitializeConvo",
			Handler:    _ClientHost_InitializeConvo_Handler,
		},
		{
			MethodName: "ConfirmConvo",
			Handler:    _ClientHost_ConfirmConvo_Handler,
		},
		{
			MethodName: "SendText",
			Handler:    _ClientHost_SendText_Handler,
		},
		{
			MethodName: "RecieveText",
			Handler:    _ClientHost_RecieveText_Handler,
		},
		{
			MethodName: "GetConversation",
			Handler:    _ClientHost_GetConversation_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "smvshost/host.proto",
}
