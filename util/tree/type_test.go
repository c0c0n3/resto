package tree

import (
	"reflect"
	"testing"
)

func TestIsLeaf(t *testing.T) {
	if IsLeaf(mk5NodeTestTree()) {
		t.Errorf("want: 5 node tree not leaf")
	}
	if !IsLeaf(Node(1)) {
		t.Errorf("want: single node tree is leaf")
	}
	if !IsLeaf[int](nil) {
		t.Errorf("want: nil is leaf")
	}
}

func TestChildrenOfNilIsEmptyList(t *testing.T) {
	got := Children[int](nil)
	want := []Tree[int]{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}
