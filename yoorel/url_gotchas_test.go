package yoorel

import (
	"net/url"
	"path"
	"testing"
)

func TestUrlWireFormatWithHostAndPort(t *testing.T) {
	want := "//jou.ma:8080" // wish := "jou.ma:8080"
	u := &url.URL{
		Host: "jou.ma:8080",
	}
	got := u.String()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestUrlHostEncoding(t *testing.T) {
	want := "http://x%2520y" // wish := "http://x%20y"
	u := &url.URL{
		Scheme: "http",
		Host:   "x%20y",
	}
	got := u.String()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestUrlParseBareHostAndPort(t *testing.T) {
	parsed, err := url.Parse("host:80")
	if err != nil {
		t.Errorf("want: URL; got: %v", err)
	}
	if parsed.Scheme != "host" {
		t.Errorf("want: host; got: %s", parsed.Scheme)
	}
	if parsed.Opaque != "80" {
		t.Errorf("want: 80; got: %s", parsed.Opaque)
	}
	if parsed.Host != "" {
		t.Errorf("want: ''; got: %s", parsed.Host)
	}
	if parsed.Hostname() != "" {
		t.Errorf("want: ''; got: %s", parsed.Hostname())
	}
	if parsed.Port() != "" {
		t.Errorf("want: ''; got: %s", parsed.Port())
	}

	// From url.Parse docs:
	//
	//   Trying to parse a hostname and path without a scheme is invalid
	//   but may not necessarily return an error, due to parsing ambiguities.
}

func TestUrlParseUntrimmedBareHostAndPort(t *testing.T) {
	parsed, err := url.Parse(" host:80 ")
	if err == nil {
		t.Errorf("want: error; got: %v", parsed)
		// err -> parse " host:80 ":
		//        first path segment in URL cannot contain colon
	}

	// From url.Parse docs:
	//
	//   Trying to parse a hostname and path without a scheme is invalid
	//   but may not necessarily return an error, due to parsing ambiguities.
}

func TestUrlParseEncodedHost(t *testing.T) {
	want := "http://x%20y"
	_, err := url.Parse(want)
	if err == nil {
		t.Fatalf("want: parse error")
		// err -> parse error: parse "http://x%20y": invalid URL escape "%20"
		//
		// this error comes from unescape (line 221) b/c as stated in the
		// comments there:
		//   https://tools.ietf.org/html/rfc3986#page-21
		//   in the host component %-encoding can only be used
		//   for non-ASCII bytes.
	}
}

func TestUrlPathEncoding(t *testing.T) {
	want := "http://jou.ma/x%2520y" // wish := "http://jou.ma/x%20y"
	u := &url.URL{
		Scheme: "http",
		Host:   "jou.ma",
		Path:   "x%20y",
	}
	got := u.String()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestUrlParseEncodedPath(t *testing.T) {
	want := "http://jou.ma/x%20y"
	u, err := url.Parse(want)
	if err != nil {
		t.Fatalf("want: %s; got: %v", want, err)
	}
	got := u.String()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
	if u.Path != "/x y" {
		t.Errorf("want: /x y; got: %s", u.Path)
	}
}

func TestPathJoinNormalizes(t *testing.T) {
	want := "/x/y/z" // wish := "/x/y/z/"
	rooted := append([]string{"/", "/x"}, "y", "/z/")
	got := path.Join(rooted...)
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}
