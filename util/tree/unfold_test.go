package tree

import (
	"fmt"
	"testing"
)

func makeTestBinTree() Tree[string] {
	return Node("0",
		Node("1",
			Node("3"),
			Node("4"),
		),
		Node("2",
			Node("5"),
			Node("6"),
		),
	)
}

func testBinTreeUnfolder(seed int) (label string, nextSeeds []int) {
	label = fmt.Sprintf("%d", seed)
	if 2*seed+1 < 7 {
		nextSeeds = []int{2*seed + 1, 2*seed + 2}
	}
	return label, nextSeeds
}

func TestUnfoldBinTree(t *testing.T) {
	want := makeTestBinTree()
	got := Unfold(testBinTreeUnfolder, 0)

	if !Equal(want, got) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}
