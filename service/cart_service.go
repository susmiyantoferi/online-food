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

type CartService interface {
	CreateCart(ctx context.Context, req *dto.CartCreateReq) (*dto.CartResponse, error)
	UpdateCart(ctx context.Context, req *dto.CartUpdateReq) (*dto.CartResponse, error)
	GetCartByUserID(ctx context.Context, userID uint) ([]*dto.CartResponse, error)
	GetCartByID(ctx context.Context, cartID uint) (*dto.CartResponse, error)
	GetAllCarts(ctx context.Context) ([]*dto.CartResponse, error)
	CheckoutCart(ctx context.Context, cartID, userID uint) (*dto.OrderResponse, error)
}

type cartServiceImpl struct {
	CartRepo repository.CartRepository
	Validate *validator.Validate
}

func NewCartServiceImpl(cartRepo repository.CartRepository, validate *validator.Validate) CartService {
	return &cartServiceImpl{
		CartRepo: cartRepo,
		Validate: validate,
	}
}

func (c *cartServiceImpl) CreateCart(ctx context.Context, req *dto.CartCreateReq) (*dto.CartResponse, error) {
	if err := c.Validate.Struct(req); err != nil {
		return nil, handling.ErrorValidation
	}

	menus := entity.Cart{
		UserID:   req.UserID,
		CartMenu: make([]entity.CartMenu, len(req.CartMenu)),
	}

	for i, v := range req.CartMenu {
		menus.CartMenu[i] = entity.CartMenu{
			MenuID: v.MenuID,
			Qty:    v.Qty,
		}
	}

	result, err := c.CartRepo.CreateCart(ctx, &menus)
	if err != nil {
		if errors.Is(err, handling.ErrEmptyItems) {
			return nil, handling.ErrEmptyItems
		}

		if errors.Is(err, handling.ErrMenuNotFound) {
			return nil, handling.ErrMenuNotFound
		}

		if errors.Is(err, handling.ErrNotEnoughStock) {
			return nil, handling.ErrNotEnoughStock
		}

		return nil, fmt.Errorf("create service: create cart: %w", err)
	}

	response := dto.ToCartResponse(result)
	return response, nil
}

func (c *cartServiceImpl) UpdateCart(ctx context.Context, req *dto.CartUpdateReq) (*dto.CartResponse, error) {
	if err := c.Validate.Struct(req); err != nil {
		return nil, handling.ErrorValidation
	}

	result, err := c.CartRepo.UpdateCart(ctx, req.CardID, req.MenuID, req.Qty)
	if err != nil {
		if errors.Is(err, handling.ErrMenuNotFound) {
			return nil, handling.ErrMenuNotFound
		}

		if errors.Is(err, handling.ErrNotEnoughStock) {
			return nil, handling.ErrNotEnoughStock
		}

		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}

		if errors.Is(err, handling.ErrCheckoutCart) {
			return nil, handling.ErrCheckoutCart
		}

		return nil, fmt.Errorf("update service: update cart: %w", err)
	}

	response := dto.ToCartResponse(result)
	return response, nil

}

func (c *cartServiceImpl) GetCartByUserID(ctx context.Context, userID uint) ([]*dto.CartResponse, error) {
	results, err := c.CartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, fmt.Errorf("get service: get cart by user id: %w", err)
	}

	responses := make([]*dto.CartResponse, 0, len(results))
	for _, v := range results {
		responses = append(responses, dto.ToCartResponse(v))
	}

	return responses, nil
}

func (c *cartServiceImpl) GetCartByID(ctx context.Context, cartID uint) (*dto.CartResponse, error) {
	result, err := c.CartRepo.GetCartByID(ctx, cartID)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, fmt.Errorf("get service: get cart by id: %w", err)
	}

	response := dto.ToCartResponse(result)
	return response, nil
}

func (c *cartServiceImpl) GetAllCarts(ctx context.Context) ([]*dto.CartResponse, error) {
	results, err := c.CartRepo.GetAllCarts(ctx)
	if err != nil {
		return nil, fmt.Errorf("get service: get all carts: %w", err)
	}

	responses := make([]*dto.CartResponse, 0, len(results))
	for _, v := range results {
		responses = append(responses, dto.ToCartResponse(v))
	}

	return responses, nil
}

func (c *cartServiceImpl) CheckoutCart(ctx context.Context, cartID, userID uint) (*dto.OrderResponse, error) {
	result, err := c.CartRepo.CheckoutCart(ctx, cartID, userID)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}

		if errors.Is(err, handling.ErrCheckoutCart) {
			return nil, handling.ErrCheckoutCart
		}
		return nil, fmt.Errorf("checkout service: checkout cart: %w", err)
	}

	response := dto.ToOrderResponse(result)
	return response, nil
}
