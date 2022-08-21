package client

import (
	"io"

	"github.com/c0c0n3/resto/hyper"
	"github.com/c0c0n3/resto/hyper/wire"
)

// RequestBody represents some data structure to be written to an
// HTTP message body. The Body function takes care of converting
// the data to a sequence of HTTP body octets.
type RequestBody interface {
	[]byte | string | *hyper.JsonBody | *hyper.StreamingBody
}

func bodyContentToSerializer[T RequestBody](data T) hyper.BodySerializer {
	var serializer hyper.BodySerializer
	switch target := any(data).(type) {
	case []byte:
		serializer = &hyper.ByteBody{Data: target}
	case string:
		serializer = &hyper.StringBody{Data: target}
	case *hyper.JsonBody:
		serializer = target
	case *hyper.StreamingBody:
		serializer = target
	}
	return serializer
}

// Body turns the given content into HTTP octets which it writes to
// the message body. Also it writes a "Content-Length" header with
// the size of the resulting octet sequence.
func Body[T RequestBody](content T) wire.RequestBuilder {
	serializer := bodyContentToSerializer(content)
	return func(msg wire.RequestWriter) error {
		return hyper.WriteBody(msg, serializer)
	}
}

// Json lets you serialise the given data structure to JSON and write
// data to the message body. Notice you use Json with the Body function
// which takes care of writing a "Content-Length" header with the size
// of the serialised JSON content---so you don't have to do that yourself.
//
// Example.
//
//     data := &MyData{greeting: "howzit!"}
//     err := Request(
//         POST("https://my.api/data"),
//         ContentType(mime.JSON),
//         Body(Json(data)),
//     ).Handle(
//         ExpectStatusCodeOneOf(200, 201),
//     )
//     fmt.Printf("error: %v\n", err)
//
func Json(data any) *hyper.JsonBody {
	return &hyper.JsonBody{Data: data}
}

// Stream the content of a reader to the message body.
//
// Example.
//
//     file, _ := os.Open("some.txt")
//     defer file.Close()
//     err := Request(
//         POST("https://my.api/data"),
//         ContentType(mime.PLAIN_TEXT),
//         Body(Stream(file)),
//     ).Handle(
//         ExpectStatusCodeOneOf(200, 201),
//     )
//     fmt.Printf("error: %v\n", err)
//
func Stream(data io.ReadCloser) *hyper.StreamingBody {
	return &hyper.StreamingBody{Data: data}
}
