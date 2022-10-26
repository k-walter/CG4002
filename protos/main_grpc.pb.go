// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.7
// source: protos/main.proto

package __

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

// RelayClient is the client API for Relay service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RelayClient interface {
	GetRound(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Relay_GetRoundClient, error)
	// glove (MPU6050) data for AI
	Gesture(ctx context.Context, opts ...grpc.CallOption) (Relay_GestureClient, error)
	// detected by infrared (KY-022 rx, KY-005 tx)
	Shoot(ctx context.Context, in *Event, opts ...grpc.CallOption) (*emptypb.Empty, error)
	Shot(ctx context.Context, in *Event, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type relayClient struct {
	cc grpc.ClientConnInterface
}

func NewRelayClient(cc grpc.ClientConnInterface) RelayClient {
	return &relayClient{cc}
}

func (c *relayClient) GetRound(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (Relay_GetRoundClient, error) {
	stream, err := c.cc.NewStream(ctx, &Relay_ServiceDesc.Streams[0], "/Relay/GetRound", opts...)
	if err != nil {
		return nil, err
	}
	x := &relayGetRoundClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Relay_GetRoundClient interface {
	Recv() (*RndResp, error)
	grpc.ClientStream
}

type relayGetRoundClient struct {
	grpc.ClientStream
}

func (x *relayGetRoundClient) Recv() (*RndResp, error) {
	m := new(RndResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *relayClient) Gesture(ctx context.Context, opts ...grpc.CallOption) (Relay_GestureClient, error) {
	stream, err := c.cc.NewStream(ctx, &Relay_ServiceDesc.Streams[1], "/Relay/Gesture", opts...)
	if err != nil {
		return nil, err
	}
	x := &relayGestureClient{stream}
	return x, nil
}

type Relay_GestureClient interface {
	Send(*Data) error
	CloseAndRecv() (*emptypb.Empty, error)
	grpc.ClientStream
}

type relayGestureClient struct {
	grpc.ClientStream
}

func (x *relayGestureClient) Send(m *Data) error {
	return x.ClientStream.SendMsg(m)
}

func (x *relayGestureClient) CloseAndRecv() (*emptypb.Empty, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(emptypb.Empty)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *relayClient) Shoot(ctx context.Context, in *Event, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Relay/Shoot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *relayClient) Shot(ctx context.Context, in *Event, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Relay/Shot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RelayServer is the server API for Relay service.
// All implementations must embed UnimplementedRelayServer
// for forward compatibility
type RelayServer interface {
	GetRound(*emptypb.Empty, Relay_GetRoundServer) error
	// glove (MPU6050) data for AI
	Gesture(Relay_GestureServer) error
	// detected by infrared (KY-022 rx, KY-005 tx)
	Shoot(context.Context, *Event) (*emptypb.Empty, error)
	Shot(context.Context, *Event) (*emptypb.Empty, error)
	mustEmbedUnimplementedRelayServer()
}

// UnimplementedRelayServer must be embedded to have forward compatible implementations.
type UnimplementedRelayServer struct {
}

func (UnimplementedRelayServer) GetRound(*emptypb.Empty, Relay_GetRoundServer) error {
	return status.Errorf(codes.Unimplemented, "method GetRound not implemented")
}
func (UnimplementedRelayServer) Gesture(Relay_GestureServer) error {
	return status.Errorf(codes.Unimplemented, "method Gesture not implemented")
}
func (UnimplementedRelayServer) Shoot(context.Context, *Event) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shoot not implemented")
}
func (UnimplementedRelayServer) Shot(context.Context, *Event) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Shot not implemented")
}
func (UnimplementedRelayServer) mustEmbedUnimplementedRelayServer() {}

// UnsafeRelayServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RelayServer will
// result in compilation errors.
type UnsafeRelayServer interface {
	mustEmbedUnimplementedRelayServer()
}

func RegisterRelayServer(s grpc.ServiceRegistrar, srv RelayServer) {
	s.RegisterService(&Relay_ServiceDesc, srv)
}

func _Relay_GetRound_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(RelayServer).GetRound(m, &relayGetRoundServer{stream})
}

type Relay_GetRoundServer interface {
	Send(*RndResp) error
	grpc.ServerStream
}

type relayGetRoundServer struct {
	grpc.ServerStream
}

func (x *relayGetRoundServer) Send(m *RndResp) error {
	return x.ServerStream.SendMsg(m)
}

func _Relay_Gesture_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(RelayServer).Gesture(&relayGestureServer{stream})
}

type Relay_GestureServer interface {
	SendAndClose(*emptypb.Empty) error
	Recv() (*Data, error)
	grpc.ServerStream
}

type relayGestureServer struct {
	grpc.ServerStream
}

func (x *relayGestureServer) SendAndClose(m *emptypb.Empty) error {
	return x.ServerStream.SendMsg(m)
}

func (x *relayGestureServer) Recv() (*Data, error) {
	m := new(Data)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _Relay_Shoot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).Shoot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Relay/Shoot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).Shoot(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _Relay_Shot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RelayServer).Shot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Relay/Shot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RelayServer).Shot(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

// Relay_ServiceDesc is the grpc.ServiceDesc for Relay service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Relay_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Relay",
	HandlerType: (*RelayServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Shoot",
			Handler:    _Relay_Shoot_Handler,
		},
		{
			MethodName: "Shot",
			Handler:    _Relay_Shot_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetRound",
			Handler:       _Relay_GetRound_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Gesture",
			Handler:       _Relay_Gesture_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protos/main.proto",
}

// VizClient is the client API for Viz service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VizClient interface {
	Update(ctx context.Context, in *State, opts ...grpc.CallOption) (*emptypb.Empty, error)
	InFov(ctx context.Context, in *Event, opts ...grpc.CallOption) (*InFovResp, error)
}

