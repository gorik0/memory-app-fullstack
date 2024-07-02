package handler

import (
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

func TestSignout(t *testing.T) {

	//	REQUIERD PARAMS
	gin.SetMode(gin.TestMode)

	//	USECASES
	userSuccess := models.User{UID: uuid.New(), Email: "go@2.ri"}
	userFail := models.User{UID: uuid.New(), Email: "go@2.ri"}

	//	MOCK CREATING
	mockTokenService := new(mocks.TokenService)

	mockTokenService.On("Signout", mock.AnythingOfType("context.backgroundCtx"), userFail.UID).Return(fmt.Errorf("...Some erroor..."))
	mockTokenService.On("Signout", mock.AnythingOfType("context.backgroundCtx"), userSuccess.UID).Return(nil)
	t.Run("success", func(t *testing.T) {
		router := gin.New()
		reader := httptest.NewRecorder()

		router.Use(func(context *gin.Context) {
			context.Set("user", &userSuccess)

		})

		NewHandler(&Config{
			Engine:        router,
			TokenServiceI: mockTokenService,
		})

		request := httptest.NewRequest(http.MethodPost, "/signout", nil)

		router.ServeHTTP(reader, request)
		mesBytes := reader.Body.Bytes()
		gi := gin.H{}
		err := json.Unmarshal(mesBytes, &gi)
		assert.NoError(t, err)

		code := reader.Code
		assert.Equal(t, http.StatusOK, code)
		_, ok := gi["mesaage"]
		assert.Equal(t, true, ok)
		fmt.Println(gi)
	})
	t.Run("token_service_error", func(t *testing.T) {
		router := gin.New()
		reader := httptest.NewRecorder()

		router.Use(func(context *gin.Context) {
			context.Set("user", &userFail)

		})

		NewHandler(&Config{
			Engine:        router,
			TokenServiceI: mockTokenService,
		})

		request := httptest.NewRequest(http.MethodPost, "/signout", nil)

		router.ServeHTTP(reader, request)
		//mesBytes := reader.Body.Bytes()
		//gi := gin.H{}
		//err := json.Unmarshal(mesBytes, &gi)
		//assert.Error(t, err)

		code := reader.Code
		assert.Equal(t, http.StatusInternalServerError, code)

	})

}
