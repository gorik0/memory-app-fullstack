package models

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name UserService
type UserService interface {
	Get(context context.Context,uid uuid.UUID) (*User, error)
}

type UserRepo interface {
	GetById(uid uuid.UUID) (*User, error)
}
