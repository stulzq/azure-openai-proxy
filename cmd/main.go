package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	registerRoute(r)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	runServer(srv)
}

func runServer(srv *http.Server) {
	go func() {
		log.Printf("Server listening at %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(errors.Errorf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server Shutdown...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
