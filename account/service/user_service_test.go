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
	"net/http"
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

func TestSignin(t *testing.T) {
	//::: USER REPO setup

	mockURep := new(mocks.UserRepository)

	//::: DATA usecase setup

	//: 		context
	ctx := context.Background()

	//: 		password
	passwordRequest := "12345"
	passwordRequestHashed, _ := generateHashPassword(passwordRequest)

	//: 		email
	email_userExist_passwordMatch := "1test@st.te"
	email_userNOTExist_ := "2test@st.te"
	email_userExist_passwordNOTMatch := "3test@st.te"

	//: 		return user
	user_psw_match := &models.User{
		Email:    email_userExist_passwordMatch,
		Password: passwordRequestHashed,
	}
	user_psw_not_match := &models.User{
		Email:    email_userExist_passwordNOTMatch,
		Password: passwordRequestHashed + "wrong_password",
	}
	//:::USER REPO methods setup
	mockURep.On("GetByEmail", mock.AnythingOfType("context.backgroundCtx"), email_userExist_passwordMatch).Return(user_psw_match, nil)
	errToReturn := apprerrors.NewNotFound("email", "email")
	mockURep.On("GetByEmail", mock.AnythingOfType("context.backgroundCtx"), email_userNOTExist_).Return(nil, errToReturn)
	mockURep.On("GetByEmail", mock.AnythingOfType("context.backgroundCtx"), email_userExist_passwordNOTMatch).Return(user_psw_not_match, nil)

	//::: USER SERVICE create
	userService := NewUserService(&UserServiceConfig{UserRepo: mockURep})

	t.Run("userExist_passwordMatch", func(t *testing.T) {
		//:::create request
		userRequest := &models.User{

			Email:    email_userExist_passwordMatch,
			Password: passwordRequest,
		}

		err := userService.Signin(ctx, userRequest)

		//:ASSERT NO ERRORS
		assert.NoError(t, err)
		//:ASSERT password MATCHES
		assert.Equal(t, userRequest.Password, passwordRequestHashed)

	})
	t.Run("userNOTExist_", func(t *testing.T) {

		userRequest := &models.User{

			Email:    email_userNOTExist_,
			Password: passwordRequest,
		}
		err := userService.Signin(ctx, userRequest)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errToReturn)

	})
	t.Run("userExist_passwordNOTMatch", func(t *testing.T) {

		userRequest := &models.User{

			Email:    email_userExist_passwordNOTMatch,
			Password: passwordRequest,
		}
		err := userService.Signin(ctx, userRequest)
		assert.Equal(t, apprerrors.Status(err), http.StatusUnauthorized)

	})

}
