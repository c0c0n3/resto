package fnc

// A function `f: X -> Y`.
type F[X, Y any] func(X) Y

// A predicate---i.e. a function `X -> bool`.
type Pred[X any] func(X) bool

// A fold operator to fold a list or other "foldable" data structures.
type FoldOp[X, Y any] func(X, Y) Y
