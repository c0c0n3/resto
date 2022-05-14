package mime

import (
	"testing"

	"github.com/c0c0n3/resto/util/list"
	"github.com/c0c0n3/resto/util/set"
)

var allTypes = []MediaType{JSON, YAML}

func TestDistinctMediaTypes(t *testing.T) {
	allStrRepr := set.FromList(list.Map(MediaType.String, allTypes))
	want := len(allTypes)
	got := allStrRepr.Size()

	if want != got {
		t.Errorf("want: %d; got: %d", want, got)
	}
}
