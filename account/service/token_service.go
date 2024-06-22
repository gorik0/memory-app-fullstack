package service

import (
	"context"
	"crypto/rsa"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
)

type TokenService struct {
	PrivKey *rsa.PrivateKey
	PublKey *rsa.PublicKey

	RefreshSecret string
}

type ConfigTokenService struct {
	PrivKey *rsa.PrivateKey
	PublKey *rsa.PublicKey

	RefreshSecret string
}

func NewTokenService(cfg *ConfigTokenService) models.TokenServiceI {
	return &TokenService{
		PrivKey:       cfg.PrivKey,
		PublKey:       cfg.PublKey,
		RefreshSecret: cfg.RefreshSecret,
	}
}

func (t TokenService) GetPairForUser(context context.Context, u *models.User, prevIdToken string) (*models.TokenPair, error) {

	//:::ID TOKEN generate
	idToken, err := generateToken(u, t.PrivKey)
	if err != nil {
		log.Printf("Coudldn't genereta token for user :::%v  with errror ::: %v \n", u, err)
		e := apprerrors.NewInternal()
		return nil, e
	}

	//:::REFRESH TOKEN generate

	refreshToken, err := generateRefreshToken(u.UID, t.RefreshSecret)
	if err != nil {
		log.Printf("Coudldn't genereta refresh token for user :::%v  with errror ::: %v \n", u, err)
		e := apprerrors.NewInternal()
		return nil, e
	}

	//:::RETURN TOKEN pair

	return &models.TokenPair{
		IdToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil

}
