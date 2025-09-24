package entity

import (
	"time"

	"gorm.io/gorm"
)

type Menu struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"size:255;notnull"`
	Stock       int            `gorm:"notnull"`
	Price       float64        `gorm:"notnull"`
	Category    string         `gorm:"type:enum('makanan','minuman');notnull"`
	Description string         `gorm:"size:255"`
	CreatedAt   time.Time      `gorm:"notnull"`
	UpdatedAt   time.Time      `gorm:"notnull"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
