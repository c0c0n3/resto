package yoorel

import (
	"reflect"
	"sort"
	"testing"
)

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

func TestBuilderFromUnsupportedScheme(t *testing.T) {
	got := BuilderFrom("ftp://h/").Build()
	if got.IsRight() {
		t.Errorf("want: unsupported scheme; got: %v", got.Right())
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

func TestBuildHostOnly(t *testing.T) {
	want := "http://h:80/"
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
	want := "http://h:80/"
	got := EmptyBuilder().Http().HostAndPort("h").Build().Right()
	gotWire := got.WireFormat()
	if gotWire != want {
		t.Errorf("want: %s; got: %v", want, gotWire)
	}
	if got.Port() != DEFAULT_HTTP_PORT {
		t.Errorf("want: %d; got: %d", DEFAULT_HTTP_PORT, got.Port())
	}

	want = "https://h:443/"
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

func TestUrlAppendPath(t *testing.T) {
	want := "/a/b"
	got := EmptyBuilder().Http().HostAndPort("h").
		JoinPath("a/").JoinPath("/b").
		Build().Right()
	if got.Path() != want {
		t.Errorf("want: %s; got: %s", want, got.Path())
	}
}

func TestBuildQuery(t *testing.T) {
	got := EmptyBuilder().Https().HostAndPort("h").
		Query("x", "1").Query("x", "2").Query("y", "3").
		Build()
	if !got.IsRight() {
		t.Errorf("want: URL; got: %v", got.Left())
	}

	result := got.Right()
	keys := result.QueryKeys()
	if keys.Size() != 2 {
		t.Errorf("want: 2; got: %d", keys.Size())
	}
	if !keys.Member("x") {
		t.Errorf("want: x in key set")
	}
	if !keys.Member("y") {
		t.Errorf("want: y in key set")
	}

	xs := result.QueryValues("x")
	sort.Strings(xs)
	if !reflect.DeepEqual(xs, []string{"1", "2"}) {
		t.Errorf("want: {1 2}; got: %v", xs)
	}

	ys := result.QueryValues("y")
	if !reflect.DeepEqual(ys, []string{"3"}) {
		t.Errorf("want: {3}; got: %v", ys)
	}

	empty := result.QueryValues("not-there")
	if !reflect.DeepEqual(empty, []string{}) {
		t.Errorf("want: {}; got: %v", empty)
	}
}

func TestExtendBaseUrl(t *testing.T) {
	got := BuilderFrom("https://h/?x=1").JoinPath("a").Query("y", "2").
		Build().Right()
	if !got.Secure() {
		t.Errorf("want: HTTPs; got: HTTP")
	}
	if got.Host() != "h" {
		t.Errorf("want: h; got: %s", got.Host())
	}
	if got.Path() != "/a" {
		t.Errorf("want: /a; got: %s", got.Path())
	}

	xs := got.QueryValues("x")
	if !reflect.DeepEqual(xs, []string{"1"}) {
		t.Errorf("want: {1}; got: %v", xs)
	}

	ys := got.QueryValues("y")
	if !reflect.DeepEqual(ys, []string{"2"}) {
		t.Errorf("want: {2}; got: %v", ys)
	}
}
