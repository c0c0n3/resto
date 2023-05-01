package msg

// HTTP status codes.
// See:
// - https://www.rfc-editor.org/rfc/rfc9110.html#status.codes

// An HTTP status code, i.e. an integer between 100 and 599.
type StatusCode int

func (c StatusCode) Value() int {
	return int(c)
}
