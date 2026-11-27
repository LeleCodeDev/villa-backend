package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type VoucherRepository struct {
	db *gorm.DB
}

func NewVoucherRepository(db *gorm.DB) *VoucherRepository {
	return &VoucherRepository{db: db}
}

func (r *VoucherRepository) WithTx(tx *gorm.DB) *VoucherRepository {
	return &VoucherRepository{db: tx}
}

func (r *VoucherRepository) GetAll(ctx context.Context, query dto.VoucherQuery) ([]model.Voucher, int64, error) {
	var vouchers []model.Voucher
	var total int64

	db := r.db.WithContext(ctx).Model(&model.Voucher{})

	if query.Search != "" {
		db = db.Where("code LIKE ?", "%"+strings.ToUpper(query.Search)+"%")
	}

	if query.IsActive != nil {
		db.Where("is_active = ?", *query.IsActive)
	}

	if err := db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	db = db.Order("created_at " + string(query.Sort))

	if !query.Unpage {
		db = db.Offset(query.GetOffset()).Limit(query.Size)
	}

	if err := db.Find(&vouchers).Error; err != nil {
		return nil, 0, err
	}

	return vouchers, total, nil
}

func (r *VoucherRepository) GetByID(ctx context.Context, id uint) (*model.Voucher, error) {
	var voucher model.Voucher

	if err := r.db.WithContext(ctx).First(&voucher, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &voucher, nil
}

func (r *VoucherRepository) GetByCode(ctx context.Context, code string) (*model.Voucher, error) {
	var voucher model.Voucher

	if err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&voucher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &voucher, nil
}

func (r *VoucherRepository) ExistByCode(ctx context.Context, code string) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&model.Voucher{}).
		Where("code = ?", code).
		Count(&count).Error

	return count > 0, err
}

func (r *VoucherRepository) IncrementUsage(ctx context.Context, id uint) (bool, error) {
	res := r.db.WithContext(ctx).Model(&model.Voucher{}).
		Where("id = ? AND (usage_limit IS NULL OR used_count < usage_limit)", id).
		UpdateColumn("used_count", gorm.Expr("used_count + 1"))
	return res.RowsAffected == 1, res.Error
}

func (r *VoucherRepository) DecrementUsage(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Voucher{}).
		Where("id = ? AND used_count > 0", id).
		UpdateColumn("used_count", gorm.Expr("used_count - 1")).Error
}

func (r *VoucherRepository) CreateIfAbsent(ctx context.Context, voucher *model.Voucher) (bool, error) {
	res := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "source_booking_id"}},
			DoNothing: true,
		}).
		Create(voucher)

	return res.RowsAffected == 1, res.Error
}

func (r *VoucherRepository) Create(ctx context.Context, voucher *model.Voucher) error {
	return r.db.WithContext(ctx).Create(voucher).Error
}

func (r *VoucherRepository) Update(ctx context.Context, voucher *model.Voucher) error {
	return r.db.WithContext(ctx).Save(voucher).Error
}

func (r *VoucherRepository) Delete(ctx context.Context, voucher *model.Voucher) error {
	return r.db.WithContext(ctx).Delete(voucher).Error
}
