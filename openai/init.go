package openai

import (
	"fmt"
	"github.com/stulzq/azure-openai-proxy/constant"
	"log"
	"os"
)

func Init() {
	AzureOpenAIAPIVer = os.Getenv(constant.ENV_AZURE_OPENAI_API_VER)
	AzureOpenAIDeploy = os.Getenv(constant.ENV_AZURE_OPENAI_DEPLOY)
	AzureOpenAIEndpoint = os.Getenv(constant.ENV_AZURE_OPENAI_ENDPOINT)

	log.Println("AzureOpenAIAPIVer: ", AzureOpenAIAPIVer)
	log.Println("AzureOpenAIDeploy: ", AzureOpenAIDeploy)
	log.Println("AzureOpenAIEndpoint: ", AzureOpenAIEndpoint)

	ChatCompletionsUrl = fmt.Sprintf("%s/openai/deployments/%s/chat/completions?api-version=%s", AzureOpenAIEndpoint, AzureOpenAIDeploy, AzureOpenAIAPIVer)
}
