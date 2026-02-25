package service

import (
	"context"
	"sika/internal/domain"
	"sika/internal/repositoy"
)

type UserService interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
}

type userService struct {
	userRepo repositoy.UserRepository
}

func NewUserService(userRepository repositoy.UserRepository) UserService {
	return &userService{
		userRepo: userRepository,
	}
}

func (us *userService) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	return us.userRepo.Create(ctx, user)
}

func (us *userService) Get(ctx context.Context, id string) (*domain.User, error) {
	return us.userRepo.Get(ctx, id)
}
