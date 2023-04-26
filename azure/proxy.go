package azure

import (
	"bytes"
	"fmt"
	"github.com/stulzq/azure-openai-proxy/util"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"path"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

// Proxy Azure OpenAI
func Proxy(c *gin.Context) {
	if c.Request.Method == http.MethodOptions {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Status(200)
		return
	}

	director := func(req *http.Request) {
		if req.Body == nil {
			util.SendError(c, errors.New("request body is empty"))
			return
		}
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(body))

		// get model from body
		model, err := sonic.Get(body, "model")
		if err != nil {
			util.SendError(c, errors.Wrap(err, "get model error"))
			return
		}

		// get deployment from request
		deployment, err := model.String()
		if err != nil {
			util.SendError(c, errors.Wrap(err, "get deployment error"))
			return
		}
		deployment = GetDeploymentByModel(deployment)

		// get auth token from header
		rawToken := req.Header.Get("Authorization")
		token := strings.TrimPrefix(rawToken, "Bearer ")
		req.Header.Set(AuthHeaderKey, token)
		req.Header.Del("Authorization")

		originURL := req.URL.String()
		req.Host = AzureOpenAIEndpointParse.Host
		req.URL.Scheme = AzureOpenAIEndpointParse.Scheme
		req.URL.Host = AzureOpenAIEndpointParse.Host
		req.URL.Path = path.Join(fmt.Sprintf("/openai/deployments/%s", deployment), strings.Replace(req.URL.Path, "/v1/", "/", 1))
		req.URL.RawPath = req.URL.EscapedPath()

		query := req.URL.Query()
		query.Add("api-version", AzureOpenAIAPIVer)
		req.URL.RawQuery = query.Encode()

		log.Printf("proxying request [%s] %s -> %s", model, originURL, req.URL.String())
	}

	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)

	// issue: https://github.com/Chanzhaoyu/chatgpt-web/issues/831
	if c.Writer.Header().Get("Content-Type") == "text/event-stream" {
		if _, err := c.Writer.Write([]byte{'\n'}); err != nil {
			log.Printf("rewrite response error: %v", err)
		}
	}
}

func GetDeploymentByModel(model string) string {
	if v, ok := AzureOpenAIModelMapper[model]; ok {
		return v
	}

	return fallbackModelMapper.ReplaceAllString(model, "")
}
