package models

import (
	"context"
	"github.com/google/uuid"
	"mime/multipart"
	"time"
)

//go:generate mockery --name UserServiceI
type UserServiceI interface {
	Get(context context.Context, uid uuid.UUID) (*User, error)
	Signup(context.Context, *User) error
	Signin(context.Context, *User) error
	UpdateDetail(ctx context.Context, u *User) error
	ClearProfileImage(ctx context.Context, uid uuid.UUID) error
	SetProfileImage(ctx context.Context, uid uuid.UUID, imageFileHeader *multipart.FileHeader) (*User, error)
}

type TokenServiceI interface {
	GetPairForUser(context context.Context, u *User, prevIdToken string) (*TokenPair, error)
	Signout(context context.Context, uid uuid.UUID) error
	ValidateIDToken(token string) (*User, error)
	ValidateRefreshToken(token string) (*RefreshToken, error)
}
//go:generate mockery --name UserRepositoryI

type UserRepositoryI interface {
	GetById(ctx context.Context, uid uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Create(ctx context.Context, u *User) error
	Update(ctx context.Context, u *User) error
	UpdateImage(ctx context.Context, uid uuid.UUID, imageURL string) (*User, error)
}

type TokenRepository interface {
	SetRefreshToken(ctx context.Context, userId string, tokenId string, expiTime time.Duration) error

	DeleteRefreshToken(ctx context.Context, userId string, prevTokenId string) error
	DeleteUserRefreshToken(ctx context.Context, userId string) error
}

//go:generate mockery --name ImageRepository

type ImageRepository interface {
	DeleteProfile(ctx context.Context, objName string) error
	UpdateProfile(ctx context.Context, objName string, imageFile multipart.File) (string, error)
}
