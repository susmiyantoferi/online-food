package config

import (
	"fmt"
	"log"
	"online-food/entity"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Database() *gorm.DB {

	user := os.Getenv("DB_USERNAME")
	pass := os.Getenv("DB_PWD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)

	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{})
	if err != nil {
		log.Fatalf("database: %v", err)
	}

	err = db.AutoMigrate(
		&entity.User{},
		&entity.Menu{},
		&entity.Cart{},
		&entity.CartMenu{},
		&entity.Order{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}

	return db
}
