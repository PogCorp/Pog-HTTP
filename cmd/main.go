package main

import (
	"fmt"
	"net/url"
	"strings"
)

type Header map[string]string

type RequestBuilder struct {
	method  string
	url     string
	body    []byte
	header  Header
	version string
}

type RequestBuilderOpt func(*RequestBuilder)

func defaultRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		method:  "GET",
		url:     "",
		body:    nil,
		header:  make(map[string]string),
		version: "HTTP/1.1",
	}
}

func NewRequestBuilder(method, url string, opts ...RequestBuilderOpt) *RequestBuilder {
	rb := defaultRequestBuilder()
	rb.url = url
	rb.method = method
	for _, opt := range opts {
		opt(rb)
	}
	return rb
}

// FIX: no error handling
func WithVersion(ver string) RequestBuilderOpt {
	return func(rb *RequestBuilder) {
		rb.version = ver
	}
}

// FIX: no error handling
func WithBody(body []byte) RequestBuilderOpt {
	return func(rb *RequestBuilder) {
		rb.body = body
	}
}

// FIX: no error handling
func WithHeader(keyVal map[string]string) RequestBuilderOpt {
	return func(rb *RequestBuilder) {
		for k, v := range keyVal {
			rb.header[k] = v
		}
	}
}

// FIX: no error handling
func (rb *RequestBuilder) AddHeader(key, val string) {
	rb.header[key] = val
}

func getHostFromURL(url *url.URL) (string, error) {
	baseHost := url.Scheme + "://" + url.Hostname()
	if url.Port() == "" {
		return baseHost, nil
	}

	return baseHost + ":" + url.Port(), nil
}

func (rb *RequestBuilder) raw() (string, error) {
	var request strings.Builder

	parsedURL, err := url.Parse(rb.url)
	if err != nil {
		return "", err
	}

	host, err := getHostFromURL(parsedURL)
	if err != nil {
		return "", err
	}

	// General
	request.WriteString(fmt.Sprintf("%s %s %s\r\n", rb.method, parsedURL.Path, rb.version))

	// Headers
	request.WriteString(fmt.Sprintf("Host: %s\r\n", host))
	for key, value := range rb.header {
		request.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}

	// Content-Length header
	if len(rb.body) > 0 {
		request.WriteString(fmt.Sprintf("Content-Length: %d\r\n", len(rb.body)))
	}

	// Empty line to separate headers from body
	request.WriteString("\r\n")

	// Body
	if len(rb.body) > 0 {
		request.Write(rb.body)
	}

	return request.String(), nil
}

func (rb *RequestBuilder) Send() error {
	payload, err := rb.raw()
	if err != nil {
		return err
	}

	fmt.Println(payload)

	return nil
}

func main() {
	rb := NewRequestBuilder(
		"POST", "https://example.com:42069/",
		WithBody([]byte("Hello, World!")),
		WithHeader(map[string]string{
			"Content-Type":    "text/plain",
			"Accept-Encoding": "gzip",
		}))

	if err := rb.Send(); err != nil {
		fmt.Println("In Raw Request got error:", err)
		return
	}
}
