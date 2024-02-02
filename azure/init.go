package azure

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"github.com/stulzq/azure-openai-proxy/constant"
	"github.com/stulzq/azure-openai-proxy/util"
)

const (
	AuthHeaderKey = "api-key"
)

var (
	C                     Config
	ModelDeploymentConfig = map[string]DeploymentConfig{}
)

func Init() error {
	var (
		apiVersion        string
		endpoint          string
		openaiModelMapper string
		apikey            string
		err               error
	)

	apiVersion = viper.GetString(constant.ENV_AZURE_OPENAI_API_VER)
	endpoint = viper.GetString(constant.ENV_AZURE_OPENAI_ENDPOINT)
	openaiModelMapper = viper.GetString(constant.ENV_AZURE_OPENAI_MODEL_MAPPER)
	if endpoint != "" && openaiModelMapper != "" {
		if apiVersion == "" {
			apiVersion = "2023-07-01-preview"
		}
		apikey = strings.TrimSpace(viper.GetString(constant.ENV_AZURE_OPENAI_API_KEY))
		InitFromEnvironmentVariables(apiVersion, apikey, endpoint, openaiModelMapper)
	} else {
		if err = InitFromConfigFile(); err != nil {
			return err
		}
	}

	// ensure apiBase likes /v1
	viper.SetDefault("api_base", "/v1")
	apiBase := viper.GetString("api_base")
	if !strings.HasPrefix(apiBase, "/") {
		apiBase = "/" + apiBase
	}
	if strings.HasSuffix(apiBase, "/") {
		apiBase = apiBase[:len(apiBase)-1]
	}
	viper.Set("api_base", apiBase)
	log.Printf("apiBase is: %s", apiBase)
	for _, itemConfig := range C.DeploymentConfig {
		u, err := url.Parse(itemConfig.Endpoint)
		if err != nil {
			return fmt.Errorf("parse endpoint error: %w", err)
		}
		itemConfig.EndpointUrl = u
		ModelDeploymentConfig[itemConfig.ModelName] = itemConfig
	}
	return err
}

func InitFromEnvironmentVariables(apiVersion, apikey, endpoint, openaiModelMapper string) {
	log.Println("Init from environment variables")
	if openaiModelMapper != "" {
		// openaiModelMapper example:
		// gpt-3.5-turbo=deployment_name_for_gpt_model,text-davinci-003=deployment_name_for_davinci_model
		for _, pair := range strings.Split(openaiModelMapper, ",") {
			info := strings.Split(pair, "=")
			if len(info) != 2 {
				log.Fatalf("error parsing %s, invalid value %s", constant.ENV_AZURE_OPENAI_MODEL_MAPPER, pair)
			}
			modelName, deploymentName := info[0], info[1]
			u, err := url.Parse(endpoint)
			if err != nil {
				log.Fatalf("parse endpoint error: %s", err.Error())
			}
			ModelDeploymentConfig[modelName] = DeploymentConfig{
				DeploymentName: deploymentName,
				ModelName:      modelName,
				Endpoint:       endpoint,
				EndpointUrl:    u,
				ApiKey:         apikey,
				ApiVersion:     apiVersion,
			}
		}
	}
}

func InitFromConfigFile() error {
	log.Println("Init from config file")

	configFile := viper.GetString("configFile")
	if configFile == "" {
		configFile = filepath.Join(util.GetWorkdir(), "config.yaml")
	} else if !filepath.IsAbs(configFile) {
		configFile = filepath.Join(util.GetWorkdir(), configFile)
	}

	viper.SetConfigType("yaml")
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("read config file error: %+v\n", err)
		return err
	}

	if err := viper.Unmarshal(&C); err != nil {
		log.Printf("unmarshal config file error: %+v\n", err)
		return err
	}
	for _, configItem := range C.DeploymentConfig {
		ModelDeploymentConfig[configItem.ModelName] = configItem
	}

	log.Println("read config file success")
	return nil
}
