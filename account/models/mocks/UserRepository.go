package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"memory-app/models"
)

type UserRepository struct {
	mock.Mock
}

func (u *UserRepository) GetById(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	ret := u.Called(ctx, uid)

	var r0 *models.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*models.User)
	}

	var r1 error
	if ret.Get(1) != nil {
		r1 = ret.Get(1).(error)
	}

	return r0, r1

}

var _ models.UserRepositoryI = &UserRepository{}
