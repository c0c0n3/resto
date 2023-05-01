package msg

import (
	"testing"
)

func TestMethodEquality(t *testing.T) {
	if GET != Method("GET") {
		t.Errorf("want: equal; got: not equal")
	}
	if GET != "GET" {
		t.Errorf("want: equal; got: not equal")
	}
}

func TestMethodAsString(t *testing.T) {
	if GET.String() != "GET" {
		t.Errorf("want: equal; got: not equal")
	}
}
