package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSigninHandler(t *testing.T) {

	gin.SetMode(gin.TestMode)
	//	::: USERservice TOKENservice mocks

	mockUserService := new(mocks.UserServiceI)
	mockTokenService := new(mocks.TokenService)

	// ::: setup MOCKS
	user := SigninRequest{

		Email:    "egorik@ko.eu",
		Password: "1212312442",
	}

	userInvalidDataCase := user
	userInvalidDataCase.Password = "1"

	userFailServiceCase := user
	userFailServiceCase.Password = "failcaseserviceuser"

	tokens := &models.TokenPair{}

	mockTokenService.On("GetPairForUser", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: user.Email, Password: user.Password}, mock.Anything).Return(tokens, nil)
	mockUserService.On("Signin", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: user.Email, Password: user.Password}).Return(nil)

	mockUserService.On("Signin", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: userFailServiceCase.Email, Password: userFailServiceCase.Password}).Return(fmt.Errorf("eee"))
	mockTokenService.On("GetPairForUser", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: userFailServiceCase.Email, Password: userFailServiceCase.Password}, mock.Anything).Return(tokens, fmt.Errorf("ERRRR"))

	//:::HANDLER setup

	router := gin.Default()
	NewHandler(&Config{
		Engine:        router,
		UserService:   mockUserService,
		TokenServiceI: mockTokenService,

		HandlerTimeout: 5,
	})
	t.Run("SUCCESS", func(t *testing.T) {

		rr := httptest.NewRecorder()
		_ = rr
		reqBody, _ := json.Marshal(user)
		log.Println(string(reqBody))

		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))

		router := gin.Default()
		NewHandler(&Config{
			Engine:        router,
			UserService:   mockUserService,
			TokenServiceI: mockTokenService,

			HandlerTimeout: 5,
		})

		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		log.Println(rr.Body.Bytes())
		assert.Equal(t, http.StatusOK, rr.Code)
		mockTokenService.AssertCalled(t, "GetPairForUser", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: user.Email, Password: user.Password}, mock.Anything)
		mockUserService.AssertCalled(t, "Signin", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: user.Email, Password: user.Password})

	})
	t.Run("FAIL service", func(t *testing.T) {

		rr := httptest.NewRecorder()
		_ = rr
		reqBody, _ := json.Marshal(userFailServiceCase)
		log.Println(string(reqBody))

		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))

		router := gin.Default()
		NewHandler(&Config{
			Engine:        router,
			UserService:   mockUserService,
			TokenServiceI: mockTokenService,

			HandlerTimeout: 5,
		})

		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		log.Println(rr.Body.Bytes())
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockTokenService.AssertNotCalled(t, "GetPairForUser", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: userFailServiceCase.Email, Password: userFailServiceCase.Password}, mock.Anything)
		mockUserService.AssertCalled(t, "Signin", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: userFailServiceCase.Email, Password: userFailServiceCase.Password})

	})
	t.Run("INVALI BODUY", func(t *testing.T) {

		rr := httptest.NewRecorder()
		_ = rr
		reqBody, _ := json.Marshal(userInvalidDataCase)
		log.Println(string(reqBody))

		req := httptest.NewRequest(http.MethodPost, "/signin", bytes.NewBuffer(reqBody))

		router := gin.Default()
		NewHandler(&Config{
			Engine:        router,
			UserService:   mockUserService,
			TokenServiceI: mockTokenService,

			HandlerTimeout: 5,
		})

		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, req)
		log.Println(rr.Body.Bytes())
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockTokenService.AssertNotCalled(t, "GetPairForUser", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: userInvalidDataCase.Email, Password: userInvalidDataCase.Password}, mock.Anything)
		mockUserService.AssertNotCalled(t, "Signin", mock.AnythingOfType("context.backgroundCtx"), &models.User{Email: userInvalidDataCase.Email, Password: userInvalidDataCase.Password})

	})

}
