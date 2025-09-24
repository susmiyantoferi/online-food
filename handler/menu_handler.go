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

type MenuHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindAll(ctx *gin.Context)
}

type menuHandlerImpl struct {
	MenuService service.MenuService
}

func NewMenuHandlerImpl(menuService service.MenuService) *menuHandlerImpl {
	return &menuHandlerImpl{
		MenuService: menuService,
	}
}

func (m *menuHandlerImpl) Create(ctx *gin.Context) {
	req := dto.MenuCreateReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	result, err := m.MenuService.Create(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusCreated, "Created", "menu created successfully", result)
}

func (m *menuHandlerImpl) Update(ctx *gin.Context) {
	req := dto.MenuUpdateReq{}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	menuId := ctx.Param("menuId")
	id, err := strconv.Atoi(menuId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	req.ID = uint(id)

	result, err := m.MenuService.Update(ctx.Request.Context(), &req)
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Updated", "menu updated successfully", result)
}

func (m *menuHandlerImpl) Delete(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	id, err := strconv.Atoi(menuId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	if err := m.MenuService.Delete(ctx.Request.Context(), uint(id)); err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Deleted", "menu deleted successfully", nil)
}

func (m *menuHandlerImpl) FindByID(ctx *gin.Context) {
	menuId := ctx.Param("menuId")
	id, err := strconv.Atoi(menuId)
	if err != nil {
		response.ToResponseJson(ctx, http.StatusBadRequest, "Bad Request", "invalid input type id", nil)
		return
	}

	result, err := m.MenuService.FindByID(ctx.Request.Context(), uint(id))
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "menu find id successfully", result)
}

func (m *menuHandlerImpl) FindAll(ctx *gin.Context) {
	result, err := m.MenuService.FindAll(ctx.Request.Context())
	if err != nil {
		handling.HandleError(ctx, err)
		return
	}

	response.ToResponseJson(ctx, http.StatusOK, "Success", "menu find successfully", result)
}
