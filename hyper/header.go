package hyper

import (
	"fmt"
	"strings"

	"github.com/c0c0n3/resto/hyper/wire"
	"github.com/c0c0n3/resto/mime"
)

// Write a "Content-Type" header with the specified MIME type.
func WriteContentType(msg wire.MessageWriter, mediaType mime.MediaType) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	return msg.Header("Content-Type", mediaType.String())
	// TODO implement charset & friends?
}

// Write a "Content-Length" header with the specified body size.
func WriteContentLength(msg wire.MessageWriter, bodySize uint64) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	sz := fmt.Sprintf("%d", bodySize)
	return msg.Header("Content-Length", sz)
}

// Write an "Accept" header with the specified MIME types.
func WriteAccept(msg wire.MessageWriter, mediaType ...mime.MediaType) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	ts := []string{}
	for _, mt := range mediaType {
		ts = append(ts, mt.String())
	}
	if len(ts) > 0 {
		return msg.Header("Accept", strings.Join(ts, ", "))
	}
	return nil
	// TODO implement weights too?
}

// Write an "Authorization" header with the specified value.
func WriteAuthorization(msg wire.MessageWriter, value string) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	return msg.Header("Authorization", value)
}

// BearerTokenProvider is a function that retrieves a Bearer token or
// an error if it couldn't retrieve it.
type BearerTokenProvider func() (string, error)

// Write an "Authorization" header with a value of "Bearer t" where `t`
// is the token value returned by the given BearerTokenProvider.
func WriteBearerToken(msg wire.MessageWriter, acquireToken BearerTokenProvider) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	if acquireToken == nil {
		return NilBearerTokenProviderErr()
	}
	if token, err := acquireToken(); err != nil {
		return err
	} else {
		authValue := fmt.Sprintf("Bearer %s", token)
		return WriteAuthorization(msg, authValue)
	}
}
