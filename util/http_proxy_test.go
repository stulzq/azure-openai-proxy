package util

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHttpProxy(t *testing.T) {
	proxyAddress := "http://127.0.0.1:1087"
	transport, err := NewHttpProxy(proxyAddress)

	assert.NoError(t, err)
	assert.NotNil(t, transport)

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get("https://www.google.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestSocksProxy(t *testing.T) {
	proxyAddress := "socks5://127.0.0.1:1080"
	transport, err := NewSocksProxy(proxyAddress)

	assert.NoError(t, err)
	assert.NotNil(t, transport)

	client := &http.Client{
		Transport: transport,
	}

	resp, err := client.Get("https://www.google.com")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
}
