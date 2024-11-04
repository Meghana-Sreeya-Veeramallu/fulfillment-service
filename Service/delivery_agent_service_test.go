package Service

import (
	client "Fulfillment/Client"
	"Fulfillment/Model"
	"Fulfillment/Repository"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

// setupTestDB function to initialize a new in-memory database for testing.
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	err := db.AutoMigrate(&Model.DeliveryAgent{})
	if err != nil {
		return nil
	}
	return db
}

// setupService function to initialize a new DeliveryAgentService with the provided database.
func setupService(db *gorm.DB) *DeliveryAgentService {
	repo := Repository.NewDeliveryAgentRepository(db)
	return NewDeliveryAgentService(repo)
}

// Test AddDeliveryAgent
func TestAddDeliveryAgentSuccessfully(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	name := "Ketan"
	city := "Hyderabad"

	agent, err := service.AddDeliveryAgent(name, city)

	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, name, agent.Name)
	assert.Equal(t, city, agent.City)
	assert.Equal(t, Model.AVAILABLE, agent.AvailabilityStatus)
	assert.Nil(t, agent.OrderID)
}

// Test AddDeliveryAgent with empty name
func TestAddDeliveryAgentWithEmptyName(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	city := "Hyderabad"

	agent, err := service.AddDeliveryAgent("", city)

	assert.Error(t, err)
	assert.Nil(t, agent)
	assert.Equal(t, "name cannot be empty", err.Error())
}

// Test AddDeliveryAgent with empty city
func TestAddDeliveryAgentWithEmptyCity(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	name := "Ketan"

	agent, err := service.AddDeliveryAgent(name, "")

	assert.Error(t, err)
	assert.Nil(t, agent)
	assert.Equal(t, "city cannot be empty", err.Error())
}

// Test AddDeliveryAgent with empty name and city
func TestAssignAgentToOrderSuccessfully(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	name := "Ketan"
	city := "Hyderabad"

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return true, nil
	}

	// Replace the real HTTP client functions with mocks
	client.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	agent, _ := service.AddDeliveryAgent(name, city)

	orderID := 1
	err := service.AssignAgentToOrder(agent.Id, orderID)

	assert.NoError(t, err)

	var updatedAgent Model.DeliveryAgent
	db.First(&updatedAgent, agent.Id)

	assert.Equal(t, orderID, *updatedAgent.OrderID)
	assert.Equal(t, Model.UNAVAILABLE, updatedAgent.AvailabilityStatus)
}

// Test AssignAgentToOrder when agent not found
func TestAssignAgentToOrderWhenDeliveryAgentNotFound(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)

	err := service.AssignAgentToOrder(999, 1)

	assert.Error(t, err)
	assert.Equal(t, "delivery agent not found", err.Error())
}

// Test AssignAgentToOrder when agent is not available
func TestAssignAgentToOrderWhenAlreadyAssigned(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	name := "Ketan"
	city := "Hyderabad"

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return true, nil
	}

	// Replace the real HTTP client functions with mocks
	client.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	agent, _ := service.AddDeliveryAgent(name, city)

	orderID1 := 1
	err := service.AssignAgentToOrder(agent.Id, orderID1)
	assert.NoError(t, err)

	orderID2 := 2
	err = service.AssignAgentToOrder(agent.Id, orderID2)

	assert.Error(t, err)
	assert.Equal(t, "delivery agent is not available", err.Error())
}

// Test AssignAgentToOrder when order is not found
func TestAssignAgentToOrderWhenOrderNotFound(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	name := "Ketan"
	city := "Hyderabad"

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return false, nil
	}

	// Replace the real HTTP client functions with mocks
	client.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	agent, _ := service.AddDeliveryAgent(name, city)

	err := service.AssignAgentToOrder(agent.Id, 123)

	assert.Error(t, err)
	assert.Equal(t, "order does not exist", err.Error())
}

// Test AssignAgentToOrder when order exists but cannot be assigned
func TestAssignAgentToOrderWhenOrderExistsButCannotBeAssigned(t *testing.T) {
	db := setupTestDB()
	service := setupService(db)
	name := "Ketan"
	city := "Hyderabad"

	// Mock CheckAndUpdateOrderStatus to simulate a successful check
	var mockCheckAndUpdateOrderStatus = func(orderID int) (bool, error) {
		return true, fmt.Errorf("order cannot be assigned")
	}

	// Replace the real HTTP client functions with mocks
	client.CheckAndUpdateOrderStatus = mockCheckAndUpdateOrderStatus

	agent, _ := service.AddDeliveryAgent(name, city)

	err := service.AssignAgentToOrder(agent.Id, 123)

	assert.Error(t, err)
	assert.Equal(t, "order cannot be assigned", err.Error())
}
