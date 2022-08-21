package client

import (
	"github.com/c0c0n3/resto/hyper"
	"github.com/c0c0n3/resto/hyper/wire"
)

// Client represents an HTTP client.
type Client struct {
	send wire.Sender
}

// Make a new HTTP client.
// If you pass no arguments, you get a Client backed by the standard lib's
// http.DefaultClient. Otherwise, you can pass in a wire.Sender of your
// liking to make Client use that to exchange HTTP messages. See wire.NewSender
// for a few options to create a wire.Sender with no effort. If that doesn't
// do, you can always provide your own wire.Sender implementation.
func New(transport ...wire.Sender) *Client {
	if len(transport) == 0 || transport[0] == nil {
		return &Client{
			send: wire.NewSender[wire.DefaultClient](),
		}
	}
	return &Client{send: transport[0]}
}

// Holds the server response to a given request.
type Response struct {
	requestError error
	reader       wire.ResponseReader
}

// Request makes an HTTP request using the given builders and returns
// the server's response.
//
// Each wire.RequestBuilder writes some fields of the HTTP request,
// possibly returning an error if something goes wrong. Request lets
// you chain builders to write a full HTTP request, each builder contributes
// its bit and the request building process stops at the first error.
// If there's an error in the request writing or sending phase that
// error gets propagated on to the Response so calling Response.Handle
// will return that error immediately and stop there. Basically
// a poor man's monomorphic either+IO monad stack---ask Google.
func (p *Client) Request(fields ...wire.RequestBuilder) *Response {
	for _, f := range fields {
		if f == nil {
			return &Response{requestError: hyper.NilRequestBuilderErr()}
		}
	}
	request := makeRequestBuilder(fields...)
	reader, err := p.send(request)

	return &Response{err, reader}
}

func makeRequestBuilder(builders ...wire.RequestBuilder) wire.RequestBuilder {
	return func(request wire.RequestWriter) error {
		for _, build := range builders {
			if err := build(request); err != nil {
				return err
			}
		}
		return nil
	}
}

// Handle processes the HTTP response returned by Client.Request.
//
// If the request had an error, Handle returns it and exits. Otherwise
// it feeds the response to each supplied wire.ResponseHandler in turn
// and in the same order as in the input list, stopping at the first
// one that errors out and returning that error. If all the handlers
// are successful, the returned error will be nil. If the response
// contains a body, Handle automatically closes the associated reader
// just before returning, so handlers don't have to do that.
//
// So you can chain handlers to do more than one thing with the response
// so the code stays modular---single responsibility principle, anyone?
// Since the response processing chain stops at the first error, basically
// we've got a poor man's monomorphic either+IO monad stack---ask Google.
func (p Response) Handle(handlers ...wire.ResponseHandler) error {
	if p.requestError != nil {
		return p.requestError
	}

	defer p.reader.Body().Close()
	for _, handle := range handlers {
		if handle == nil {
			return hyper.NilResponseHandlerErr()
		}
		if err := handle(p.reader); err != nil {
			return err
		}
	}
	return nil
}

// Request is a convenience function to send an HTTP request using
// http.DefaultClient. The given builders write the request as explained
// in Client.Request.
func Request(fields ...wire.RequestBuilder) *Response {
	return New().Request(fields...)
}
