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

	if ok := BindData(ctx, &req); !ok {
		return
	}
	log.Println(req)
	user := &models.User{Email: req.Email, Password: req.Password}

	err := h.UserService.Signup(ctx, user)
	if err != nil {
		log.Printf("Fail to create user signup ::: %s \n", err)
		ctx.JSON(apprerrors.Status(err), gin.H{"error": err})
		return
	}

	token, err := h.TokenService.GetPairForUser(ctx, user, "")

	if err != nil {
		log.Printf("Fail to create token for user (%v) ::: %v \n", user, err.Error())
		ctx.JSON(apprerrors.Status(err), gin.H{"error": err})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"token": token})

}
