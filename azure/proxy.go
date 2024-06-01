package azure

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stulzq/azure-openai-proxy/util"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

const cognitiveservicesScope = "https://cognitiveservices.azure.com/.default"

func ProxyWithConverter(requestConverter RequestConverter) gin.HandlerFunc {
	return func(c *gin.Context) {
		Proxy(c, requestConverter)
	}
}

type DeploymentInfo struct {
	Data   []map[string]interface{} `json:"data"`
	Object string                   `json:"object"`
}

func ModelProxy(c *gin.Context) {
	// Create a channel to receive the results of each request
	results := make(chan []map[string]interface{}, len(ModelDeploymentConfig))

	// Send a request for each deployment in the map
	for _, deployment := range ModelDeploymentConfig {
		go func(deployment DeploymentConfig) {
			// Create the request
			req, err := http.NewRequest(http.MethodGet, deployment.Endpoint+"/openai/deployments?api-version=2022-12-01", nil)
			if err != nil {
				log.Printf("error parsing response body for deployment %s: %v", deployment.DeploymentName, err)
				results <- nil
				return
			}

			// Set the auth header
			req.Header.Set(APIKeyHeaderKey, deployment.ApiKey)

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("error sending request for deployment %s: %v", deployment.DeploymentName, err)
				results <- nil
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Printf("unexpected status code %d for deployment %s", resp.StatusCode, deployment.DeploymentName)
				results <- nil
				return
			}

			// Read the response body
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("error reading response body for deployment %s: %v", deployment.DeploymentName, err)
				results <- nil
				return
			}

			// Parse the response body as JSON
			var deploymentInfo DeploymentInfo
			err = json.Unmarshal(body, &deploymentInfo)
			if err != nil {
				log.Printf("error parsing response body for deployment %s: %v", deployment.DeploymentName, err)
				results <- nil
				return
			}
			results <- deploymentInfo.Data
		}(deployment)
	}

	// Wait for all requests to finish and collect the results
	var allResults []map[string]interface{}
	for i := 0; i < len(ModelDeploymentConfig); i++ {
		result := <-results
		if result != nil {
			allResults = append(allResults, result...)
		}
	}
	var info = DeploymentInfo{Data: allResults, Object: "list"}
	combinedResults, err := json.Marshal(info)
	if err != nil {
		log.Printf("error marshalling results: %v", err)
		util.SendError(c, err)
		return
	}

	// Set the response headers and body
	c.Header("Content-Type", "application/json")
	c.String(http.StatusOK, string(combinedResults))
}

// Proxy Azure OpenAI
func Proxy(c *gin.Context, requestConverter RequestConverter) {
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Status(200)
		return
	}

	// preserve request body for error logging
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	bodyBytes, err := io.ReadAll(tee)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return
	}
	c.Request.Body = io.NopCloser(&buf)

	director := func(req *http.Request) {
		if req.Body == nil {
			util.SendError(c, errors.New("request body is empty"))
			return
		}
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		// get model from url params or body
		model := c.Param("model")
		if model == "" {
			_model, err := sonic.Get(body, "model")
			if err != nil {
				util.SendError(c, fmt.Errorf("get model error: %w", err))
				return
			}
			_modelStr, err := _model.String()
			if err != nil {
				util.SendError(c, fmt.Errorf("get model name error: %w", err))
				return
			}
			model = _modelStr
		}

		// get deployment from request
		deployment, err := GetDeploymentByModel(model)
		if err != nil {
			util.SendError(c, err)
			return
		}

		// get auth token from header or deployment config
		token := deployment.ApiKey
		tokenFound := false
		if token == "" && token != "msi" {
			rawToken := req.Header.Get("Authorization")
			token = strings.TrimPrefix(rawToken, "Bearer ")
			req.Header.Set(APIKeyHeaderKey, token)
			req.Header.Del("Authorization")
			tokenFound = true
		}
		// get azure token using managed identity
		var azureToken azcore.AccessToken
		if token == "" || token == "msi" {
			cred, err := azidentity.NewManagedIdentityCredential(nil)
			if err != nil {
				util.SendError(c, fmt.Errorf("failed to create managed identity credential: %w", err))
			}

			azureToken, err = cred.GetToken(context.TODO(), policy.TokenRequestOptions{
				Scopes: []string{cognitiveservicesScope},
			})
			if err != nil {
				util.SendError(c, fmt.Errorf("failed to get token: %w", err))
			}

			req.Header.Del(APIKeyHeaderKey)
			req.Header.Set(AuthHeaderKey, "Bearer "+azureToken.Token)
			tokenFound = true
		}

		if !tokenFound {
			util.SendError(c, errors.New("token is empty"))
			return
		}

		originURL := req.URL.String()
		req, err = requestConverter.Convert(req, deployment)
		if err != nil {
			util.SendError(c, fmt.Errorf("convert request error: %w", err))
			return
		}
		log.Printf("proxying request [%s] %s -> %s", model, originURL, req.URL.String())
	}

	proxy := &httputil.ReverseProxy{Director: director}
	transport, err := util.NewProxyFromEnv()
	if err != nil {
		util.SendError(c, fmt.Errorf("get proxy error: %w", err))
		return
	}
	if transport != nil {
		proxy.Transport = transport
	}

	proxy.ServeHTTP(c.Writer, c.Request)

	// issue: https://github.com/Chanzhaoyu/chatgpt-web/issues/831
	if c.Writer.Header().Get("Content-Type") == "text/event-stream" {
		if _, err := c.Writer.Write([]byte{'\n'}); err != nil {
			log.Printf("rewrite response error: %v", err)
		}
	}

	if c.Writer.Status() != 200 {
		log.Printf("encountering error with body: %s", string(bodyBytes))
	}
}

func GetDeploymentByModel(model string) (*DeploymentConfig, error) {
	deploymentConfig, exist := ModelDeploymentConfig[model]
	if !exist {
		return nil, fmt.Errorf("deployment config for %s not found", model)
	}
	return &deploymentConfig, nil
}
