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

	// if viper get cors is true, then apply corsMiddleware
	if viper.GetBool("cors") {
		log.Printf("CORS supported! \n")
		r.Use(corsMiddleware())
	}

	registerRoute(r)

	srv := &http.Server{
		Addr:    viper.GetString("listen"),
		Handler: r,
	}

	runServer(srv)
}

// corsMiddleware sets up the CORS headers for all responses
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Clear any previously set headers
		if c.Request.Method != "POST" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Stainless-OS, X-STAINLESS-LANG, X-STAINLESS-PACKAGE-VERSION, X-STAINLESS-RUNTIME, X-STAINLESS-RUNTIME-VERSION, X-STAINLESS-ARCH")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
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
	pflag.BoolP("cors", "s", false, "cors support")
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
