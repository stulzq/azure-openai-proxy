package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stulzq/azure-openai-proxy/azure"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var (
	version   = ""
	buildDate = ""
	gitCommit = ""
)

func main() {
	viper.AutomaticEnv()
	parseFlag()

	err := azure.Init()
	if err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	registerRoute(r)

	srv := &http.Server{
		Addr:    viper.GetString("listen"),
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Server Shutdown...")
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func parseFlag() {
	pflag.StringP("configFile", "c", "config.yaml", "config file")
	pflag.StringP("listen", "l", ":8080", "listen address")
	pflag.BoolP("version", "v", false, "version information")
	pflag.Parse()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}
	if viper.GetBool("version") {
		fmt.Println("version:", version)
		fmt.Println("buildDate:", buildDate)
		fmt.Println("gitCommit:", gitCommit)
		os.Exit(0)
	}
}
