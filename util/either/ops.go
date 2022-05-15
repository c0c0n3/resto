package either

import (
	"github.com/c0c0n3/resto/util/fnc"
)

// Apply f to e's right value if e has it otherwise propagate e's left
// value. In pseudo-code:
//
//     e = R v  ==>  R (f(v))
//     e = L v  ==>  L v
//
func Map[L, R, Y any](f fnc.F[R, Y], e Either[L, R]) Either[L, Y] {
	if e.IsRight() {
		y := f(e.Right())
		return MakeRight[L](y)
	}
	return MakeLeft[L, Y](e.Left())
}

// Apply f to e's right value if e has it otherwise propagate e's left
// value. In pseudo-code:
//
//     e = R v  ==>  f(v)
//     e = L v  ==>  L v
//
func Bind[L, R, Y any](f fnc.F[R, Either[L, Y]], e Either[L, R]) Either[L, Y] {
	if e.IsRight() {
		return f(e.Right())
	}
	return MakeLeft[L, Y](e.Left())
}
