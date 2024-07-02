package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
)

func (h *Handler) Signout(ctx *gin.Context) {
	userGin := ctx.MustGet("user")

	context := ctx.Request.Context()
	user, ok := userGin.(*models.User)
	if !ok {
		fmt.Println("!!!! nooooo")
		return
	}
	err := h.TokenService.Signout(context, user.UID)
	if err != nil {

		ctx.JSON(apprerrors.Status(err), gin.H{

			"error": err,
		})
		return
	}
	ctx.JSON(200, gin.H{"mesaage": "Token signout succecc!!!"})
}
