package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"net/http"
)

type SigninRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,gte=5"`
}

func (h *Handler) Signin(ctx *gin.Context) {

	var req SigninRequest
	if ok := BindData(ctx, &req); !ok {
		log.Printf("eRRR while binding data")
		e := apprerrors.NewBadRequest("invalid field")
		ctx.JSON(e.Status(), gin.H{"error": e})
		return
	}

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
	err := h.UserService.Signin(ctx.Request.Context(), user)
	if err != nil {
		log.Printf("Coudln't signin ::: %s", err.Error())

		ctx.JSON(apprerrors.Status(err), gin.H{"error": err})
		return
	}
	tokens, err := h.TokenService.GetPairForUser(ctx.Request.Context(), user, "")
	if err != nil {
		log.Printf("Coudln't get token pair ::: %s", err.Error())

		ctx.JSON(apprerrors.Status(err), gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"tokens": tokens})

}
