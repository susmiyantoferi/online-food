package main

import (
	"log"
	"online-food/config"
	"online-food/handler"
	"online-food/repository"
	"online-food/routes"
	"online-food/service"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error load env")
	}

	database := config.Database()
	//redis:= config.RedisCLient()
	validate := validator.New()
	userRepo := repository.NewUserRepositoryImpl(database)
	userService := service.NewUserServiceImpl(userRepo, validate)
	userHandler := handler.NewUserHandlerImpl(userService)

	routes := routes.SetupRouter(userHandler)

	port := os.Getenv("APP_PORT")
	routes.Run(port)
	log.Println("server running on port " + port)
}
