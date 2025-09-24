package service

import (
	"context"
	"errors"
	"fmt"
	"online-food/dto"
	"online-food/entity"
	"online-food/repository"
	"online-food/utils/constanta"
	"online-food/utils/handling"
	"online-food/utils/hashing"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	Create(ctx context.Context, req *dto.UserCreateReq) (*dto.UserResponse, error)
	Update(ctx context.Context, id uint, req *dto.UserUpdateReq) (*dto.UserResponse, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*dto.UserResponse, error)
	FindAll(ctx context.Context) ([]*dto.UserResponse, error)
	FindByEmail(ctx context.Context, email string) (*dto.UserResponse, error)
}

type userServiceImpl struct {
	UserRepo repository.UserRepository
	Validate *validator.Validate
}

func NewUserServiceImpl(userRepo repository.UserRepository, validate *validator.Validate) *userServiceImpl {
	return &userServiceImpl{
		UserRepo: userRepo,
		Validate: validate,
	}
}

func (u *userServiceImpl) Create(ctx context.Context, req *dto.UserCreateReq) (*dto.UserResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	if err := u.Validate.Struct(req); err != nil {
		return nil, handling.ErrorValidation
	}

	pass, err := hashing.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing: %w", err)
	}

	user := entity.User{
		Name:     req.Name,
		Email:    email,
		Password: pass,
		Role:     constanta.Customer,
		Hp:       req.Hp,
		Address:  req.Address,
	}

	result, err := u.UserRepo.Create(ctx, &user)
	if err != nil {
		if errors.Is(err, handling.ErrorEmailExist) {
			return nil, handling.ErrorEmailExist
		}
		return nil, fmt.Errorf("user service: create: %w", err)
	}

	response := dto.ToUserResponse(result)

	return response, nil
}

func (u *userServiceImpl) Update(ctx context.Context, id uint, req *dto.UserUpdateReq) (*dto.UserResponse, error) {
	if err := u.Validate.Struct(req); err != nil {
		return nil, handling.ErrorValidation
	}

	user := entity.User{}

	if req.Name != nil {
		user.Name = *req.Name
	}

	if req.Password != nil {
		pass, err := hashing.HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("hashing: %w", err)
		}
		user.Password = pass
	}

	if req.Hp != nil {
		user.Hp = *req.Hp
	}

	if req.Address != nil {
		user.Address = *req.Address
	}

	result, err := u.UserRepo.Update(ctx, id, &user)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, fmt.Errorf("user service: update: %w", err)
	}

	response := dto.ToUserResponse(result)

	return response, nil

}

func (u *userServiceImpl) Delete(ctx context.Context, id uint) error {
	err := u.UserRepo.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return handling.ErrorIdNotFound
		}
		return fmt.Errorf("user service: delete: %w", err)
	}

	return nil
}

func (u *userServiceImpl) FindByID(ctx context.Context, id uint) (*dto.UserResponse, error) {
	user, err := u.UserRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, handling.ErrorIdNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, fmt.Errorf("user service: find by id: %w", err)
	}

	response := dto.ToUserResponse(user)

	return response, nil
}

func (u *userServiceImpl) FindAll(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := u.UserRepo.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("user service: find all: %w", err)
	}

	var responses []*dto.UserResponse
	for _, user := range users {
		responses = append(responses, dto.ToUserResponse(user))
	}

	return responses, nil
}

func (u *userServiceImpl) FindByEmail(ctx context.Context, email string) (*dto.UserResponse, error) {
	user, err := u.UserRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, handling.ErrorEmailNotFound) {
			return nil, handling.ErrorEmailNotFound
		}
		return nil, fmt.Errorf("user service: find by email: %w", err)
	}

	response := dto.ToUserResponse(user)

	return response, nil
}
