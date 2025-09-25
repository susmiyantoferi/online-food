package routes

import (
	"online-food/handler"
	"online-food/middleware"

	"github.com/gin-gonic/gin"
)

func MenuRouter(router *gin.Engine, MenuHandler handler.MenuHandler) {
	menu := router.Group("/api/v1")
	menu.Use(middleware.Authentication())
	{
		admin := menu.Group("/menus")
		admin.Use(middleware.RoleAccessMiddleware("admin"))
		{
			admin.POST("/", MenuHandler.Create)
			admin.PUT("/:menuId", MenuHandler.Update)
			admin.DELETE("/:menuId", MenuHandler.Delete)
			admin.GET("/:menuId", MenuHandler.FindByID)

		}

		cust := menu.Group("/menus")
		cust.Use(middleware.RoleAccessMiddleware("customer", "admin"))
		{
			cust.GET("/", MenuHandler.FindAll)
		}
	}

}
