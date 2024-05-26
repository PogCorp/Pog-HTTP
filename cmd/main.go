package main

import (
	"fmt"
	"pog-http/pkg/requests"
)

func main() {
	req, err := requests.NewRequest(
		requests.GET, "http://example.com/",
		nil,
		requests.WithHeader(map[string]string{
			"Accept": "*/*",
		}),
	)
	if err != nil {
		fmt.Println("With New Request got error:", err)
		return
	}

	client := requests.NewClient()
	res, err := client.Send(req)
	if err != nil {
		fmt.Println("Got error while making request, got err:", err)
		return
	}
	fmt.Printf("Got response %s\n", string(res))
}
