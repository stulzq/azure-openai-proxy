package util

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

func NewHttpProxy(proxyAddress string) (*http.Transport, error) {
	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy URL: %v", err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}

	if proxyURL.User != nil {
		proxyAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxyURL.User.String()))

		transport.ProxyConnectHeader = http.Header{
			"Proxy-Authorization": []string{proxyAuth},
		}
	}

	return transport, nil
}

func NewSocksProxy(proxyAddress string) (*http.Transport, error) {
	// proxyAddress: socks5://user:password@127.0.0.1:1080
	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		return nil, fmt.Errorf("error parsing proxy URL: %v", err)
	}

	dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("error creating proxy dialer: %v", err)
	}

	transport := &http.Transport{
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			return dialer.Dial(network, address)
		},
	}

	return transport, nil
}
