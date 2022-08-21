package wire

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/c0c0n3/resto/yoorel"
)

func getNowhere(req RequestWriter) error {
	url := yoorel.BuilderFrom("http://nowhere/").Build().Right()
	return req.RequestLine(GET, url)
}

func assertGetWithMock(t *testing.T, send Sender, wantCode StatusCode) {
	response, err := send(getNowhere)

	if err != nil {
		t.Fatalf("want: response; got: %v", err)
	}
	got, _ := response.StatusLine()
	if got != wantCode {
		t.Errorf("want: %d; got: %v", wantCode, got)
	}
}

func assertGetWithClient(t *testing.T, send Sender) {
	response, err := send(getNowhere)

	if response != nil {
		t.Errorf("want: err; got: %v", response)
	}
	if _, ok := err.(*url.Error); !ok {
		t.Errorf("want: url error; got: %v", err)
	}
}

func TestSenderWithRawClientFunc(t *testing.T) {
	want := 234
	client := func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: want}, nil
	}
	send := NewSender(client)

	assertGetWithMock(t, send, StatusCode(want))
}

func TestSenderWithStdLibSender(t *testing.T) {
	want := 234
	var client StdLibSender = func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: want}, nil
	}
	send := NewSender(client)

	assertGetWithMock(t, send, StatusCode(want))
}

func TestSenderWithDefaultClient(t *testing.T) {
	send := NewSender[DefaultClient]()
	assertGetWithClient(t, send)
}

func TestSenderWithCustomClient(t *testing.T) {
	client := &http.Client{Timeout: time.Second * 10}
	send := NewSender(client)
	assertGetWithClient(t, send)
}
