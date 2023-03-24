package apis

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stulzq/azure-openai-proxy/openai"
	"io"
	"strings"
)

// ChatCompletions xxx
// Path: /v1/chat/completions
func ChatCompletions(c *gin.Context) {
	// get auth token from header
	rawToken := c.GetHeader("Authorization")
	token := strings.TrimPrefix(rawToken, "Bearer ")

	reqContent, err := io.ReadAll(c.Request.Body)
	if err != nil {
		SendError(c, errors.Wrap(err, "failed to read request body"))
		return
	}

	oaiResp, err := openai.ChatCompletions(token, reqContent)
	if err != nil {
		SendError(c, errors.Wrap(err, "failed to call Azure OpenAI"))
		return
	}

	// pass-through header
	extraHeaders := map[string]string{}
	for k, v := range oaiResp.Header {
		if _, ok := ignoreHeaders[k]; ok {
			continue
		}

		extraHeaders[k] = strings.Join(v, ",")
	}

	c.DataFromReader(oaiResp.StatusCode, oaiResp.ContentLength, oaiResp.Header.Get("Content-Type"), oaiResp.Response.Body, extraHeaders)

	_, _ = c.Writer.Write([]byte{'\n'}) // add a newline to the end of the response https://github.com/Chanzhaoyu/chatgpt-web/issues/831
}
