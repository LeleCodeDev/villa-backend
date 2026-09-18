package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithTx(tx *gorm.DB) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) ExistByEmail(ctx context.Context, email string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Count(&count).Error

	return count > 0, err
}
