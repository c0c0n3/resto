package msg

import (
	"reflect"
	"testing"
)

func TestEmptyFields(t *testing.T) {
	fs := FieldsFromMap(nil)
	got := len(fs.Names())
	if got != 0 {
		t.Errorf("want: 0; got: %d", got)
	}
}

func TestMergedFields(t *testing.T) {
	m := map[string][]string{
		"a": {"1", "2"},
		"b": {"3"},
		"A": {"4"},
	}
	fs := FieldsFromMap(m)

	gotNames := fs.Names()
	if len(gotNames) != 2 {
		t.Fatalf("want: 2 names; got: %v", gotNames)
	}
	if got, ok := fs.Lookup("a"); ok {
		if got.Size() != 3 {
			t.Errorf("want: sz(a) = 3; got: %d", got.Size())
		}
		vs := make([]string, 3)
		for k := 0; k < 3; k++ {
			vs[k] = got.Get(k)
		}
		if !reflect.DeepEqual(vs, []string{"1", "2", "4"}) &&
			!reflect.DeepEqual(vs, []string{"4", "1", "2"}) {
			t.Errorf("want: merge; got: %v", vs)
		}
	} else {
		t.Errorf("want: a; got: %v", got)
	}
	if got, ok := fs.Lookup("B"); ok {
		if got.Size() != 1 {
			t.Errorf("want: sz(b) = 1; got: %d", got.Size())
		}
		if got.First() != "3" {
			t.Errorf("want: b[0] = 3; got: %s", got.First())
		}
	} else {
		t.Errorf("want: b; got: %v", got)
	}
}

func TestEmptyField(t *testing.T) {
	got := NewField("A", nil)
	if !got.Name().Eq("a") {
		t.Errorf("want: a = A")
	}
	if got.Size() != 0 {
		t.Errorf("want: sz = 0; got: %d", got.Size())
	}
	if !got.Empty() {
		t.Errorf("want: empty")
	}
	if got.Get(0) != "" {
		t.Errorf("want: f[0] = ''; got: %s", got.Get(0))
	}
	if got.Get(1) != "" {
		t.Errorf("want: f[1] = ''; got: %s", got.Get(1))
	}
}
