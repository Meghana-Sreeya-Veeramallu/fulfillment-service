package DeliveryAgentService

import (
	"Fulfillment/DeliveryAgent"
	"errors"
	"gorm.io/gorm"
)

func AddDeliveryAgent(db *gorm.DB, name string, city string) (*DeliveryAgent.DeliveryAgent, error) {
	return DeliveryAgent.NewDeliveryAgent(db, name, city)
}

func AssignAgentToOrder(db *gorm.DB, deliveryAgentID uint, orderID int) error {
	var agent DeliveryAgent.DeliveryAgent

	if err := db.First(&agent, deliveryAgentID).Error; err != nil {
		return errors.New("delivery agent not found")
	}

	if agent.AvailabilityStatus != DeliveryAgent.AVAILABLE {
		return errors.New("delivery agent is not available")
	}

	agent.OrderID = &orderID
	agent.AvailabilityStatus = DeliveryAgent.UNAVAILABLE

	if err := db.Save(&agent).Error; err != nil {
		return err
	}

	return nil
}
