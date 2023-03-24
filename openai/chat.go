package openai

import "github.com/imroc/req/v3"

func ChatCompletions(token string, body []byte) (*req.Response, error) {
	return client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader(AuthHeaderKey, token).
		SetBodyBytes(body).
		Post(ChatCompletionsUrl)
}
