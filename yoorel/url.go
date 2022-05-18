package yoorel

import (
	"net/url"
	"strconv"

	"github.com/c0c0n3/resto/util/set"
)

// A valid reference to a Web resource, e.g. "http://some/api".
// It's an absolute URI you can only get through a Builder.
type HttpUrl interface {
	// Is this an HTTPs or HTTP URL?
	Secure() bool
	// Host part of the URL.
	Host() string
	// Port part of the URL.
	Port() int
	// Path part of the URL.
	Path() string
	// Query key set. Will be empty if this URL doesn't have a query part.
	QueryKeys() set.Set[string]
	// Look up query values associated to the given key. Will be an empty
	// slice if there are no values for the key.
	QueryValues(key string) []string
	// Convert the URL to its string representation you can use in an HTTP
	// message. Path and query parts get URL-escaped as needed.
	WireFormat() string

	// NOTE. Keep this interface read-only. Builder takes care of mutability.
}

type httpUrl struct {
	ref   *url.URL
	query url.Values
}

func (p httpUrl) Secure() bool {
	return p.ref.Scheme == "https"
}

func (p httpUrl) Host() string {
	return p.ref.Hostname()
}

func (p httpUrl) Port() int {
	if port, err := strconv.Atoi(p.ref.Port()); err == nil {
		return port
	}
	if p.Secure() {
		return DEFAULT_HTTPS_PORT
	}
	return DEFAULT_HTTP_PORT
}

func (p httpUrl) Path() string {
	if p.ref.Path == "" {
		return "/"
		// see comments in builderImpl.JoinPath
	}
	return p.ref.Path
}

func (p httpUrl) QueryKeys() set.Set[string] {
	keys := make([]string, len(p.query))

	i := 0
	for k := range p.query {
		keys[i] = k
		i++
	}
	return set.FromList(keys)
}

func (p httpUrl) QueryValues(key string) []string {
	if vs, ok := p.query[key]; ok {
		return vs
	}
	return []string{}
}

func (p httpUrl) WireFormat() string {
	p.ref.ForceQuery = false
	refRepr := p.ref.String()

	if len(p.query) == 0 {
		return refRepr
	}
	queryRepr := p.query.Encode()
	return refRepr + "?" + queryRepr
}
