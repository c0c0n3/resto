package yoorel

import (
	"testing"
)

func TestToURLFromNil(t *testing.T) {
	got := ToURL(nil)
	if got == nil {
		t.Fatalf("want: url; got: nil")
	}
	if got.Host != "" || got.Path != "" {
		t.Errorf("want: empty; got: %v", got)
	}
}

func TestToURLFromValue(t *testing.T) {
	value := EmptyBuilder().Https().HostAndPort("h").JoinPath("p").
		Query("x", "1").Query("y", "2").
		Build().Right()
	want := "https://h:443/p?x=1&y=2"
	got := ToURL(value).String()

	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}
