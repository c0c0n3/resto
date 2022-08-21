package client

import (
	"net/http"
	"testing"

	"github.com/c0c0n3/resto/hyper/wire"
)

type mockClient struct {
	capturedReq *http.Request
	resToSend   *http.Response
	errToSend   error
}

func (p *mockClient) Send(req *http.Request) (*http.Response, error) {
	p.capturedReq = req
	if p.errToSend != nil {
		return &http.Response{}, p.errToSend
	}
	return p.resToSend, nil
}

func (p *mockClient) Sender() wire.Sender {
	return wire.NewSender(p.Send)
}

func TestNewClientFromNil(t *testing.T) {
	got := New(nil)
	if got == nil {
		t.Errorf("want: client; got: nil")
	}
}

func TestNewClientFromEmpty(t *testing.T) {
	senders := []wire.Sender{}
	got := New(senders...)
	if got == nil {
		t.Errorf("want: client; got: nil")
	}
	if New() == nil {
		t.Errorf("want: client; got: nil")
	}
}
