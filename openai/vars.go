package openai

import "github.com/imroc/req/v3"

const (
	AuthHeaderKey = "api-key"
)

var (
	AzureOpenAIEndpoint = ""
	AzureOpenAIAPIVer   = ""
	AzureOpenAIDeploy   = ""

	ChatCompletionsUrl = ""

	client = req.C()
)
