package wire

import (
	"io"
	"net/http"

	"github.com/c0c0n3/resto/util/bytez"
	"github.com/c0c0n3/resto/yoorel"
)

// Request writing buffer. It implements RequestWriter and conversion to
// http.Request.
type reqBuf struct {
	method  string
	url     string
	headers http.Header
	body    io.ReadCloser
}

func emptyReqBuf() *reqBuf {
	return &reqBuf{
		headers: make(http.Header),
	}
}

func (p *reqBuf) Header(name string, content string) error {
	p.headers.Set(name, content)
	return nil
}

func (p *reqBuf) Body(content io.ReadCloser) error {
	p.body = content
	return nil
}

func (p *reqBuf) RequestLine(verb Method, resource yoorel.HttpUrl) error {
	p.method = verb.String()
	p.url = resource.WireFormat()
	return nil
}

func (p *reqBuf) toHttpRequest() (*http.Request, error) {
	req, err := http.NewRequest(p.method, p.url, p.body)
	if err != nil {
		return nil, err
	}
	req.Header = p.headers
	return req, nil
}

// Response reading. It implements ResponseReader by sourcing data from
// http.Response.
type resReader struct {
	res *http.Response
}

func (p *resReader) Header(name string) string {
	return p.res.Header.Get(name)
}

func (p *resReader) Headers() map[string][]string {
	return p.res.Header
}

func (p *resReader) Body() io.ReadCloser {
	if p.res.Body == nil {
		return bytez.NewBuffer()
	}
	return p.res.Body
}

func (p *resReader) StatusLine() (code StatusCode, reason string) {
	return StatusCode(p.res.StatusCode), p.res.Status
}

// StdLibSender is a function to send an HTTP request and receive a
// response from the server using the http package request and response
// types. It returns an error if something goes wrong---e.g. a network
// failure. Notice http.Client.Do is a StdLibSender.
type StdLibSender func(*http.Request) (*http.Response, error)

// DefaultClient is a convenience type to tell NewSender to create a
// Sender backed by http.DefaultClient.
type DefaultClient int

// StdLibClient declares the type of arguments NewSender accepts.
type StdLibClient interface {
	*http.Client | DefaultClient |
		StdLibSender | func(*http.Request) (*http.Response, error)
	// NOTE. No implicit type cast.
	// The Go compiler won't accept ~StdLibSender, so we've got add the
	// actual function type too otherwise NewSender's callers won't be
	// able to pass in a plain function of type
	//   func(*http.Request) (*http.Response, error)
	// without a type cast.
}

func getStdLibRequestSender[T StdLibClient](client ...T) StdLibSender {
	if len(client) > 0 {
		switch impl := any(client[0]).(type) {
		case *http.Client:
			return impl.Do
		case StdLibSender:
			return impl
		case func(*http.Request) (*http.Response, error):
			return impl // see NOTE in StdLibClient about implicit casts.
		}
	}
	return http.DefaultClient.Do
}

func sendRequest(build RequestBuilder, send StdLibSender) (ResponseReader, error) {
	buf := emptyReqBuf()
	if err := build(buf); err != nil {
		return nil, err
	}
	request, err := buf.toHttpRequest()
	if err != nil {
		return nil, err
	}
	response, err := send(request)
	if err != nil {
		return nil, err
	}
	return &resReader{response}, nil
}

// Build a Sender to make HTTP requests.
//
// The Sender uses the http package under the bonnet to make HTTP
// requests and read responses. To get a basic Sender, use
//
//     send := NewSender[DefaultClient]()
//
// This will make the send function use http.DefaultClient to exchange
// messages. But you can also make the returned Sender use your own
// instance of http.Client as in this example
//
//     client := &http.Client{Timeout: time.Second * 10}
//     send := NewSender(client)
//
// where the given client is what the send function will use to make
// an HTTP request and read the response each time it gets called.
// Also notice the Sender can actually exchange messages using any
// function of type StdLibSender
//
//     func(*http.Request) (*http.Response, error)
//
// which is any function that can send an http.Request and pack the
// server's reply in an http.Response, returning an error if something
// goes wrong. So if needed, you can even make the Sender use your own
// StdLibSender function. This comes in handy especially when you want
// to unit-test your code without having to jump through hoops since
// you can easily swap out any HTTP call with a stub. Here's an example
//
//     clientStub := func(req *http.Request) (*http.Response, error) {
//         return &http.Response{StatusCode: 200}, nil
//     }
//     getRequest := func(req RequestWriter) error {
//         url := yoorel.BuilderFrom("http://nowhere/").Build().Right()
//         return req.RequestLine(GET, url)
//     }
//     send := NewSender(clientStub)
//
//     response, _ := send(getRequest)
//     code, _ := response.StatusLine() // code == 200
//
func NewSender[T StdLibClient](client ...T) Sender {
	send := getStdLibRequestSender(client...)
	return func(build RequestBuilder) (ResponseReader, error) {
		return sendRequest(build, send)
	}
}
