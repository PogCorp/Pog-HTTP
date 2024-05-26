package requests

type Version string

const (
	HTTP10 Version = "HTTP/1.0"
	HTTP11 Version = "HTTP/1.1"
	HTTP20 Version = "HTTP/2.0"
)

type Method string

const (
	GET     Method = "GET"
	POST    Method = "POST"
	PUT     Method = "PUT"
	DELETE  Method = "DELETE"
	OPTIONS Method = "OPTIONS"
)

type HttpClient interface {
	Send(req *Request) ([]byte, error)
}
