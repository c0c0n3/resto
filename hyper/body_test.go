package hyper

import (
	"encoding/json"
	"io"
	"reflect"
	"testing"

	e "github.com/c0c0n3/resto/util/err"
)

var writeBodyBytesFixtures = [][]byte{
	{}, {4}, {4, 2},
}

func TestWriteBodyBytes(t *testing.T) {
	for _, want := range writeBodyBytesFixtures {
		msg := newMsgWriter(t)
		WriteBody(msg, &ByteBody{want})

		got := msg.body
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v; got: %v", want, got)
		}
	}
}

var readBodyBytesFixtures = []string{
	"", "1", "12345678",
}

func TestReadBodyBytes(t *testing.T) {
	for _, want := range readBodyBytesFixtures {
		msg := newMsgReader()
		msg.body = want
		output := &ByteBody{}
		ReadBody(msg, output)

		got := string(output.Data)
		if want != got {
			t.Errorf("want: %s; got: %s", want, got)
		}
	}
}

var writeStringBodyFixtures = []string{
	"", "4", "42", "wada wada",
}

func TestWriteStringBody(t *testing.T) {
	for _, want := range writeStringBodyFixtures {
		msg := newMsgWriter(t)
		WriteBody(msg, &StringBody{want})

		wantBytes := []byte(want)
		got := msg.body
		if !reflect.DeepEqual(wantBytes, got) {
			t.Errorf("want: %v; got: %v", wantBytes, got)
		}
	}
}

var readStringBodyFixtures = []string{
	"", "1", "12345678",
}

func TestReadStringBody(t *testing.T) {
	for _, want := range readStringBodyFixtures {
		msg := newMsgReader()
		msg.body = want
		output := &StringBody{}
		ReadBody(msg, output)

		got := output.Data
		if want != got {
			t.Errorf("want: %s; got: %s", want, got)
		}
	}
}

var writeStreamingBodyFixtures = [][]byte{
	{}, {1}, {1, 2, 3, 4, 5, 6, 7, 8},
}

func TestWriteStreamingBody(t *testing.T) {
	for _, want := range writeStreamingBodyFixtures {
		msg := newMsgWriter(t)
		WriteBody(msg, &StreamingBody{newByteStreamer(want)})

		got := msg.body
		if !reflect.DeepEqual(want, got) {
			t.Errorf("want: %v; got: %v", want, got)
		}
	}
}

var readStreamingBodyFixtures = []string{
	"", "1", "12345678",
}

func TestReadStreamingBody(t *testing.T) {
	for _, want := range readStreamingBodyFixtures {
		msg := newMsgReader()
		msg.body = want
		output := &StreamingBody{}
		ReadBody(msg, output)

		gotBytes, _ := io.ReadAll(output.Data)
		got := string(gotBytes)
		if want != got {
			t.Errorf("want: %s; got: %s", want, got)
		}
	}
}

type MyData struct{ Greeting string }
type Unknown struct {
	X interface{}
}

func TestWriteJsonBody(t *testing.T) {
	msg := newMsgWriter(t)
	greeting := &JsonBody{&MyData{"howzit!"}}
	WriteBody(msg, greeting)

	want := `{"Greeting":"howzit!"}`
	got := msg.stringBody()
	if want != got {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestWriteNilJsonBody(t *testing.T) {
	msg := newMsgWriter(t)
	noData := &JsonBody{Data: nil}
	WriteBody(msg, noData)

	want := `null`
	got := msg.stringBody()
	if want != got {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestWriteBodyJsonMarshalError(t *testing.T) {
	msg := newMsgWriter(t)
	notSerializable := &JsonBody{Data: func() {}}
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
	unknown := &JsonBody{Data: &Unknown{X: x}}
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
	ReadBody(msg, &JsonBody{output})

	if output.Greeting != "howzit!" {
		t.Errorf("want: howzit!; got: %v", output)
	}
}

func TestReadNilJsonBody(t *testing.T) {
	msg := newMsgReader()
	msg.body = `null`
	output := &MyData{}
	err := ReadBody(msg, &JsonBody{output})

	if err != nil {
		t.Errorf("want: empty read; got: %v", err)
	}
	if output.Greeting != "" {
		t.Errorf("want: empty read got: %v", output)
	}
}

func TestWriteBodyContentLengthError(t *testing.T) {
	msg := newMsgWriter(t)
	msg.headerWritingErr = UnexpectedResponseErr("boom")
	err := WriteBody(msg, &StringBody{"42"})

	if _, ok := err.(e.Err[UnexpectedResponse]); !ok {
		t.Errorf("want: unexpected response err; got: %v", err)
	}
	msg.assertNoBody()
}

func TestWriteBodyNilWriterError(t *testing.T) {
	err := WriteBody(nil, &StringBody{"42"})

	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestWriteBodyNilSerializerError(t *testing.T) {
	msg := newMsgWriter(t)
	err := WriteBody(msg, nil)

	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestReadBodyWithNilReaderError(t *testing.T) {
	err := ReadBody(nil, &StringBody{})

	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestReadBodyWithNilDeserializerError(t *testing.T) {
	msg := newMsgReader()
	msg.body = `{"Greeting":"howzit!"}`
	err := ReadBody(msg, nil)

	if _, ok := err.(e.Err[NilPtr]); !ok {
		t.Errorf("want: nil ptr err; got: %v", err)
	}
}

func TestReadBodyEmptyReadOnNilMessageBody(t *testing.T) {
	msg := newMsgReader()
	msg.returnNilBody = true
	output := &ByteBody{}
	err := ReadBody(msg, output)

	if err != nil {
		t.Errorf("want: empty read; got: %v", err)
	}
	if len(output.Data) > 0 {
		t.Errorf("want: empty read; got: %v", output)
	}
}
