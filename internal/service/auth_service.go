package service

import (
	"context"

	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/model"
	"github.com/lelecodedev/villa-backend/internal/repository"
	"github.com/lelecodedev/villa-backend/pkg/errors"
	"github.com/lelecodedev/villa-backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	txManager *repository.TxManager
	userRepo  *repository.UserRepository
}

func NewAuthService(txManager *repository.TxManager, userRepo *repository.UserRepository) *AuthService {
	return &AuthService{txManager: txManager, userRepo: userRepo}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (dto.AuthResponse, error) {
	var registeredUser *model.User

	if err := s.txManager.Transaction(ctx, func(tx *gorm.DB) error {
		txRepo := s.userRepo.WithTx(tx)

		exist, err := txRepo.ExistByEmail(ctx, req.Email)
		if err != nil {
			return err
		}
		if exist {
			return errors.AlreadyExist("User email already exist!")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := mapper.ToRegisteUserModel(req, string(hashedPassword))
		if err := txRepo.Create(ctx, user); err != nil {
			return err
		}

		registeredUser = user

		return nil
	}); err != nil {
		return dto.AuthResponse{}, err
	}

	token, err := jwt.GenerateToken(registeredUser.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return mapper.ToAuthResponse(registeredUser, token), nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	if user == nil {
		return dto.AuthResponse{}, errors.Unauthorized("Incorrect email or password!")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return dto.AuthResponse{}, errors.Unauthorized("Incorrect email or password!")
	}

	token, err := jwt.GenerateToken(user.ID)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	return mapper.ToAuthResponse(user, token), nil
}
