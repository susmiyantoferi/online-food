package dto

import (
	"online-food/entity"
	"time"
)

type CreateMenuItem struct {
	MenuID uint `validate:"required" json:"menu_id"`
	Qty    int  `validate:"required,gt=0" json:"qty"`
}

type CartCreateReq struct {
	UserID   uint             `validate:"required" json:"user_id"`
	CartMenu []CreateMenuItem `validate:"required" json:"cart_menu"`
}

type CartUpdateReq struct {
	UserID uint `validate:"required" json:"user_id"`
	CardID uint `validate:"required" json:"card_id"`
	MenuID uint `validate:"required" json:"menu_id"`
	Qty    int  `validate:"required" json:"qty"`
}

type MenuDetails struct {
	MenuID    uint    `json:"menu_id"`
	Name      string  `json:"name"`
	Qty       int     `json:"qty"`
	UnitPrice float64 `json:"unit_price"`
}

type UserDetails struct {
	Name    string `json:"name"`
	Hp      string `json:"hp"`
	Address string `json:"address"`
}

type CartResponse struct {
	CartID    uint          `json:"cart_id"`
	User      UserDetails   `json:"user"`
	Amount    float64       `json:"amount"`
	Status    string        `json:"status"`
	Menus     []MenuDetails `json:"menus"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

func ToCartResponse(cart *entity.Cart) *CartResponse {
	menus := make([]MenuDetails, 0, len(cart.CartMenu))
	for _, v := range cart.CartMenu {
		menus = append(menus, MenuDetails{
			MenuID:    v.Menu.ID,
			Name:      v.Menu.Name,
			Qty:       v.Qty,
			UnitPrice: v.UnitPrice,
		})
	}
	return &CartResponse{
		CartID: cart.ID,
		User: UserDetails{
			Name:    cart.User.Name,
			Hp:      cart.User.Hp,
			Address: cart.User.Address,
		},
		Amount:    cart.Amount,
		Status:    cart.Status,
		Menus:     menus,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
	}
}

type OrderResponse struct {
	OrderID   uint          `json:"order_id"`
	OrderDate time.Time     `json:"order_date"`
	User      UserDetails   `json:"user"`
	AmountPay float64       `json:"amount_pay"`
	Menus     []MenuDetails `json:"menus"`
	Status    string        `json:"status"`
}

func ToOrderResponse(order *entity.Order) *OrderResponse {
	menus := make([]MenuDetails, 0, len(order.Cart.CartMenu))
	for _, v := range order.Cart.CartMenu {
		menus = append(menus, MenuDetails{
			MenuID:    v.Menu.ID,
			Name:      v.Menu.Name,
			Qty:       v.Qty,
			UnitPrice: v.UnitPrice,
		})
	}
	return &OrderResponse{
		OrderID:   order.ID,
		OrderDate: order.OrderDate,
		User: UserDetails{
			Name:    order.User.Name,
			Hp:      order.User.Hp,
			Address: order.User.Address,
		},
		AmountPay: order.AmountPay,
		Menus:     menus,
		Status:    order.Status,
	}
}
