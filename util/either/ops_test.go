package either

import (
	"fmt"
	"testing"
)

func intToString(x int) string {
	return fmt.Sprintf("%d", x)
}

func intToStringE(x int) Either[bool, string] {
	return MakeRight[bool](intToString(x))
}

func TestMapPropagatesLeft(t *testing.T) {
	e := MakeLeft[bool, int](true)
	got := Map(intToString, e)

	if got.IsRight() {
		t.Errorf("want: left; got: %v", got)
	}
	if !got.Left() {
		t.Errorf("want: true; got: %v", got.Left())
	}
}

func TestMapRight(t *testing.T) {
	e := MakeRight[bool](2)
	got := Map(intToString, e)

	if !got.IsRight() {
		t.Errorf("want: right; got: %v", got)
	}
	if got.Right() != "2" {
		t.Errorf("want: '2'; got: %v", got.Right())
	}
}

func TestBindPropagatesLeft(t *testing.T) {
	e := MakeLeft[bool, int](true)
	got := Bind(intToStringE, e)

	if got.IsRight() {
		t.Errorf("want: left; got: %v", got)
	}
	if !got.Left() {
		t.Errorf("want: true; got: %v", got.Left())
	}
}

func TestBindRight(t *testing.T) {
	e := MakeRight[bool](2)
	got := Bind(intToStringE, e)

	if !got.IsRight() {
		t.Errorf("want: right; got: %v", got)
	}
	if got.Right() != "2" {
		t.Errorf("want: '2'; got: %v", got.Right())
	}
}

func TestBindRightNewLeft(t *testing.T) {
	e := MakeRight[int](2)
	got := Bind(MakeLeft[int, int], e)

	if got.IsRight() {
		t.Errorf("want: left; got: %v", got)
	}
	if got.Left() != 2 {
		t.Errorf("want: 2; got: %v", got.Left())
	}
}
