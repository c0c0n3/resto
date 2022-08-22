/*

Utils to make your life slightly easier than working directly with
the HTTP client from "net/http".


Request building

You can put together a request assembly line by mixing and matching
little, discrete, reusable pieces of functionality encapsulated by
wire.RequestBuilder functions. Request building reads like an HTTP
request on the wire and is more type-safe than doing it the Go way.
Example:

    response := Request(
        POST("http://you.api/greeting"),
        ContentType(mime.JSON),
        Body("howzit!"),
    )

There's quite a bit of built-in builders, but it's fairly easy
to write your own since a builder is just a function

    func(wire.RequestWriter) error

that uses the given wire.RequestWriter to write something to the
HTTP request to be sent to the server and returns an error if it
couldn't. For example, here's a builder to write a custom header
to specify the content's language is South African English

    func ContentLanguage(request wire.RequestWriter) error {
        return request.Header("Content-Language", "en-ZA")
    }

This custom builder writes a header of "Content-Language: en-ZA"
you could add to the previous request like this

    response := Request(
        POST("http://you.api/greeting"),
        ContentType(mime.JSON),
        ContentLanguage,
        Body("howzit!"),
    )

Surely, it'd probably be more useful if you could actually pass
in the language code. That's easy to do too,

    func ContentLanguage(code string) wire.RequestBuilder {
        return func(request wire.RequestWriter) error {
            return request.Header("Content-Language", code)
        }
    }

Yep, it's functions all the way down. And here's how you'd use
it

    response := Request(
        POST("http://you.api/greeting"),
        ContentType(mime.JSON),
        ContentLanguage("en-US"),
        Body("howdy, stranger"),
    )


Response handling

Similarly, the returned Response lets you build a response processing
pipeline out of wire.ResponseHandler functions. You do that by calling
the Response.Handle method as in the example below

	output := &hyper.ByteBody{}           // read the body in a byte slice
    err := response.Handle(
        ExpectStatusCodeOneOf(200, 201),  // fail if status not 200 or 201
        ReadResponse(output),             // otherwise read in the body
    )
    fmt.Printf("response body: %v\nerror: %v\n", output.Data, err)

If the request had an error, Handle returns it and exits. Otherwise
it feeds the response to each supplied wire.ResponseHandler in turn
and in the same order as in the input list, stopping at the first one
that errors out and returning that error. If all the handlers are
successful, the returned error will be nil. If the response contains
a body, Handle automatically closes the associated reader just before
returning, so handlers don't have to do that.

So you can chain handlers to do more than one thing with the response
so the code stays modular---single responsibility principle, anyone?

There's a fair bit of built-in handlers. If you need to write your
own, keep in mind, like for request builders, a handler is just a
function

    func(wire.ResponseReader) error

It uses the given wire.ResponseReader to read from the server's response,
returning an error if it couldn't. For example, this handler prints
all the response headers to stdout

    func PrintHeaders(response wire.ResponseReader) error {
        fmt.Printf("response headers:\n%v\n", response.Headers())
        return nil
    }

and you can easily slot it in

	output := &hyper.ByteBody{}
    err := response.Handle(
        PrintHeaders,
        ExpectStatusCodeOneOf(200, 201),
        ReadResponse(output),
    )
    fmt.Printf("response body: %v\nerror: %v\n", output.Data, err)

Notice it makes sense for PrintHeaders to be the first handler
in the response processing pipeline since if there was another
handler before that and that handler failed, then PrintHeaders
wouldn't run.


Message exchange

So we've seen how to use request builders to send an HTTP request
and response handlers to consume the response. Besides modularity,
one nice thing about them is that you can pull them together in an
HTTP request-reply message flow where execution stops at the first
error---so you don't have to litter your code with `if err ...`
statements. Here's a complete example that POSTs some JSON to the
HttpBin service at https://httpbin.org/ and prints the response.
Notice the code to build a custom header (to specify a language)
and how the framework handles writing/reading Go data to/from the
request/response body under the bonnet for you.

    package main

    import (
        "fmt"

        "github.com/c0c0n3/resto/hyper"
        //lint:ignore ST1001 using dot import to make EDSL read better
        . "github.com/c0c0n3/resto/hyper/client"
        "github.com/c0c0n3/resto/hyper/wire"
        "github.com/c0c0n3/resto/mime"
    )

    func Run() {
        greeting := &Greeting{Message: "howzit!"}
        output := &hyper.StringBody{}
        err := Request(
            POST("https://httpbin.org/post"),
            ContentType(mime.JSON),
            ContentLanguage("en-ZA"),
            Body(Json(greeting)),
        ).Handle(
            PrintHeaders,
            ExpectSuccess,
            ReadResponse(output),
        )
        fmt.Printf("body: %s\nerror: %v\n", output.Data, err)
    }

    type Greeting struct {
        Message string
    }

    func ContentLanguage(code string) wire.RequestBuilder {
        return func(request wire.RequestWriter) error {
            return request.Header("Content-Language", code)
        }
    }

    func PrintHeaders(response wire.ResponseReader) error {
        fmt.Println("headers:")
        for name, value := range response.Headers() {
            fmt.Printf("\t%s: %v\n", name, value)
        }
        return nil
    }

    func main() {
        Run()
    }

HttpBin replies with a JSON object containing the JSON data we POSTed
as well as the headers we sent along in the request. If you run this
code, you should see an output similar to the one below.

    headers:
        Date: [Mon, 22 Aug 2022 09:39:41 GMT]
        Content-Type: [application/json]
        Content-Length: [491]
        Server: [gunicorn/19.9.0]
        Access-Control-Allow-Origin: [*]
        Access-Control-Allow-Credentials: [true]
    body: {
        "args": {},
        "data": "{\"Message\":\"howzit!\"}",
        "files": {},
        "form": {},
        "headers": {
            "Accept-Encoding": "gzip",
            "Content-Language": "en-ZA",
            "Content-Length": "21",
            "Content-Type": "application/json",
            "Host": "httpbin.org",
            "User-Agent": "Go-http-client/2.0",
            "X-Amzn-Trace-Id": "Root=1-63034edd-5312e02638edf96215f21a47"
        },
        "json": {
            "Message": "howzit!"
        },
        "origin": "83.79.100.214",
        "url": "https://httpbin.org/post"
    }

    error: <nil>


HTTP client

The Request function we've been using so far is just a convenience
function to send an HTTP request using http.DefaultClient. If you
need more flexibility, use the New function instead which lets you
specify your own function to exchange HTTP messages. In fact, New
accepts a Sender which is a function of type

    func(wire.RequestBuilder) (wire.ResponseReader, error)

capable of writing an HTTP request with the given RequestBuilder
and returning a ResponseReader to let the caller read the server's
reply. If you need total control over how the framework exchanges
HTTP messages, then implement your own Sender. But you're probably
already happy with the "net/http" package, so why reinvent the wheel?
Well, the wire package actually comes with its own Sender implementation
backed by "net/http".

So one option is to tweak that implementation to suit your needs.
This is fairly easy since there's a wire.NewSender function that
can build a Sender out of your own instance of http.Client. Here's
how you could do it

    client := &http.Client{Timeout: time.Second * 10}
    sender := wire.NewSender(client)
    hyperc := New(sender)

    err := hyperc.Request(
        GET("http://some/stuff"),
    ).Handle()

    fmt.Printf("error: %v", err)


Unit testing

As you might've guessed the only place where IO actually happens is
in the Sender function. If you provided your own Sender, you could
swap out actual HTTP calls with stubs. This would make it super easy
to unit-test your code. But is there a cheap way to write a Sender
stub for unit testing? Turns out there is. In fact, wire.NewSender
can also build a fully-fledged Sender out of a function

    func(*http.Request) (*http.Response, error)

and it's kinda easy to implement stubs with that signature. Let's
have a look at an example. We're going to write a stub that returns
a 200 with a version number in the body and use that stub to test
a call to the version endpoint of a fictitious API.

    func TestVersion(t *testing.T) {
        version := "v2"
        send := func(req *http.Request) (*http.Response, error) {
            return &http.Response{
                StatusCode: 200,
                Body:       bytez.NewBufferFrom([]byte(version)),
            }, nil
        }
        sender := wire.NewSender(send)

        output := &hyper.StringBody{}
        err := New(sender).Request(
            GET("http://you.api/version"),
        ).Handle(
            ExpectSuccess,
            ReadResponse(output),
        )

        if err != nil {
            t.Errorf("want: success; got: %v", err)
        }
        if output.Data != version {
            t.Errorf("want: %s; got: %s", version, output.Data)
        }
    }

*/
package client

