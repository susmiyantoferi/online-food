package dto

import (
	"online-food/entity"
	"time"
)

type UserCreateReq struct {
	Name     string `validate:"required,min=1,max=100" json:"name"`
	Email    string `validate:"required,email,min=1,max=100" json:"email"`
	Password string `validate:"required,min=8,max=255" json:"password"`
	Hp       string `validate:"required,numeric" json:"hp"`
	Address  string `validate:"required" json:"address"`
}

type UserUpdateReq struct {
	Name     *string `validate:"omitempty,min=1,max=100" json:"name,omitempty"`
	Password *string `validate:"omitempty,min=8,max=255" json:"password,omitempty"`
	Hp       *string `validate:"omitempty,numeric" json:"hp,omitempty"`
	Address  *string `validate:"omitempty,min=1" json:"address,omitempty"`
}

type UserResponse struct {
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Hp        string    `json:"hp"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserLoginReq struct {
	Email    string `validate:"required,email,min=1,max=100" json:"email"`
	Password string `validate:"required,min=8,max=100" json:"password"`
}

type UserRefreshTokenReq struct {
	TokenRefresh string `validate:"required" json:"refresh_token"`
}

func ToUserResponse(user *entity.User) *UserResponse {
	return &UserResponse{
		Name:      user.Name,
		Email:     user.Email,
		Hp:        user.Hp,
		Address:   user.Address,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
