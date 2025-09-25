package routes

import (
	"online-food/handler"
	"online-food/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(router *gin.Engine, UserHandler handler.UserHandler) {
	public := router.Group("/api/v1/auth")
	{
		public.POST("/login", UserHandler.Login)
		public.POST("/refresh-token", UserHandler.RefreshToken)
		public.POST("/register", UserHandler.Create)
	}

	user := router.Group("/api/v1")
	user.Use(middleware.Authentication())
	{
		customers := user.Group("/users")
		customers.Use(middleware.RoleAccessMiddleware("customer", "admin"))
		{
			customers.PUT("/me", UserHandler.Update)
			customers.GET("/me", UserHandler.Profile)
		}

		admin := user.Group("/users")
		admin.Use(middleware.RoleAccessMiddleware("admin"))
		{
			admin.DELETE("/:userId", UserHandler.Delete)
			admin.GET("/", UserHandler.FindAll)
			admin.GET("/email/:email", UserHandler.FindByEmail)
			admin.GET("/:userId", UserHandler.FindByID)
		}
	}
}
