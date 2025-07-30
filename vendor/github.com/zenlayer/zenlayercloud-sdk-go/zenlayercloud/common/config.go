package common

import (
	"net/http"
	"time"
)

type Config struct {
	AutoRetry     bool `default:"false"`
	MaxRetryTime  int  `default:"3"`
	RetryDuration time.Duration
	HttpTransport *http.Transport   `default:""`
	Transport     http.RoundTripper `default:""`
	Proxy         string            `default:""`
	Scheme        string            `default:"HTTPS"`
	Domain        string            `default:""`
	// timeout in seconds
	Timeout int `default:"300"`
	Debug   *bool
}

func NewConfig() (config *Config) {
	config = &Config{}
	InitStructWithDefaultTag(config)
	return
}

func (c *Config) WithAutoRetry(isAutoRetry bool) *Config {
	c.AutoRetry = isAutoRetry
	return c
}

func (c *Config) WithMaxRetryTime(maxRetryTime int) *Config {
	c.MaxRetryTime = maxRetryTime
	return c
}

func (c *Config) WithTimeout(timeout int) *Config {
	c.Timeout = timeout
	return c
}

func (c *Config) WithHttpTransport(httpTransport *http.Transport) *Config {
	c.HttpTransport = httpTransport
	return c
}

func (c *Config) WithScheme(scheme string) *Config {
	c.Scheme = scheme
	return c
}
