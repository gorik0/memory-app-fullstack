package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"net/http"
)

func (h *Handler) MeAbout(ctx *gin.Context) {

	user, exists := ctx.Get("user")
	if !exists {
		err := apprerrors.NewInternal()
		log.Printf("User not found for context")
		ctx.JSON(err.Status(), gin.H{
			"error": err,
		})
		return
	}

	uid := user.(*models.User).UID

	ginCtx := ctx.Request.Context()
	user, err := h.UserService.Get(ginCtx, uid)
	if err != nil {
		log.Printf("Coudln't get user for id ::: %v", uid.String())

		e := apprerrors.NewNotFound("user", uid.String())
		ctx.JSON(e.Status(), gin.H{
			"error": e,
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{"user": user})

}
