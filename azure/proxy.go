package azure

import (
	"bytes"
	"fmt"
	"github.com/stulzq/azure-openai-proxy/util"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ProxyWithConverter(requestConverter RequestConverter) gin.HandlerFunc {
	return func(c *gin.Context) {
		Proxy(c, requestConverter)
	}
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
				util.SendError(c, errors.Wrap(err, "get model error"))
				return
			}
			_modelStr, err := _model.String()
			if err != nil {
				util.SendError(c, errors.Wrap(err, "get model name error"))
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

		// get auth token from header or deployemnt config
		token := deployment.ApiKey
		if token == "" {
			rawToken := req.Header.Get("Authorization")
			token = strings.TrimPrefix(rawToken, "Bearer ")
		}
		if token == "" {
			util.SendError(c, errors.New("token is empty"))
			return
		}
		req.Header.Set(AuthHeaderKey, token)
		req.Header.Del("Authorization")

		originURL := req.URL.String()
		req, err = requestConverter.Convert(req, deployment)
		if err != nil {
			util.SendError(c, errors.Wrap(err, "convert request error"))
			return
		}
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

func GetDeploymentByModel(model string) (*DeploymentConfig, error) {
	deploymentConfig, exist := ModelDeploymentConfig[model]
	if !exist {
		return nil, errors.New(fmt.Sprintf("deployment config for %s not found", model))
	}
	return &deploymentConfig, nil
}
