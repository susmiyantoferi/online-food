package routes

import (
	"online-food/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(UserHandler handler.UserHandler, MenuHandler handler.MenuHandler) *gin.Engine {

	router := gin.Default()
	UserRouter(router, UserHandler)
	MenuRouter(router, MenuHandler)

	return router
}
