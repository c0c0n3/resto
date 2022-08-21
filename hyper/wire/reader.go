package wire

import (
	"io"
)

// Read the HTTP core message parts that both requests and responses
// share.
type MessageReader interface {
	// Read a message header.
	Header(name string) (content string)
	// Read all the headers.
	Headers() map[string][]string
	// Get a stream to read the message body.
	Body() io.ReadCloser
}

// Read an HTTP request.
type RequestReader interface {
	MessageReader
	// Read the request line.
	RequestLine() (verb Method, path string)
}

// Read an HTTP response.
type ResponseReader interface {
	MessageReader
	// Read the status line.
	StatusLine() (code StatusCode, reason string)
}
