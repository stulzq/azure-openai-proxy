package azure

import (
	"github.com/stulzq/azure-openai-proxy/constant"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	AuthHeaderKey = "api-key"
)

var (
	AzureOpenAIEndpoint      = ""
	AzureOpenAIEndpointParse *url.URL

	AzureOpenAIAPIVer = ""

	AzureOpenAIModelMapper = map[string]string{
		"gpt-3.5-turbo": "gpt-35-turbo",
	}
	fallbackModelMapper = regexp.MustCompile(`[.:]`)
)

func Init() {
	AzureOpenAIAPIVer = os.Getenv(constant.ENV_AZURE_OPENAI_API_VER)
	AzureOpenAIEndpoint = os.Getenv(constant.ENV_AZURE_OPENAI_ENDPOINT)

	if AzureOpenAIAPIVer == "" {
		AzureOpenAIAPIVer = "2023-03-15-preview"
	}

	var err error
	AzureOpenAIEndpointParse, err = url.Parse(AzureOpenAIEndpoint)
	if err != nil {
		log.Fatal("parse AzureOpenAIEndpoint error: ", err)
	}

	if v := os.Getenv(constant.ENV_AZURE_OPENAI_MODEL_MAPPER); v != "" {
		for _, pair := range strings.Split(v, ",") {
			info := strings.Split(pair, "=")
			if len(info) != 2 {
				log.Fatalf("error parsing %s, invalid value %s", constant.ENV_AZURE_OPENAI_MODEL_MAPPER, pair)
			}

			AzureOpenAIModelMapper[info[0]] = info[1]
		}
	}

	log.Println("AzureOpenAIAPIVer: ", AzureOpenAIAPIVer)
	log.Println("AzureOpenAIEndpoint: ", AzureOpenAIEndpoint)
	log.Println("AzureOpenAIModelMapper: ", AzureOpenAIModelMapper)
}
