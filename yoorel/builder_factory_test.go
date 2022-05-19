package yoorel

import (
	"net/url"
	"testing"
)

func TestBuilderFactoryWithNoSeed(t *testing.T) {
	assert := func(buf Builder) {
		got := buf.Http().HostAndPort("h").Build()
		if !got.IsRight() {
			t.Errorf("want: URL; got: %v", got.Left())
		}
	}
	assert(BuilderFrom[url.URL]())
	assert(BuilderFrom[*url.URL]())
	assert(BuilderFrom[*url.URL](nil))
	assert(BuilderFrom[string]())
}

func TestBuilderFactoryPropagatesUrlParseErr(t *testing.T) {
	builder := BuilderFrom("\n")
	got := builder.Http().HostAndPort("h").Build()
	if got.IsRight() {
		t.Errorf("want: error; got: %v", got)
	}

	got = builder.Https().HostAndPort("k").JoinPath("a/b").Query("x", "1").
		Build()
	if got.IsRight() {
		t.Errorf("want: error; got: %v", got)
	}
}

var builderFactoryFixtures = []Builder{
	BuilderFrom("http://h"),
	BuilderFrom(url.URL{Scheme: "http", Host: "h"}),
	BuilderFrom(&url.URL{Scheme: "http", Host: "h"}),
}

func TestBuilderFactory(t *testing.T) {
	want := "http://h:80/"
	for k, builder := range builderFactoryFixtures {
		got := builder.Build()
		if !got.IsRight() {
			t.Errorf("[%d] want: URL; got: %v", k, got.Left())
		}
		gotWire := got.Right().WireFormat()
		if gotWire != want {
			t.Errorf("[%d] want: %s; got: %s", k, want, gotWire)
		}
	}
}
