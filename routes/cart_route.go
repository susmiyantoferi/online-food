package routes

import (
	"online-food/handler"
	"online-food/middleware"

	"github.com/gin-gonic/gin"
)

func CartRouter(router *gin.Engine, CartHandler handler.CartHandler) {
	cart := router.Group("/api/v1/")
	cart.Use(middleware.Authentication())
	{
		cust := cart.Group("/carts")
		cust.Use(middleware.RoleAccessMiddleware("customer", "admin"))
		{
			cust.POST("/", CartHandler.CreateCart)
			cust.PUT("/:cartId", CartHandler.UpdateCart)
			cust.GET("/users", CartHandler.GetCartByUserID)
			cust.POST("/checkout/:cartId", CartHandler.CheckoutCart)
		}

		admin := cart.Group("/carts")
		admin.Use(middleware.RoleAccessMiddleware("admin"))
		{
			cust.GET("/:cartId", CartHandler.GetCartByID)
			cust.GET("/", CartHandler.GetAllCarts)
		}
	}
}
