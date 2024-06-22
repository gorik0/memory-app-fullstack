package service

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
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
type IDTokenCustomClaims struct {
	User *models.User `json:"user"`
	jwt.RegisteredClaims
}

func generateRefreshToken(uid uuid.UUID, secret string) (*RefreshToken, error) {
	//:::PRE setup
	currentTime := time.Now()
	tokenExp := currentTime.AddDate(0, 0, 3)

	tokenID, _ := uuid.NewRandom()

	//:::CLAIMS setup
	customClaims := RefreshTokenCustomClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenExp),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ID:        tokenID.String(),
		},
	}

	//:::TOKEN setup
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	signedToken, err := token.SignedString(secret)

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	//:::RETURN
	return &RefreshToken{
		SS:      signedToken,
		ID:      tokenID.String(),
		Expired: tokenExp.Sub(currentTime),
	}, nil
}

func generateToken(u *models.User, key *rsa.PrivateKey) (string, error) {
	{
		//:::PRE setup
		currentTime := time.Now().Unix()
		tokenExp := currentTime + 60*15

		tokenID, _ := uuid.NewRandom()

		//:::CLAIMS setup
		customClaims := IDTokenCustomClaims{
			User: u,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Unix(tokenExp, 0)),
				IssuedAt:  jwt.NewNumericDate(time.Unix(currentTime, 0)),
				ID:        tokenID.String(),
			},
		}

		//:::TOKEN setup
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, customClaims)
		signedToken, err := token.SignedString(key)

		if err != nil {
			log.Println("Failed to sign refresh token string")
			return "", err
		}

		//:::RETURN
		return signedToken, nil
	}
}
