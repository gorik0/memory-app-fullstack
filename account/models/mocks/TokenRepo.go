package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"time"
)

type TokenRepo struct {
	mock.Mock
}

func (t *TokenRepo) DeleteUserRefreshToken(ctx context.Context, userId string) error {
	ret := t.Called(ctx, userId)

	var r1 error

	if ret.Get(0) != nil {
		r1 = ret.Get(0).(error)
	}
	return r1

}

func (t *TokenRepo) SetRefreshToken(ctx context.Context, userId string, tokenId string, expiTime time.Duration) error {
	called := t.Called(ctx, userId, tokenId, expiTime)

	var e error
	if called.Get(0) != nil {
		e := called.Get(0).(error)
		return e
	}
	return e
}

func (t *TokenRepo) DeleteRefreshToken(ctx context.Context, userId string, prevTokenId string) error {
	called := t.Called(ctx, userId, prevTokenId)

	var e error
	if called.Get(0) != nil {
		e := called.Get(0).(error)
		return e
	}
	return e
}

var _ models.TokenRepository = &TokenRepo{}
