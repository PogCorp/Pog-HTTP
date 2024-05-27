package requests

import (
	"errors"
	"strings"
)

type Url struct {
	Authorithy string
	Scheme     string
	PathQuery  string
}

func NewUrl(u string) (*Url, error) {
	if hasControlCharacter(u) {
		return nil, errors.New("Invalid Control Character included in the url")
	}
	if u == "" {
		return nil, errors.New("Empty Url")
	}

	url := new(Url)
	scheme, rest, err := schemeExtractor(u)
	if err != nil {
		return nil, err
	}

	url.Scheme = scheme

	if len(rest) <= 0 || rest[0] != '/' || rest[1] != '/' {
		return nil, errors.New("authority was not provided")
	}
	authority, rest, err := authorityExtractor(rest[2:])
	if err != nil {
		return nil, err
	}
	url.Authorithy = authority
	url.PathQuery = rest

	return url, nil
}

func authorityExtractor(u string) (authority string, rest string, err error) {
	var i int
	if i = strings.Index(u, "/"); i == -1 {
		return "", "", errors.New("no authority found in url")
	}
	return u[:i], u[i:], nil
}

// (Scheme must match [a-zA-Z][a-zA-Z0-9+.-]*)
func schemeExtractor(u string) (scheme string, rest string, err error) {
	first := u[0]
	if first == ':' {
		return "", "", errors.New("No scheme was provided")
	} else if (first < 'a' || first > 'z') && (first < 'A' || first > 'Z') {
		return "", "", errors.New("First character of scheme is not alphabetical")
	}
	for i := 1; i < len(u); i++ {
		v := u[i]
		switch {
		case (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z'):
		case ('0' <= v && v <= '9') || v == '+' || v == '-' || v == '.':
			break
		case v == ':':
			return u[:i], u[i+1:], nil

		default:
			return "", "", errors.New("No scheme was provided")
		}
	}
	return
}

func hasControlCharacter(url string) bool {
	for _, v := range url {
		if v < 0x20 || v == 0x7f {
			return true
		}
	}
	return false
}
