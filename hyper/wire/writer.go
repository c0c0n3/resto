package wire

import (
	"io"

	"github.com/c0c0n3/resto/yoorel"
)

// Write the HTTP core message parts that both requests and responses
// share. This interface doesn't force the implementation to buffer the
// whole message in memory before starting to write data on the wire.
// So a streaming implementation is possible too.
// (Notice unlike Go's standard lib, this interface doesn't force you
// to keep all the headers in memory.)
type MessageWriter interface {
	// Write a message header.
	Header(name string, content string) error
	// Write the message body.
	Body(content io.ReadCloser) error
}

// Write an HTTP request.
type RequestWriter interface {
	MessageWriter
	// Write the request line.
	RequestLine(verb Method, resource yoorel.HttpUrl) error
}

// Write an HTTP response.
type ResponseWriter interface {
	MessageWriter
	// Write the status line.
	StatusLine(code StatusCode, reason string) error
}
