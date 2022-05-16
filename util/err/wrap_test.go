package err

import (
	"fmt"
	"testing"
)

type TooBig string
type TooSmall string
type Bad string
type Wrong string

func checkBounds(x int) error {
	if x < 1 {
		return Mk[TooSmall]("got: %d", x)
	}
	if 10 < x {
		return Mk[TooBig]("got: %d", x)
	}
	return nil
}

func whichErrorType(e error) string {
	switch e.(type) {
	case Err[Bad]:
		return "Bad"
	case Err[Wrong]:
		return "Wrong"
	case Err[TooBig]:
		return "TooBig"
	case Err[TooSmall]:
		return "TooSmall"
	default:
		return "unknown"
	}
}

func TestErrMsg(t *testing.T) {
	want := "something bad happened!"
	bad := Mk[Bad](want)

	if got := fmt.Sprintf("%v", bad); got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
	if got := bad.Error(); got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestErrMsgWithParams(t *testing.T) {
	want := "was expecting 4 but got: 2"
	wrong := Mk[Wrong]("was expecting 4 but got: %d", 2)

	if got := fmt.Sprintf("%v", wrong); got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
	if got := wrong.Error(); got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestSwitchOnErrType(t *testing.T) {
	bad := Mk[Bad]("")
	wrong := Mk[Wrong]("")

	if got := whichErrorType(bad); got != "Bad" {
		t.Errorf("want: Bad; got: %s", got)
	}
	if got := whichErrorType(wrong); got != "Wrong" {
		t.Errorf("want: Wrong; got: %s", got)
	}
	if got := whichErrorType(checkBounds(-1)); got != "TooSmall" {
		t.Errorf("want: TooSmall; got: %s", got)
	}
	if got := whichErrorType(checkBounds(100)); got != "TooBig" {
		t.Errorf("want: TooBig; got: %s", got)
	}
}
