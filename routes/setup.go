package routes

import (
	"online-food/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	UserHandler handler.UserHandler,
	MenuHandler handler.MenuHandler,
	CartHandler handler.CartHandler,
) *gin.Engine {

	router := gin.Default()
	UserRouter(router, UserHandler)
	MenuRouter(router, MenuHandler)
	CartRouter(router, CartHandler)

	return router
}
