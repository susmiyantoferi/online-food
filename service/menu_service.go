package service

import (
	"context"
	"errors"
	"fmt"
	"online-food/dto"
	"online-food/entity"
	"online-food/repository"
	"online-food/utils/handling"

	"github.com/go-playground/validator/v10"
)

type MenuService interface {
	Create(ctx context.Context, req *dto.MenuCreateReq) (*dto.MenuResponse, error)
	Update(ctx context.Context, req *dto.MenuUpdateReq) (*dto.MenuResponse, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*dto.MenuResponse, error)
	FindAll(ctx context.Context) ([]*dto.MenuResponse, error)
}

type menuServiceImpl struct {
	MenuRepo repository.MenuRepository
	Validate *validator.Validate
}

func NewMenuServiceImpl(menuRepo repository.MenuRepository, validate *validator.Validate) *menuServiceImpl {
	return &menuServiceImpl{
		MenuRepo: menuRepo,
		Validate: validate,
	}
}

func (m *menuServiceImpl) Create(ctx context.Context, req *dto.MenuCreateReq) (*dto.MenuResponse, error) {
	if err := m.Validate.Struct(req); err != nil {
		return nil, handling.ErrorValidation
	}

	menu := &entity.Menu{
		Name:        req.Name,
		Stock:       req.Stock,
		Price:       req.Price,
		Category:    req.Category,
		Description: req.Description,
	}

	result, err := m.MenuRepo.Create(ctx, menu)
	if err != nil {
		return nil, fmt.Errorf("menu service: create: %w", err)
	}

	response := dto.ToMenuResponse(result)
	return response, nil
}

func (m *menuServiceImpl) Update(ctx context.Context, req *dto.MenuUpdateReq) (*dto.MenuResponse, error) {
	if err := m.Validate.Struct(req); err != nil {
		return nil, handling.ErrorValidation
	}

	menu, err := m.MenuRepo.FindByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, fmt.Errorf("menu service: update: find id: %w", err)
	}

	if req.Name != nil {
		menu.Name = *req.Name
	}

	if req.Stock != nil {
		menu.Stock = *req.Stock
	}

	if req.Price != nil {
		menu.Price = *req.Price
	}

	if req.Category != nil {
		menu.Category = *req.Category
	}

	if req.Description != nil {
		menu.Description = *req.Description
	}

	result, err := m.MenuRepo.Update(ctx, menu)
	if err != nil {
		return nil, fmt.Errorf("menu service: update: %w", err)
	}

	respone := dto.ToMenuResponse(result)
	return respone, nil
}

func (m *menuServiceImpl) Delete(ctx context.Context, id uint) error {
	if err := m.MenuRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return handling.ErrorIdNotFound
		}
		return fmt.Errorf("menu service: delete: %w", err)
	}

	return nil
}

func (m *menuServiceImpl) FindByID(ctx context.Context, id uint) (*dto.MenuResponse, error) {
	result, err := m.MenuRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, fmt.Errorf("menu service: find id: %w", err)
	}

	response := dto.ToMenuResponse(result)
	return response, nil
}

func (m *menuServiceImpl) FindAll(ctx context.Context) ([]*dto.MenuResponse, error) {
	result, err := m.MenuRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("menu service: find all: %w", err)
	}

	var responses []*dto.MenuResponse
	for _, v := range result {
		responses = append(responses, dto.ToMenuResponse(v))
	}
	return responses, nil
}
