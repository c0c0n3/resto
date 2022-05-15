package either

import (
	"testing"
)

func TestLeftReturnsRightZeroValue(t *testing.T) {
	e := MakeLeft[string, int]("left")
	if e.Right() != 0 {
		t.Errorf("want: 0; got: %v", e.Right())
	}
}

func TestLeftReturnsRightGivenValue(t *testing.T) {
	e := MakeLeft[string, int]("left")
	got := e.RightOr(2)
	if got != 2 {
		t.Errorf("want: 2; got: %v", got)
	}
}

func TestRightReturnsLeftZeroValue(t *testing.T) {
	e := MakeRight[string](2)
	if e.Left() != "" {
		t.Errorf("want: empty string; got: %v", e.Left())
	}
}

func TestRightIgnoresRightGivenValue(t *testing.T) {
	e := MakeRight[string](2)
	got := e.RightOr(3)
	if got != 2 {
		t.Errorf("want: 2; got: %v", got)
	}
}
