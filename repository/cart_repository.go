package repository

import (
	"context"
	"errors"
	"fmt"
	"online-food/entity"
	"online-food/utils/constanta"
	"online-food/utils/handling"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CartRepository interface {
	CreateCart(ctx context.Context, cart *entity.Cart) (*entity.Cart, error)
	UpdateCart(ctx context.Context, cartID, menuID, userID uint, qty int) (*entity.Cart, error)
	//DeleteCart(ctx context.Context, cartID uint) error
	GetCartByUserID(ctx context.Context, userID uint) ([]*entity.Cart, error)
	GetCartByID(ctx context.Context, cartID uint) (*entity.Cart, error)
	GetAllCarts(ctx context.Context) ([]*entity.Cart, error)
	CheckoutCart(ctx context.Context, cartID, userID uint) (*entity.Order, error)
}

type cartRepositoryImpl struct {
	Db *gorm.DB
}

func NewCartRepositoryImpl(db *gorm.DB) CartRepository {
	return &cartRepositoryImpl{
		Db: db,
	}
}

func (c *cartRepositoryImpl) CreateCart(ctx context.Context, cart *entity.Cart) (*entity.Cart, error) {
	if cart == nil {
		return nil, fmt.Errorf("cart is nil")
	}

	if len(cart.CartMenu) == 0 {
		return nil, handling.ErrEmptyItems
	}

	err := c.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		cart.Status = constanta.Uncheckout

		//create on table cart
		if err := tx.Omit("CartMenu").Create(cart).Error; err != nil {
			return fmt.Errorf("create cart: %w", err)
		}

		//create on table cart_menu
		for i := range cart.CartMenu {
			cart.CartMenu[i].CartID = cart.ID

			//check menu id exist
			var menu entity.Menu
			if err := tx.First(&menu, cart.CartMenu[i].MenuID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return handling.ErrMenuNotFound
				}
				return fmt.Errorf("find menu: %w", err)
			}

			cart.CartMenu[i].UnitPrice = menu.Price

			//pengurangan stock menu
			stock := tx.Model(&entity.Menu{}).Where("id = ? AND stock >= ?", cart.CartMenu[i].MenuID, cart.CartMenu[i].Qty).
				UpdateColumn("stock", gorm.Expr("stock - ?", cart.CartMenu[i].Qty))
			if stock.Error != nil {
				return fmt.Errorf("update stock: %w", stock.Error)
			}

			if stock.RowsAffected == 0 {
				return handling.ErrNotEnoughStock
			}
		}

		if err := tx.Create(&cart.CartMenu).Error; err != nil {
			return fmt.Errorf("create cart menu: %w", err)
		}

		pay := 0.0
		for _, v := range cart.CartMenu {
			pay += v.UnitPrice * float64(v.Qty)
		}

		if err := tx.Model(cart).Update("amount", pay).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err := c.Db.WithContext(ctx).Preload("User").Preload("CartMenu").Preload("CartMenu.Menu").
		First(cart, cart.ID).Error; err != nil {
		return nil, err
	}

	return cart, nil
}

