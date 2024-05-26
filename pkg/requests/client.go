package requests

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

type Client struct {
}

func NewClient() HttpClient {
	return &Client{}
}

func (c *Client) Send(req *Request) ([]byte, error) {
	ips, err := net.LookupIP(req.url.Authorithy)
	if err != nil {
		return nil, err
	}
	if len(ips) <= 0 {
		return nil, errors.New("no ips where found")
	}
	log.Printf("ips: %+v", ips)
	conn, err := net.Dial("tcp", fmt.Sprintf("[%s]:80", ips[0].String()))
	if err != nil {
		return nil, err
	}
	log.Println("after Dial")

	fmt.Printf("raw request:\n%s", string(req.raw()))

	_, err = conn.Write(req.raw())
	if err != nil {
		return nil, err
	}
	log.Println("after Write")
	resp, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
