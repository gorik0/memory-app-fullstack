package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
)

func (h *Handler) ImageDelete(ctx *gin.Context) {

	//:::: USER FORM CONTEXT
	user := ctx.MustGet("user").(*models.User)

	context := ctx.Request.Context()
	err := h.UserService.ClearProfileImage(context, user.UID)
	if err != nil {

		fmt.Println("fail to clear profile image with err: ", err)
		ctx.JSON(apprerrors.Status(err), gin.H{
			"err": err,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "success on delete image",
	})
}
