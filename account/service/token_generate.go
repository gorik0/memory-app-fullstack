package service

import (
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"memory-app/account/models"
	"strconv"
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

func generateRefreshToken(uid uuid.UUID, secret string, exp string) (*RefreshToken, error) {
	//:::PRE setup

	exp64, err := strconv.ParseInt(exp, 0, 64)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp64) * time.Second)

	tokenID, _ := uuid.NewRandom()
	println(secret)
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
	signedToken, err := token.SignedString([]byte(secret))

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

func generateToken(u *models.User, key *rsa.PrivateKey, exp string) (string, error) {
	{
		//:::PRE setup
		exp64, err := strconv.ParseInt(exp, 0, 64)
		if err != nil {
			return "", err
		}
		currentTime := time.Now().Unix()
		tokenExp := currentTime + exp64

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
