package client

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/c0c0n3/resto/hyper"
	"github.com/c0c0n3/resto/mime"
	"github.com/c0c0n3/resto/util/bytez"
	e "github.com/c0c0n3/resto/util/err"
)

type MyData struct {
	Greeting string
}

func TestGetJsonGreeting(t *testing.T) {
	dataToSendBack := `{"Greeting": "howzit!"}`
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 201,
			Status:     "OK",
			Body:       bytez.NewBufferFrom([]byte(dataToSendBack)),
		},
	}
	client := New(mock.Sender())

	output := &MyData{}
	err := client.Request(
		GET("https://my.api/data"),
		Accept(mime.JSON),
	).Handle(
		ExpectSuccess,
		ReadJsonResponse(output),
	)

	if err != nil {
		t.Errorf("want: server reply; got: %v", err)
	}
	if output.Greeting != "howzit!" {
		t.Errorf("want: howzit!; got: %v", output)
	}
}

func TestPostJsonGreeting(t *testing.T) {
	dataToPost := &MyData{Greeting: "howzit!"}
	dataToSendBack := "welcome, stranger"
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 201,
			Status:     "OK",
			Body:       bytez.NewBufferFrom([]byte(dataToSendBack)),
		},
	}
	client := New(mock.Sender())

	responseBody := &hyper.StringBody{}
	err := client.Request(
		POST("https://my.api/data"),
		ContentType(mime.JSON),
		Body(Json(dataToPost)),
	).Handle(
		ExpectStatusCodeOneOf(200, 201),
		ReadResponse(responseBody),
	)

	if err != nil {
		t.Errorf("want: server reply; got: %v", err)
	}
	if responseBody.Data != dataToSendBack {
		t.Errorf("want: %s; got: %s", dataToSendBack, responseBody.Data)
	}
	reqLen := mock.capturedReq.Header.Get("Content-Length")
	if reqLen != "22" {
		t.Errorf("want: 22; got: %s", reqLen)
	}
}

func TestPostForm(t *testing.T) {
	creds := func() (string, error) {
		return "t0k3n", nil
	}
	dataToPost := "field1=value1&field2=value2"
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 201,
			Status:     "OK",
		},
	}
	client := New(mock.Sender())

	err := client.Request(
		POST("https://my.api/data"),
		BearerToken(creds),
		ContentType(mime.URL_ENCODED),
		ContentLength(1234), // Body overrides it with the right value
		Body(dataToPost),
	).Handle(
		ExpectStatusCodeOneOf(200, 201),
	)

	if err != nil {
		t.Errorf("want: server reply; got: %v", err)
	}
	reqLen := mock.capturedReq.Header.Get("Content-Length")
	if reqLen != "27" {
		t.Errorf("want: 27; got: %s", reqLen)
	}
	bearer := mock.capturedReq.Header.Get("Authorization")
	if bearer != "Bearer t0k3n" {
		t.Errorf("want: 'Bearer t0k3n'; got: %s", bearer)
	}
}

func TestPutBytes(t *testing.T) {
	creds := "t0k3n"
	dataToPut := []byte{1, 2, 3, 4}
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 201,
			Status:     "OK",
		},
	}
	client := New(mock.Sender())

	err := client.Request(
		PUT("https://my.api/data"),
		Authorization(creds),
		ContentType(mime.OCTET_STREAM),
		Body(dataToPut),
	).Handle(
		ExpectStatusCodeOneOf(200, 201),
	)

	if err != nil {
		t.Errorf("want: server reply; got: %v", err)
	}
	reqLen := mock.capturedReq.Header.Get("Content-Length")
	if reqLen != "4" {
		t.Errorf("want: 4; got: %s", reqLen)
	}
	auth := mock.capturedReq.Header.Get("Authorization")
	if auth != creds {
		t.Errorf("want: %s; got: %s", creds, auth)
	}
}

func TestPostStream(t *testing.T) {
	dataToStream := Stream(bytez.NewBufferFrom([]byte{1, 2, 3, 4}))
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 201,
			Status:     "OK",
		},
	}
	client := New(mock.Sender())

	err := client.Request(
		POST("https://my.api/data"),
		ContentType(mime.OCTET_STREAM),
		Body(dataToStream),
	).Handle(
		ExpectStatusCodeOneOf(200, 201),
	)

	if err != nil {
		t.Errorf("want: server reply; got: %v", err)
	}
	reqLen := mock.capturedReq.Header.Get("Content-Length")
	if reqLen != "" {
		t.Errorf("want: no len; got: %s", reqLen)
	}
}

func TestPropagateRequestBuildingError(t *testing.T) {
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 200,
			Status:     "OK",
		},
	}
	client := New(mock.Sender())

	err := client.Request(
		GET("@://https://my.api/data"), // invalid URL
		Accept(mime.JSON),
	).Handle(
		ExpectSuccess,
	)

	if _, ok := err.(*url.Error); !ok {
		t.Errorf("want: unexpected response err; got: %v", err)
	}
}

func TestPropagateNilRequestBuilderError(t *testing.T) {
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 200,
			Status:     "OK",
		},
	}
	client := New(mock.Sender())

	err := client.Request(
		GET("https://my.api/data"),
		Accept(mime.JSON),
		nil,
	).Handle(
		ExpectSuccess,
	)

	if _, ok := err.(e.Err[hyper.NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestPropagateNilResponseHandlerError(t *testing.T) {
	mock := &mockClient{
		resToSend: &http.Response{
			StatusCode: 200,
			Status:     "OK",
		},
	}
	client := New(mock.Sender())

	err := client.Request(
		GET("https://my.api/data"),
		Accept(mime.JSON),
	).Handle(
		nil,
		ExpectSuccess,
	)

	if _, ok := err.(e.Err[hyper.NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestDnsLookupError(t *testing.T) {
	err := Request(
		GET("http://nohwere/"),
		Accept(mime.JSON),
	).Handle(
		ExpectSuccess,
	)
	if _, ok := err.(*url.Error); !ok {
		t.Errorf("want: url error; got: %v", err)
	}
}
