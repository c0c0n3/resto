package tree

import (
	"fmt"
	"reflect"
	"testing"
)

func mk5NodeTestTree() Tree[int] {
	return Node(1,
		Node(2,
			Node(3),
			Node(4),
		),
		Node(5),
	)
}

func stringStacker(label int, result string) string {
	return result + fmt.Sprintf("%d", label)
}

func TestFold5NodeTreeWithStringStacker(t *testing.T) {
	u := mk5NodeTestTree()
	got := Fold(stringStacker, "", u)
	want := "12345"

	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func TestFoldBreadth5NodeTreeWithStringStacker(t *testing.T) {
	u := mk5NodeTestTree()
	got := FoldBreadth(stringStacker, "", u)
	want := "12534"

	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

func Test5NodeTreeToList(t *testing.T) {
	got := ToList(mk5NodeTestTree())
	want := []int{1, 2, 3, 4, 5}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func Test5NodeTreeToBreadthList(t *testing.T) {
	got := ToBreadthList(mk5NodeTestTree())
	want := []int{1, 2, 5, 3, 4}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func TestToListReturnsEmptyOnNilInput(t *testing.T) {
	got := ToList[int](nil)
	want := []int{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}

func TestToBreadthListReturnsEmptyOnNilInput(t *testing.T) {
	got := ToBreadthList[int](nil)
	want := []int{}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v; got: %v", want, got)
	}
}
