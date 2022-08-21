package hyper

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/c0c0n3/resto/hyper/wire"
	"github.com/c0c0n3/resto/util/bytez"
)

// BodyContent represents some data structure to be written (read) to
// (from) an HTTP message body. The WriteBody (ReadBody) function takes
// care of converting the data to (from) a sequence of HTTP body octets.
type BodyContent interface {
	[]byte | string | JSON
}

func bodyContentToReader[T BodyContent](data T) (
	reader io.ReadCloser, size int, err error) {
	switch target := any(data).(type) {
	case []byte:
		size = len(target)
		reader = bytez.Reader(target)
	case string:
		buf := []byte(target)
		size = len(buf)
		reader = bytez.Reader(buf)
	case JSON:
		reader, size, err = target.serialize()
	}
	return reader, size, err
}

/*
TODO why this function can't compile?
Getting IncompatibleAssign errors, e.g. for the []byte case

    cannot use io.ReadAll(data) (value of type []byte) as T value
	in assignment: cannot assign []byte to string (in T)

func bodyContentFromReader[T BodyContent](data io.ReadCloser) (out T, err error) {
	switch any(out).(type) {
	case []byte:
		out, err = io.ReadAll(data)
	case string:
		buf, err := io.ReadAll(data)
		out = string(buf)
	case JSON:
		err = out.deserialize(data)
	}
	return out, err
}

If this, or something similar, worked then we could have a single

    func ReadBody[T BodyContent](msg wire.MessageReader, output *T) error

instead of the many Read* ones we've got at the moment.
*/

// Write a message body with the given content.
// Also write a "Content-Length" header with the size of the content.
func WriteBody[T BodyContent](msg wire.MessageWriter, content T) error {
	contentReader, bodySize, err := bodyContentToReader(content)
	if err != nil {
		return err
	}
	if err := WriteContentLength(msg, uint64(bodySize)); err != nil {
		return err
	}
	return msg.Body(contentReader)
}

// TODO also implement streaming body? most of the standard libs aren't built
// w/ streaming in mind, so in practice you'll likely have the whole body in
// memory most of the time for common cases---e.g. JSON, YAML.

// Deserialise a JSON message body into the given output data structure.
func ReadJsonBody[T any](msg wire.MessageReader, output *T) error {
	if output == nil {
		return NilJsonOutDataErr()
	}
	out := JSON{output}
	return out.deserialize(msg.Body())
}

// Read the message body into the given buffer.
func ReadBody(msg wire.MessageReader, output *bytes.Buffer) error {
	if output == nil {
		return NilBytesBufferErr()
	}
	_, err := io.Copy(output, msg.Body())
	return err
}

// JSON holds a data structure that needs to be (de-)serialized (from)
// to an HTTP body octet stream containing JSON data.
type JSON struct {
	Data any
}

func (p JSON) serialize() (io.ReadCloser, int, error) {
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary // (*)
	buf, err := json.Marshal(p.Data)
	return bytez.NewBufferFrom(buf), len(buf), err

	// (*) json-iterator lib.
	// Should we use it in place of json from Go's standard lib? Pros: it
	// can handle serialisation corner cases. Cons: resto should only depend
	// on Go's standard lib.
	//
	// json-iterator can handle the serialisation of fields of type
	//    map[interface {}]interface{}
	// where the built-in json module will blow up w/
	//    json: unsupported type: map[interface {}]interface{}
	// (See TestWriteUntypedJsonBody for an example.)
	// If you're reading in YAML and then writing it out as JSON you could get
	// bitten by this. For example, say you use "gopkg.in/yaml.v2" to read
	// some YAML that has a field containing arbitrary JSON into a struct
	// with a field X of type interface{}---you don't know what the JSON looks
	// like, but later on you still want to be able to write it out.
	// The YAML lib will read the JSON into X with a type of
	//     map[interface {}]interface{}
	// but when you call the built-in json.Marshal, it'll blow up in your face
	// b/c it doesn't know how to handle that type.
	// See:
	// - https://stackoverflow.com/questions/35377477
}

func (p JSON) deserialize(reader io.ReadCloser) error {
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary // (*)
	decoder := json.NewDecoder(reader)
	return decoder.Decode(p.Data)

	// (*) json-iterator lib.
	// We could use it in serialize() to work around encoding/json's
	// inability to serialise map[interface {}]interface{} types. Here
	// we're parsing JSON into a data structure and AFAICT the built-in
	// json lib can parse pretty much any valid JSON you throw at it.
	// So the only reason to use json-iterator in place of encoding/json
	// would be performance: json-iterator is way faster than encoding/json.
	// But ideally resto shouldn't depend on external libs...
	// TODO if we switch to json-iterator in serialize(), then use
	// json-iterator here too.
}
