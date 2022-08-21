/*
Package wire defines the interfaces to read and write HTTP messages
as well as the interfaces to exchange them. Plus, it provides a default
implementation backed by the standard lib's "net/http" package.
Notice this package is low-level stuff you probably don't want to
use, but it lets us easily build EDSLs on top of it. (It takes quite
a bit of boilerplate to use this package directly, look at the example
below.)


Reading and writing a message

The RequestReader and ResponseReader interfaces provide methods to
read the parts that make up an HTTP request and response, respectively.
They both extend MessageReader which factors out the reading of headers
and body since both requests and responses can have these parts.

Symmetrically, the RequestWriter and ResponseWriter interfaces have
methods to write the parts that make up an HTTP request and response,
respectively. They both extend MessageWriter which factors out the
writing of headers and body.

All these reading/writing interfaces don't force the implementation
into imperative programming and buffering the whole message content
into memory before sending it out. In fact, an implementation could
use immutable objects and streaming for example.


Exchanging messages

The Sender and RequestHandler function types encapsulate the sending
of a request to a server and the server-side processing of it, respectively.
You can implement your own functions or get a default implementation
backed by the "net/http" package from the standard library. (At the
moment, there's a fully-fledged Sender implementation you get through
NewSender, but there's no server-side RequestHandler implementation.)


Example - POSTing to HttpBin

The code below sends a POST request to https://httpbin.org/post with
a custom header of "greeting: howzit!" and a body with the '*' char.
Then it prints out HttpBin's response which echoes the request data
in a JSON object.

	package main

	import (
		"fmt"
		"io"

		"github.com/c0c0n3/resto/util/bytez"
		"github.com/c0c0n3/resto/yoorel"
	)

	func postRequest(req RequestWriter) error {
		url := yoorel.EmptyBuilder().Https().HostAndPort("httpbin.org").
			JoinPath("/post").Build().Right()
		err := req.RequestLine(POST, url)
		if err != nil {
			return err
		}
		err = req.Header("greeting", "howzit!")
		if err != nil {
			return err
		}
		content := bytez.Reader([]byte{42}) // ASCII 42 = '*'
		return req.Body(content)
	}

	func Run() {
		send := NewSender[DefaultClient]()
		responseReader, err := send(postRequest)
		if err != nil {
			panic(err)
		}

		code, reason := responseReader.StatusLine()
		fmt.Printf("code: %d; reason: %s\n", code, reason)

		if body, err := io.ReadAll(responseReader.Body()); err != nil {
			fmt.Printf("body: %v\n", err)
		} else {
			fmt.Printf("body: %s\n", string(body))
		}
	}

	func main() {
		Run()
	}

Prints:

	code: 200; reason: 200 OK
	body: {
		"args": {},
		"data": "*",
		"files": {},
		"form": {},
		"headers": {
			"Accept-Encoding": "gzip",
			"Greeting": "howzit!",
			"Host": "httpbin.org",
			"Transfer-Encoding": "chunked",
			"User-Agent": "Go-http-client/2.0",
			"X-Amzn-Trace-Id": "Root=1-62fe61fc-47a0076c71834e966ff458cf"
		},
		"json": null,
		"origin": "83.79.100.214",
		"url": "https://httpbin.org/post"
	}

*/
package wire

/*

// code for the above example, comment back in to test.

import (
	"fmt"
	"io"

	"github.com/c0c0n3/resto/util/bytez"
	"github.com/c0c0n3/resto/yoorel"
)

func postRequest(req RequestWriter) error {
	url := yoorel.EmptyBuilder().Https().HostAndPort("httpbin.org").
		JoinPath("/post").Build().Right()
	err := req.RequestLine(POST, url)
	if err != nil {
		return err
	}
	err = req.Header("greeting", "howzit!")
	if err != nil {
		return err
	}
	content := bytez.Reader([]byte{42}) // ASCII 42 = '*'
	return req.Body(content)
}

func Run() {
	send := NewSender[DefaultClient]()
	responseReader, err := send(postRequest)
	if err != nil {
		panic(err)
	}

	code, reason := responseReader.StatusLine()
	fmt.Printf("code: %d; reason: %s\n", code, reason)

	if body, err := io.ReadAll(responseReader.Body()); err != nil {
		fmt.Printf("body: %v\n", err)
	} else {
		fmt.Printf("body: %s\n", string(body))
	}
}

*/
