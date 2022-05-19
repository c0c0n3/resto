package yoorel

import (
	"net/url"

	"github.com/c0c0n3/resto/util/err"
)

// Build an HttpUrl.
// Call the provided methods to build the various URL parts, then call
// Build when done to get an HttpUrl or an error if the parts entered
// so far don't add up to a valid HttpUrl.
// Escaping. Don't enter escaped strings, they'll be double escaped.
// Just enter the string as is, e.g.
//     JoinPath("a b") ==> /a%20b
type Builder interface {
	// Set the scheme to HTTP.
	Http() Builder
	// Set the scheme to HTTPs.
	Https() Builder
	// Set the host and port part of the URL.
	// You can either specify just the host without the port or a host and
	// port part separated by a colon as in "host:8080". Any extra white
	// space gets trimmed from both ends of the input.
	HostAndPort(hp string) Builder
	// Append one or more path components to the current URL.
	// Examples.
	//     JoinPath("a")                   ==> /a
	//     JoinPath("a", "b/c")            ==> /a/b/c
	//     JoinPath("/a/b").JoinPath("/c") ==> /a/b/c
	JoinPath(ps ...string) Builder
	// Add the given key-value pair to the query part of the URL.
	Query(key, value string) Builder
	// Build an HttpUrl from the parts entered so far or return an error
	// if the parts don't add up to a valid HttpUrl.
	Build() err.ErrOr[HttpUrl]
}

// Convenience union type to specify all the kind of inputs BuilderFrom
// accepts.
type InitialBuilderUrl interface {
	string | *url.URL | url.URL
}

// Pre-populate a Builder with the data in the given URL seed.
func BuilderFrom[T InitialBuilderUrl](seed ...T) Builder {
	if len(seed) > 0 {
		switch partialUrl := any(seed[0]).(type) {
		case url.URL:
			return builderFromUrl(&partialUrl)
		case *url.URL:
			return builderFromUrl(partialUrl)
		case string:
			return builderFromRawUrl(partialUrl)
		}
	}
	return EmptyBuilder()
}

// Create an empty Builder.
func EmptyBuilder() Builder {
	return builderBuffer{
		path:  make([]string, 0),
		query: make(url.Values),
	}
}

func builderFromUrl(partialUrl *url.URL) Builder {
	if partialUrl == nil {
		return EmptyBuilder()
	}
	return builderBuffer{
		scheme:      partialUrl.Scheme,
		hostAndPort: partialUrl.Host,
		path:        []string{partialUrl.Path},
		query:       partialUrl.Query(),
	}
}

func builderFromRawUrl(rawUrl string) Builder {
	if parsed, e := url.Parse(rawUrl); e == nil {
		return builderBuffer{
			scheme:      parsed.Scheme,
			hostAndPort: parsed.Host,
			path:        []string{parsed.Path},
			query:       parsed.Query(),
		}
	} else {
		return builderErr{e}
	}
}
