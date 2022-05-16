package err

import (
	"testing"
)

func TestErrStackFromMostToLeastRecent(t *testing.T) {
	got := Stack().Push(Mk[Bad]("e1")).Push(Mk[Wrong]("e2")).Error()
	want := "e2\ne1\n"

	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestErrStackPushEmpty(t *testing.T) {
	got := Stack().Push().Error()
	if got != "" {
		t.Errorf("want: empty; got: %s", got)
	}
}

func TestErrStackPushMany(t *testing.T) {
	got := Stack().Push(Mk[Bad]("e1"), Mk[Wrong]("e2")).Error()
	want := "e2\ne1\n"

	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}
