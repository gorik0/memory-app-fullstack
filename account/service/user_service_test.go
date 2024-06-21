package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"memory-app/account/models/mocks"
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

func TestSignup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {

		uid, _ := uuid.NewRandom()

		userToCreate := &models.User{
			Email:    "goriko@ko.ru",
			Password: "111",
		}

		mockUserRepo := new(mocks.UserRepository)
		mockUserRepo.On("Create", mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("*models.User")).
			Run(func(args mock.Arguments) {

				user := args.Get(1).(*models.User)
				user.UID = uid
			}).Return(nil)
		userService := NewUserService(&UserServiceConfig{
			UserRepo: mockUserRepo,
		})
		ctx := context.TODO()
		err := userService.Signup(ctx, userToCreate)

		assert.NoError(t, err)

		assert.Equal(t, uid, userToCreate.UID)
		mockUserRepo.AssertExpectations(t)
	})
	t.Run("Error", func(t *testing.T) {

		userToCreate := &models.User{
			Email:    "goriko@ko.ru",
			Password: "111",
		}

		mockError := apprerrors.NewConflict("email", "gogogoo")

		mockUserRepo := new(mocks.UserRepository)
		mockUserRepo.On("Create", mock.AnythingOfType("context.todoCtx"), mock.AnythingOfType("*models.User")).Return(mockError)
		userService := NewUserService(&UserServiceConfig{
			UserRepo: mockUserRepo,
		})
		ctx := context.TODO()
		err := userService.Signup(ctx, userToCreate)

		assert.Error(t, err)
		assert.Equal(t, mockError, err)

		mockUserRepo.AssertExpectations(t)
	})

}
