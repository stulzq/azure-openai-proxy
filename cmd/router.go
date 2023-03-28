package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stulzq/azure-openai-proxy/azure"
)

// registerRoute registers all routes
func registerRoute(r *gin.Engine) {
	// https://platform.openai.com/docs/api-reference
	r.Any("*path", azure.Proxy)
}