type vizClient struct {
	cc grpc.ClientConnInterface
}

func NewVizClient(cc grpc.ClientConnInterface) VizClient {
	return &vizClient{cc}
}

func (c *vizClient) Update(ctx context.Context, in *State, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/Viz/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *vizClient) InFov(ctx context.Context, in *Event, opts ...grpc.CallOption) (*InFovResp, error) {
	out := new(InFovResp)
	err := c.cc.Invoke(ctx, "/Viz/InFov", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VizServer is the server API for Viz service.
// All implementations must embed UnimplementedVizServer
// for forward compatibility
type VizServer interface {
	Update(context.Context, *State) (*emptypb.Empty, error)
	InFov(context.Context, *Event) (*InFovResp, error)
	mustEmbedUnimplementedVizServer()
}

// UnimplementedVizServer must be embedded to have forward compatible implementations.
type UnimplementedVizServer struct {
}

func (UnimplementedVizServer) Update(context.Context, *State) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedVizServer) InFov(context.Context, *Event) (*InFovResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InFov not implemented")
}
func (UnimplementedVizServer) mustEmbedUnimplementedVizServer() {}

// UnsafeVizServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VizServer will
// result in compilation errors.
type UnsafeVizServer interface {
	mustEmbedUnimplementedVizServer()
}

func RegisterVizServer(s grpc.ServiceRegistrar, srv VizServer) {
	s.RegisterService(&Viz_ServiceDesc, srv)
}

func _Viz_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(State)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VizServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Viz/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VizServer).Update(ctx, req.(*State))
	}
	return interceptor(ctx, in, info, handler)
}

func _Viz_InFov_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VizServer).InFov(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Viz/InFov",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VizServer).InFov(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

// Viz_ServiceDesc is the grpc.ServiceDesc for Viz service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Viz_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Viz",
	HandlerType: (*VizServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Update",
			Handler:    _Viz_Update_Handler,
		},
		{
			MethodName: "InFov",
			Handler:    _Viz_InFov_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/main.proto",
}

// PynqClient is the client API for Pynq service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PynqClient interface {
	Emit(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Event, error)
	Poll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Event, error)
}

type pynqClient struct {
	cc grpc.ClientConnInterface
}

func NewPynqClient(cc grpc.ClientConnInterface) PynqClient {
	return &pynqClient{cc}
}

func (c *pynqClient) Emit(ctx context.Context, in *Data, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/Pynq/Emit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pynqClient) Poll(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*Event, error) {
	out := new(Event)
	err := c.cc.Invoke(ctx, "/Pynq/Poll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PynqServer is the server API for Pynq service.
// All implementations must embed UnimplementedPynqServer
// for forward compatibility
type PynqServer interface {
	Emit(context.Context, *Data) (*Event, error)
	Poll(context.Context, *emptypb.Empty) (*Event, error)
	mustEmbedUnimplementedPynqServer()
}

// UnimplementedPynqServer must be embedded to have forward compatible implementations.
type UnimplementedPynqServer struct {
}

func (UnimplementedPynqServer) Emit(context.Context, *Data) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Emit not implemented")
}
func (UnimplementedPynqServer) Poll(context.Context, *emptypb.Empty) (*Event, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Poll not implemented")
}
func (UnimplementedPynqServer) mustEmbedUnimplementedPynqServer() {}

// UnsafePynqServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PynqServer will
// result in compilation errors.
type UnsafePynqServer interface {
	mustEmbedUnimplementedPynqServer()
}

func RegisterPynqServer(s grpc.ServiceRegistrar, srv PynqServer) {
	s.RegisterService(&Pynq_ServiceDesc, srv)
}

func _Pynq_Emit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Data)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PynqServer).Emit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pynq/Emit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PynqServer).Emit(ctx, req.(*Data))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pynq_Poll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PynqServer).Poll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Pynq/Poll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PynqServer).Poll(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Pynq_ServiceDesc is the grpc.ServiceDesc for Pynq service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pynq_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Pynq",
	HandlerType: (*PynqServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Emit",
			Handler:    _Pynq_Emit_Handler,
		},
		{
			MethodName: "Poll",
			Handler:    _Pynq_Poll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/main.proto",
}
