package handler

import (
	"net/http"
	"online-food/dto"
	"online-food/service"
	"online-food/utils/handling"
	"online-food/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CartHandler interface {
	CreateCart(ctx *gin.Context)
	UpdateCart(ctx *gin.Context)
	GetCartByUserID(ctx *gin.Context)
	GetCartByID(ctx *gin.Context)
	GetAllCarts(ctx *gin.Context)
	CheckoutCart(ctx *gin.Context)
}

type cartHandlerImpl struct {
	CartService service.CartService
}

func NewCartHandlerImpl(cartService service.CartService) CartHandler {
	return &cartHandlerImpl{
		CartService: cartService,
	}
}

func (c *cartHandlerImpl) CreateCart(ctx *gin.Context) {
	req := dto.CartCreateReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	userClaims, _ := ctx.Get("user")
	user := userClaims.(*dto.TokenClaim)

	req.UserID = user.UserID

	result, err := c.CartService.CreateCart(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusCreated, "Created", "cart created successfully", result)

}

func (c *cartHandlerImpl) UpdateCart(ctx *gin.Context) {
	req := dto.CartUpdateReq{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	cartId := ctx.Param("cartId")
	id, err := strconv.Atoi(cartId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	req.CardID = uint(id)

	result, err := c.CartService.UpdateCart(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "cart updated successfully", result)

}

func (c *cartHandlerImpl) GetCartByUserID(ctx *gin.Context) {
	userClaims, _ := ctx.Get("user")
	user := userClaims.(*dto.TokenClaim)

	result, err := c.CartService.GetCartByUserID(ctx.Request.Context(), user.UserID)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "get cart by user id successfully", result)
}

func (c *cartHandlerImpl) GetCartByID(ctx *gin.Context) {
	cartId := ctx.Param("cartId")
	id, err := strconv.Atoi(cartId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	result, err := c.CartService.GetCartByID(ctx.Request.Context(), uint(id))
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "get cart by id successfully", result)
}

func (c *cartHandlerImpl) GetAllCarts(ctx *gin.Context) {
	results, err := c.CartService.GetAllCarts(ctx.Request.Context())
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "get all carts successfully", results)
}

func (c *cartHandlerImpl) CheckoutCart(ctx *gin.Context) {
	cartId := ctx.Param("cartId")
	id, err := strconv.Atoi(cartId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	userClaims, _ := ctx.Get("user")
	user := userClaims.(*dto.TokenClaim)

	result, err := c.CartService.CheckoutCart(ctx.Request.Context(), uint(id), user.UserID)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "checkout cart successfully", result)
}
