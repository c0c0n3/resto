package msg

import (
	"testing"
)

func TestStatusCodeEquality(t *testing.T) {
	if StatusCode(201) == StatusCode(200) {
		t.Errorf("want: not equal; got: equal")
	}
	if StatusCode(200) != 200 {
		t.Errorf("want: equal; got: not equal")
	}
}

func TestStatusCodeAsInt(t *testing.T) {
	if StatusCode(200).Value() != 200 {
		t.Errorf("want: equal; got: not equal")
	}
}
