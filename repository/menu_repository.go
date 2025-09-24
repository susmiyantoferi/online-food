package repository

import (
	"context"
	"errors"
	"online-food/entity"
	"online-food/utils/handling"

	"gorm.io/gorm"
)

type MenuRepository interface {
	Create(ctx context.Context, menu *entity.Menu) (*entity.Menu, error)
	Update(ctx context.Context,  menu *entity.Menu) (*entity.Menu, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*entity.Menu, error)
	FindAll(ctx context.Context) ([]*entity.Menu, error)
}

type menuRepositoryImpl struct {
	Db *gorm.DB
}

func NewMenuRepositoryImpl(db *gorm.DB) *menuRepositoryImpl {
	return &menuRepositoryImpl{
		Db: db,
	}
}

func (m *menuRepositoryImpl) Create(ctx context.Context, menu *entity.Menu) (*entity.Menu, error) {
	if err := m.Db.WithContext(ctx).Create(menu).Error; err != nil {
		return nil, err
	}

	return menu, nil
}

func (m *menuRepositoryImpl) Update(ctx context.Context, menu *entity.Menu) (*entity.Menu, error) {
	data := entity.Menu{
		Name:        menu.Name,
		Stock:       menu.Stock,
		Price:       menu.Price,
		Category:    menu.Category,
		Description: menu.Description,
	}

	if err := m.Db.WithContext(ctx).First(menu, menu.ID).Updates(data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, err
	}

	return menu, nil
}

func (m *menuRepositoryImpl) Delete(ctx context.Context, id uint) error {
	result := m.Db.WithContext(ctx).Delete(&entity.Menu{}, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return handling.ErrorIdNotFound
	}

	return nil
}

func (m *menuRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.Menu, error) {
	menus := entity.Menu{}
	if err := m.Db.WithContext(ctx).First(&menus, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, handling.ErrorIdNotFound
		}
		return nil, err
	}

	return &menus, nil
}

func (m *menuRepositoryImpl) FindAll(ctx context.Context) ([]*entity.Menu, error) {
	var menus []*entity.Menu
	if err := m.Db.WithContext(ctx).Find(&menus).Error; err != nil {
		return nil, err
	}

	return menus, nil
}
