package file

import (
	"fmt"

	"github.com/c0c0n3/resto/util/err"
)

// A nil pointer error.
type NilPtr string

// An error happened while trying to make a path absolute.
type PathResolve string

// An error to signal a file system node isn't a directory.
type NotADirectory string

func nilVisitorErr() err.Err[NilPtr] {
	return err.Mk[NilPtr]("nil visitor")
}

func pathResolveWhitespaceErr() err.Err[PathResolve] {
	msg := "must be a non-empty, non-whitespace-only string"
	return err.Mk[PathResolve](msg)
}

func pathResolveAbsErr(cause error) err.Err[PathResolve] {
	return err.Mk[PathResolve](cause.Error())
}

func notADirErr(path AbsPath) err.Err[NotADirectory] {
	msg := "not a directory: %v"
	return err.Mk[NotADirectory](msg, path.Value())
}

// VisitError wraps any error that happened while traversing the target
// directory with an additional path to indicate where the error happened.
type VisitError struct {
	AbsPath string
	Err     error
}

// Error implements the standard error interface.
func (e VisitError) Error() string {
	return fmt.Sprintf("%s: %v", e.AbsPath, e.Err)
}

// Unwrap implements Go's customary error unwrapping.
func (e VisitError) Unwrap() error { return e.Err }
