package either

// Either an L-value or an R-value.
// We keep the Haskell convention of calling the L-value "left" and the
// R-value "right".
type Either[L, R any] interface {
	// Return the L-value if this is an L-value, otherwise the L zero value.
	Left() L
	// Return the R-value if this is an R-value, otherwise the R zero value.
	Right() R
	// Return the R-value if this is an R-value, otherwise the given argument.
	RightOr(r R) R
	// Is this an R-value?
	IsRight() bool
}

type left[L, R any] struct {
	value L
}

func (p left[L, R]) Left() L {
	return p.value
}

func (p left[L, R]) Right() R {
	var r R
	return r
}

func (p left[L, R]) RightOr(r R) R {
	return r
}

func (p left[L, R]) IsRight() bool {
	return false
}

type right[L, R any] struct {
	value R
}

func (p right[L, R]) Left() L {
	var l L
	return l
}

func (p right[L, R]) Right() R {
	return p.value
}

func (p right[L, R]) RightOr(r R) R {
	return p.value
}

func (p right[L, R]) IsRight() bool {
	return true
}

// Make an Either holding a left value of l.
func MakeLeft[L, R any](l L) Either[L, R] {
	return left[L, R]{l}
}

// Make an Either holding a right value of r.
func MakeRight[L, R any](r R) Either[L, R] {
	return right[L, R]{r}
}
