package hyper

import (
	"io"
	"testing"

	"github.com/c0c0n3/resto/util/bytez"
)

type msgWriter struct {
	headers          map[string]string
	body             []byte
	headerWritingErr error
	t                *testing.T
}

func newMsgWriter(t *testing.T) *msgWriter {
	return &msgWriter{
		headers: make(map[string]string, 0),
		t:       t,
	}
}

func (p *msgWriter) Header(name string, content string) error {
	if p.headerWritingErr != nil {
		return p.headerWritingErr
	}
	p.headers[name] = content
	return nil
}

func (p *msgWriter) Body(content io.ReadCloser) error {
	data, err := io.ReadAll(content)
	p.body = data
	return err
}

func (p *msgWriter) getHeader(name string) (string, bool) {
	val, present := p.headers[name]
	return val, present
}

func (p *msgWriter) stringBody() string {
	return string(p.body)
}

func (p *msgWriter) assertNoHeader(name string) {
	if val, present := p.getHeader(name); present {
		p.t.Errorf("want no header: %s; got: %s", name, val)
	}
}

func (p *msgWriter) assertHeader(name, want string) {
	if got, present := p.getHeader(name); !present {
		p.t.Errorf("want header %s; got nothing", name)
	} else {
		if got != want {
			p.t.Errorf("want: %s = %s; got: %s = %s", name, want, name, got)
		}
	}
}

func (p *msgWriter) assertNoBody() {
	got := p.stringBody()
	if len(got) > 1 {
		p.t.Errorf("want: empty body; got: %s", got)
	}
}

type msgReader struct {
	headers       map[string]string
	body          string
	returnNilBody bool
}

func newMsgReader() *msgReader {
	return &msgReader{
		headers:       make(map[string]string, 0),
		returnNilBody: false,
	}
}

func (p *msgReader) Header(name string) string {
	if val, present := p.headers[name]; present {
		return val
	}
	return ""
}

func (p *msgReader) Headers() map[string][]string {
	return nil
}

func (p *msgReader) Body() io.ReadCloser {
	if p.returnNilBody {
		return nil
	}
	buf := []byte(p.body)
	return bytez.NewBufferFrom(buf)
}

type byteStreamer struct {
	dataToStream []byte
	streamed     int
}

func newByteStreamer(content []byte) *byteStreamer {
	return &byteStreamer{
		dataToStream: content,
		streamed:     0,
	}
}

func (p *byteStreamer) Read(buf []byte) (n int, err error) {
	if len(p.dataToStream) == p.streamed {
		return 0, io.EOF
	}
	buf[0] = p.dataToStream[p.streamed]
	p.streamed = p.streamed + 1
	return 1, nil
}

func (p *byteStreamer) Close() error {
	return nil
}
