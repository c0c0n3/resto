package set

import (
	"testing"
)

func TestEmptySet(t *testing.T) {
	empty := From[int]()

	if empty.Size() != 0 {
		t.Errorf("want: 0; got: %d", empty.Size())
	}
	if empty.Member(1) {
		t.Errorf("want: no members; got: 1 is a member")
	}
	vs := empty.List()
	if len(vs) != 0 {
		t.Errorf("want: no members; got: %v", vs)
	}
}

func TestRemoveDups(t *testing.T) {
	set := From(1, 2, 2, 3, 4, 3)

	if set.Size() != 4 {
		t.Errorf("want: 4; got: %d", set.Size())
	}
	for k := 1; k < 5; k++ {
		if !set.Member(k) {
			t.Errorf("want: %d is a member; got: false", k)
		}
	}
	if set.Member(5) {
		t.Errorf("want: 5 not a member; got: 5 is a member")
	}
	vs := set.List()
	if len(vs) != 4 {
		t.Errorf("want: 4 members; got: %v", vs)
	}
}
