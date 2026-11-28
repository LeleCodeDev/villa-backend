package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/config"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"gorm.io/gorm"
)

type VoucherService struct {
	txManager   *repository.TxManager
	repo        *repository.VoucherRepository
	bookingRepo *repository.BookingRepository
}

func NewVoucherService(
	txManager *repository.TxManager,
	repo *repository.VoucherRepository,
	bookingRepo *repository.BookingRepository,
) *VoucherService {
	return &VoucherService{
		txManager:   txManager,
		repo:        repo,
		bookingRepo: bookingRepo,
	}
}

func (s *VoucherService) GetAll(ctx context.Context, query dto.VoucherQuery) ([]dto.VoucherResponse, int64, error) {
	vouchers, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	voucherResponses := make([]dto.VoucherResponse, 0, len(vouchers))
	for _, voucher := range vouchers {
		voucherResponses = append(voucherResponses, mapper.ToVoucherResponse(&voucher))
	}

	return voucherResponses, total, nil
}

func (s *VoucherService) GetByID(ctx context.Context, id uint) (dto.VoucherResponse, error) {
	voucher, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.VoucherResponse{}, err
	}
	if voucher == nil {
		return dto.VoucherResponse{}, appError.BadRequest(fmt.Sprintf("Voucher not found with ID: %d", id))
	}

	return mapper.ToVoucherResponse(voucher), nil
}

func (s *VoucherService) Validate(ctx context.Context, currentUser model.User, req dto.VoucherValidateRequest) (dto.VoucherValidateResponse, error) {
	voucher, err := s.repo.GetByCode(ctx, req.Code)
	if err != nil {
		return dto.VoucherValidateResponse{}, err
	}
	if voucher == nil {
		return dto.VoucherValidateResponse{}, appError.BadRequest("Voucher is not valid")
	}

	if voucher.UsageLimit != nil && voucher.UsedCount >= *voucher.UsageLimit {
		return dto.VoucherValidateResponse{}, appError.BadRequest("Voucher usage limit reached")
	}

	discount, err := calculateVoucher(voucher, req.Amount, time.Now())
	if err != nil {
		return dto.VoucherValidateResponse{}, err
	}

	if voucher.PerUserLimit != nil {
		used, err := s.bookingRepo.CountVoucherUsage(ctx, currentUser.ID, voucher.ID)
		if err != nil {
			return dto.VoucherValidateResponse{}, err
		}

		if used >= int64(*voucher.PerUserLimit) {
			return dto.VoucherValidateResponse{}, appError.BadRequest("You have reached the maximum used of this voucher")
		}
	}

	finalAmount := req.Amount - discount

	response := dto.VoucherValidateResponse{
		Code:           voucher.Code,
		DiscountAmount: discount,
		FinalAmount:    finalAmount,
	}

	return response, nil
}

func parseVoucherDates(start, exp *string) (*time.Time, *time.Time, error) {
	loc := config.Env.Location

	var startAt *time.Time
	if start != nil {
		s, err := time.ParseInLocation("2006-01-02 15:04", *start, loc)
		if err != nil {
			return nil, nil, appError.BadRequest("invalid start_at, expected YYYY-MM-DD HH:MM")
		}

		startAt = &s
	}

	var expAt *time.Time
	if exp != nil {
		e, err := time.ParseInLocation("2006-01-02 15:04", *exp, loc)
		if err != nil {
			return nil, nil, appError.BadRequest("invalid expires_at, expected YYYY-MM-DD HH:MM")
		}

		if startAt != nil && !e.After(*startAt) {
			return nil, nil, appError.BadRequest("expires_at must be after start_at")
		}

		expAt = &e
	}

	return startAt, expAt, nil
}

