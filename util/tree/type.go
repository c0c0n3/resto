package tree

// A multi-way tree where each node can have zero or more children.
// A Tree value is a node in the tree. It can be a leaf node if it
// has no children otherwise it's an inner node.
type Tree[X any] interface {
	// This node's label.
	Label() X
	// This node's children.
	Children() []Tree[X]
}

type mwTree[X any] struct {
	label    X
	children []Tree[X]
}

func (t *mwTree[X]) Label() X {
	return t.label
}

func (t *mwTree[X]) Children() []Tree[X] {
	return t.children
}

// Make a node with the given label and children.
func Node[X any](label X, ns ...Tree[X]) Tree[X] {
	return &mwTree[X]{
		label:    label,
		children: ns,
	}
}

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
