package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/models"
	"memory-app/models/mocks"
	"testing"
)

func TestGet(t *testing.T) {
	t.Run("RepositoryHaveUser", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		userExpected := &models.User{
			UID: uid,
		}

		mockUserRepo := new(mocks.UserRepository)
		mockUserRepo.On("GetById", mock.Anything, uid).Return(userExpected, nil)
		userService := NewUserService(&UserServiceConfig{
			UserRepo: mockUserRepo,
		})
		ctx := context.TODO()
		user, err := userService.Get(ctx, uid)

		assert.NoError(t, err)

		assert.Equal(t, userExpected, user)
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("RepositoryHaveNOuser", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		mockUserRepo := new(mocks.UserRepository)
		mockUserRepo.On("GetById", mock.Anything, uid).Return(nil, fmt.Errorf("has no user"))
		userService := NewUserService(&UserServiceConfig{
			UserRepo: mockUserRepo,
		})
		ctx := context.TODO()
		user, err := userService.Get(ctx, uid)

		assert.Nil(t, user)
		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})

}
