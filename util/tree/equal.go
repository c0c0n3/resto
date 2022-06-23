package tree

import "reflect"

// Are the given two trees semantically equivalent? i.e. do they have
// the same structure and do corresponding labels are equal?
// We use reflect.DeepEqual to compare labels.
func Equal[X any](t, u Tree[X]) bool {
	if t == nil && u == nil {
		return true
	}
	if t == nil || u == nil {
		return false
	}
	if !reflect.DeepEqual(t.Label(), u.Label()) {
		return false
	}
	tc, uc := t.Children(), u.Children()
	if len(tc) != len(uc) {
		return false
	}
	for k, c := range tc {
		if !Equal(c, uc[k]) {
			return false
		}
	}
	return true
}

// NOTE. Why not use reflect.DeepEqual to compare the trees? In principle,
// it should work. In practice, it doesn't and I've got no clue why.
