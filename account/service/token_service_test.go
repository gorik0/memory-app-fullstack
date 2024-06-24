package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"memory-app/account/models"
	"memory-app/account/models/mocks"
	"testing"
	"time"
)

func TestNewPairFromUser(t *testing.T) {
	privString, _ := ioutil.ReadFile("../rsa_private_.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privString)
	publicString, _ := ioutil.ReadFile("../rsa_public_.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicString)
	secret := "egoriktestownservice"

	//uid, _ := uuid.NewRandom()
	userIDsuccess, _ := uuid.NewRandom()
	userIDfail, _ := uuid.NewRandom()
	user := &models.User{
		UID:      userIDsuccess,
		Email:    "gorik@ko.ru",
		Password: "213123124",
	}

	tokenRepoMock := new(mocks.TokenRepo)
	argsSuccess := mock.Arguments{
		mock.AnythingOfType("context.backgroundCtx"),
		userIDsuccess.String(),
		mock.Anything,
		mock.AnythingOfType("time.Duration"),
	}
	argsFail := mock.Arguments{
		mock.AnythingOfType("context.backgroundCtx"),
		userIDfail.String(),
		mock.Anything,
		mock.AnythingOfType("time.Duration"),
	}
	tokenRepoMock.On("SetRefreshToken", argsSuccess...).Return(nil)
	tokenRepoMock.On("SetRefreshToken", argsFail...).Return(fmt.Errorf("SOME errr...."))

	tokenService := NewTokenService(&ConfigTokenService{
		PrivKey:         privKey,
		PublKey:         publicKey,
		RefreshSecret:   secret,
		RefreshExp:      "20000",
		IdExp:           "900",
		TokenRepository: tokenRepoMock,
	})
	t.Run("RETURN corerct token", func(t *testing.T) {

		ctx := context.Background()
		//user.UID = userIDsuccess
		//::TOKEN pair creating
		tokenPair, err := tokenService.GetPairForUser(ctx, user, "")
		assert.NoError(t, err)

		//::ID token parsing & asserting

		customClaims := &IDTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.IdToken, customClaims, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		assert.NoError(t, err)

		expectedClaims := []interface{}{
			user.UID,
			user.Email,
			user.Name,
		}
		actualClaims := []interface{}{
			customClaims.User.UID,
			customClaims.User.Email,
			customClaims.User.Name,
		}
		assert.ElementsMatch(t, actualClaims, expectedClaims)
		//	::ASSERT expired dates

		expiredExpected := time.Now().Add(time.Minute * 15)
		expiredActual := customClaims.ExpiresAt

		assert.WithinDuration(t, expiredActual.Time, expiredExpected, time.Second*5)

		//::REFRESH parsing & asserting
		customClaimsRefresh := &RefreshTokenCustomClaims{}
		_, err = jwt.ParseWithClaims(tokenPair.IdToken, customClaimsRefresh, func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		})
		assert.NoError(t, err)
		//	::ASSERT expired dates

		expiredExpectedRefr := time.Now().Add(time.Minute * 15)
		expiredActualRefr := customClaims.ExpiresAt

		assert.WithinDuration(t, expiredActualRefr.Time, expiredExpectedRefr, time.Second*5)

	})
	t.Run("Fail to refresh in redis ", func(t *testing.T) {

		ctx := context.Background()
		user.UID = userIDfail
		//::TOKE_ = tokN pair creating
		tokenPair, err := tokenService.GetPairForUser(ctx, user, "")
		assert.Error(t, err)
		//::ID token parsing & asserting

		tokenRepoMock.AssertNotCalled(t, "DeleteRefreshToken")
		_ = tokenPair

	})

}
