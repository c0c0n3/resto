package yoorel

import (
	"net/url"
	"path"

	"github.com/c0c0n3/resto/util/err"
)

type Builder interface {
	Http() Builder
	Https() Builder
	HostAndPort(hp string) Builder
	JoinPath(ps ...string) Builder // starts from /; appends wada/wada/
	Query(key, value string) Builder
	Build() err.ErrOr[HttpUrl]
}

type builderImpl struct {
	refBuf   *url.URL
	queryBuf url.Values
	buildErr error
}

type InitialBuilderUrl interface {
	~string | *url.URL | url.URL
}

func BuilderFrom[T InitialBuilderUrl](seed ...T) Builder {
	if len(seed) == 0 {
		return EmptyBuilder()
	}
	switch partialUrl := any(seed[0]).(type) {
	case url.URL:
		return builderFromUrl(&partialUrl)
	case *url.URL:
		return builderFromUrl(partialUrl)
	case string:
		return builderFromRawUrl(partialUrl)
	default:
		return EmptyBuilder()
	}
}

func EmptyBuilder() builderImpl {
	return builderImpl{
		refBuf:   &url.URL{},
		queryBuf: make(url.Values),
	}
}

func builderFromUrl(partialUrl *url.URL) builderImpl {
	if partialUrl == nil {
		return EmptyBuilder()
	}
	return builderImpl{
		refBuf:   partialUrl,
		queryBuf: partialUrl.Query(),
	}
}

func builderFromRawUrl(rawUrl string) builderImpl {
	parsed, err := url.Parse(rawUrl)
	return builderImpl{
		refBuf:   parsed,
		queryBuf: parsed.Query(),
		buildErr: err,
	}
}

func (p builderImpl) Http() Builder {
	p.refBuf.Scheme = "http"
	return p
}

func (p builderImpl) Https() Builder {
	p.refBuf.Scheme = "https"
	return p
}

func (p builderImpl) HostAndPort(hp string) Builder {
	p.refBuf.Host = hp
	return p
}

func (p builderImpl) JoinPath(ps ...string) Builder {
	rooted := append([]string{"/", p.refBuf.Path}, ps...) // (*)
	p.refBuf.Path = path.Join(rooted...)
	return p

	// NOTE. Start and end slashes.
	// We add a start slash since if it it wasn't there we've got to have
	// it. If refBuf.Path already starts with a slash, path.Join will ignore
	// it, so no harm done. But when it comes to end slashes, path.Join
	// insists on removing them---e.g. /a/b/ becomes /a/b. Not sure if
	// this is correct, but will keep it for now.
}

func (p builderImpl) Query(key, value string) Builder {
	p.queryBuf.Add(key, value)
	return p
}

func (p builderImpl) Build() err.ErrOr[HttpUrl] {
	r := httpUrl{p.refBuf, p.queryBuf}
	e := p.validate()
	return err.FromResult[HttpUrl](r, e)
}

func (p builderImpl) validate() error {
	if p.buildErr != nil {
		return p.buildErr
	}
	if !p.refBuf.IsAbs() {
		return err.Mk[InvalidUrl]("not an absolute URL: '%s'", p.refBuf)
	}
	if p.refBuf.Scheme != "http" && p.refBuf.Scheme != "https" {
		// This can happen if you instantiate the builder with a raw URL
		// that isn't HTTP and then you don't call the Http* methods.
		return err.Mk[InvalidUrl]("not an HTTP URL: '%s'", p.refBuf)
	}
	if HasPort(p.refBuf.Host) {
		if _, err := ParseHostAndPort(p.refBuf.Host); err != nil {
			// ParseHostAndPort implements extra checks currently missing
			// from url.Parse.
			return err
		}
	} else {
		if err := IsHostname(p.refBuf.Host); err != nil {
			// ParseHostAndPort calls IsHostname, so we check this in
			// the case there's no port.
			return err
		}
	}
	return nil
}
