package client

import (
	"net/url"
	"testing"

	e "github.com/c0c0n3/resto/util/err"
	"github.com/c0c0n3/resto/yoorel"
)

func checkParsedUrl(t *testing.T, want string, got e.ErrOr[yoorel.HttpUrl]) {
	if !got.IsRight() {
		t.Fatalf("want: parsed url; got: %v", got)
	}
	if got.Right().WireFormat() != want {
		t.Errorf("want: %s; got: %s", want, got.Right().WireFormat())
	}
}

func TestToHttpUrlFromString(t *testing.T) {
	want := "http://cape.town:8080/clifton"
	got := toHttpUrl(want)
	checkParsedUrl(t, want, got)
}

func TestToHttpUrlFromURL(t *testing.T) {
	want := "http://cape.town:8080/clifton"
	u := url.URL{
		Scheme: "http",
		Host:   "cape.town:8080",
		Path:   "/clifton",
	}
	got := toHttpUrl(u)
	checkParsedUrl(t, want, got)
}

func TestToHttpUrlFromURLPtr(t *testing.T) {
	want := "http://cape.town:8080/clifton"
	u, _ := url.Parse(want)
	got := toHttpUrl(u)
	checkParsedUrl(t, want, got)
}

func TestToHttpUrlFromWellFormedUrl(t *testing.T) {
	want := "http://cape.town:8080/clifton"
	u := yoorel.BuilderFrom(want).Build()
	got := toHttpUrl(WellFormedUrl{u.Right()})
	checkParsedUrl(t, want, got)
}

func TestToHttpUrlFromNil(t *testing.T) {
	var target *url.URL = nil
	got := toHttpUrl(target)
	if got.IsRight() {
		t.Errorf("want: error; got: %v", got)
	}
}
