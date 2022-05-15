package fnc

// A function `f: X -> Y`.
type F[X, Y any] func(X) Y

// A procedure you feed with an X-value and get back either a result of
// type Y or an error. Lots of Go procedures look like this because of
// error handling conventions.
type Procedure[X, Y any] func(X) (Y, error)

// A predicate---i.e. a function `X -> bool`.
type Pred[X any] func(X) bool

// A fold operator to fold a list or other "foldable" data structures.
type FoldOp[X, Y any] func(X, Y) Y
