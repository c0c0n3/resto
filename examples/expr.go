// Solving the expression problem with tagless final interpreters in Go.
//
// Can you solve the expression problem in Go 1.18 using typed tagless
// final interpreters? It looks like you can. The program below is just
// a simple port to Go of the brilliant Haskell solution in these excellent
// lecture notes:
//
// * http://okmij.org/ftp/tagless-final/course/lecture.pdf
//
// The oke who wrote that really has a knack for making complex stuff
// understandable, read it first before trying to make sense of the Go
// code...
//
// We use tagless final interpreters in a few places in resto, so I'm
// parking this code here for future reference.

package examples

import "fmt"

// initial expression type
type expr[R any] interface {
	lit(v int) R
	add(x R, y R) R
}

// arithmetic expr

type num int

func numEval() num {
	return num(0)
}

//lint:ignore U1000 actually used but the linter isn't smart enough to see that.
func (p num) lit(v int) num {
	return num(v)
}

//lint:ignore U1000 actually used but the linter isn't smart enough to see that.
func (p num) add(x num, y num) num {
	return x + y
}

// expr pretty printing as sexpr

type str string

func strEval() str {
	return str("")
}

//lint:ignore U1000 actually used but the linter isn't smart enough to see that.
func (p str) lit(v int) str {
	return str(fmt.Sprintf("%d", v))
}

//lint:ignore U1000 actually used but the linter isn't smart enough to see that.
func (p str) add(x str, y str) str {
	return str(fmt.Sprintf("(+ %s %s)", x, y))
}

// builds a generic expr you can evaluate in different ways
func exampleExpr[R expr[R]](r R) R {
	return r.add(r.lit(1), r.add(r.lit(2), r.lit(3)))
}

// expr extension

type mult[R any] interface {
	mult(x R, y R) R
}

//lint:ignore U1000 actually used but the linter isn't smart enough to see that.
func (p num) mult(x num, y num) num {
	return x * y
}

//lint:ignore U1000 actually used but the linter isn't smart enough to see that.
func (p str) mult(x str, y str) str {
	return str(fmt.Sprintf("(* %s %s)", x, y))
}

type extendedExpr[R any] interface {
	expr[R]
	mult[R]
}

// builds a generic extendedExpr you can evaluate in different ways
func exampleExtendedExpr[R extendedExpr[R]](r R) R {
	return r.mult(r.lit(-4), r.add(r.lit(1), r.add(r.lit(2), r.lit(3))))
}

func RunExprProblem() {
	fmt.Println("expression problem w/ tagless final")

	numExpr := exampleExpr(numEval())
	strExpr := exampleExpr(strEval())

	fmt.Printf("arithmetic evaluation: %v\n", numExpr)
	fmt.Printf("pretty printing: %v\n", strExpr)

	extendedNumExpr := exampleExtendedExpr(numEval())
	extendedStrExpr := exampleExtendedExpr(strEval())

	fmt.Printf("extended arithmetic evaluation: %v\n", extendedNumExpr)
	fmt.Printf("extended pretty printing: %v\n", extendedStrExpr)
}
