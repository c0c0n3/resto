package yoorel

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/c0c0n3/resto/util/err"
)

type builderBuffer struct {
	scheme      string
	hostAndPort string
	path        []string
	query       url.Values
}

func (p builderBuffer) Http() Builder {
	p.scheme = string(Http)
	return p
}

func (p builderBuffer) Https() Builder {
	p.scheme = string(Https)
	return p
}

func (p builderBuffer) HostAndPort(hp string) Builder {
	p.hostAndPort = strings.TrimSpace(hp)
	return p
}

func (p builderBuffer) JoinPath(ps ...string) Builder {
	p.path = append(p.path, ps...)
	return p
}

func (p builderBuffer) Query(key, value string) Builder {
	p.query.Add(key, value)
	return p
}

func (p builderBuffer) Build() err.ErrOr[HttpUrl] {
	r := &httpUrl{
		path:  p.path,
		query: p.query,
	}
	er := p.buildScheme(r)
	return err.Bind(p.buildHostAndPort, er)
}

func (p builderBuffer) buildScheme(r *httpUrl) err.ErrOr[*httpUrl] {
	if Http.unwrap().Eq(p.scheme) {
		r.scheme = Http
	}
	if Https.unwrap().Eq(p.scheme) {
		r.scheme = Https
	}

	var e error = nil
	if p.scheme == "" {
		e = err.Mk[InvalidUrl]("not an absolute URL: '%v'", p)
	} else if !(Http.unwrap().Eq(p.scheme) || Https.unwrap().Eq(p.scheme)) {
		e = err.Mk[InvalidUrl]("not a valid scheme: '%v'", p)
	}

	return err.FromResult(r, e)
}

func (p builderBuffer) buildHostAndPort(r *httpUrl) (HttpUrl, error) {
	raw := p.hostAndPort
	if !HasPort(raw) {
		port := DefaultPort(r.scheme)
		raw = fmt.Sprintf("%s:%d", raw, port)
	}

	hp, e := ParseHostAndPort(raw)
	r.hostAndPort = hp

	return r, e
}

type builderErr struct {
	initUrlErr error
}

func (p builderErr) Http() Builder {
	return p
}

func (p builderErr) Https() Builder {
	return p
}

func (p builderErr) HostAndPort(hp string) Builder {
	return p
}

func (p builderErr) JoinPath(ps ...string) Builder {
	return p
}

func (p builderErr) Query(key, value string) Builder {
	return p
}

func (p builderErr) Build() err.ErrOr[HttpUrl] {
	return err.FromResult[HttpUrl](nil, p.initUrlErr)
}
