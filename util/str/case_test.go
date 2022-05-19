package str

import (
	"testing"
)

type StrWrap string

func TestEqWithNil(t *testing.T) {
	target := CaseInsensitive("")
	if target.Eq(nil) {
		t.Errorf("want: %v != nil", target)
	}
}

func TestEqWithWrapper(t *testing.T) {
	target := CaseInsensitive("Ab3Z")
	other := StrWrap("aB3z")
	if !target.Eq(other) {
		t.Errorf("want: %v == %v", target, other)
	}

	other = StrWrap("aB3z.")
	if target.Eq(other) {
		t.Errorf("want: %v != %v", target, other)
	}
}

func TestEqWithString(t *testing.T) {
	target := CaseInsensitive("Ab3Z")
	other := "aB3z"
	if !target.Eq(other) {
		t.Errorf("want: %v == %v", target, other)
	}

	other = "aB3z."
	if target.Eq(other) {
		t.Errorf("want: %v != %v", target, other)
	}
}

func TestEqWithOtherType(t *testing.T) {
	target := CaseInsensitive("3")
	other := 3
	if !target.Eq(other) {
		t.Errorf("want: %v == %v", target, other)
	}

	other = 4
	if target.Eq(other) {
		t.Errorf("want: %v != %v", target, other)
	}
}
