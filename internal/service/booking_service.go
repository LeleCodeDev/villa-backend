package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type BookingService struct {
	txManager        *repository.TxManager
	repo             *repository.BookingRepository
	bookingAddonRepo *repository.BookingAddonRepository
}

func NewBookingService(
	txManager *repository.TxManager,
	repo *repository.BookingRepository,
	bookingAddonRepo *repository.BookingAddonRepository,
) *BookingService {
	return &BookingService{
		txManager:        txManager,
		repo:             repo,
		bookingAddonRepo: bookingAddonRepo,
	}
}

func (s *BookingService) GetAll(ctx context.Context, query dto.BookingQuery) ([]dto.BookingResponse, int64, error) {
	bookings, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, nil
	}

	bookingIds := make([]uint, 0, len(bookings))
	for _, booking := range bookings {
		bookingIds = append(bookingIds, booking.ID)
	}

	bookingAddons, err := s.bookingAddonRepo.GetAllInBookingIds(ctx, bookingIds)
	if err != nil {
		return nil, 0, err
	}

	bookingAddonMap := make(map[uint][]model.BookingAddon, len(bookings))
	for _, bookingAddon := range bookingAddons {
		bookingAddonMap[bookingAddon.BookingID] = append(bookingAddonMap[bookingAddon.BookingID], bookingAddon)
	}

	responses := make([]dto.BookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		responses = append(responses, mapper.ToBookingResponse(&booking, bookingAddonMap[booking.ID]))
	}

	return responses, total, nil
}

func (s *BookingService) GetByID(ctx context.Context, id uint) (dto.BookingResponse, error) {
	booking, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.BookingResponse{}, err
	}
	if booking == nil {
		return dto.BookingResponse{}, appError.NotFound(fmt.Sprintf("Booking not found with ID: %d", id))
	}

	bookingAddons, err := s.bookingAddonRepo.GetAllByBookingID(ctx, id)
	if err != nil {
		return dto.BookingResponse{}, err
	}

	return mapper.ToBookingResponse(booking, bookingAddons), nil
}
