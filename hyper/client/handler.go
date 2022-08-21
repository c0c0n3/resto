package client

import (
	"github.com/c0c0n3/resto/hyper"
	"github.com/c0c0n3/resto/hyper/wire"
	"github.com/c0c0n3/resto/util/set"
)

// ExpectSuccess is a wire.ResponseHandler that checks for successful
// responses. If the response code is in the range 200-299 (both inclusive),
// ExpectSuccess does nothing, otherwise it returns an error---that
// stops any following wire.ResponseHandler from running.
//
// Example.
//
//     err := Request(
//         GET("https://google.com"),
//     ).Handle(
//         ExpectSuccess,
//     )
//     fmt.Printf("error: %v\n", err)
//
func ExpectSuccess(response wire.ResponseReader) error {
	code, _ := response.StatusLine()
	if code < 200 || code > 299 {
		return hyper.UnexpectedResponseErr(
			"expected successful response, got: %v", code)
	}
	return nil
}

// TODO. Implement expect for other status code ranges too?
// Informational responses (100–199)
// Successful responses (200–299)     --> DONE
// Redirects (300–399)
// Client errors (400–499)
// Server errors (500–599)

// ExpectStatusCodeOneOf builds a wire.ResponseHandler to check the
// response status code is among the given ones.
// If the response code is in the given list, the returned handler
// does nothing. Otherwise it returns an error---that stops any following
// wire.ResponseHandler from running.
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
func ExpectStatusCodeOneOf(expectedStatusCode ...int) wire.ResponseHandler {
	return func(response wire.ResponseReader) error {
		allowed := set.From(expectedStatusCode...)
		code, _ := response.StatusLine()

		if !allowed.Member(code.Value()) {
			return hyper.UnexpectedResponseErr("%v", code)
		}
		return nil
	}
}

// ReadJsonResponse builds a wire.ResponseHandler to deserialise a
// JSON response body, returning any error that stopped it from
// deserializing the body.
//
// Example.
//
//     data := &MyData{}
//     err := Request(
//         GET("https://my.api/data"),
//         Accept(mime.JSON),
//     ).Handle(
//         ExpectSuccess,
//         ReadJsonResponse(data),
//     )
//     fmt.Printf("data: %v\nerror: %v\n", data, err)
//
func ReadJsonResponse[T any](output *T) wire.ResponseHandler {
	return func(response wire.ResponseReader) error {
		deserializer := &hyper.JsonBody{Data: output}
		return hyper.ReadBody(response, deserializer)
	}
}

// ReadResponse builds a wire.ResponseHandler to read in a response body,
// returning any error that stopped it from reading the body. You use an
// hyper.BodyDeserializer to have ReadResponse convert the HTTP body octets
// to a data type instance of your liking.
//
// Example.
//
//     output := &hyper.StringBody{}  // reads the body in a string,
//     // output := &hyper.ByteBody{} // reads the body in a byte slice, etc.
//     err := Request(
//         GET("https://my.api/data"),
//         Accept(mime.OCTET_STREAM),
//     ).Handle(
//         ExpectSuccess,
//         ReadResponse(output),
//     )
//     fmt.Printf("data: %s\nerror: %v\n", output.Data, err)
//
func ReadResponse(output hyper.BodyDeserializer) wire.ResponseHandler {
	return func(response wire.ResponseReader) error {
		return hyper.ReadBody(response, output)
	}
}
