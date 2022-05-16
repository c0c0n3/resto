package err

import (
	"fmt"
)

// Generic error of type T where T is string-like.
// If T and U are different types so are E[T] and E[U]. Since E implements
// the error interface, you can make typed errors by just declaring a string
// wrapper type.
//
// Example.
//
//     TooBig string
//     TooSmall string
//
//     func checkBounds(x int) error {
//         if x < 1 {
//             return Mk[TooSmall]("got: %d", x)
//         }
//         if 10 < x {
//             return Mk[TooBig]("got: %d", x)
//         }
//         return nil
//     }
//
//     func whichErrorType(e error) string {
//         switch e.(type) {
//         case Err[TooBig]:
//             return "TooBig"
//         case Err[TooSmall]:
//             return "TooSmall"
//         default:
//             return "unknown"
//         }
//     }
//
type Err[T ~string] struct {
	wrapped T
}

// Make a new error of type T where T is string-like.
// The error message gets formatted according to the given specifier and
// arguments as if you called Sprintf.
func Mk[T ~string](format string, args ...any) Err[T] {
	msg := fmt.Sprintf(format, args...)
	return Err[T]{T(msg)}
}

func (e Err[T]) Error() string {
	return fmt.Sprintf("%v", e.wrapped)
}
