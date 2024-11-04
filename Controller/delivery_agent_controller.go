package Controller

import (
	"Fulfillment/Service"
	pb "Fulfillment/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// DeliveryAgentServer represents the gRPC server for delivery agents.
type DeliveryAgentServer struct {
	pb.UnimplementedDeliveryAgentServiceServer
	Service *Service.DeliveryAgentService
}

// NewDeliveryAgentServer creates a new DeliveryAgentServer.
func NewDeliveryAgentServer(service *Service.DeliveryAgentService) *DeliveryAgentServer {
	return &DeliveryAgentServer{Service: service}
}

// AddDeliveryAgent handles the gRPC request for adding a delivery agent.
func (s *DeliveryAgentServer) AddDeliveryAgent(ctx context.Context, req *pb.AddDeliveryAgentRequest) (*pb.AddDeliveryAgentResponse, error) {
	_, err := s.Service.AddDeliveryAgent(req.Name, req.City)
	if err != nil {
		if err.Error() == "name cannot be empty" || err.Error() == "city cannot be empty" {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, grpc.Errorf(codes.Internal, "Failed to add delivery agent: "+err.Error())
	}

	return &pb.AddDeliveryAgentResponse{
		Message: "Delivery agent added successfully",
	}, nil
}

// AssignAgentToOrder handles the gRPC request for assigning an agent to an order.
func (s *DeliveryAgentServer) AssignAgentToOrder(ctx context.Context, req *pb.AssignAgentToOrderRequest) (*pb.AssignAgentToOrderResponse, error) {
	agentID := req.AgentId
	orderID := req.OrderId

	err := s.Service.AssignAgentToOrder(uint(agentID), int(orderID))
	if err != nil {
		if err.Error() == "delivery agent not found" || err.Error() == "delivery agent is not available" || err.Error() == "order does not exist" || err.Error() == "order cannot be assigned" {
			return nil, grpc.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, grpc.Errorf(codes.Internal, "An error occurred: "+err.Error())
	}

	return &pb.AssignAgentToOrderResponse{
		Message: "Delivery agent assigned to order successfully",
	}, nil
}
