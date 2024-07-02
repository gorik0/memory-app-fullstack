package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"net/http"
)

type DetailsRequest struct {
	Name    string `json:"name" binding:"omitempty,max=50"`
	Email   string `json:"email" binding:"omitempty,email"`
	Website string `json:"website" binding:"omitempty,url"`
}

func (h *Handler) Details(ctx *gin.Context) {

	userRequest := ctx.MustGet("user").(*models.User)

	var req DetailsRequest

	ok := BindData(ctx, &req)

	if !ok {
		return
	}
	fmt.Println("REQQQQ", req)
	user := models.User{
		UID:     userRequest.UID,
		Email:   req.Email,
		Name:    req.Name,
		Website: req.Website,
	}

	fmt.Println("USER", user)
	err := h.UserService.UpdateDetail(ctx, &user)
	if err != nil {
		ctx.JSON(apprerrors.Status(err), gin.H{
			"error": err,
		})
		return
	}
	fmt.Println("USER", user)

	ctx.JSON(http.StatusOK, gin.H{"user": user})

}
