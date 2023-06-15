package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/stulzq/azure-openai-proxy/azure"
)

// registerRoute registers all routes
func registerRoute(r *gin.Engine) {
	// https://platform.openai.com/docs/api-reference
	r.HEAD("/", func(c *gin.Context) {
		c.Status(200)
	})
	r.Any("/health", func(c *gin.Context) {
		c.Status(200)
	})
	apiBase := viper.GetString("api_base")
	stripPrefixConverter := azure.NewStripPrefixConverter(apiBase)
	templateConverter := azure.NewTemplateConverter("/openai/deployments/{{.DeploymentName}}/embeddings")
	apiBasedRouter := r.Group(apiBase)
	{
		apiBasedRouter.Any("/engines/:model/embeddings", azure.ProxyWithConverter(templateConverter))
		apiBasedRouter.Any("/completions", azure.ProxyWithConverter(stripPrefixConverter))
		apiBasedRouter.Any("/chat/completions", azure.ProxyWithConverter(stripPrefixConverter))
		apiBasedRouter.Any("/embeddings", azure.ProxyWithConverter(stripPrefixConverter))
	}
}
