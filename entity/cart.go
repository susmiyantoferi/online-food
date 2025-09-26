package entity

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	UserID    uint           `gorm:"notnull"`
	User      User           `gorm:"foreignKey:UserID;references:ID"`
	CartMenu  []CartMenu     `gorm:"foreignKey:CartID"`
	Amount    float64        `gorm:"default:null"`
	Status    string         `gorm:"type:enum('uncheckout','checkout');default:'uncheckout';notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
