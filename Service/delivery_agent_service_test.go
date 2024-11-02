package Service

import (
	"Fulfillment/Model"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

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

func TestCreateDeliveryAgentSuccessfully(t *testing.T) {
	db := setupTestDB()
	name := "Ketan"
	city := "Hyderabad"

	agent, err := AddDeliveryAgent(db, name, city)

	assert.NoError(t, err)
	assert.NotNil(t, agent)
	assert.Equal(t, name, agent.Name)
	assert.Equal(t, city, agent.City)
	assert.Equal(t, Model.AVAILABLE, agent.AvailabilityStatus)
	assert.Nil(t, agent.OrderID)
}

func TestCreateDeliveryAgentWithEmptyName(t *testing.T) {
	db := setupTestDB()
	city := "Hyderabad"

	agent, err := AddDeliveryAgent(db, "", city)

	assert.Error(t, err)
	assert.Nil(t, agent)
	assert.Equal(t, "name cannot be empty", err.Error())
}

func TestCreateDeliveryAgentWithEmptyCity(t *testing.T) {
	db := setupTestDB()
	name := "Ketan"

	agent, err := AddDeliveryAgent(db, name, "")

	assert.Error(t, err)
	assert.Nil(t, agent)
	assert.Equal(t, "city cannot be empty", err.Error())
}

func TestAssignAgentToOrderSuccessfully(t *testing.T) {
	db := setupTestDB()
	name := "Ketan"
	city := "Hyderabad"

	agent, _ := AddDeliveryAgent(db, name, city)

	orderID := 123
	err := AssignAgentToOrder(db, agent.Id, orderID)

	assert.NoError(t, err)

	var updatedAgent Model.DeliveryAgent
	db.First(&updatedAgent, agent.Id)

	assert.Equal(t, orderID, *updatedAgent.OrderID)
	assert.Equal(t, Model.UNAVAILABLE, updatedAgent.AvailabilityStatus)
}

func TestAssignAgentToOrderWhenDeliveryAgentNotFound(t *testing.T) {
	db := setupTestDB()

	err := AssignAgentToOrder(db, 999, 123)

	assert.Error(t, err)
	assert.Equal(t, "delivery agent not found", err.Error())
}

func TestAssignAgentToOrderWhenAlreadyAssigned(t *testing.T) {
	db := setupTestDB()
	name := "Ketan"
	city := "Hyderabad"

	agent, _ := AddDeliveryAgent(db, name, city)

	orderID1 := 123
	err := AssignAgentToOrder(db, agent.Id, orderID1)
	assert.NoError(t, err)

	orderID2 := 456
	err = AssignAgentToOrder(db, agent.Id, orderID2)

	assert.Error(t, err)
	assert.Equal(t, "delivery agent is not available", err.Error())
}
