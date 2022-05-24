package tree

import (
	"sort"

	"github.com/c0c0n3/resto/util/fnc"
	"github.com/c0c0n3/resto/util/list"
)

// Does this node have any children?
func IsLeaf[X any](t Tree[X]) bool {
	if t == nil {
		return true
	}
	return len(t.Children()) == 0
}

// Convenience function to get t.Label().
func Label[X any](t Tree[X]) X {
	var x X
	if t != nil {
		x = t.Label()
	}
	return x
}

// Convenience function to get t.Children().
func Children[X any](t Tree[X]) []Tree[X] {
	if t == nil {
		return []Tree[X]{}
	}
	return t.Children()
}

// Traverse the tree in depth-first order, computing `op` at each step.
// In pseudo-code:
//
//     x1                            r5 where
//      |- x2            Fold            r1 = op(x1, seed)
//      |   |- x3    -------------\      r2 = op(x2, r1)
//      |   |- x4    -------------/      r3 = op(x3, r2)
//      |- x5                            r4 = op(x4, r3)
//                                       r5 = op(x5, r4)
//
func Fold[X, Y any](op fnc.FoldOp[X, Y], seed Y, t Tree[X]) Y {
	if t == nil {
		return seed
	}
	result := op(t.Label(), seed)
	for _, kid := range t.Children() {
		result = Fold(op, result, kid) // (*)
	}
	return result

	// NOTE. Recursion. Easy, but probably bad news in a language like
	// Go---stack overflow, anyone? Should we use iteration instead?
	// Notice recursion is a problem only when the tree height is huge...
}

// Traverse the tree in breadth-first order, computing `op` at each step.
// In pseudo-code:
//
//     x1                            r5 where
//      |- x2            Fold            r1 = op(x1, seed)
//      |   |- x3    -------------\      r2 = op(x2, r1)
//      |   |- x4    -------------/      r3 = op(x5, r2)
//      |- x5                            r4 = op(x3, r3)
//                                       r5 = op(x4, r4)
//
func FoldBreadth[X, Y any](op fnc.FoldOp[X, Y], seed Y, t Tree[X]) Y {
	return list.Fold(op, seed, ToBreadthList(t))
}

// Collect the tree labels in breadth-first order into a map keyed by
// depth level---root's level = 0, root children's level = 1, etc.
// In pseudo-code:
//
//     x1
//      |- x2           Levels           m[0] = [x1]
//      |   |- x3    -------------\      m[1] = [x2, x5]
//      |   |- x4    -------------/      m[2] = [x3, x4]
//      |- x5
//
func Levels[X any](t Tree[X]) map[int][]X {
	levelMap := make(map[int][]X)
	if t != nil {
		collectLevels(levelMap, 0, t)
	}
	return levelMap
}

func collectLevels[X any](levelMap map[int][]X, level int, ts ...Tree[X]) {
	if len(ts) == 0 {
		return
	}
	levelMap[level] = list.Map(Label[X], ts)
	collectLevels(levelMap, level+1, list.ConcatMap(Children[X], ts)...) // (*)

	// NOTE. Recursion. Easy, but probably bad news in a language like
	// Go---stack overflow, anyone? Should we use iteration instead?
	// Notice recursion is a problem only when the tree height is huge...
}

// NOTE. Expressing tree functions as folds.
// You could, but it's not worth my while in Go. See similar implementation
// note about list functions as folds in the list package.

func ToList[X any](t Tree[X]) []X {
	return nil // TODO
}

func ToBreadthList[X any](t Tree[X]) []X {
	if t == nil {
		return []X{}
	}

	levelMap := Levels(t)
	keys := make([]int, 0, len(levelMap))
	for k := range levelMap {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	labels := make([]X, 0, len(levelMap)) // there's at least that many children
	for k := range keys {
		labels = append(labels, levelMap[k]...)
	}
	return labels
}
