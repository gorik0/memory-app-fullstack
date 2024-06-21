package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"net/http"
)

type SignUpRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6,lte=20"`
}

func (h *Handler) Signup(ctx *gin.Context) {

	var req SignUpRequest

	if ok := BindData(ctx, req); !ok {
		return
	}

	user := &models.User{Email: req.Email, Password: req.Password}

	err := h.UserService.Signup(ctx, user)
	if err != nil {
		log.Printf("Fail to craete user signup ::: %s \n", err)
		ctx.JSON(apprerrors.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ok": true})

}
