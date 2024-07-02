package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"memory-app/account/models/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDetails(t *testing.T) {
	//	::: USEECASES

	userFromContext := &models.User{
		UID: uuid.New(),
	}

	successUser := DetailsRequest{
		Email:   "success@ui.ru",
		Name:    "qwe",
		Website: "https://sfd.rt",
	}
	successUserForService := &models.User{
		UID:     userFromContext.UID,
		Email:   successUser.Email,
		Name:    successUser.Name,
		Website: successUser.Website,
	}
	failUser := DetailsRequest{
		Email:   "fail@ui.ru",
		Name:    "qwe",
		Website: "https://sfd.rt",
	}
	failUserInvalidData := DetailsRequest{
		Email:   "fail@",
		Name:    "qwe",
		Website: "https://sfd.rt",
	}
	failUserForServiceInvalidData := &models.User{
		UID:     userFromContext.UID,
		Email:   failUserInvalidData.Email,
		Name:    failUserInvalidData.Name,
		Website: failUserInvalidData.Website,
	}
	failUserForService := &models.User{
		UID:     userFromContext.UID,
		Email:   failUser.Email,
		Name:    failUser.Name,
		Website: failUser.Website,
	}

	//	::: MOCKS CREATING

	moUserService := new(mocks.UserServiceI)
	moUserService.On("UpdateDetail", mock.AnythingOfType("*gin.Context"), successUserForService).Return(nil)
	moUserService.On("UpdateDetail", mock.AnythingOfType("*gin.Context"), failUserForService).Return(fmt.Errorf("some error "))
	moUserService.On("UpdateDetail", mock.AnythingOfType("*gin.Context"), failUserForServiceInvalidData).Return(fmt.Errorf("some error "))
	//	::: REQUIRED PARAMS
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.Use(func(context *gin.Context) {
		context.Set("user", userFromContext)
	})

	NewHandler(&Config{Engine: router, UserService: moUserService})

	//	::: TESTS
	t.Run("Success", func(t *testing.T) {
		userBytes, err := json.Marshal(successUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(userBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		moUserService.AssertCalled(t, "UpdateDetail", mock.AnythingOfType("*gin.Context"), successUserForService)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
	t.Run("Fail", func(t *testing.T) {
		userBytes, err := json.Marshal(failUser)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(userBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		moUserService.AssertCalled(t, "UpdateDetail", mock.AnythingOfType("*gin.Context"), failUserForService)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
	t.Run("FailValidateData", func(t *testing.T) {
		userBytes, err := json.Marshal(failUserInvalidData)
		assert.NoError(t, err)
		req := httptest.NewRequest(http.MethodPut, "/details", bytes.NewBuffer(userBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		moUserService.AssertNotCalled(t, "UpdateDetail", mock.AnythingOfType("*gin.Context"), failUserForServiceInvalidData)
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
}
