package yoorel

import (
	"fmt"
	"testing"
)

var invalidHostnameFixtures = []string{
	"", "\n", ":", ":80", "some.host:", "some host", "some host.com",
	"what?is.this", "em@il", "what.the.h*ll",
	"x1234567890123456789012345678901234567890123456789012345678901234.com",
}

func TestInvalidHostname(t *testing.T) {
	for k, d := range invalidHostnameFixtures {
		if err := IsHostname(d); err == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var validHostnameFixtures = []string{
	"::123", "1.2.3.4", "_h.com", "a-b.some_where", "some.host",
	"x12345678901234567890123456789012345678901234567890123456789012.com",
}

func TestValidHostname(t *testing.T) {
	for k, d := range validHostnameFixtures {
		if err := IsHostname(d); err != nil {
			t.Errorf("[%d] want: valid; got: %v", k, err)
		}
	}
}

var invalidHostnameAndPortFixtures = []string{
	"", "\n", ":", ":80", "some.host:", "some host:80", "some.host:123456789",
}

func TestInvalidHostnameAndPort(t *testing.T) {
	for k, d := range invalidHostnameAndPortFixtures {
		if _, err := ParseHostAndPort(d); err == nil {
			t.Errorf("[%d] want: error; got: valid", k)
		}
	}
}

var parseHostAndPortFixtures = []struct {
	in       string
	wantHost string
	wantPort int
}{
	{"h:0", "h", 0}, {"h:1", "h", 1}, {"h:65535", "h", 65535},
	{"[::123]:0", "::123", 0}, {"[::123]:1", "::123", 1},
	{"[::123]:65535", "::123", 65535},
	{"1.2.3.4:0", "1.2.3.4", 0}, {"1.2.3.4:1", "1.2.3.4", 1},
	{"1.2.3.4:65535", "1.2.3.4", 65535},
}

func TestParseHostAndPort(t *testing.T) {
	for k, d := range parseHostAndPortFixtures {
		if hp, err := ParseHostAndPort(d.in); err != nil {
			t.Errorf("[%d] want: valid parse; got: %v", k, err)
		} else {
			if d.wantHost != hp.Host() || d.wantPort != hp.Port() {
				t.Errorf("[%d] want: %s:%d; got: %v",
					k, d.wantHost, d.wantPort, hp)
			}

			repr := fmt.Sprintf("%s:%d", d.wantHost, d.wantPort)
			if repr != hp.String() {
				t.Errorf("[%d] want string repr: %s; got: %v", k, repr, hp)
			}
		}
	}
}
