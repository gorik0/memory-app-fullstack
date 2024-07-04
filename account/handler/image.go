package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"net/http"
)

func (h *Handler) Image(ctx *gin.Context) {
	//:::: USER FORM CONTEXT
	user := ctx.MustGet("user").(*models.User)

	//::: SET MAX SIZE bytes
	ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, h.MaxBytesSize)

	//::: GET IMAGE
	file, err := ctx.FormFile("file")
	/*
		//::ERROR ON
		errors:
				- max size
				- simple internal
				- empty image

	*/

	if err != nil {

		fmt.Println("FormFile err:", err)
		if err.Error() == "http: request body too large" {
			ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{
				"err": fmt.Sprintf("Too large file for max size %v", h.MaxBytesSize),
			})
			return
		}
		e := apprerrors.NewBadRequest(fmt.Sprintf("Unable  to parse file"))
		ctx.JSON(e.Status(), gin.H{"err": e})
	}

	if file == nil {
		e := apprerrors.NewBadRequest("Must contain content!!!")
		ctx.JSON(e.Status(), gin.H{"err": e})

		return
	}

	//::: VALIDATE IMAGE TYPE (MIME)
	mimeType := file.Header.Get("Content-Type")
	if valid := validateImageMimeType(mimeType); !valid {

		e := apprerrors.NewBadRequest("Invalid contetnt type!!! Must be JPG or PNG")
		ctx.JSON(e.Status(), gin.H{"err": e})

		return

	}
	//:: UPDATE TO USER SERVICE

	context := ctx.Request.Context()
	image, err := h.UserService.SetProfileImage(context, user.UID, file)

	ctx.JSON(http.StatusOK, gin.H{"image": image})

}

func validateImageMimeType(mimeType string) bool {

	switch mimeType {
	case "image/jpeg":
		return true
	case "image/png":
		return true
	default:
		return false
	}
}
