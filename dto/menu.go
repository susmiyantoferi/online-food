package dto

import (
	"online-food/entity"
	"time"
)

type MenuCreateReq struct {
	Name        string  `validate:"required,min=1,max=100" json:"name"`
	Stock       int     `validate:"required,gt=0" json:"stock"`
	Price       float64 `validate:"required" json:"price"`
	Category    string  `validate:"required,oneof=makanan minuman" json:"category"`
	Description string  `validate:"required" json:"description"`
}

type MenuUpdateReq struct {
	ID          uint     `validate:"required"`
	Name        *string  `validate:"omitempty,min=1,max=100" json:"name,omitempty"`
	Stock       *int     `validate:"omitempty,gt=0" json:"stock,omitempty"`
	Price       *float64 `validate:"omitempty" json:"price"`
	Category    *string  `validate:"omitempty,oneof=makanan minuman" json:"category,omitempty"`
	Description *string  `validate:"omitempty" json:"description,omitempty"`
}

type MenuResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Stock       int       `json:"stock"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToMenuResponse(menu *entity.Menu) *MenuResponse {
	return &MenuResponse{
		ID:          menu.ID,
		Name:        menu.Name,
		Stock:       menu.Stock,
		Price:       menu.Price,
		Category:    menu.Category,
		Description: menu.Description,
		CreatedAt:   menu.CreatedAt,
		UpdatedAt:   menu.UpdatedAt,
	}
}
