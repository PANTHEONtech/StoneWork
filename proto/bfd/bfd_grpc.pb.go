// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.1.0
// - protoc             v3.17.3
// source: bfd/bfd.proto

package bfd

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

// BFDWatcherClient is the client API for BFDWatcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BFDWatcherClient interface {
	// WatchBFDEvents allows to subscribe for BFD events.
	WatchBFDEvents(ctx context.Context, in *WatchBFDEventsRequest, opts ...grpc.CallOption) (BFDWatcher_WatchBFDEventsClient, error)
}

type bFDWatcherClient struct {
	cc grpc.ClientConnInterface
}

func NewBFDWatcherClient(cc grpc.ClientConnInterface) BFDWatcherClient {
	return &bFDWatcherClient{cc}
}

func (c *bFDWatcherClient) WatchBFDEvents(ctx context.Context, in *WatchBFDEventsRequest, opts ...grpc.CallOption) (BFDWatcher_WatchBFDEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &BFDWatcher_ServiceDesc.Streams[0], "/bfd.BFDWatcher/WatchBFDEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &bFDWatcherWatchBFDEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type BFDWatcher_WatchBFDEventsClient interface {
	Recv() (*BFDEvent, error)
	grpc.ClientStream
}

type bFDWatcherWatchBFDEventsClient struct {
	grpc.ClientStream
}

func (x *bFDWatcherWatchBFDEventsClient) Recv() (*BFDEvent, error) {
	m := new(BFDEvent)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// BFDWatcherServer is the server API for BFDWatcher service.
// All implementations must embed UnimplementedBFDWatcherServer
// for forward compatibility
type BFDWatcherServer interface {
	// WatchBFDEvents allows to subscribe for BFD events.
	WatchBFDEvents(*WatchBFDEventsRequest, BFDWatcher_WatchBFDEventsServer) error
	mustEmbedUnimplementedBFDWatcherServer()
}

// UnimplementedBFDWatcherServer must be embedded to have forward compatible implementations.
type UnimplementedBFDWatcherServer struct {
}

func (UnimplementedBFDWatcherServer) WatchBFDEvents(*WatchBFDEventsRequest, BFDWatcher_WatchBFDEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchBFDEvents not implemented")
}
func (UnimplementedBFDWatcherServer) mustEmbedUnimplementedBFDWatcherServer() {}

// UnsafeBFDWatcherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BFDWatcherServer will
// result in compilation errors.
type UnsafeBFDWatcherServer interface {
	mustEmbedUnimplementedBFDWatcherServer()
}

func RegisterBFDWatcherServer(s grpc.ServiceRegistrar, srv BFDWatcherServer) {
	s.RegisterService(&BFDWatcher_ServiceDesc, srv)
}

func _BFDWatcher_WatchBFDEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchBFDEventsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BFDWatcherServer).WatchBFDEvents(m, &bFDWatcherWatchBFDEventsServer{stream})
}

type BFDWatcher_WatchBFDEventsServer interface {
	Send(*BFDEvent) error
	grpc.ServerStream
}

type bFDWatcherWatchBFDEventsServer struct {
	grpc.ServerStream
}

func (x *bFDWatcherWatchBFDEventsServer) Send(m *BFDEvent) error {
	return x.ServerStream.SendMsg(m)
}

// BFDWatcher_ServiceDesc is the grpc.ServiceDesc for BFDWatcher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BFDWatcher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bfd.BFDWatcher",
	HandlerType: (*BFDWatcherServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchBFDEvents",
			Handler:       _BFDWatcher_WatchBFDEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "bfd/bfd.proto",
}
