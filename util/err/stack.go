package err

import (
	"fmt"
	"strings"

	"github.com/c0c0n3/resto/util/list"
)

// A stack to collect errors from most (top of the stack) to least recent
// (bottom of the stack). This type implements the error interface by
// listing errors one per line, from most to least recent ones.
type ErrStack struct {
	stack []error
}

// Create a new error stack.
func Stack() ErrStack {
	return ErrStack{stack: make([]error, 0)}
}

// Push one or more errors to the stack.
func (p ErrStack) Push(es ...error) ErrStack {
	p.stack = append(p.stack, es...)
	return p
}

func (p ErrStack) Error() string {
	buf := strings.Builder{}
	for _, e := range list.Reverse(p.stack) {
		line := fmt.Sprintf("%v\n", e)
		buf.WriteString(line)
	}
	return buf.String()
}
