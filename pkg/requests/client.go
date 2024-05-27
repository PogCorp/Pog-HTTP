package requests

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
}

func NewClient() HttpClient {
	return &Client{}
}

func (c *Client) Send(req *Request) ([]byte, error) {
	var ip string
	var err error
	if i := strings.Index(req.url.Authorithy, ":"); i < 0 {
		ip, err = c.extractIp(req.url.Authorithy, "")
		if err != nil {
			return nil, err
		}
	} else {
		ip, err = c.extractIp(req.url.Authorithy[:i], req.url.Authorithy[i+1:])
		if err != nil {
			return nil, err
		}
	}
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		return nil, err
	}

	log.Printf("Request:\n%s\n", string(req.raw()))

	_, err = conn.Write(req.raw())
	if err != nil {
		return nil, err
	}
	resp, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) extractIp(host string, port string) (string, error) {
	ips, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}
	if len(ips) <= 0 {
		return "", errors.New("no ips where found")
	}
	if port == "" {
		return fmt.Sprintf("[%s]:%s", ips[0].String(), "80"), nil
	}
	return fmt.Sprintf("[%s]:%s", ips[0].String(), port), nil

}
