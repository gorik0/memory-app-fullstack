package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"memory-app/account/models/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_MeAbout(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("userExist", func(t *testing.T) {

		mockUid, err := uuid.NewRandom()
		assert.NoError(t, err)

		mockUser := &models.User{
			UID:  mockUid,
			Name: "gorik",
		}

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		userServiceMock.On("Get", mock.AnythingOfType("*gin.Context"), mockUid).Return(mockUser, nil)

		router := gin.Default()
		router.Use(func(context *gin.Context) {
			context.Set("user", mockUser)
		})

		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:      router,
			UserService: userServiceMock,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		userMockBytes, err := json.Marshal(gin.H{
			"user": mockUser,
		})
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, userMockBytes, rr.Body.Bytes())

		userServiceMock.AssertExpectations(t)

	})
	t.Run("userNotFound", func(t *testing.T) {

		mockUid, err := uuid.NewRandom()
		assert.NoError(t, err)

		//::USER SERVICE mock via mockery
		userServiceMock := mocks.NewUserServiceI(t)
		userServiceMock.On("Get", mock.Anything, mockUid).Return(nil, fmt.Errorf("erro while getting user from down"))

		router := gin.Default()
		router.Use(func(context *gin.Context) {
			context.Set("user", &models.User{
				UID: mockUid,
			})
		})

		rr := httptest.NewRecorder()

		NewHandler(&Config{
			Engine:      router,
			UserService: userServiceMock,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		respError := apprerrors.NewNotFound("user", mockUid.String())
		respBody, err := json.Marshal(gin.H{
			"error": respError,
		})

		assert.NoError(t, err)

		assert.Equal(t, respError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		userServiceMock.AssertExpectations(t)

	})
	t.Run("GetWithNoContext", func(t *testing.T) {

		//::USER SERVICE mock via mockery
		mockUserService := mocks.NewUserServiceI(t)
		//mockUserService.On("Get", mock.Anything, mock.Anything).Return(nil, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// do not append user to context
		router := gin.Default()
		NewHandler(&Config{
			Engine:      router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/me", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		mockUserService.AssertNotCalled(t, "Get", mock.Anything)
		println("passed")
		time.Sleep(time.Second)
	})

}
