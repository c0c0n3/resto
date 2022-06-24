package tgz

import (
	"github.com/c0c0n3/resto/util/err"
)

// A nil pointer error.
type NilPtr string

// An attempt to access a closed resource handle, like a file or a stream.
type ClosedHandle string

func nilSinkErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil sink")
}

func nilSourceErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil source")
}

func nilEntryReaderErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil entry reader")
}

func closedReaderErr() err.Err[ClosedHandle] {
	return err.Mk[ClosedHandle]("closed reader")
}
