package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"memory-app/account/models/fixture"
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

		mockUserRepo := new(mocks.UserRepositoryI)
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

		mockUserRepo := new(mocks.UserRepositoryI)
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

		mockUserRepo := new(mocks.UserRepositoryI)
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

		mockUserRepo := new(mocks.UserRepositoryI)
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

	mockURep := new(mocks.UserRepositoryI)

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

func TestUpdate(t *testing.T) {
	//	::: USECASES
	userFail := models.User{Email: "sd"}
	userSuccess := models.User{Email: "s0"}
	//	::; MOCK CREATING
	mo := new(mocks.UserRepositoryI)
	mo.On("Update", mock.AnythingOfType("context.backgroundCtx"), &userSuccess).Return(nil)
	mo.On("Update", mock.AnythingOfType("context.backgroundCtx"), &userFail).Return(fmt.Errorf("Somw error"))

	//	::: REQUIRED PARAM

	ctx := context.Background()
	service := NewUserService(&UserServiceConfig{UserRepo: mo})

	//	::: TEST

	t.Run("userFail", func(t *testing.T) {

		err := service.UpdateDetail(ctx, &userFail)
		assert.Error(t, err)

	})
	t.Run("userSuccess", func(t *testing.T) {
		err := service.UpdateDetail(ctx, &userSuccess)
		assert.NoError(t, err)
	})
}

func TestSetProfileImage(t *testing.T) {
	//	::: USECASES
	userFailGetUser := models.User{UID: uuid.New()}
	userFailUpdateImage := models.User{UID: uuid.New()}
	imageFileNameSuccess := "succeess"
	imageFileNameFailImageRepo := "fail"
	imageFileNameSuccessWithExt := imageFileNameSuccess + ".png"
	imageFileName_FailImageRepo_WithExt := imageFileNameFailImageRepo + ".png"
	userFailImageRepo := models.User{UID: uuid.New(), ImageURL: imageFileNameFailImageRepo}
	user := models.User{UID: uuid.New(), ImageURL: imageFileNameSuccess}
	//imageFileNameFailImageRepo:="fail"

	//	::: REQUIRED PARAM

	errGetUser := fmt.Errorf("fail to return user")
	errImageUpdate := fmt.Errorf("Some errror")
	errImageRepo := fmt.Errorf("Some eerrror ")
	ctx := context.Background()
	image := fixture.NewMultipartImage("egor.png", "image/png")
	imageFileHeader := image.GetFormFile()
	//imageFile,err:=imageFileHeader.Open()
	//assert.NoError(t, err)

	//	::; MOCK CREATING
	moUser := new(mocks.UserRepositoryI)
	moImage := new(mocks.ImageRepository)
	moUser.On("GetById", mock.AnythingOfType("context.backgroundCtx"), user.UID).Return(&user, nil)
	moUser.On("GetById", mock.AnythingOfType("context.backgroundCtx"), userFailImageRepo.UID).Return(&userFailImageRepo, nil)
	moUser.On("GetById", mock.AnythingOfType("context.backgroundCtx"), userFailUpdateImage.UID).Return(&user, nil)
	moUser.On("GetById", mock.AnythingOfType("context.backgroundCtx"), userFailGetUser.UID).Return(&user, errGetUser)
	moUser.On("UpdateImage", mock.AnythingOfType("context.backgroundCtx"), user.UID, imageFileNameSuccessWithExt).Return(&user, nil)
	moUser.On("UpdateImage", mock.AnythingOfType("context.backgroundCtx"), userFailImageRepo.UID, imageFileNameSuccessWithExt).Return(&user, nil)
	moUser.On("UpdateImage", mock.AnythingOfType("context.backgroundCtx"), userFailUpdateImage.UID, imageFileNameSuccessWithExt).Return(&user, errImageUpdate)
	moImage.On("UpdateProfile", mock.AnythingOfType("context.backgroundCtx"), imageFileNameSuccessWithExt, mock.AnythingOfType("multipart.sectionReadCloser")).Return(imageFileNameSuccessWithExt, nil)
	moImage.On("UpdateProfile", mock.AnythingOfType("context.backgroundCtx"), imageFileName_FailImageRepo_WithExt, mock.AnythingOfType("multipart.sectionReadCloser")).Return("", errImageRepo)
	userService := NewUserService(&UserServiceConfig{UserRepo: moUser, ImageRepo: moImage})
	//	::: TEST

	t.Run("Success", func(t *testing.T) {

		userReturn, err := userService.SetProfileImage(ctx, user.UID, imageFileHeader)
		assert.NoError(t, err)
		fmt.Println(user)
		fmt.Println(userReturn)
		assert.Equal(t, user, *userReturn)

	})
	t.Run("FailGetUser", func(t *testing.T) {

		_, err := userService.SetProfileImage(ctx, userFailGetUser.UID, imageFileHeader)
		assert.Error(t, err)
		assert.Equal(t, errGetUser, err)

	})

	t.Run("FailUpdateImageUser", func(t *testing.T) {

		_, err := userService.SetProfileImage(ctx, userFailUpdateImage.UID, imageFileHeader)
		assert.Error(t, err)
		assert.Equal(t, errImageUpdate, err)

	})
	t.Run("FailImageRepo", func(t *testing.T) {

		_, err := userService.SetProfileImage(ctx, userFailImageRepo.UID, imageFileHeader)
		assert.Error(t, err)
		assert.Equal(t, errImageRepo, err)

	})

}

//success
//errorGetUser
//errorUpdateImage
//errorImageRepo
