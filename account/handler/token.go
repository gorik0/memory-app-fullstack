package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"memory-app/account/models/apprerrors"
	"net/http"
)

type TokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) Tokens(ctx *gin.Context) {
	var req = TokenRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		log.Println("Error while bindind refresh token request")
		ctx.JSON(apprerrors.Status(err), gin.H{
			"error": err,
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"say": "Tokens"})

	//	:: VALIDATE REFRESH token
	token, err := h.TokenService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(apprerrors.Status(err), gin.H{
			"error": err,
		})

	}
	//	:: GET USER for token

	user, err := h.UserService.Get(ctx, token.UID)
	if err != nil {

		ctx.JSON(apprerrors.Status(err), gin.H{
			"error": err,
		})
	}
	//	:: GET TOKEN pair for user

	newTokenPair, err := h.TokenService.GetPairForUser(ctx, user, token.ID.String())
	if err != nil {

		ctx.JSON(apprerrors.Status(err), gin.H{
			"error": err,
		})
		пше
	}
	//	:: TOKEN pair return
	ctx.JSON(http.StatusOK, gin.H{
		"tokens": newTokenPair,
	})
}
