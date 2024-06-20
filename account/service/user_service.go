package service

import (
	"context"
	"github.com/google/uuid"
	"memory-app/models"
)

type UserService struct {
	UserRepository models.UserRepositoryI
}

type UserServiceConfig struct {
	UserRepo models.UserRepositoryI
}

func NewUserService(cfg *UserServiceConfig) models.UserServiceI {
	return &UserService{
		UserRepository: cfg.UserRepo,
	}
}

func (u *UserService) Get(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	return u.UserRepository.GetById(ctx, uid)

}

var _ models.UserServiceI = &UserService{}
