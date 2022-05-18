package list

import (
	"github.com/c0c0n3/resto/util/fnc"
)

// Fold a list. In pseudo-code:
//
//     Fold(op, y, [x1, x2, ...]) = op(...(op(x2, op(x1, y))...)
//
func Fold[X, Y any, L ~[]X](op fnc.FoldOp[X, Y], seed Y, xs L) Y {
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
// This operation doesn't modify the input list.
func Map[X, Y any, L ~[]X](f fnc.F[X, Y], xs L) []Y {
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

// Keep only the list elements x such that p(x). In pseudo-code:
//
//     Filter(p, [x1, x2, ...]) = [x2, ...]
//
// where p(x1) = false, p(x2) = true, etc.
// This operation doesn't modify the input list.
func Filter[X any, L ~[]X](p fnc.Pred[X], xs L) []X {
	ys := make([]X, 0, len(xs))
	for _, x := range xs {
		if p(x) {
			ys = append(ys, x)
		}
	}
	return ys
}

// Reverse the input list. In pseudo-code:
//
//     Reverse([x1, x2, ...]) = [..., x2, x1]
//
// This operation doesn't modify the input list.
func Reverse[X any, L ~[]X](xs L) []X {
	lx := len(xs)
	ys := make([]X, lx)
	for k, x := range xs {
		ys[lx-k-1] = x
	}
	return ys
}

// Get the first element of the input list if the list isn't empty,
// otherwise return def if provided. If the list is empty and no def
// was passed in, just return X's zero value.
func Head[X any, L ~[]X](xs L, def ...X) X {
	if len(xs) > 0 {
		return xs[0]
	}
	if len(def) > 0 {
		return def[0]
	}
	var zero X
	return zero
}
