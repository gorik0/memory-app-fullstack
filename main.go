package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//	::: ROUTER setup

	gi := gin.New()
	gi.GET("/api/account", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"ant creatures": "are undeerrrated"})
	})

	//	::: SERVER setup

	server := http.Server{
		Addr:    ":9000",
		Handler: gi,
	}
	//	::: SERVER run

	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error while runnong server :::: %s", err)
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	//	:::SERVER shutdown

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("While shut down serve :::: %s", err)
	}

}
