package utils

import (
	"net/http"
	"net/url"
	"time"

	"github.com/bestruirui/mihomo-check/config"
	"golang.org/x/net/proxy"
)

// NewHTTPClient 根据配置创建并返回一个 HTTP 客户端
func NewHTTPClient() *http.Client {
	var client *http.Client

	if config.GlobalConfig.Proxy.Type == "http" {
		proxyURL, err := url.Parse(config.GlobalConfig.Proxy.Address)
		if err != nil {
			client = &http.Client{Timeout: 30 * time.Second}
		} else {
			transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
			client = &http.Client{Transport: transport, Timeout: 30 * time.Second}
		}
	} else if config.GlobalConfig.Proxy.Type == "socks" {
		socksDialer, err := proxy.SOCKS5("tcp", config.GlobalConfig.Proxy.Address, nil, proxy.Direct)
		if err != nil {
			client = &http.Client{Timeout: 30 * time.Second}
		} else {
			transport := &http.Transport{Dial: socksDialer.Dial}
			client = &http.Client{Transport: transport, Timeout: 30 * time.Second}
		}
	} else {
		client = &http.Client{Timeout: 30 * time.Second}
	}

	return client
}
