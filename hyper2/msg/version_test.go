package msg

import (
	"fmt"
	"testing"

	"github.com/c0c0n3/resto/hyper2"
	e "github.com/c0c0n3/resto/util/err"
)

var validParseVersionFixtures = []struct {
	input string
	want  ProtocolVersion
}{
	{"HTTP/1", ProtocolVersion{1, 0}},
	{"Http/1", ProtocolVersion{1, 0}},
	{"HttP/1.1", ProtocolVersion{1, 1}},
	{"HTTP/21.32", ProtocolVersion{21, 32}},
	{"hTTp/21.32", ProtocolVersion{21, 32}},
	{"HTTP/2.32", ProtocolVersion{2, 32}},
	{"http/21.3", ProtocolVersion{21, 3}},
}

func TestParseValidVersion(t *testing.T) {
	for _, d := range validParseVersionFixtures {
		got, err := ParseProtocolVersion(d.input)
		if err != nil {
			t.Errorf("input: %s; want: %v; got error: %s", d.input, d.want, err)
		}
		if got != d.want {
			t.Errorf("input: %s; want: %v; got: %s", d.input, d.want, got)
		}
	}
}

var invalidParseVersionFixtures = []string{
	"", "H", "HTTP", "HTTP/", "HTTP/1.", "HTTPs/1.1", "HTTP/1.2.3",
}

func TestParseInvalidVersion(t *testing.T) {
	for _, raw := range invalidParseVersionFixtures {
		_, err := ParseProtocolVersion(raw)
		if _, ok := err.(e.Err[hyper2.UnparsableProtocolVersion]); !ok {
			t.Errorf("want: UnparsableProtocolVersion; got: %v", err)
		}
	}
}

func TestParseVersionThenStringIsId(t *testing.T) {
	for k := 0; k < 20; k++ {
		want := fmt.Sprintf("HTTP/%d.%d", k, 10*k+1)
		parsed, _ := ParseProtocolVersion(want)
		got := parsed.String()
		if got != want {
			t.Errorf("want: %s; got: %s", want, got)
		}
	}
}

func TestStringVersionThenParseIsId(t *testing.T) {
	for k := 0; k < 20; k++ {
		want := ProtocolVersion{k, 10*k + 1}
		raw := want.String()
		got, _ := ParseProtocolVersion(raw)
		if got != want {
			t.Errorf("want: %v; got: %v", want, got)
		}
	}
}

func TestVersion_1_1(t *testing.T) {
	got := HTTP_1_1()
	if got.Major() != 1 {
		t.Errorf("want: 1; got: %d", got)
	}
	if got.Minor() != 1 {
		t.Errorf("want: 1; got: %d", got)
	}
}

func TestVersion_2_0(t *testing.T) {
	got := HTTP_2_0()
	if got.Major() != 2 {
		t.Errorf("want: 2; got: %d", got)
	}
	if got.Minor() != 0 {
		t.Errorf("want: 0; got: %d", got)
	}
}

func TestVersion_3_0(t *testing.T) {
	got := HTTP_3_0()
	if got.Major() != 3 {
		t.Errorf("want: 3; got: %d", got)
	}
	if got.Minor() != 0 {
		t.Errorf("want: 0; got: %d", got)
	}
}
