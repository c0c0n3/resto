package list

import (
	"github.com/c0c0n3/resto/util/fnc"
)

// Fold a list. In pseudo-code:
//
//     Fold(op, y, [x1, x2, ...]) = op(...(op(x2, op(x1, y))...)
//
func Fold[X, Y any](op fnc.FoldOp[X, Y], seed Y, xs []X) Y {
	result := seed
	for _, x := range xs {
		result = op(x, result)
	}
	return result
}

// Apply f to each element of the list. In pseudo-code:
//
//     Map(f, [x1, x2, ...]) = [f(x1), f(x2), ...]
//
func Map[X, Y any](f fnc.F[X, Y], xs []X) []Y {
	ys := make([]Y, len(xs))
	for k, x := range xs {
		ys[k] = f(x)
	}
	return ys
}

// NOTE. Expressing list functions as folds.
// Most list functions can be implemented as a fold---this has to do with
// fold capturing the pattern of primitive recursion, but I'm digressing.
// I'm not doing this in Go though since the code gets seriously ugly.
// Here's a *working* implementation of Map in terms of Fold:
//
//      func Map[X, Y any](f fnc.F[X, Y], xs []X) []Y {
// 	        seed := make([]Y, 0, len(xs))
// 	        op := func(x X, ys []Y) []Y {
// 		        ys = append(ys, f(x))
// 		        return ys
// 	        }
// 	        return Fold(op, seed, xs)
//      }
//
// For a comparison, this is how you could define map in terms of fold
// in Haskell
//
//     map f = foldr ((:) . f) []
//

// Keep only the list elements x such that p(x). In pseudo-code
//
//     Filter(p, [x1, x2, ...]) = [x2, ...]
//
// where p(x1) = false, p(x2) = true, etc.
func Filter[X any](p fnc.Pred[X], xs []X) []X {
	ys := make([]X, 0, len(xs))
	for _, x := range xs {
		if p(x) {
			ys = append(ys, x)
		}
	}
	return ys
}
