package proxy

import (
	"net"
	"net/http"
	"net/url"
	"time"
)

// NewHTTPClient 创建一个 HTTP 客户端。
// proxyAddr 为代理地址（如 "127.0.0.1:6450"），为空时不使用代理。
func NewHTTPClient(proxyAddr string) *http.Client {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          200,
		MaxIdleConnsPerHost:   100,
		MaxConnsPerHost:       200,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if proxyAddr != "" {
		proxyURL, err := url.Parse("http://" + proxyAddr)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}
