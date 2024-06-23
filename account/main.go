package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	//	::: DATA SOURCE setup

	ds, err := initDS()
	if err != nil {
		log.Fatalf("dta source init ::: %s", err.Error())
	}
	//:::INJECTING DATA SOURCE / ROUTER setup
	gi, err := inject(ds)
	if err != nil {
		log.Fatalf("injecting data source  ::: %s", err.Error())

	}
	//	::: SERVER setup

	server := http.Server{
		Addr:    ":9000",
		Handler: gi,
	}
	//	::: SERVER run

	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Error while running server :::: %s", err)
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
