package repository

import (
	"context"
	"errors"
	"online-food/entity"
	"online-food/utils/handling"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	Update(ctx context.Context, id uint, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindAll(ctx context.Context) ([]*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userRepositoryImpl struct {
	Db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *userRepositoryImpl {
	return &userRepositoryImpl{
		Db: db,
	}
}

func (u *userRepositoryImpl) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := u.Db.WithContext(ctx).Create(user).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, handling.ErrorEmailExist
		}
		return nil, err
	}

	return user, nil
}

func (u *userRepositoryImpl) Update(ctx context.Context, id uint, user *entity.User) (*entity.User, error) {
	var dataUser entity.User
	if err := u.Db.WithContext(ctx).First(&dataUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, handling.ErrorIdNotFound
		}

		return nil, err
	}

	update := map[string]interface{}{}

	if user.Name != "" {
		update["name"] = user.Name
	}

	if user.Password != "" {
		update["password"] = user.Password
	}

	if user.Hp != "" {
		update["hp"] = user.Hp
	}

	if user.Address != "" {
		update["address"] = user.Address
	}

	result := u.Db.WithContext(ctx).Model(&dataUser).Updates(update)
	if result.Error != nil {
		return nil, result.Error
	}

	return &dataUser, nil
}

func (u *userRepositoryImpl) Delete(ctx context.Context, id uint) error {
	var user entity.User
	result := u.Db.WithContext(ctx).Delete(&user, id)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return handling.ErrorIdNotFound
	}

	return nil
}

func (u *userRepositoryImpl) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	if err := u.Db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, handling.ErrorIdNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (u *userRepositoryImpl) FindAll(ctx context.Context) ([]*entity.User, error) {
	var user []*entity.User
	if err := u.Db.WithContext(ctx).Find(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := u.Db.WithContext(ctx).Where("email = ?", email).Take(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, handling.ErrorEmailNotFound
		}

		return nil, err
	}

	return &user, nil
}
