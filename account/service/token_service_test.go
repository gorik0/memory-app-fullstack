package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"memory-app/account/models"
	"testing"
	"time"
)

func TestNewPairFromUser(t *testing.T) {
	privString, _ := ioutil.ReadFile("../rsa_private_.pem")
	privKey, _ := jwt.ParseRSAPrivateKeyFromPEM(privString)
	publicString, _ := ioutil.ReadFile("../rsa_public_.pem")
	publicKey, _ := jwt.ParseRSAPublicKeyFromPEM(publicString)
	secret := "egoriktestownservice"

	tokenService := NewTokenService(&ConfigTokenService{
		PrivKey:       privKey,
		PublKey:       publicKey,
		RefreshSecret: secret,
	})

	uid, _ := uuid.NewRandom()
	user := &models.User{
		UID:      uid,
		Email:    "gorik@ko.ru",
		Password: "213123124",
	}

	t.Run("RETURN corerct token", func(t *testing.T) {

		ctx := context.TODO()

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

}
