package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"memory-app/account/models/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Signup(t *testing.T) {

	gin.SetMode(gin.TestMode)

	//t.Run("success", func(t *testing.T) {
	//
	//	//		:::CREATING request
	//	reqBody, err := json.Marshal(gin.H{
	//		"Email":    "gorik@ri.ro",
	//		"Password": "gorikwee23e",
	//	})
	//	assert.NoError(t, err)
	//
	//	//::USER SERVICE mock via mockery
	//	userServiceMock := mocks.NewUserServiceI(t)
	//	userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(nil)
	//
	//	//:::Router setup
	//	router := gin.Default()
	//	rr := httptest.NewRecorder()
	//
	//	NewHandler(&Config{
	//		Engine:      router,
	//		UserService: userServiceMock,
	//	})
	//	//:::Request setup
	//	request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
	//	assert.NoError(t, err)
	//	request.Header.Set("Content-Type", "application/json")
	//	//::MAKE request
	//	router.ServeHTTP(rr, request)
	//
	//	//::ASSERTION
	//	assert.Equal(t, http.StatusOK, rr.Code)
	//
	//	userServiceMock.AssertExpectations(t)
	//
	//})
	t.Run("wrongEmail", func(t *testing.T) {
		//		:::CREATING request
		reqBody, err := json.Marshal(gin.H{
			"Email":    "gorik0riro",
			"Password": "gorikwee23e",
		})
		assert.NoError(t, err)

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		//userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(nil)

		//:::Router setup
		router := gin.Default()
		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:      router,
			UserService: userServiceMock,
		})
		//:::Request setup
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		userServiceMock.AssertNotCalled(t, "Signup")

	})
	t.Run("wrongPassword", func(t *testing.T) {

		//		:::CREATING request
		reqBody, err := json.Marshal(gin.H{
			"Email":    "gorik@dd.ri",
			"Password": "go",
		})
		assert.NoError(t, err)

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		//userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(nil)

		//:::Router setup
		router := gin.Default()
		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:      router,
			UserService: userServiceMock,
		})
		//:::Request setup
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		userServiceMock.AssertNotCalled(t, "Signup")

	})
	t.Run("fieldsFRequired", func(t *testing.T) {

		//		:::CREATING request
		reqBody, err := json.Marshal(gin.H{
			"Email": "gorik@dd.ri",
		})
		assert.NoError(t, err)

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		//userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(nil)

		//:::Router setup
		router := gin.Default()
		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:      router,
			UserService: userServiceMock,
		})
		//:::Request setup
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		userServiceMock.AssertNotCalled(t, "Signup")

	})
	t.Run("ErrorFromUserService", func(t *testing.T) {

		//		:::CREATING request
		reqBody, err := json.Marshal(gin.H{
			"Email":    "gorik@dd.ri",
			"Password": "sdasdsddg0eqwrw",
		})
		assert.NoError(t, err)

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(fmt.Errorf("Something was going wrong..."))

		//:::Router setup
		router := gin.Default()
		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:      router,
			UserService: userServiceMock,
		})
		//:::Request setup
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		userServiceMock.AssertExpectations(t)

	})
	t.Run("GetTokenSuccess", func(t *testing.T) {

		//		:::CREATING request
		reqBody, err := json.Marshal(gin.H{
			"Email":    "gorik@dd.ri",
			"Password": "sdasdsddg0eqwrw",
		})
		assert.NoError(t, err)

		//::TOKEN SERVICE mock via mockery

		respToken := &models.TokenPair{
			IdToken:      "99999",
			RefreshToken: "000000",
		}

		var tokenServiceMock = new(mocks.TokenService)
		tokenServiceMock.On("GetPairForUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User"), mock.Anything).Return(respToken, nil)

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(nil)

		//:::Router setup
		router := gin.Default()
		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:        router,
			UserService:   userServiceMock,
			TokenServiceI: tokenServiceMock,
		})
		//:::Request setup
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")
		expResp, _ := json.Marshal(gin.H{"token": respToken})

		router.ServeHTTP(rr, request)
		assert.Equal(t, expResp, rr.Body.Bytes())
		assert.Equal(t, http.StatusCreated, rr.Code)

		userServiceMock.AssertExpectations(t)

	})
	t.Run("GetTokenFail", func(t *testing.T) {

		//		:::CREATING request
		reqBody, err := json.Marshal(gin.H{
			"Email":    "gorik@dd.ri",
			"Password": "sdasdsddg0eqwrw",
		})
		assert.NoError(t, err)

		//::TOKEN SERVICE mock via mockery

		var tokenServiceMock = new(mocks.TokenService)
		tokenServiceMock.On("GetPairForUser", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User"), mock.Anything).Return(nil, fmt.Errorf("Something going wromg"))

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		userServiceMock.On("Signup", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("*models.User")).Return(nil)

		//:::Router setup
		router := gin.Default()
		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:        router,
			UserService:   userServiceMock,
			TokenServiceI: tokenServiceMock,
		})
		//:::Request setup
		request, err := http.NewRequest(http.MethodPost, "/signup", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		userServiceMock.AssertExpectations(t)
		tokenServiceMock.AssertExpectations(t)

	})

}
