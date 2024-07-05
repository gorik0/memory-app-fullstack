package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"memory-app/account/models"
	"memory-app/account/models/fixture"
	"memory-app/account/models/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestImage(t *testing.T) {
	//	:: USECASE
	user := &models.User{UID: uuid.New()}
	userServiceFail := &models.User{UID: uuid.New()}
	imageSuccess := fixture.NewMultipartImage("egor.png", "image/png")
	imageSuccess2 := fixture.NewMultipartImage("egor.png", "image/png")

	imageMimeTypeFail := fixture.NewMultipartImage("egor.oleg", "image/oleg")

	//	:: MOCKS

	mo := new(mocks.UserServiceI)

	imageSuccessUrl := uuid.New()
	mo.On("SetProfileImage", mock.AnythingOfType("context.backgroundCtx"), user.UID, mock.AnythingOfType("*multipart.FileHeader")).Return(&models.User{ImageURL: imageSuccessUrl.String()}, nil)
	mo.On("SetProfileImage", mock.AnythingOfType("context.backgroundCtx"), userServiceFail.UID, mock.AnythingOfType("*multipart.FileHeader")).Return(&models.User{}, fmt.Errorf("some errorr "))
	//	:: REQUIRED PARAMS
	gin.SetMode(gin.TestMode)

	//	:: TESTS

	t.Run("Success", func(t *testing.T) {
		router := gin.Default()

		w := httptest.NewRecorder()
		router.Use(func(context *gin.Context) {

			context.Set("user", user)
		})
		NewHandler(&Config{UserService: mo, Engine: router, MaxBytesSize: 10 * 1024})
		req := httptest.NewRequest(http.MethodPost, "/image", imageSuccess.MultipartBody)
		req.Header.Set("Content-Type", imageSuccess.ContentType)
		router.ServeHTTP(w, req)
		println("!!!!!!!!!!!!")
		println(w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		respBody, _ := json.Marshal(gin.H{"image": imageSuccessUrl.String()})
		assert.Equal(t, respBody, w.Body.Bytes())
	})
	t.Run("FailOnSrvice", func(t *testing.T) {
		router := gin.Default()

		w := httptest.NewRecorder()
		router.Use(func(context *gin.Context) {

			context.Set("user", userServiceFail)
		})
		NewHandler(&Config{UserService: mo, Engine: router, MaxBytesSize: 10 * 1024})
		req := httptest.NewRequest(http.MethodPost, "/image", imageSuccess2.MultipartBody)
		req.Header.Set("Content-Type", imageSuccess2.ContentType)
		router.ServeHTTP(w, req)
		println("!!!!!!!!!!!!")
		println(w.Body.String())

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("FailOnMimeType", func(t *testing.T) {
		router := gin.Default()

		w := httptest.NewRecorder()
		router.Use(func(context *gin.Context) {

			context.Set("user", user)
		})
		NewHandler(&Config{UserService: mo, Engine: router, MaxBytesSize: 10 * 1024})
		req := httptest.NewRequest(http.MethodPost, "/image", imageMimeTypeFail.MultipartBody)
		req.Header.Set("Content-Type", imageMimeTypeFail.ContentType)
		router.ServeHTTP(w, req)
		println("!!!!!!!!!!!!")

		println(w.Body.String())

		assert.Equal(t, http.StatusBadRequest, w.Code)
		//respBody, _ := json.Marshal(gin.H{"image": imageSuccessUrl.String()})
		//assert.Equal(t, respBody, w.Body.Bytes())
	})
	t.Run("EmptyContent", func(t *testing.T) {
		router := gin.Default()

		w := httptest.NewRecorder()
		router.Use(func(context *gin.Context) {

			context.Set("user", user)
		})
		NewHandler(&Config{UserService: mo, Engine: router, MaxBytesSize: 10 * 1024})
		req := httptest.NewRequest(http.MethodPost, "/image", nil)
		req.Header.Set("Content-Type", "multipart/form-data")
		router.ServeHTTP(w, req)
		println("!!!!!!!")
		println(w.Body.String())

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

}

//
//wrongMimeType
//emptyBody
//maxSizeExceed
//Success
//FailUserService
