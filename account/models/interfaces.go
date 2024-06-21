package models

import (
	"context"
	"github.com/google/uuid"
)

//go:generate mockery --name UserServiceI
type UserServiceI interface {
	Get(context context.Context, uid uuid.UUID) (*User, error)
	Signup(context.Context, *User) error
}

type TokenServiceI interface {
	GetPairForUser(context context.Context, u *User, prevIdToken string) (*TokenPair, error)
}
type UserRepositoryI interface {
	GetById(ctx context.Context, uid uuid.UUID) (*User, error)
	Create(ctx context.Context, u *User) error
}
