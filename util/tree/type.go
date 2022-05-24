package tree

type Tree[X any] interface {
	Node(label X, ns ...Tree[X]) Tree[X]
	Label() X
	Children() []Tree[X]
}
