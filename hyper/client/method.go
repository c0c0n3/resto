package client

import (
	"net/url"

	"github.com/c0c0n3/resto/hyper/wire"
	e "github.com/c0c0n3/resto/util/err"
	"github.com/c0c0n3/resto/yoorel"
)

// WellFormedUrl lets you pass a yoorel.HttpUrl to the GET, POST, etc.
// builders.
// Because of some limitations of Go's union constraint element declarations,
// you can't pass a yoorel.HttpUrl directly to those functions. So we
// need this wrapper---see NOTE about it in TargetUrl for the details.
type WellFormedUrl struct {
	Url yoorel.HttpUrl
}

// The type of arguments you can pass to the GET, POST, etc. builders.
type TargetUrl interface {
	string | *url.URL | url.URL | WellFormedUrl // (*)

	// NOTE. Invalid union workaround. Ideally, we should be able to add
	// yoorel.HttpUrl directly to the union without a flipping wrapper.
	// But Go won't let you do that b/c the interface has methods.
	// See: InvalidUnion error explanation.
}

func toHttpUrl[T TargetUrl](target T) e.ErrOr[yoorel.HttpUrl] {
	var result e.ErrOr[yoorel.HttpUrl]
	switch value := any(target).(type) {
	case WellFormedUrl:
		result = e.FromResult(value.Url, nil)
	case string:
		result = yoorel.BuilderFrom(value).Build()
	case *url.URL:
		result = yoorel.BuilderFrom(value).Build()
	case url.URL:
		result = yoorel.BuilderFrom(value).Build()
	}
	return result
}

func makeRequestLineBuilder[U TargetUrl](verb wire.Method, resource U) wire.RequestBuilder {
	return func(req wire.RequestWriter) error {
		errOrUrl := toHttpUrl(resource)
		if errOrUrl.IsRight() {
			return req.RequestLine(verb, errOrUrl.Right())
		}
		return errOrUrl.Left()
		// write := func(target yoorel.HttpUrl) (R, error) {
		// 	return req, req.RequestLine(verb, target)
		// }
		// return e.Bind(write, toHttpUrl(resource))
	}
}

// GET writes the request line of a GET HTTP request to the specified
// resource.
func GET[U TargetUrl](resource U) wire.RequestBuilder {
	return makeRequestLineBuilder(wire.GET, resource)
}

// POST writes the request line of a POST HTTP request to the specified
// resource.
func POST[U TargetUrl](resource U) wire.RequestBuilder {
	return makeRequestLineBuilder(wire.POST, resource)
}

// PUT writes the request line of a PUT HTTP request to the specified
// resource.
func PUT[U TargetUrl](resource U) wire.RequestBuilder {
	return makeRequestLineBuilder(wire.PUT, resource)
}

// PATCH writes the request line of a PATCH HTTP request to the specified
// resource.
func PATCH[U TargetUrl](resource U) wire.RequestBuilder {
	return makeRequestLineBuilder(wire.PATCH, resource)
}

// DELETE writes the request line of a DELETE HTTP request to the specified
// resource.
func DELETE[U TargetUrl](resource U) wire.RequestBuilder {
	return makeRequestLineBuilder(wire.DELETE, resource)
}
