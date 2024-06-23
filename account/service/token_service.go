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
	RefreshExp    string
	IdExp         string

	TokenRepository models.TokenRepository
}

type ConfigTokenService struct {
	PrivKey *rsa.PrivateKey
	PublKey *rsa.PublicKey

	RefreshSecret string

	RefreshExp string
	IdExp      string

	TokenRepository models.TokenRepository
}

func NewTokenService(cfg *ConfigTokenService) models.TokenServiceI {
	return &TokenService{
		PrivKey:         cfg.PrivKey,
		PublKey:         cfg.PublKey,
		RefreshSecret:   cfg.RefreshSecret,
		RefreshExp:      cfg.RefreshExp,
		IdExp:           cfg.IdExp,
		TokenRepository: cfg.TokenRepository,
	}
}

func (t TokenService) GetPairForUser(ctx context.Context, u *models.User, prevIdToken string) (*models.TokenPair, error) {

	//:::ID TOKEN generate
	idToken, err := generateToken(u, t.PrivKey, t.IdExp)
	if err != nil {
		log.Printf("Coudldn't genereta token for user :::%v  with errror ::: %v \n", u, err)
		e := apprerrors.NewInternal()
		return nil, e
	}

	//:::REFRESH TOKEN generate

	refreshToken, err := generateRefreshToken(u.UID, t.RefreshSecret, t.RefreshExp)
	if err != nil {
		log.Printf("Coudldn't genereta refresh token for user :::%v  with errror ::: %v \n", u, err)
		e := apprerrors.NewInternal()
		return nil, e
	}
	//::: REDIS job

	err = t.TokenRepository.SetRefreshToken(ctx, u.UID.String(), refreshToken.ID, refreshToken.Expired)
	if err != nil {
		log.Printf("Coudldn't set refresh in REDIS :::%v  with errror ::: %v \n", u, err)
		e := apprerrors.NewInternal()
		return nil, e
	}

	if prevIdToken != "" {
		err = t.TokenRepository.DeleteRefreshToken(ctx, u.UID.String(), prevIdToken)
		if err != nil {
			log.Printf("Coudldn't delete refresh  in REDIS :::%v  with errror ::: %v \n", u, err)
			e := apprerrors.NewInternal()
			return nil, e
		}
	}

	//:::RETURN TOKEN pair

	return &models.TokenPair{
		IdToken:      idToken,
		RefreshToken: refreshToken.SS,
	}, nil

}
