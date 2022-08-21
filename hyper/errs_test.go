package hyper

import (
	"testing"

	e "github.com/c0c0n3/resto/util/err"
)

func TestErrors(t *testing.T) {
	var err e.Err[NilPtr]
	if NilRequestBuilderErr() == err {
		t.Errorf("want: nil req builder err; got: nil")
	}
	if NilResponseHandlerErr() == err {
		t.Errorf("want: nil res handler err; got: nil")
	}
}
