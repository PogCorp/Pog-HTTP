package requests

import (
	"errors"
	"fmt"
)

type Request struct {
	method  Method
	url     string
	body    []byte
	header  Header
	version Version
}

type RequestOpt func(*Request)

func defaultRequest() *Request {
	return &Request{
		method:  GET,
		url:     "",
		body:    nil,
		header:  make(Header),
		version: HTTP10,
	}
}

func NewRequestBuilder(method Method, url string, body []byte, opts ...RequestOpt) *Request {
	r := defaultRequest()
	r.url = url
	r.method = method
	r.body = body
	r.header["Content-Length"] = string((body))
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func WithVersion(ver Version) RequestOpt {
	return func(r *Request) {
		r.version = ver
	}
}

func WithHeader(keyVal map[string]string) RequestOpt {
	return func(r *Request) {
		for k, v := range keyVal {
			r.header[k] = v
		}
	}
}

func (r *Request) AddHeader(key, val string) error {
	if key != "Content-Length" {
		r.header[key] = val
		return nil
	}
	return errors.New("Trying to write value into Content-Length")
}

func (r *Request) SetVersion(ver Version) {
	r.version = ver
}

func (r *Request) AddBody(body []byte) {
	r.body = body
	r.header["Content-Length"] = fmt.Sprint(len(body))
}
