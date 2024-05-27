package requests

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Request struct {
	method  Method
	url     *Url
	rawUrl  string
	body    []byte
	header  Header
	version Version
}

type RequestOpt func(*Request)

// Init
func defaultRequest() *Request {
	return &Request{
		method:  GET,
		rawUrl:  "",
		body:    nil,
		header:  make(Header),
		version: HTTP10,
	}
}

func NewRequest(method Method, url string, body []byte, opts ...RequestOpt) (*Request, error) {
	var err error
	r := defaultRequest()
	r.rawUrl = url
	r.method = method
	r.body = body
	r.url, err = NewUrl(url)
	if err != nil {
		return nil, err
	}
	if body != nil {
		r.header["Content-Length"] = fmt.Sprint(len(body))
	}
	for _, opt := range opts {
		opt(r)
	}
	return r, nil
}

func WithBasicAuth(auth string) RequestOpt {
	return func(r *Request) {
		r.BasicAuth(auth)
	}
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
	if body != nil {
		r.header["Content-Length"] = fmt.Sprint(len(body))
	}
}

func (r *Request) BasicAuth(auth string) {
	r.header["Authorization"] = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}

func (r *Request) raw() []byte {
	raw := strings.Builder{}
	raw.WriteString(fmt.Sprintf("%s %s %s\r\n", string(r.method), r.url.PathQuery, string(r.version)))
	raw.WriteString(fmt.Sprintf("Host: %s\r\n", r.url.Authorithy))
	for k, v := range r.header {
		header := fmt.Sprintf("%s: %s\r\n", k, v)
		raw.WriteString(header)
	}
	raw.WriteString("\r\n") // NOTE: Necessary format to separate headers from body
	raw.WriteString(string(r.body))
	return []byte(raw.String())
}
