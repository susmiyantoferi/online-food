package entity

import (
	"time"

	"gorm.io/gorm"
)

type CartMenu struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	CartID    uint           `gorm:"notnull;uniqueIndex:idx_cart_menu"`
	Cart      Cart           `gorm:"foreignKey:CartID;references:ID;onDelete:RESTRICT"`
	MenuID    uint           `gorm:"notnull;uniqueIndex:idx_cart_menu"`
	Menu      Menu           `gorm:"foreignKey:MenuID;references:ID;onDelete:RESTRICT"`
	UnitPrice float64        `gorm:"notnull"`
	Qty       int            `gorm:"notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
