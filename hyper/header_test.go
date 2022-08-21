package hyper

import (
	"testing"

	"github.com/c0c0n3/resto/mime"
	e "github.com/c0c0n3/resto/util/err"
)

func TestWriteContentType(t *testing.T) {
	msg := newMsgWriter(t)
	WriteContentType(msg, mime.JSON)
	msg.assertHeader("Content-Type", "application/json")
}

func TestWriteContentLength(t *testing.T) {
	msg := newMsgWriter(t)
	WriteContentLength(msg, 42)
	msg.assertHeader("Content-Length", "42")
}

func TestWriteAcceptWithNoParams(t *testing.T) {
	msg := newMsgWriter(t)
	WriteAccept(msg)
	msg.assertNoHeader("Accept")
}

func TestWriteAcceptWithParams(t *testing.T) {
	msg := newMsgWriter(t)
	WriteAccept(msg, mime.YAML, mime.JSON)
	msg.assertHeader("Accept", "application/yaml, application/json")
}

func TestWriteAuthorization(t *testing.T) {
	msg := newMsgWriter(t)
	creds := "foo bar"
	WriteAuthorization(msg, creds)
	msg.assertHeader("Authorization", creds)
}

func TestWriteBearerTokenNilProvider(t *testing.T) {
	msg := newMsgWriter(t)
	err := WriteBearerToken(msg, nil)
	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
	msg.assertNoHeader("Authorization")
}

func TestWriteBearerTokenProviderError(t *testing.T) {
	msg := newMsgWriter(t)
	provider := func() (string, error) {
		return "", UnexpectedResponseErr("boom")
	}
	err := WriteBearerToken(msg, provider)
	if _, ok := err.(e.Err[UnexpectedResponse]); !ok {
		t.Errorf("want: unexpected response err; got: %v", err)
	}
	msg.assertNoHeader("Authorization")
}

func TestWriteBearerToken(t *testing.T) {
	msg := newMsgWriter(t)
	provider := func() (string, error) {
		return "foo bar", nil
	}
	WriteBearerToken(msg, provider)
	msg.assertHeader("Authorization", "Bearer foo bar")
}
