package tree

import "testing"

func TestNilsAreEqual(t *testing.T) {
	if !Equal[int](nil, nil) {
		t.Errorf("want: nil == nil; got: false")
	}
}

func TestNilIsNotEqualToNonNil(t *testing.T) {
	if Equal(nil, Node(1)) {
		t.Errorf("want: nil != Node(1); got: false")
	}
	if Equal(Node(1), nil) {
		t.Errorf("want: Node(1) != nil; got: false")
	}
}

func TestTreesWithDifferentLabelsAreNotEqual(t *testing.T) {
	if Equal(Node(1), Node(2)) {
		t.Errorf("want: Node(1) != Node(2); got: false")
	}

	u := Node(0, Node(1), Node(2))
	w := Node(0, Node(2), Node(1))
	if Equal(u, w) {
		t.Errorf("want: u != w; got: false")
	}
}

func TestTreesWithDifferentStructureAreNotEqual(t *testing.T) {
	u := Node(0, Node(1), Node(2))
	w := Node(0, Node(1), Node(2, Node(3), Node(4)))
	if Equal(u, w) {
		t.Errorf("want: u != w; got: false")
	}
}
