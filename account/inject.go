package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"memory-app/account/handler"
	"memory-app/account/models/repository"
	"memory-app/account/service"
	"os"
	"strconv"
	"time"
)

// ::: USER SERVICE / TOKEN SEWRVICE / USER REPO
func inject(ds *dataSources) (*gin.Engine, error) {
	//::: USER  SERVICE
	userRepo := repository.PgUSerRepo{
		DB: ds.DB,
	}

	imageRepo := repository.ImageRepository{PathToFile: "./images/"}
	userService := service.NewUserService(&service.UserServiceConfig{
		UserRepo:  userRepo,
		ImageRepo: &imageRepo,
	})
	//::: TOKEN  SERVICE

	refreshExp := os.Getenv("ID_TOKEN_EXPIRED")
	idExp := os.Getenv("REFRESH_TOKEN_EXPIRED")

	privateString, err := os.ReadFile(os.Getenv("PRIVATE_KEY_FILE"))
	if err != nil {
		return nil, fmt.Errorf("while reading private key file ::: %w", err)
	}
	private := privateString
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	if err != nil {
		return nil, fmt.Errorf("while parsing private key  ::: %w", err)

	}

	publicString, err := os.ReadFile(os.Getenv("PUBLIC_KEY_FILE"))
	if err != nil {
		return nil, fmt.Errorf("whil reading  public key fiule  ::: %w", err)

	}
	public := publicString
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(public)
	if err != nil {
		return nil, fmt.Errorf("while parsing  public key   ::: %w", err)

	}
	secret := os.Getenv("REFRESH_SECRET")

	tokenRepo := repository.NewTokenRepo(ds.Redis)
	tokenService := service.NewTokenService(&service.ConfigTokenService{
		PrivKey:         privKey,
		PublKey:         publicKey,
		RefreshSecret:   secret,
		RefreshExp:      refreshExp,
		IdExp:           idExp,
		TokenRepository: tokenRepo,
	})

	//::: USER  SERVICE
	gi := gin.Default()

	baseUrl := os.Getenv("ACCOUNT_API_URL")
	handlerTimeoutString := os.Getenv("HANDLER_TIMEOUT")
	handlerTimeout, err := strconv.Atoi(handlerTimeoutString)
	if err != nil {
		return nil, fmt.Errorf("parsing timeout  public key fiule  ::: %w", err)
	}
	handler.NewHandler(&handler.Config{
		Engine:         gi,
		UserService:    userService,
		TokenServiceI:  tokenService,
		BaseURL:        baseUrl,
		HandlerTimeout: time.Duration(handlerTimeout) * time.Second,
		MaxBytesSize:   10 * 1024 * 1024,
	})

	return gi, nil

}
