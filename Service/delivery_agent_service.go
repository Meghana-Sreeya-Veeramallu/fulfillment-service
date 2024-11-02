package Service

import (
	"Fulfillment/Model"
	"errors"
	"gorm.io/gorm"
)

func AddDeliveryAgent(db *gorm.DB, name string, city string) (*Model.DeliveryAgent, error) {
	return Model.NewDeliveryAgent(db, name, city)
}

func AssignAgentToOrder(db *gorm.DB, deliveryAgentID uint, orderID int) error {
	var agent Model.DeliveryAgent

	if err := db.First(&agent, deliveryAgentID).Error; err != nil {
		return errors.New("delivery agent not found")
	}

	if agent.AvailabilityStatus != Model.AVAILABLE {
		return errors.New("delivery agent is not available")
	}

	agent.OrderID = &orderID
	agent.AvailabilityStatus = Model.UNAVAILABLE

	if err := db.Save(&agent).Error; err != nil {
		return err
	}

	return nil
}
