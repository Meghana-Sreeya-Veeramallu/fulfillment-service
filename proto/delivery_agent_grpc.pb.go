// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: delivery_agent.proto

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
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DeliveryAgentService_AddDeliveryAgent_FullMethodName   = "/proto.DeliveryAgentService/AddDeliveryAgent"
	DeliveryAgentService_AssignAgentToOrder_FullMethodName = "/proto.DeliveryAgentService/AssignAgentToOrder"
)

// DeliveryAgentServiceClient is the client API for DeliveryAgentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeliveryAgentServiceClient interface {
	AddDeliveryAgent(ctx context.Context, in *AddDeliveryAgentRequest, opts ...grpc.CallOption) (*AddDeliveryAgentResponse, error)
	AssignAgentToOrder(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AssignAgentToOrderResponse, error)
}

type deliveryAgentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDeliveryAgentServiceClient(cc grpc.ClientConnInterface) DeliveryAgentServiceClient {
	return &deliveryAgentServiceClient{cc}
}

func (c *deliveryAgentServiceClient) AddDeliveryAgent(ctx context.Context, in *AddDeliveryAgentRequest, opts ...grpc.CallOption) (*AddDeliveryAgentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddDeliveryAgentResponse)
	err := c.cc.Invoke(ctx, DeliveryAgentService_AddDeliveryAgent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deliveryAgentServiceClient) AssignAgentToOrder(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AssignAgentToOrderResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AssignAgentToOrderResponse)
	err := c.cc.Invoke(ctx, DeliveryAgentService_AssignAgentToOrder_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeliveryAgentServiceServer is the server API for DeliveryAgentService service.
// All implementations must embed UnimplementedDeliveryAgentServiceServer
// for forward compatibility.
type DeliveryAgentServiceServer interface {
	AddDeliveryAgent(context.Context, *AddDeliveryAgentRequest) (*AddDeliveryAgentResponse, error)
	AssignAgentToOrder(context.Context, *emptypb.Empty) (*AssignAgentToOrderResponse, error)
	mustEmbedUnimplementedDeliveryAgentServiceServer()
}

// UnimplementedDeliveryAgentServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDeliveryAgentServiceServer struct{}

func (UnimplementedDeliveryAgentServiceServer) AddDeliveryAgent(context.Context, *AddDeliveryAgentRequest) (*AddDeliveryAgentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddDeliveryAgent not implemented")
}
func (UnimplementedDeliveryAgentServiceServer) AssignAgentToOrder(context.Context, *emptypb.Empty) (*AssignAgentToOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AssignAgentToOrder not implemented")
}
func (UnimplementedDeliveryAgentServiceServer) mustEmbedUnimplementedDeliveryAgentServiceServer() {}
func (UnimplementedDeliveryAgentServiceServer) testEmbeddedByValue()                              {}

// UnsafeDeliveryAgentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeliveryAgentServiceServer will
// result in compilation errors.
type UnsafeDeliveryAgentServiceServer interface {
	mustEmbedUnimplementedDeliveryAgentServiceServer()
}

func RegisterDeliveryAgentServiceServer(s grpc.ServiceRegistrar, srv DeliveryAgentServiceServer) {
	// If the following call pancis, it indicates UnimplementedDeliveryAgentServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DeliveryAgentService_ServiceDesc, srv)
}

func _DeliveryAgentService_AddDeliveryAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddDeliveryAgentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeliveryAgentServiceServer).AddDeliveryAgent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeliveryAgentService_AddDeliveryAgent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeliveryAgentServiceServer).AddDeliveryAgent(ctx, req.(*AddDeliveryAgentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeliveryAgentService_AssignAgentToOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeliveryAgentServiceServer).AssignAgentToOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeliveryAgentService_AssignAgentToOrder_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeliveryAgentServiceServer).AssignAgentToOrder(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// DeliveryAgentService_ServiceDesc is the grpc.ServiceDesc for DeliveryAgentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeliveryAgentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DeliveryAgentService",
	HandlerType: (*DeliveryAgentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddDeliveryAgent",
			Handler:    _DeliveryAgentService_AddDeliveryAgent_Handler,
		},
		{
			MethodName: "AssignAgentToOrder",
			Handler:    _DeliveryAgentService_AssignAgentToOrder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "delivery_agent.proto",
}