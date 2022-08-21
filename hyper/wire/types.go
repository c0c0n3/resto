package wire

// One of the HTTP methods: GET, POST, etc.
type Method string

const (
	GET     = Method("GET")
	HEAD    = Method("HEAD")
	POST    = Method("POST")
	PUT     = Method("PUT")
	PATCH   = Method("PATCH")
	DELETE  = Method("DELETE")
	CONNECT = Method("CONNECT")
	OPTIONS = Method("OPTIONS")
	TRACE   = Method("TRACE")
)

func (m Method) String() string {
	return string(m)
}

// An HTTP status code, i.e. an integer between 100 and 599.
type StatusCode int

func (c StatusCode) Value() int {
	return int(c)
}
