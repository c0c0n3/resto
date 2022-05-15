package either

import (
	"fmt"
	"testing"
)

func doubleOdd(x int) (int, error) {
	if x%2 == 1 {
		return x * 2, nil
	}
	return 0, fmt.Errorf("not an odd integer: %d", x)
}

func digitToString(x int) (string, error) {
	if 0 <= x && x < 10 {
		return fmt.Sprintf("%d", x), nil
	}
	return "", fmt.Errorf("not a digit: %d", x)
}

func TestBindERightValues(t *testing.T) {
	e1 := FromResult(doubleOdd(3))
	e2 := BindE(digitToString, e1)
	got := e2.Right()

	if !e2.IsRight() {
		t.Errorf("want: right; got: %v", e2)
	}
	if got != "6" {
		t.Errorf("want: '6'; got: %v", got)
	}
}

func TestBindEPropagatesError(t *testing.T) {
	e1 := FromResult(doubleOdd(2))
	e2 := BindE(digitToString, e1)

	if e2.IsRight() {
		t.Errorf("want: left; got: %v", e2)
	}
	if e2.Left() == nil {
		t.Errorf("want: error; got: %v", e2)
	}
}

func TestWrapProcPropagatesError(t *testing.T) {
	e1 := FromResult(doubleOdd(7))
	e2 := BindE(digitToString, e1)

	if e2.IsRight() {
		t.Errorf("want: left; got: %v", e2)
	}
	if e2.Left() == nil {
		t.Errorf("want: error; got: %v", e2)
	}
}
