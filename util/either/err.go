package either

import (
	"github.com/c0c0n3/resto/util/fnc"
)

// Convert the output of a typical Go procedure returning an error-value
// pair to an Either value. If there's an error, it gets mapped to a left
// whereas the return value gets mapped to a right.
func FromResult[V any](value V, err error) Either[error, V] {
	if err != nil {
		return MakeLeft[error, V](err)
	}
	return MakeRight[error](value)
}

// Convert a typical Go procedure returning an error-value pair to one
// returning an Either value instead. If the procedure returns an error,
// the wrapper returns that error as a left value. Otherwise the wrapper
// returns the procedure's output as a right value.
func WrapProc[X, Y any](f fnc.Procedure[X, Y]) func(X) Either[error, Y] {
	return func(x X) Either[error, Y] {
		if y, err := f(x); err != nil {
			return MakeLeft[error, Y](err)
		} else {
			return MakeRight[error](y)
		}
	}
}

// Specialise Bind to the case where the function to bind is a typical
// Go procedure returning an error-value pair. The error becomes the
// left and the output the right.
func BindE[V, U any](f fnc.Procedure[V, U], e Either[error, V]) Either[error, U] {
	return Bind(WrapProc(f), e)
}
