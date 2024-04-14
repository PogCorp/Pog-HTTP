package requests

type Request struct {
	method  string
	url     string
	body    []byte
	header  Header
	version string
}
