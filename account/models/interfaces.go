package models

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name UserService
type UserServiceI interface {
	Get(context context.Context, uid uuid.UUID) (*User, error)
}

type UserRepositoryI interface {
	GetById(ctx context.Context, uid uuid.UUID) (*User, error)
}
