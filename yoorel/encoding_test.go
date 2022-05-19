package yoorel

import (
	"testing"
)

func TestWireFormatEscapes(t *testing.T) {
	want := "http://h:80/a%20b?ke+y=valu+e"
	got := EmptyBuilder().Http().HostAndPort(" h:80 ").
		JoinPath("a b").
		Query("ke y", "valu e").
		Build().Right().WireFormat()
	if got != want {
		t.Errorf("want: %s; got: %s", want, got)
	}
}

var rebuildFromWireFormatIsIdentityFixtures = []Builder{
	EmptyBuilder().Http().HostAndPort(" h ").JoinPath("a b").
		Query("ke y", "valu e"),
	EmptyBuilder().Https().HostAndPort("h:9090").JoinPath("a b/c").
		Query("key", "valu e").Query("x", "1"),
}

func TestRebuildFromWireFormatIsIdentity(t *testing.T) {
	for k, builder := range rebuildFromWireFormatIsIdentityFixtures {
		built := builder.Build().Right().WireFormat()
		parsed := BuilderFrom(built).Build().Right().WireFormat()
		if built != parsed {
			t.Errorf("[%d] built: %s; parsed: %s", k, built, parsed)
		}
	}
}
