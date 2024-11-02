package Model

import (
	"errors"
	"gorm.io/gorm"
)

type DeliveryAgent struct {
	Id                 uint               `gorm:"primaryKey;autoIncrement"`
	Name               string             `gorm:"size:100"`
	City               string             `gorm:"size:100"`
	AvailabilityStatus AvailabilityStatus `gorm:"size:20"`
	OrderID            *int
}

func NewDeliveryAgent(db *gorm.DB, name string, city string) (*DeliveryAgent, error) {
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if city == "" {
		return nil, errors.New("city cannot be empty")
	}
	agent := DeliveryAgent{
		Name:               name,
		City:               city,
		AvailabilityStatus: AVAILABLE,
		OrderID:            nil,
	}

	if err := db.Create(&agent).Error; err != nil {
		return nil, err
	}

	return &agent, nil
}
