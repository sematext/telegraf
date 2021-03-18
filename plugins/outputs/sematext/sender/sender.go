package sender

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	// defaultConnectTimeout is the maximum amount of time a dial will wait for a connect to complete.
	defaultConnectTimeout = time.Second * 3

	// defaultTimeout specifies a time limit for requests made by this client
	defaultTimeout = time.Second * 10

	headerAgent       = "User-Agent"
	headerContentType = "Content-Type"
	headerProxyAuth   = "Proxy-Authorization"
)

// Config contains sender configuration
type Config struct {
	baseReceiverUrl string
	databaseName    string

	proxyURL *url.URL
	username string
	password string
}

// Sender is a simple wrapper around standard HTTP client
type Sender struct {
	client    *http.Client
	proxyAuth string
}

// NewSender constructs a new HTTP sender that should be used to send requests to Sematext
func NewSender(config *Config) *Sender {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: defaultConnectTimeout,
		}).DialContext,
		TLSHandshakeTimeout: defaultConnectTimeout,
	}

	if config.proxyURL != nil {
		transport.Proxy = http.ProxyURL(config.proxyURL)
	}

	c := &http.Client{
		Timeout:   defaultTimeout,
		Transport: transport,
	}

	return &Sender{client: c, proxyAuth: getProxyHeader(config)}
}

// getProxyHeader creates proxy authentication header based on config settings
func getProxyHeader(config *Config) string {
	var proxyAuth string
	if config.proxyURL != nil && len(config.username) > 0 && len(config.password) > 0 {
		auth := fmt.Sprintf("%s:%s", config.username, config.password)
		proxyAuth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
	}
	return proxyAuth
}

// Request emits an HTTP request targeting given HTTP method and URL.
func (c *Sender) Request(method, url, contentType string, body []byte) (*http.Response, error) {
	req, err := c.createRequest(method, url, contentType, body)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	res, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Sender) createRequest(method, url, contentType string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set(headerContentType, contentType)
	req.Header.Set(headerAgent, "telegraf")

	if len(c.proxyAuth) > 0 {
		req.Header.Add(headerProxyAuth, c.proxyAuth)
	}

	return req, nil
}
