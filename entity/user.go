package entity

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	Name      string         `gorm:"size:100;notnull"`
	Email     string         `gorm:"size:100;unique;notnull"`
	Password  string         `gorm:"size:255;notnull"`
	Role      string         `gorm:"type:enum('customer','admin');default:'customer';notnull"`
	Hp        string         `gorm:"notnull"`
	Address   string         `gorm:"notnull"`
	CreatedAt time.Time      `gorm:"notnull"`
	UpdatedAt time.Time      `gorm:"notnull"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
