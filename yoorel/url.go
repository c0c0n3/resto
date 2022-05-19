package yoorel

import (
	"net/url"
	"path"

	"github.com/c0c0n3/resto/util/set"
	"github.com/c0c0n3/resto/util/str"
)

// Enumerates the URL schemes for a Web resource.
type Scheme str.CaseInsensitive

const (
	Http  Scheme = "http"
	Https Scheme = "https"
)

func (p Scheme) unwrap() str.CaseInsensitive {
	return str.CaseInsensitive(p)
}

// Return the default port for the given scheme: 80 for HTTP, 443 for HTTPs.
func DefaultPort(scheme Scheme) int {
	if scheme.unwrap().Eq(Https) {
		return DEFAULT_HTTPS_PORT
	}
	return DEFAULT_HTTP_PORT
}

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
	scheme      Scheme
	hostAndPort *HostAndPort
	path        []string
	query       url.Values
}

func (p httpUrl) Secure() bool {
	return p.scheme.unwrap().Eq(Https)
}

func (p httpUrl) Host() string {
	return p.hostAndPort.Host()
}

func (p httpUrl) Port() int {
	return p.hostAndPort.Port()
}

func (p httpUrl) Path() string {
	rooted := append([]string{"/"}, p.path...) // (*)
	return path.Join(rooted...)

	// NOTE. Start and end slashes.
	// We add a start slash since if it it wasn't there we've got to have
	// it. If p.path already starts with a slash, path.Join will ignore
	// it, so no harm done. But when it comes to end slashes, path.Join
	// insists on removing them---e.g. /a/b/ becomes /a/b. Not sure if
	// this is correct, but will keep it for now.
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
	ref := url.URL{
		Scheme:     string(p.scheme),
		Host:       p.hostAndPort.String(),
		Path:       p.Path(),
		ForceQuery: false,
	}
	refRepr := ref.String()

	if len(p.query) == 0 {
		return refRepr
	}
	queryRepr := p.query.Encode()
	return refRepr + "?" + queryRepr
}
