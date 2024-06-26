package service

import (
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log"
	"memory-app/account/models"
	"strconv"
	"time"
)

type RefreshTokenData struct {
	SS        string
	ID        uuid.UUID
	ExpiresIN time.Duration

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

func generateRefreshToken(uid uuid.UUID, secret string, exp string) (*RefreshTokenData, error) {
	//:::PRE setup

	exp64, err := strconv.ParseInt(exp, 0, 64)
	if err != nil {
		return nil, err
	}

	currentTime := time.Now()
	tokenExp := currentTime.Add(time.Duration(exp64) * time.Second)

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
	signedToken, err := token.SignedString([]byte(secret))

	if err != nil {
		log.Println("Failed to sign refresh token string")
		return nil, err
	}

	//:::RETURN
	return &RefreshTokenData{
		SS:      signedToken,
		ID:      tokenID,
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

func validateIDtoken(tokenString string, key *rsa.PublicKey) (*IDTokenCustomClaims, error) {
	claims := &IDTokenCustomClaims{}
	println(tokenString)
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		log.Println("Error while parsing tokenString with pub key ::: ", err.Error())
		println(err)
		return nil, fmt.Errorf("parsing tokenString with pub kye")
	}

	if !token.Valid {
		log.Println("Token is invalid!!! ::: ", err.Error())
		return nil, fmt.Errorf("invalid token")

	}

	claims, ok := token.Claims.(*IDTokenCustomClaims)
	if !ok {
		log.Println("Error while type assertion token-claims ::: ", err.Error())
		return nil, fmt.Errorf("Error while type assertion token-claims")
	}

	return claims, nil

}

func validateRefreshToken(tokenToValidate string, secret string) (*RefreshTokenCustomClaims, error) {
	claims := &RefreshTokenCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenToValidate, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		log.Println("Error while parsing tokenString with secret key ::: ", err.Error())
		println(err)
		return nil, fmt.Errorf("parsing tokenString with secrte kye")
	}

	if !token.Valid {
		log.Println("Token is invalid!!! ::: ", err.Error())
		return nil, fmt.Errorf("invalid token")

	}

	claims, ok := token.Claims.(*RefreshTokenCustomClaims)
	if !ok {
		log.Println("Error while type assertion token-claims ::: ", err.Error())
		return nil, fmt.Errorf("Error while type assertion token-claims")
	}

	return claims, nil
}