/*

// Message exchange example above.
// Copy-paste this code in main.go, then run it with `go run main.go`.

package main

import (
	"fmt"

	"github.com/c0c0n3/resto/hyper"
	//lint:ignore ST1001 using dot import to make EDSL read better
	. "github.com/c0c0n3/resto/hyper/client"
	"github.com/c0c0n3/resto/hyper/wire"
	"github.com/c0c0n3/resto/mime"
)

func Run() {
	greeting := &Greeting{Message: "howzit!"}
	output := &hyper.StringBody{}
	err := Request(
		POST("https://httpbin.org/post"),
		ContentType(mime.JSON),
		ContentLanguage("en-ZA"),
		Body(Json(greeting)),
	).Handle(
		PrintHeaders,
		ExpectSuccess,
		ReadResponse(output),
	)
	fmt.Printf("body: %s\nerror: %v\n", output.Data, err)
}

type Greeting struct {
	Message string
}

func ContentLanguage(code string) wire.RequestBuilder {
	return func(request wire.RequestWriter) error {
		return request.Header("Content-Language", code)
	}
}

func PrintHeaders(response wire.ResponseReader) error {
	fmt.Println("headers:")
	for name, value := range response.Headers() {
		fmt.Printf("\t%s: %v\n", name, value)
	}
	return nil
}

func main() {
	Run()
}

*/
