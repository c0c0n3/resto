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
//     type TooBig string
//     type TooSmall string
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
// arguments as if you called Sprintf and prefixed with the error type.
//
// Example.
//
//     package bounds
//
//     type TooBig string
//
//     func check(x int) {
//         if x > 10 {
//             err := Mk[TooBig]("check: %d", x)
//             fmt.Printf(err)
//         }
//     }
//     // check(15) prints: "bounds.TooBig: check: 15"
//
func Mk[T ~string](format string, args ...any) Err[T] {
	msg := fmt.Sprintf(format, args...)
	var t T
	taggedMsg := fmt.Sprintf("%T: %s", t, msg)

	return Err[T]{T(taggedMsg)}
}

func (e Err[T]) Error() string {
	return fmt.Sprintf("%v", e.wrapped)
}
