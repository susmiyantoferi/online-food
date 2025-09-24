package routes

import (
	"online-food/handler"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, UserHandler handler.UserHandler){

	regist := router.Group("/api/v1/users")
	{
		regist.POST("/register", UserHandler.Create)
		regist.PUT("/:userId", UserHandler.Update)
		regist.DELETE("/:userId", UserHandler.Delete)
		regist.GET("/:userId", UserHandler.FindByID)
		regist.GET("/", UserHandler.FindAll)
		regist.GET("/email/:email", UserHandler.FindByEmail)
	}
}
