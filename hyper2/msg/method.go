package msg

// HTTP methods.
// See:
// - https://www.rfc-editor.org/rfc/rfc9110.html#name-methods

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
