package tree

// Tryna port Haskell's excellent Data.Tree over to Go, but I don't
// think I'm happy w/ the result. I wish Go was more Haskell-like...

// Unfolder generates the data to build a Tree in breadth-first order.
// Given a seed, it returns the label of the current node and the seeds
// to use for building its children. If the list of seeds is nil or empty,
// then the current node is a leaf.
type Unfolder[X, S any] func(seed S) (X, []S)

// Use the given tree data generator to build a Tree in breadth-first order.
// Example.
//
//     func sevenNodeBinTree(seed int) (label string, nextSeeds []int) {
//         label = fmt.Sprintf("%d", seed)
//         if 2*seed+1 < 7 {
//             nextSeeds = []int{2*seed + 1, 2*seed + 2}
//         }
//         return label, nextSeeds
//     }
//
//     func mkTree() Tree[string] {
//         return Unfold(sevenNodeBinTree, 0)
//     }
//
//     mkTree() ~~>      0
//                      / \
//                     1   2
//                   / |   | \
//                  3  4   5  6
//
func Unfold[X, S any](f Unfolder[X, S], seed S) Tree[X] {
	label, childSeeds := f(seed)
	children := unfoldForest(f, childSeeds)
	return Node(label, children...)
}

func unfoldForest[X, S any](f Unfolder[X, S], seeds []S) []Tree[X] {
	if seeds == nil {
		seeds = []S{}
	}
	nodes := make([]Tree[X], 0, len(seeds))
	for _, seed := range seeds {
		n := Unfold(f, seed) // (*)
		nodes = append(nodes, n)
	}
	return nodes

	// NOTE. Mutual recursion. Easy, but probably bad news in a language
	// like Go---stack overflow, anyone? Should we use iteration instead?
	// Notice recursion is a problem only when the tree height is huge...
}