func (c *cartRepositoryImpl) UpdateCart(ctx context.Context, cartID, menuID, userID uint, qty int) (*entity.Cart, error) {
	var result *entity.Cart

	err := c.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var cart entity.Cart
		if err := tx.Where("user_id", userID).First(&cart, cartID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return handling.ErrorIdNotFound
			}
		}

		if cart.Status == constanta.Checkout {
			return handling.ErrCheckoutCart
		}

		var cartMenu entity.CartMenu
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("cart_id = ? AND menu_id = ?", cartID, menuID).First(&cartMenu).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			//insert new menu
			if qty <= 0 {
				return fmt.Errorf("invalid qty to add: %d", qty)
			}

			//cek menu exist
			var menu entity.Menu
			if err := tx.First(&menu, menuID).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return handling.ErrMenuNotFound
				}
				return fmt.Errorf("find menu: %w", err)
			}

			//reduce stock menu
			stock := tx.Model(&entity.Menu{}).Where("id = ? AND stock >= ?", menuID, qty).
				UpdateColumn("stock", gorm.Expr("stock - ?", qty))
			if stock.Error != nil {
				return fmt.Errorf("update stock: %w", stock.Error)
			}

			if stock.RowsAffected == 0 {
				return handling.ErrNotEnoughStock
			}

			newCartMenu := entity.CartMenu{
				CartID:    cartID,
				MenuID:    menuID,
				Qty:       qty,
				UnitPrice: menu.Price,
			}

			if err := tx.Create(&newCartMenu).Error; err != nil {
				return fmt.Errorf("create cart menu: %w", err)
			}

		} else if err != nil {
			return err
		} else {

			if qty == 0 {

			} else if qty > 0 {
				stock := tx.Model(&entity.Menu{}).Where("id = ? AND stock >= ?", menuID, qty).
					UpdateColumn("stock", gorm.Expr("stock - ?", qty))
				if stock.Error != nil {
					return fmt.Errorf("update stock: %w", stock.Error)
				}

				if stock.RowsAffected == 0 {
					return handling.ErrNotEnoughStock
				}

				if err := tx.Model(&entity.CartMenu{}).Where("id = ?", cartMenu.ID).
					UpdateColumn("qty", gorm.Expr("qty + ?", qty)).Error; err != nil {
					return fmt.Errorf("increment cart menu qty: %w", err)
				}
			} else {
				remove := -qty
				if remove > cartMenu.Qty {
					return fmt.Errorf("invalid qty: remove %d > existing %d", remove, cartMenu.Qty)
				}

				if err := tx.Model(&entity.Menu{}).Where("id = ?", menuID).
					UpdateColumn("stock", gorm.Expr("stock + ?", remove)).Error; err != nil {
					return fmt.Errorf("restore stock: %w", err)
				}

				newQty := cartMenu.Qty - remove
				if newQty == 0 {
					if err := tx.Delete(&entity.CartMenu{}, cartMenu.ID).Error; err != nil {
						return fmt.Errorf("delete cart menu: %w", err)
					}
				} else {
					if err := tx.Model(&entity.CartMenu{}).Where("id = ?", cartMenu.ID).
						Update("qty", newQty).Error; err != nil {
						return fmt.Errorf("update cart menu qty: %w", err)
					}
				}
			}
		}

		if err := tx.Preload("User").Preload("CartMenu").
			Preload("CartMenu.Menu").First(&cart, cartID).Error; err != nil {
			return err
		}

		total := 0.0
		for _, v := range cart.CartMenu {
			total += v.UnitPrice * float64(v.Qty)
		}

		if err := tx.Model(&cart).Update("amount", total).Error; err != nil {
			return err
		}

		result = &cart
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// func (c *cartRepositoryImpl) DeleteCart(ctx context.Context, cartID uint) error {

// }

func (c *cartRepositoryImpl) GetCartByUserID(ctx context.Context, userID uint) ([]*entity.Cart, error) {
	var carts []*entity.Cart
	if err := c.Db.WithContext(ctx).Preload("User").Preload("CartMenu").Preload("CartMenu.Menu").
		Where("user_id = ?", userID).Find(&carts).Error; err != nil {
		return nil, err
	}

	return carts, nil
}

func (c *cartRepositoryImpl) GetCartByID(ctx context.Context, cartID uint) (*entity.Cart, error) {
	var cart entity.Cart
	if err := c.Db.WithContext(ctx).Preload("User").Preload("CartMenu").Preload("CartMenu.Menu").
		First(&cart, cartID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, err
	}

	return &cart, nil
}

func (c *cartRepositoryImpl) GetAllCarts(ctx context.Context) ([]*entity.Cart, error) {
	var carts []*entity.Cart
	if err := c.Db.WithContext(ctx).Preload("User").Preload("CartMenu").Preload("CartMenu.Menu").
		Find(&carts).Error; err != nil {
		return nil, err
	}

	return carts, nil
}

func (c *cartRepositoryImpl) CheckoutCart(ctx context.Context, cartID, userID uint) (*entity.Order, error) {
	var order entity.Order
	err := c.Db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var cart entity.Cart
		if err := tx.Preload("CartMenu").Preload("CartMenu.Menu").First(&cart, cartID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return handling.ErrorIdNotFound
			}
			return fmt.Errorf("find cart: %w", err)
		}

		if cart.Status == constanta.Checkout {
			return handling.ErrCheckoutCart
		}

		if err := tx.Model(&cart).UpdateColumn("status", constanta.Checkout).Error; err != nil {
			return fmt.Errorf("update cart status: %w", err)
		}

		order = entity.Order{
			CartID:    cartID,
			UserID:    userID,
			AmountPay: cart.Amount,
			OrderDate: time.Now().UTC(),
			Status:    constanta.Pending,
		}

		if err := tx.Create(&order).Error; err != nil {
			return fmt.Errorf("create order: %w", err)
		}

		if err := tx.Preload("User").Preload("Cart").
			Preload("Cart.CartMenu").
			Preload("Cart.CartMenu.Menu").
			First(&order, order.ID).Error; err != nil {
			return fmt.Errorf("preload order: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}
