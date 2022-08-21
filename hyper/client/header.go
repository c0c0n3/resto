package client

import (
	"github.com/c0c0n3/resto/hyper"
	"github.com/c0c0n3/resto/hyper/wire"
	"github.com/c0c0n3/resto/mime"
)

// ContentType writes a "Content-Type" header with the specified MIME
// type.
func ContentType(mediaType mime.MediaType) wire.RequestBuilder {
	return func(msg wire.RequestWriter) error {
		return hyper.WriteContentType(msg, mediaType)
	}
}

// ContentLength writes a "Content-Length" header with the specified
// body size.
func ContentLength(bodySize uint64) wire.RequestBuilder {
	return func(msg wire.RequestWriter) error {
		return hyper.WriteContentLength(msg, bodySize)
	}
}

// Accept writes an "Accept" header with the specified MIME types.
func Accept(mediaType ...mime.MediaType) wire.RequestBuilder {
	return func(msg wire.RequestWriter) error {
		return hyper.WriteAccept(msg, mediaType...)
	}
}

// Authorization writes an "Authorization" header with the specified
// value.
func Authorization(value string) wire.RequestBuilder {
	return func(msg wire.RequestWriter) error {
		return hyper.WriteAuthorization(msg, value)
	}
}

// BearerToken writes an "Authorization" header with a value of
// "Bearer t" where `t` is the token value returned by the given
// BearerTokenProvider.
func BearerToken(acquireToken hyper.BearerTokenProvider) wire.RequestBuilder {
	return func(msg wire.RequestWriter) error {
		return hyper.WriteBearerToken(msg, acquireToken)
	}
}
