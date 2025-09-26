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

	//user
	userRepo := repository.NewUserRepositoryImpl(database)
	userService := service.NewUserServiceImpl(userRepo, validate)
	userHandler := handler.NewUserHandlerImpl(userService)

	//menu
	menuRepo := repository.NewMenuRepositoryImpl(database)
	menuService := service.NewMenuServiceImpl(menuRepo, validate)
	menuHandler := handler.NewMenuHandlerImpl(menuService)

	//cart
	cartRepo := repository.NewCartRepositoryImpl(database)
	cartService := service.NewCartServiceImpl(cartRepo, validate)
	cartHandler := handler.NewCartHandlerImpl(cartService)

	routes := routes.SetupRouter(userHandler, menuHandler, cartHandler)

	port := os.Getenv("APP_PORT")
	routes.Run(port)
	log.Println("server running on port " + port)
}
