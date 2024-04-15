package requests

import (
	"errors"
	"fmt"
	"log"
)

type Request struct {
	method  Method
	url     string
	body    []byte
	header  Header
	version Version
}

type RequestOpt func(*Request)

// Init
func defaultRequest() *Request {
	return &Request{
		method:  GET,
		url:     "",
		body:    nil,
		header:  make(Header),
		version: HTTP10,
	}
}

func NewRequest(method Method, url string, body []byte, opts ...RequestOpt) *Request {
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

// SetVersion
func WithVersion(ver Version) RequestOpt {
	return func(r *Request) {
		r.version = ver
	}
}

// alternative to AddHeader
func WithHeader(keyVal map[string]string) RequestOpt {
	return func(r *Request) {
		for k, v := range keyVal {
			err := r.AddHeader(k, v)
			if err != nil {
				log.Println("[Warn]: trying to insert Content-Length", err)
			}
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
