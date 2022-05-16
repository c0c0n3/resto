package list

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/c0c0n3/resto/util/fnc"
)

type IntList []int

func mapAsFold[X, Y any, L ~[]X](f fnc.F[X, Y], xs L) []Y {
	seed := make([]Y, 0, len(xs))
	op := func(x X, ys []Y) []Y {
		ys = append(ys, f(x))
		return ys
	}
	return Fold(op, seed, xs)
}

func intToString(x int) string {
	return fmt.Sprintf("%d", x)
}

var testMapFixtures = []struct {
	input []int
	want  []string
}{
	{[]int{}, []string{}},
	{[]int{1}, []string{"1"}},
	{[]int{1, 2}, []string{"1", "2"}},
	{IntList{}, []string{}},
	{IntList{1}, []string{"1"}},
	{IntList{1, 2}, []string{"1", "2"}},
}

func TestMap(t *testing.T) {
	for k, fixture := range testMapFixtures {
		got := Map(intToString, fixture.input)
		if !reflect.DeepEqual(fixture.want, got) {
			t.Errorf("[%d] want: %v; got: %v", k, fixture.want, got)
		}
	}
}

func TestMapAsFold(t *testing.T) {
	for k, fixture := range testMapFixtures {
		want := Map(intToString, fixture.input)
		got := mapAsFold(intToString, fixture.input)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("[%d] want: %v; got: %v", k, want, got)
		}
	}
}

func TestReduceAsFold(t *testing.T) {
	sum := func(x, y int) int { return x + y }
	want := 1 + 2 + 3
	got := Fold(sum, 0, []int{1, 2, 3})

	if want != got {
		t.Errorf("want: %d; got: %d", want, got)
	}
}

var testFilterFixtures = []struct {
	input []int
	want  []int
}{
	{[]int{}, []int{}},
	{[]int{1}, []int{}},
	{[]int{1, 2}, []int{2}},
	{[]int{1, 2, 3}, []int{2}},
	{[]int{1, 2, 3, 4}, []int{2, 4}},
	{IntList{}, []int{}},
	{IntList{1}, []int{}},
	{IntList{1, 2}, []int{2}},
	{IntList{1, 2, 3}, []int{2}},
	{IntList{1, 2, 3, 4}, []int{2, 4}},
}

func TestFilter(t *testing.T) {
	even := func(x int) bool { return x%2 == 0 }
	for k, fixture := range testFilterFixtures {
		got := Filter(even, fixture.input)
		if !reflect.DeepEqual(fixture.want, got) {
			t.Errorf("[%d] want: %v; got: %v", k, fixture.want, got)
		}
	}
}

var testReverseFixtures = []struct {
	input []int
	want  []int
}{
	{[]int{}, []int{}},
	{[]int{1}, []int{1}},
	{[]int{1, 2}, []int{2, 1}},
	{[]int{1, 2, 3}, []int{3, 2, 1}},
	{IntList{}, []int{}},
	{IntList{1}, []int{1}},
	{IntList{1, 2}, []int{2, 1}},
	{IntList{1, 2, 3}, []int{3, 2, 1}},
}

func TestReverse(t *testing.T) {
	for k, fixture := range testReverseFixtures {
		got := Reverse(fixture.input)
		if !reflect.DeepEqual(fixture.want, got) {
			t.Errorf("[%d] want: %v; got: %v", k, fixture.want, got)
		}
	}
}
