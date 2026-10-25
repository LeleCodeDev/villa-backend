package service

import (
	"context"
	"fmt"

	appError "github.com/lelecodedev/villa-backend/internal/apperror"
	"github.com/lelecodedev/villa-backend/internal/dto"
	"github.com/lelecodedev/villa-backend/internal/mapper"
	"github.com/lelecodedev/villa-backend/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
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
