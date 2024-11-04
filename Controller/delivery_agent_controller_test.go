package Controller

import (
	orderClient "Fulfillment/Client"
	"Fulfillment/Model"
	"Fulfillment/Repository"
	"Fulfillment/Service"
	pb "Fulfillment/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *gorm.DB
var lis *bufconn.Listener
var server *grpc.Server
var client pb.DeliveryAgentServiceClient

const bufSize = 1024 * 1024 // 1MB

// Setup function to initialize the test database and server
func setup() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&Model.DeliveryAgent{})
	if err != nil {
		panic(err)
	}

	repo := Repository.NewDeliveryAgentRepository(db)
	service := Service.NewDeliveryAgentService(repo)

	lis = bufconn.Listen(bufSize)
	server = grpc.NewServer()
	pb.RegisterDeliveryAgentServiceServer(server, NewDeliveryAgentServer(service))

	go func() {
		if err := server.Serve(lis); err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(
		func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	client = pb.NewDeliveryAgentServiceClient(conn)
}

// Teardown function to stop the server
func teardown() {
	server.Stop()
}

// Test AddDeliveryAgent
func TestAddDeliveryAgentSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	req := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: "Hyderabad"}
	resp, err := client.AddDeliveryAgent(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "Delivery agent added successfully", resp.Message)
}

// Test AddDeliveryAgent with invalid name
func TestAddDeliveryAgentInvalidName(t *testing.T) {
	setup()
	defer teardown()

	req := &pb.AddDeliveryAgentRequest{Name: "", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "name cannot be empty", status.Convert(err).Message())
}

// Test AddDeliveryAgent with invalid city
func TestAddDeliveryAgentInvalidCity(t *testing.T) {
	setup()
	defer teardown()

	req := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: ""}
	_, err := client.AddDeliveryAgent(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "city cannot be empty", status.Convert(err).Message())
}

// Test AssignAgentToOrder
func TestAssignAgentToOrderSuccessfully(t *testing.T) {
	setup()
	defer teardown()

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return true, nil
	}

	// Replace the real HTTP client functions with mocks
	orderClient.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	addReq := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), addReq)
	assert.NoError(t, err)

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 1, OrderId: 1}
	resp, err := client.AssignAgentToOrder(context.Background(), assignReq)
	assert.NoError(t, err)
	assert.Equal(t, "Delivery agent assigned to order successfully", resp.Message)
}

// Test AssignAgentToOrder when agent not found
func TestAssignAgentToOrderWhenDeliveryAgentNotFound(t *testing.T) {
	setup()
	defer teardown()

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 999, OrderId: 1}
	_, err := client.AssignAgentToOrder(context.Background(), assignReq)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "delivery agent not found", status.Convert(err).Message())
}

// Test AssignAgentToOrder when already assigned
func TestAssignAgentToOrderWhenAlreadyAssigned(t *testing.T) {
	setup()
	defer teardown()

	addReq := &pb.AddDeliveryAgentRequest{Name: "Harish", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), addReq)
	assert.NoError(t, err)

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 1, OrderId: 1}
	_, err = client.AssignAgentToOrder(context.Background(), assignReq)
	assert.NoError(t, err)

	_, err = client.AssignAgentToOrder(context.Background(), assignReq)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "delivery agent is not available", status.Convert(err).Message())
}

// Test AssignAgentToOrder when order does not exist
func TestAssignAgentToOrderWhenOrderNotFound(t *testing.T) {
	setup()
	defer teardown()

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return false, nil
	}

	// Replace the real HTTP client functions with mocks
	orderClient.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	addReq := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), addReq)
	assert.NoError(t, err)

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 1, OrderId: 999}
	_, err = client.AssignAgentToOrder(context.Background(), assignReq)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "order does not exist", status.Convert(err).Message())
}

// Test AssignAgentToOrder when order exists but cannot be assigned
func TestAssignAgentToOrderWhenOrderExistsButCannotBeAssigned(t *testing.T) {
	setup()
	defer teardown()

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return true, fmt.Errorf("order cannot be assigned")
	}

	// Replace the real HTTP client functions with mocks
	orderClient.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	addReq := &pb.AddDeliveryAgentRequest{Name: "Ketan", City: "Hyderabad"}
	_, err := client.AddDeliveryAgent(context.Background(), addReq)
	assert.NoError(t, err)

	assignReq := &pb.AssignAgentToOrderRequest{AgentId: 1, OrderId: 10}
	_, err = client.AssignAgentToOrder(context.Background(), assignReq)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
	assert.Equal(t, "order cannot be assigned", status.Convert(err).Message())
}
