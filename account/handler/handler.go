package handler

import (
	"github.com/gin-gonic/gin"
	"memory-app/account/models"
	"net/http"
	"os"
)

type Handler struct {
	UserService models.UserServiceI
}

func (h *Handler) Signin(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"say": "Signin"})

}

func (h *Handler) Signout(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"say": "Signout"})

}

func (h *Handler) Tokens(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"say": "Tokens"})

}

func (h *Handler) Image(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"say": "Image"})

}

func (h *Handler) Details(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"say": "Details"})

}

type Config struct {
	Engine      *gin.Engine
	UserService models.UserServiceI
}

func NewHandler(c *Config) {
	h := Handler{
		UserService: c.UserService,
	}

	group := c.Engine.Group(os.Getenv("ACCOUNT_API_URL"))
	group.GET("/me", h.MeAbout)
	group.POST("/signin", h.Signin)
	group.POST("/signout", h.Signout)
	group.POST("/signup", h.Signup)
	group.POST("/tokens", h.Tokens)
	group.POST("/image", h.Image)
	group.DELETE("/image", h.Image)
	group.PUT("/details", h.Details)

}
