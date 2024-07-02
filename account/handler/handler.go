package handler

import (
	"github.com/gin-gonic/gin"
	"memory-app/account/handler/middleware"
	"memory-app/account/models"
	"memory-app/account/models/apprerrors"
	"net/http"
	"time"
)

type Handler struct {
	UserService  models.UserServiceI
	TokenService models.TokenServiceI
}

func (h *Handler) Image(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{"say": "Image"})

}

type Config struct {
	Engine         *gin.Engine
	UserService    models.UserServiceI
	TokenServiceI  models.TokenServiceI
	BaseURL        string
	HandlerTimeout time.Duration
}

func NewHandler(c *Config) {
	h := Handler{
		UserService:  c.UserService,
		TokenService: c.TokenServiceI,
	}

	group := c.Engine.Group(c.BaseURL)
	if gin.Mode() != gin.TestMode {

		group.Use(middleware.Timeout(c.HandlerTimeout, apprerrors.NewTimedOut()))
		group.GET("/me", middleware.AuthUser(c.TokenServiceI), h.MeAbout)
		group.POST("/signout", middleware.AuthUser(c.TokenServiceI), h.Signout)
		group.PUT("/details", middleware.AuthUser(c.TokenServiceI), h.Details)

	} else {

		group.PUT("/details", h.Details)
		group.POST("/signout", h.Signout)
		group.GET("/me", h.MeAbout)
	}

	group.POST("/signin", h.Signin)
	group.POST("/signup", h.Signup)
	group.POST("/tokens", h.Tokens)
	group.POST("/image", h.Image)
	group.DELETE("/image", h.Image)

}
