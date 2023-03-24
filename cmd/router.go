package main

import (
	"github.com/gin-gonic/gin"
	"github.com/stulzq/azure-openai-proxy/apis"
)

func registerRoute(r *gin.Engine) {
	r.POST("/v1/chat/completions", apis.ChatCompletions)
}
