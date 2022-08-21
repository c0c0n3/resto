package hyper

import (
	"github.com/c0c0n3/resto/util/err"
)

// A nil pointer error.
type NilPtr string

func NilMessageWriterErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil MessageWriter")
}

func NilBodySerializerErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil BodySerializer")
}

func NilBodyDeserializerErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil BodyDeserializer")
}

func NilRequestBuilderErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil RequestBuilder")
}

func NilResponseHandlerErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil ResponseHandler")
}

func NilBearerTokenProviderErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil BearerTokenProvider")
}

// An unexpected server response.
type UnexpectedResponse string

func UnexpectedResponseErr(format string, args ...any) err.Err[UnexpectedResponse] {
	return err.Mk[UnexpectedResponse](format, args...)
}
