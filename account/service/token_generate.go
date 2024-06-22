package service

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"memory-app/account/models"
	"time"
)

type RefreshToken struct {
	SS string
	ID string

	Expired time.Duration
}
type RefreshTokenCustomClaims struct {
	Uid uuid.UUID `json:"uid"`
	jwt.RegisteredClaims
}

func generateRefreshToken(uid uuid.UUID, secret string) (*RefreshToken, error) {

	currentTime := time.Now()
	tokenTime := currentTime.AddDate(0, 0, 3)

	tokenID, _ := uuid.NewRandom()

	claims := RefreshTokenCustomClaims{}

	jwt.NewWithClaims()
}

func generateToken(u *models.User, key *rsa.PrivateKey) (string, error) {
	return "", err
}
