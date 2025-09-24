package routes

import (
	"online-food/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(UserHandler handler.UserHandler) *gin.Engine {

	router := gin.Default()
	UserRouter(router, UserHandler)

	return router
}
