package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
)

type TokenService struct {
	mock.Mock
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
