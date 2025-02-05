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

func twiceWhenGt10(x int) []int {
	if x > 10 {
		return []int{x, x}
	}
	return []int{}
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

var testConcatMapFixtures = []struct {
	input []int
	want  []int
}{
	{[]int{}, []int{}},
	{[]int{1}, []int{}},
	{[]int{1, 2}, []int{}},
	{[]int{1, 20}, []int{20, 20}},
	{[]int{11, 20}, []int{11, 11, 20, 20}},
	{IntList{}, []int{}},
	{IntList{1}, []int{}},
	{IntList{1, 2}, []int{}},
	{IntList{1, 20}, []int{20, 20}},
	{IntList{11, 20}, []int{11, 11, 20, 20}},
}

func TestConcatMap(t *testing.T) {
	for k, fixture := range testConcatMapFixtures {
		got := ConcatMap(twiceWhenGt10, fixture.input)
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

var testHeadFixtures = []struct {
	input []int
	want  int
}{
	{[]int{}, 0},
	{[]int{1}, 1},
	{[]int{1, 2}, 1},
	{[]int{2, 1, 3}, 2},
	{IntList{}, 0},
	{IntList{1}, 1},
	{IntList{1, 2}, 1},
	{IntList{2, 1, 3}, 2},
}

func TestHead(t *testing.T) {
	for k, fixture := range testHeadFixtures {
		got := Head(fixture.input)
		if got != fixture.want {
			t.Errorf("[%d] want: %v; got: %v", k, fixture.want, got)
		}
	}

	got := Head([]int{}, 5)
	if got != 5 {
		t.Errorf("want: 5; got: %d", got)
	}
}
