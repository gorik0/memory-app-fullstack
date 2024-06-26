package models

import (
	"context"
	"github.com/google/uuid"
	"time"
)

//go:generate mockery --name UserServiceI
type UserServiceI interface {
	Get(context context.Context, uid uuid.UUID) (*User, error)
	Signup(context.Context, *User) error
	Signin(context.Context, *User) error
}

type TokenServiceI interface {
	GetPairForUser(context context.Context, u *User, prevIdToken string) (*TokenPair, error)
	ValidateIDToken(token string) (*User, error)
	ValidateRefreshToken(token string) (*RefreshToken, error)
}
type UserRepositoryI interface {
	GetById(ctx context.Context, uid uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userId string, tokenId string, expiTime time.Duration) error
	DeleteRefreshToken(ctx context.Context, userId string, prevTokenId string) error
}
