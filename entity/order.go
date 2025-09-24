package entity

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	CartID    uint           `gorm:"notnull"`
	Cart      Cart           `gorm:"foreignKey:CartID;references:ID;onDelete:RESTRICT"`
	AmountPay float64        `gorm:"notnull"`
	OrderDate time.Time      `gorm:"notnull"`
	Status    string         `gorm:"type:enum('pending','paid','cancelled');default:'pending';notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
