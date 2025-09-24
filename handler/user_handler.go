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

type UserHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	FindByEmail(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
}

type userHandlerImpl struct {
	UserService service.UserService
}

func NewUserHandlerImpl(userService service.UserService) *userHandlerImpl {
	return &userHandlerImpl{
		UserService: userService,
	}
}

func (u *userHandlerImpl) Create(ctx *gin.Context) {
	req := dto.UserCreateReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	result, err := u.UserService.Create(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusCreated, "Created", "user created successfully", result)
}

func (u *userHandlerImpl) Update(ctx *gin.Context) {
	req := dto.UserUpdateReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	result, err := u.UserService.Update(ctx.Request.Context(), uint(id), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Updated", "user updated successfully", result)
}

func (u *userHandlerImpl) Delete(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	if err := u.UserService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Deleted", "user deleted successfully", nil)
}

func (u *userHandlerImpl) FindByID(ctx *gin.Context) {
	userId := ctx.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	result, err := u.UserService.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "find id successfully", result)
}

func (u *userHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := u.UserService.FindAll(ctx.Request.Context())
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "find all successfully", result)
}

func (u *userHandlerImpl) FindByEmail(ctx *gin.Context) {
	userEmail := ctx.Param("email")

	result, err := u.UserService.FindByEmail(ctx.Request.Context(), userEmail)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "find email successfully", result)
}

func (u *userHandlerImpl) Login(ctx *gin.Context) {
	req := dto.UserLoginReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil{
		handling.HandleError(ctx, err)
		return
	}

	result, err := u.UserService.Login(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "generate token successfully", result)
}

func (u *userHandlerImpl) RefreshToken(ctx *gin.Context) {
	req := dto.UserRefreshTokenReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil{
		handling.HandleError(ctx, err)
		return
	}

	result, err := u.UserService.RefreshToken(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "generate new token successfully", result)
}	
