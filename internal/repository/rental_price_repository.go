package repository

import (
	"context"
	"errors"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
)

type RentalPriceRepository struct {
	db *gorm.DB
}

func NewRentalPriceRepository(db *gorm.DB) *RentalPriceRepository {
	return &RentalPriceRepository{db: db}
}

func (r *RentalPriceRepository) WithTx(tx *gorm.DB) *RentalPriceRepository {
	return &RentalPriceRepository{db: tx}
}

func (r *RentalPriceRepository) GetAll(ctx context.Context, query dto.RentalPriceQuery) ([]model.RentalPrice, int64, error) {
	var rentalPrices []model.RentalPrice
	var total int64

	db := r.db.WithContext(ctx).Model(&model.RentalPrice{})

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&rentalPrices).Error; err != nil {
		return nil, 0, err
	}

	return rentalPrices, total, nil
}

func (r *RentalPriceRepository) GetAllByRentalOptionID(ctx context.Context, rentalOptionID uint) ([]model.RentalPrice, error) {
	var rentalPrices []model.RentalPrice

	if err := r.db.WithContext(ctx).
		Model(&model.RentalPrice{}).
		Where("rental_option_id = ?", rentalOptionID).
		Find(&rentalPrices).Error; err != nil {
		return nil, err
	}

	return rentalPrices, nil
}

func (r *RentalPriceRepository) GetAllByRentalOptionIDIn(ctx context.Context, rentalOptionIds []uint) ([]model.RentalPrice, error) {
	var rentalPrices []model.RentalPrice

	if err := r.db.WithContext(ctx).
		Model(&model.RentalPrice{}).
		Where("rental_option_id IN ?", rentalOptionIds).
		Find(&rentalPrices).Error; err != nil {
		return nil, err
	}

	return rentalPrices, nil
}

func (r *RentalPriceRepository) GetByID(ctx context.Context, id uint) (*model.RentalPrice, error) {
	var rentalPrice model.RentalPrice

	if err := r.db.WithContext(ctx).Preload("RentalOption").First(&rentalPrice, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &rentalPrice, nil
}

func (r *RentalPriceRepository) ExistByRentalOptionIDAndSeasonAndDayType(ctx context.Context, rentalOptionID uint, season model.SeasonType, dayType model.DayType) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Where("rental_option_id = ? AND season = ? AND day_type = ?", rentalOptionID, season, dayType).
		Count(&count).Error

	return count > 0, err
}

func (r *RentalPriceRepository) Create(ctx context.Context, rentalPrice *model.RentalPrice) error {
	return r.db.WithContext(ctx).Create(rentalPrice).Error
}

func (r *RentalPriceRepository) Update(ctx context.Context, rentalPrice *model.RentalPrice) error {
	return r.db.WithContext(ctx).Save(rentalPrice).Error
}

func (r *RentalPriceRepository) Delete(ctx context.Context, rentalPrice *model.RentalPrice) error {
	return r.db.WithContext(ctx).Delete(rentalPrice).Error
}
