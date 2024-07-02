package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
)

type TokenService struct {
	mock.Mock
}

func (t *TokenService) Signout(context context.Context, uid uuid.UUID) error {
	ret := t.Called(context, uid)

	var r1 error

	if ret.Get(0) != nil {
		r1 = ret.Get(0).(error)
	}
	return r1
}

func (t *TokenService) ValidateRefreshToken(token string) (*models.RefreshToken, error) {
	ret := t.Called(token)

	var r0 *models.RefreshToken
	var r1 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.RefreshToken)
	}
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	return r0, r1
}

func (t *TokenService) ValidateIDToken(token string) (*models.User, error) {
	ret := t.Called(token)

	var r0 *models.User
	var r1 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.User)
	}
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	return r0, r1
}

func (t *TokenService) GetPairForUser(ctx context.Context, u *models.User, prevIdToken string) (*models.TokenPair, error) {
	ret := t.Called(ctx, u, prevIdToken)

	var r0 *models.TokenPair
	var r1 error

	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.TokenPair)
	}
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}
	return r0, r1

}

var _ models.TokenServiceI = &TokenService{}
