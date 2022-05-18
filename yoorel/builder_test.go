package yoorel

import (
	"net/url"
	"testing"
)

type UrlWrap string

func TestBuilderMissingHttp(t *testing.T) {
	got := EmptyBuilder().HostAndPort("jo.burg:8080").JoinPath("ben", "oni").
		Build()
	if got.IsRight() {
		t.Errorf("want: missing scheme; got: %v", got.Right())
	}
}

func TestBuilderMissingHttps(t *testing.T) {
	got := EmptyBuilder().HostAndPort("jo.burg:8080").JoinPath("ben", "oni").
		Build()
	if got.IsRight() {
		t.Errorf("want: missing scheme; got: %v", got.Right())
	}
}

var builderInvalidHostFixtures = []string{"\n", "\n:8080"}

func TestBuilderInvalidHost(t *testing.T) {
	for k, hp := range builderInvalidHostFixtures {
		got := EmptyBuilder().Http().HostAndPort(hp).JoinPath().Build()
		if got.IsRight() {
			t.Errorf("[%d] want: invalid host; got: %v", k, got.Right())
		}
	}
}

func TestBuilderInvalidPort(t *testing.T) {
	got := EmptyBuilder().Http().HostAndPort("h:70000").JoinPath().Build()
	if got.IsRight() {
		t.Errorf("want: invalid port; got: %v", got.Right())
	}
}

func TestBuilderFromNoSeed(t *testing.T) {
	assert := func(buf Builder) {
		got := buf.Http().HostAndPort("h").Build()
		if !got.IsRight() {
			t.Errorf("want: http://h/; got: %v", got.Left())
		}
	}
	assert(BuilderFrom[url.URL]())
	assert(BuilderFrom[*url.URL]())
	assert(BuilderFrom[string]())
	assert(BuilderFrom[UrlWrap]())
}

func TestBuildHostOnly(t *testing.T) {
	want := "http://h"
	got := EmptyBuilder().Http().HostAndPort("h").Build()
	if !got.IsRight() {
		t.Errorf("want: %s; got: %v", want, got.Left())
	}
	gotWire := got.Right().WireFormat()
	if gotWire != want {
		t.Errorf("want: %s; got: %v", want, gotWire)
	}
}

func TestBuildHostOnlyUseDefaultPort(t *testing.T) {
	want := "http://h"
	got := EmptyBuilder().Http().HostAndPort("h").Build().Right()
	gotWire := got.WireFormat()
	if gotWire != want {
		t.Errorf("want: %s; got: %v", want, gotWire)
	}
	if got.Port() != DEFAULT_HTTP_PORT {
		t.Errorf("want: %d; got: %d", DEFAULT_HTTP_PORT, got.Port())
	}

	want = "https://h"
	got = EmptyBuilder().Https().HostAndPort("h").Build().Right()
	gotWire = got.WireFormat()
	if gotWire != want {
		t.Errorf("want: %s; got: %v", want, gotWire)
	}
	if got.Port() != DEFAULT_HTTPS_PORT {
		t.Errorf("want: %d; got: %d", DEFAULT_HTTPS_PORT, got.Port())
	}
}

func TestUrlRootPath(t *testing.T) {
	got := EmptyBuilder().Http().HostAndPort("h").Build().Right()
	if got.Path() != "/" {
		t.Errorf("want: /; got: %s", got.Path())
	}

	got = EmptyBuilder().Http().HostAndPort("h").JoinPath("/").Build().Right()
	if got.Path() != "/" {
		t.Errorf("want: /; got: %s", got.Path())
	}
}

// var buildFromBaseUrlFixtures = []interface{}{}

// func TestBuildFromBaseUrl(t *testing.T) {
// 	for k, base := range buildFromBaseUrlFixtures {
// 		got := BuilderFrom[string](base) // TODO
// 		if !got.IsRight() {
// 			t.Errorf("[%d] want: ??; got: %v", k, got.Right())
// 		}
// 	}
// }