func (s *VoucherService) Create(ctx context.Context, req dto.VoucherRequest) (dto.VoucherResponse, error) {
	startAt, expAt, err := parseVoucherDates(req.StartAt, req.ExpiresAt)
	if err != nil {
		return dto.VoucherResponse{}, err
	}

	var createdVoucher *model.Voucher

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		if startAt != nil {
			now := time.Now().In(config.Env.Location)
			y, m, d := now.Date()
			today := time.Date(y, m, d, 0, 0, 0, 0, config.Env.Location)

			if startAt.Before(today) {
				return appError.BadRequest("start_at cannot be in the past")
			}
		}

		if req.DiscountType == model.DiscountTypePercent && req.DiscountValue > 100 {
			return appError.BadRequest("Percent discount must be between 1 and 100")
		}

		if req.DiscountType == model.DiscountTypeFixed && req.MaxDiscount != nil {
			return appError.BadRequest("max_discount only applies to percent vouchers")
		}

		code := strings.ToUpper(strings.TrimSpace(req.Code))
		exist, err := txRepo.ExistByCode(ctx, code)
		if err != nil {
			return err
		}
		if exist {
			return appError.AlreadyExist("Voucher code already exist")
		}

		voucher := mapper.ToVoucherModel(req, startAt, expAt)
		if err := txRepo.Create(ctx, voucher); err != nil {
			return err
		}

		createdVoucher = voucher

		return nil
	}); err != nil {
		return dto.VoucherResponse{}, err
	}

	return mapper.ToVoucherResponse(createdVoucher), nil
}

func (s *VoucherService) CreateRewardVoucher(ctx context.Context, booking *model.Booking) (*model.Voucher, error) {
	const voucherAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return &model.Voucher{}, err
	}

	for i := range b {
		b[i] = voucherAlphabet[int(b[i])%len(voucherAlphabet)]
	}

	code := "VILLA" + string(b)

	rewardPercent := 5
	expiresAt := time.Now().AddDate(0, 3, 0)
	usageLimit := 1
	perUserLimit := 1

	voucher := &model.Voucher{
		Code:            code,
		DiscountType:    model.DiscountTypePercent,
		DiscountValue:   int64(rewardPercent),
		UsageLimit:      &usageLimit,
		PerUserLimit:    &perUserLimit,
		ExpiresAt:       &expiresAt,
		IsActive:        true,
		SourceBookingID: &booking.ID,
	}

	created, err := s.repo.CreateIfAbsent(ctx, voucher)
	if err != nil {
		return nil, err
	}
	if !created {
		return nil, nil
	}

	return voucher, nil
}

func (s *VoucherService) Update(ctx context.Context, req dto.VoucherRequest, id uint) (dto.VoucherResponse, error) {
	startAt, expAt, err := parseVoucherDates(req.StartAt, req.ExpiresAt)
	if err != nil {
		return dto.VoucherResponse{}, err
	}

	var updatedVoucher *model.Voucher

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		voucher, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if voucher == nil {
			return appError.NotFound(fmt.Sprintf("Voucher not found with ID: %d", id))
		}

		if req.DiscountType == model.DiscountTypePercent && req.DiscountValue > 100 {
			return appError.BadRequest("Percent discount must be between 1 and 100")
		}

		if req.DiscountType == model.DiscountTypeFixed && req.MaxDiscount != nil {
			return appError.BadRequest("max_discount only applies to percent vouchers")
		}

		code := strings.ToUpper(strings.TrimSpace(req.Code))
		if code != voucher.Code {
			exist, err := txRepo.ExistByCode(ctx, code)
			if err != nil {
				return err
			}
			if exist {
				return appError.AlreadyExist("Voucher code already exist")
			}
		}

		mapper.UpdateVoucherModel(voucher, req, startAt, expAt)
		if err := txRepo.Update(ctx, voucher); err != nil {
			return err
		}

		updatedVoucher = voucher

		return nil
	}); err != nil {
		return dto.VoucherResponse{}, err
	}

	return mapper.ToVoucherResponse(updatedVoucher), nil
}

func (s *VoucherService) Delete(ctx context.Context, id uint) error {
	return s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)
		txBookingRepo := s.bookingRepo.WithTx(tx)

		voucher, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if voucher == nil {
			return appError.NotFound(fmt.Sprintf("Voucher not found with ID: %d", id))
		}

		if err := txBookingRepo.SetNullVoucherID(ctx, voucher.ID); err != nil {
			return err
		}

		return txRepo.Delete(ctx, voucher)
	})
}
