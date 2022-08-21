package hyper

import (
	"encoding/json"
	"io"

	"github.com/c0c0n3/resto/hyper/wire"
	"github.com/c0c0n3/resto/util/bytez"
)

// BodyTransfer says if the body content should be streamed.
type BodyTransfer interface {
	// Return true if the body content should be streamed, false otherwise.
	Streaming() bool
}

// BodySerializer turns data into an HTTP body octet sequence.
type BodySerializer interface {
	BodyTransfer
	// Serialize turns the data it holds into an HTTP body octet sequence.
	// The returned io.ReadCloser hods the octets whereas the integer output
	// counts them---i.e. it's the body size. The body size will be ignored
	// in the case of a streaming body.
	Serialize() (body io.ReadCloser, bodySize int, err error)
}

// BodyDeserializer reads an HTTP body octet sequence into some data
// structure or stream.
type BodyDeserializer interface {
	BodyTransfer
	// Deserialize reads the given HTTP body octets into its own data
	// structure. Deserialize doesn't close the body stream, the caller
	// is responsible for that.
	Deserialize(body io.ReadCloser) error
}

// Write a message body with the given content.
// Also write a "Content-Length" header with the size of the content
// if the content objects knows upfront how many bytes it'll write to
// the body---i.e. non-streaming content.
func WriteBody(msg wire.MessageWriter, content BodySerializer) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	if content == nil {
		return NilBodySerializerErr()
	}

	contentReader, bodySize, err := content.Serialize()
	if err != nil {
		return err
	}
	if !content.Streaming() {
		if err := WriteContentLength(msg, uint64(bodySize)); err != nil {
			return err
		}
	}
	return msg.Body(contentReader)
}

// Read the message body into the given buffer.
func ReadBody(msg wire.MessageReader, content BodyDeserializer) error {
	if msg == nil {
		return NilMessageWriterErr()
	}
	if content == nil {
		return NilBodyDeserializerErr()
	}
	return content.Deserialize(msg.Body())
}

func ensureReader(r io.ReadCloser) io.ReadCloser {
	if r == nil {
		return bytez.NewBuffer()
	}
	return r
}

// JsonBody holds a data structure that needs to be (de-)serialized
// (from) to an HTTP body octet stream containing JSON data.
type JsonBody struct {
	Data any
}

func (p *JsonBody) Streaming() bool {
	return false
}

func (p *JsonBody) Serialize() (io.ReadCloser, int, error) {
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

func (p *JsonBody) Deserialize(reader io.ReadCloser) error {
	// var json = jsoniter.ConfigCompatibleWithStandardLibrary // (*)
	decoder := json.NewDecoder(ensureReader(reader))
	return decoder.Decode(p.Data)

	// (*) json-iterator lib.
	// We could use it in Serialize() to work around encoding/json's
	// inability to serialise map[interface {}]interface{} types. Here
	// we're parsing JSON into a data structure and AFAICT the built-in
	// json lib can parse pretty much any valid JSON you throw at it.
	// So the only reason to use json-iterator in place of encoding/json
	// would be performance: json-iterator is way faster than encoding/json.
	// But ideally resto shouldn't depend on external libs...
	// TODO if we switch to json-iterator in Serialize(), then use
	// json-iterator here too.
}

// StringBody holds a string that needs to be (de-)serialized (from) to
// an HTTP body octet stream containing text.
type StringBody struct {
	Data string
}

func (p *StringBody) Streaming() bool {
	return false
}

func (p *StringBody) Serialize() (io.ReadCloser, int, error) {
	buf := []byte(p.Data)
	return bytez.Reader(buf), len(buf), nil
}

func (p *StringBody) Deserialize(reader io.ReadCloser) error {
	buf := bytez.NewBuffer()
	_, err := io.Copy(buf, ensureReader(reader))
	p.Data = string(buf.Bytes())
	return err
}

// ByteBody holds a byte slice that needs to be written/read to/from
// an HTTP body octet stream.
type ByteBody struct {
	Data []byte
}

func (p *ByteBody) Streaming() bool {
	return false
}

func (p *ByteBody) Serialize() (io.ReadCloser, int, error) {
	return bytez.Reader(p.Data), len(p.Data), nil
}

func (p *ByteBody) Deserialize(reader io.ReadCloser) error {
	buf := bytez.NewBuffer()
	_, err := io.Copy(buf, ensureReader(reader))
	p.Data = buf.Bytes()
	return err
}

// StreamingBody produces/consumes an HTTP body octet stream in constant
// space.
type StreamingBody struct {
	Data io.ReadCloser
}

func (p *StreamingBody) Streaming() bool {
	return true
}

func (p *StreamingBody) Serialize() (io.ReadCloser, int, error) {
	return p.Data, 0, nil
}

func (p *StreamingBody) Deserialize(reader io.ReadCloser) error {
	p.Data = ensureReader(reader)
	return nil
}
