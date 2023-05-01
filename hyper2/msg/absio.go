package msg

import "io"

// Abstract reading and writing of HTTP messages.
// See:
// - https://www.rfc-editor.org/rfc/rfc9110.html#name-message-abstraction

// Write an HTTP message. Can write both requests and responses. The
// type T is the type of the tagless final interpreter that produces
// the actual output.
// This interface doesn't force the implementation to buffer the whole
// message in memory before starting to write data on the wire. So a
// streaming implementation is possible too.
// (Notice unlike Go's standard lib, this interface doesn't force you
// to keep all the headers in memory.)
type Writer[T any] interface {
	// Write the request start line.
	WriteStartLine(line RequestControlData) T
	// Write the response status line.
	WriteStatusLine(line ResponseControlData) T
	// Write a header.
	WriteHeader(name, value string) T
	// Write the message body.
	WriteContent(body io.ReadCloser) T
	// Write a trailer.
	WriteTrailer(name, value string) T
}

type Reader[T any] interface {
	ReadStartLine() (RequestControlData, T)
	ReadStatusLine() (ResponseControlData, T)
	ReadHeaders() (Fields, T)
	ReadContent() (io.ReadCloser, T)
	ReadTrailers() (Fields, T)
}

type Builder[T any] func(msg Writer[T]) T

type Sender[W, R any] func(build Builder[W]) Reader[R]
type Receiver[R, W any] func(msg Reader[R]) Builder[W]
