package client

import (
	"net/http"
	"testing"

	"github.com/c0c0n3/resto/hyper"
	e "github.com/c0c0n3/resto/util/err"
)

func TestExpectSuccess(t *testing.T) {
	mock := &mockClient{
		resToSend: &http.Response{},
	}
	client := New(mock.Sender())

	for code := 200; code < 300; code++ {
		mock.resToSend.StatusCode = code

		err := client.Request(
			DELETE("https://my.api/data"),
		).Handle(
			ExpectSuccess,
		)

		if err != nil {
			t.Errorf("want: server reply (%d); got: %v", code, err)
		}
	}
	for _, code := range []int{100, 199, 300, 400, 500} {
		mock.resToSend.StatusCode = code

		err := client.Request(
			PATCH("https://my.api/data"),
		).Handle(
			ExpectSuccess,
		)

		if _, ok := err.(e.Err[hyper.UnexpectedResponse]); !ok {
			t.Errorf("want: unexpected response err (%d); got: %v", code, err)
		}
	}
}

func TestExpectStatusCodeOneOf(t *testing.T) {
	mock := &mockClient{
		resToSend: &http.Response{},
	}
	client := New(mock.Sender())

	want := []int{200, 201, 404}
	for _, code := range want {
		mock.resToSend.StatusCode = code

		err := client.Request(
			PUT("https://my.api/data"),
		).Handle(
			ExpectStatusCodeOneOf(want...),
		)

		if err != nil {
			t.Errorf("want: server reply (%d); got: %v", code, err)
		}
	}
	for _, code := range []int{100, 199, 300, 400, 500} {
		mock.resToSend.StatusCode = code

		err := client.Request(
			PUT("https://my.api/data"),
		).Handle(
			ExpectStatusCodeOneOf(want...),
		)

		if _, ok := err.(e.Err[hyper.UnexpectedResponse]); !ok {
			t.Errorf("want: unexpected response err (%d); got: %v", code, err)
		}
	}
}

func TestExpectStatusCodeNone(t *testing.T) {
	mock := &mockClient{
		resToSend: &http.Response{StatusCode: 200},
	}
	client := New(mock.Sender())

	err := client.Request(
		PUT("https://my.api/data"),
	).Handle(
		ExpectStatusCodeOneOf(),
	)

	if _, ok := err.(e.Err[hyper.UnexpectedResponse]); !ok {
		t.Errorf("want: unexpected response err; got: %v", err)
	}
}
