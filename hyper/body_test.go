package hyper

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	e "github.com/c0c0n3/resto/util/err"
)

func TestWriteBodyBytes(t *testing.T) {
	msg := newMsgWriter(t)
	want := []byte{4, 2}
	WriteBody(msg, want)

	got := msg.body
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func TestWriteStringBody(t *testing.T) {
	msg := newMsgWriter(t)
	WriteBody(msg, "42")

	want := []byte("42")
	got := msg.body
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func TestWriteBodyContentLengthError(t *testing.T) {
	msg := newMsgWriter(t)
	msg.headerWritingErr = UnexpectedResponseErr("boom")
	err := WriteBody(msg, "42")

	if _, ok := err.(e.Err[UnexpectedResponse]); !ok {
		t.Errorf("want: unexpected response err; got: %v", err)
	}
	msg.assertNoBody()
}

type MyData struct{ Greeting string }
type Unknown struct {
	X interface{}
}

func TestWriteJsonBody(t *testing.T) {
	msg := newMsgWriter(t)
	greeting := JSON{Data: &MyData{"howzit!"}}
	WriteBody(msg, greeting)

	want := `{"Greeting":"howzit!"}`
	got := msg.stringBody()
	if want != got {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestWriteNilJsonBody(t *testing.T) {
	msg := newMsgWriter(t)
	noData := JSON{Data: nil}
	WriteBody(msg, noData)

	want := `null`
	got := msg.stringBody()
	if want != got {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestWriteBodyJsonMarshalError(t *testing.T) {
	msg := newMsgWriter(t)
	notSerializable := JSON{Data: func() {}}
	err := WriteBody(msg, notSerializable)

	if _, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("want: json.UnsupportedTypeError; got: %v", err)
	}
	msg.assertNoBody()
}

func TestWriteUntypedJsonBody(t *testing.T) {
	msg := newMsgWriter(t)
	x := make(map[interface{}]interface{})
	x[1] = 2
	unknown := JSON{Data: &Unknown{X: x}}
	err := WriteBody(msg, unknown)

	if _, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("want: json.UnsupportedTypeError; got: %v", err)
	}
	msg.assertNoBody()
}

func TestReadJsonBody(t *testing.T) {
	msg := newMsgReader()
	msg.body = `{"Greeting":"howzit!"}`
	output := &MyData{}
	ReadJsonBody(msg, output)

	if output.Greeting != "howzit!" {
		t.Errorf("want: howzit!; got: %v", output)
	}
}

func TestReadNilJsonBody(t *testing.T) {
	msg := newMsgReader()
	msg.body = `null`
	output := &MyData{}
	err := ReadJsonBody(msg, output)

	if err != nil {
		t.Errorf("want: empty read; got: %v", err)
	}
	if output.Greeting != "" {
		t.Errorf("want: empty read got: %v", output)
	}
}

func TestReadJsonBodyNilOutputPtrError(t *testing.T) {
	msg := newMsgReader()
	msg.body = `{"Greeting":"howzit!"}`
	var output *MyData = nil
	err := ReadJsonBody(msg, output)

	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestReadBodyBuffer(t *testing.T) {
	msg := newMsgReader()
	want := "wada wada"
	msg.body = want
	var output bytes.Buffer
	ReadBody(msg, &output)

	got := output.String()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestReadBodyNilOutputPtrError(t *testing.T) {
	msg := newMsgReader()
	msg.body = "wada wada"
	var output *bytes.Buffer
	err := ReadBody(msg, output)

	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}
