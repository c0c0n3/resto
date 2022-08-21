package wire

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/c0c0n3/resto/util/bytez"
	"github.com/c0c0n3/resto/yoorel"
)

type echoMock struct {
	capturedRequest *http.Request
	errorToReturn   error
}

func (p *echoMock) send(req *http.Request) (*http.Response, error) {
	p.capturedRequest = req

	if p.errorToReturn != nil {
		return nil, p.errorToReturn
	}

	response := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header:     req.Header,
		Body:       req.Body,
	}
	return response, nil
}

func postRequest(req RequestWriter) error {
	url := yoorel.EmptyBuilder().Https().HostAndPort("httpbin.org").
		JoinPath("/post").Build().Right()
	err := req.RequestLine(POST, url)
	if err != nil {
		return err
	}
	err = req.Header("greeting", "howzit!")
	if err != nil {
		return err
	}
	content := bytez.Reader([]byte{42}) // ASCII 42 = '*'
	return req.Body(content)
}

func bogusRequest(req RequestWriter) error {
	url := yoorel.EmptyBuilder().Https().HostAndPort("httpbin.org").
		JoinPath("/get").Build().Right()
	bogusMethod := Method("I shouldn't be able to build this...")
	return req.RequestLine(bogusMethod, url)
}

func TestSuccessfulRequestReplyExchange(t *testing.T) {
	mock := &echoMock{}
	send := NewSender(mock.send)
	response, err := send(postRequest)
	request := mock.capturedRequest

	if err != nil {
		t.Fatalf("want: response; got: %v", err)
	}

	wantRequestUrl := "https://httpbin.org:443/post"
	wantHeaderContent := "howzit!"
	wantBodyContent := "*"
	wantCode := StatusCode(200)
	wantReason := "200 OK"

	if request.URL.String() != wantRequestUrl {
		t.Errorf("want: %s; got: %v", wantRequestUrl, request.URL)
	}
	if request.Header.Get("greeting") != wantHeaderContent {
		t.Errorf("want: %s; got: %s", wantHeaderContent, request.Header.Get("greeting"))
	}

	code, reason := response.StatusLine()
	if code != wantCode {
		t.Errorf("want: %v; got: %v", wantCode, code)
	}
	if reason != wantReason {
		t.Errorf("want: %s; got: %s", wantReason, reason)
	}

	gotHeader := response.Header("greeting")
	if gotHeader != wantHeaderContent {
		t.Errorf("want: %s; got: %s", wantHeaderContent, gotHeader)
	}
	if len(response.Headers()) != 1 {
		t.Errorf("want: 1 header; got: %d", len(response.Headers()))
	}
	if len(response.Headers()["Greeting"]) != 1 {
		t.Errorf("want: 1 header content; got: %v", response.Headers())
	}

	if body, err := io.ReadAll(response.Body()); err != nil {
		t.Errorf("want: body; got: %v", err)
	} else {
		if string(body) != wantBodyContent {
			t.Errorf("want: %s; got: %s", wantBodyContent, string(body))
		}
	}
}

func TestNetworkError(t *testing.T) {
	mock := &echoMock{
		errorToReturn: fmt.Errorf("net err"),
	}
	send := NewSender(mock.send)
	_, err := send(postRequest)
	request := mock.capturedRequest

	if err == nil {
		t.Fatalf("want: error; got: nil")
	}

	wantRequestUrl := "https://httpbin.org:443/post"
	wantHeaderContent := "howzit!"
	wantBodyContent := "*"

	if request.URL.String() != wantRequestUrl {
		t.Errorf("want: %s; got: %v", wantRequestUrl, request.URL)
	}
	if request.Header.Get("greeting") != wantHeaderContent {
		t.Errorf("want: %s; got: %s", wantHeaderContent, request.Header.Get("greeting"))
	}
	if body, err := io.ReadAll(request.Body); err != nil {
		t.Errorf("want: body; got: %v", err)
	} else {
		if string(body) != wantBodyContent {
			t.Errorf("want: %s; got: %s", wantBodyContent, string(body))
		}
	}
}

func TestFailToBuildHttpRequest(t *testing.T) {
	send := NewSender[DefaultClient]()
	_, err := send(bogusRequest)

	if err == nil {
		t.Fatalf("want: error; got: nil")
	}
	wantPrefix := "net/http: invalid method"
	got := err.Error()
	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("want prefix: %s; got: %s", wantPrefix, got)
	}
}

func TestRequestBuilderFail(t *testing.T) {
	send := NewSender[DefaultClient]()
	rogueBuilder := func(RequestWriter) error {
		return fmt.Errorf("boom!")
	}
	_, err := send(rogueBuilder)

	if err == nil {
		t.Fatalf("want: error; got: nil")
	}
	got := err.Error()
	if got != "boom!" {
		t.Errorf("want: boom!; got: %s", got)
	}
}
