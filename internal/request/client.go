package request

import (
	"net/http"
	"net/url"
	"time"

	"go.uber.org/ratelimit"
)

var rl = ratelimit.New(10)

var defaultClient *Client = nil

func GetDefaultClient() *Client {
	if defaultClient == nil {
		defaultClient = NewClient(NewClientOptions{})
	}

	return defaultClient
}

func SetDefaultClient(client *Client) {
	defaultClient = client
}

type Client struct {
	*http.Client
	Header  http.Header
	Cookies []*http.Cookie
}

type NewClientOptions struct {
	Timeout   time.Duration
	RateLimit int // requests per second
	ProxyURL  *url.URL

	Header  http.Header
	Cookies []*http.Cookie
}

func NewClient(opts NewClientOptions) *Client {
	if opts.Timeout == 0 {
		opts.Timeout = 10 * time.Second
	}

	if opts.RateLimit > 0 {
		rl = ratelimit.New(opts.RateLimit)
	}

	if opts.Header == nil {
		opts.Header = http.Header{}
	}

	if opts.Cookies == nil {
		opts.Cookies = []*http.Cookie{}
	}

	var proxy func(*http.Request) (*url.URL, error)
	if opts.ProxyURL != nil && opts.ProxyURL.String() != "" {
		proxy = http.ProxyURL(opts.ProxyURL)
	} else {
		proxy = http.ProxyFromEnvironment
	}

	return &Client{
		&http.Client{
			Timeout: opts.Timeout,

			Transport: &http.Transport{
				Proxy: proxy,

				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
			},
		},
		opts.Header,
		opts.Cookies,
	}
}

func (c *Client) WithHeader(header http.Header) *Client {
	c.Header = header
	return c
}

func (c *Client) WithCookies(cookies []*http.Cookie) *Client {
	c.Cookies = cookies
	return c
}
