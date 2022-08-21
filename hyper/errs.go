package hyper

import (
	"github.com/c0c0n3/resto/util/err"
)

// A nil pointer error.
type NilPtr string

func NilRequestBuilderErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil RequestBuilder")
}

func NilResponseHandlerErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil ResponseHandler")
}

func NilBearerTokenProviderErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil BearerTokenProvider")
}

func NilJsonOutDataErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil JSON output data structure")
}

func NilBytesBufferErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil bytes.Buffer")
}

// An unexpected server response.
type UnexpectedResponse string

func UnexpectedResponseErr(format string, args ...any) err.Err[UnexpectedResponse] {
	return err.Mk[UnexpectedResponse](format, args...)
}
