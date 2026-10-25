package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	txManager *repository.TxManager
	repo      *repository.UserRepository
}

func NewUserService(txManager *repository.TxManager, repo *repository.UserRepository) *UserService {
	return &UserService{
		txManager: txManager,
		repo:      repo,
	}
}

func (s *UserService) GetAll(ctx context.Context, query dto.UserQuery) ([]dto.UserResponse, int64, error) {
	users, total, err := s.repo.GetAll(ctx, query)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		responses = append(responses, mapper.ToUserResponse(&user))
	}

	return responses, total, nil
}

func (s *UserService) GeyByID(ctx context.Context, id uint) (dto.UserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if user == nil {
		return dto.UserResponse{}, appError.NotFound(fmt.Sprintf("User not found with ID: %d", id))
	}

	return mapper.ToUserResponse(user), nil
}

func (s *UserService) Create(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	var createdUser *model.User

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		exist, err := txRepo.ExistByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		if exist {
			return appError.AlreadyExist("User email already exist!")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := mapper.ToUserModel(req, string(hashedPassword))
		if err := txRepo.Create(ctx, user); err != nil {
			return err
		}

		createdUser = user

		return nil
	}); err != nil {
		return dto.UserResponse{}, err
	}

	return mapper.ToUserResponse(createdUser), nil
}

func (s *UserService) Update(ctx context.Context, id uint, req dto.UpdateUserRequest) (dto.UserResponse, error) {
	var updatedUser *model.User

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.repo.WithTx(tx)

		user, err := txRepo.GetByID(ctx, id)
		if err != nil {
			return err
		}
		if user == nil {
			return appError.NotFound(fmt.Sprintf("User not found with ID: %d", id))
		}

		mapper.UpdateUserModel(user, req)
		if err := txRepo.Update(ctx, user); err != nil {
			return err
		}

		updatedUser = user

		return nil
	}); err != nil {
		return dto.UserResponse{}, err
	}

	return mapper.ToUserResponse(updatedUser), nil
}
