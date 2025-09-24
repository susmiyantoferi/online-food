package routes

import (
	"online-food/handler"

	"github.com/gin-gonic/gin"
)

func MenuRouter(router *gin.Engine, MenuHandler handler.MenuHandler) {
	menu := router.Group("/api/v1/menus")
	{
		menu.POST("/", MenuHandler.Create)
		menu.PUT("/:menuId", MenuHandler.Update)
		menu.DELETE("/:menuId", MenuHandler.Delete)
		menu.GET("/:menuId", MenuHandler.FindByID)
		menu.GET("/", MenuHandler.FindAll)
	}
}
