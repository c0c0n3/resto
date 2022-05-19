package str

import (
	"fmt"
	"strings"
)

// A case-insensitive string value. E.g. "Ab3Z" is the same as "aB3z".
type CaseInsensitive string

// Compare, disregarding case, the string representation of other with
// its own value.
//
// Examples.
//
//     CaseInsensitive("2").Eq(2) == true
//     CaseInsensitive("Ab3Z").Eq("aB3z") == true
//
func (p CaseInsensitive) Eq(other interface{}) bool {
	if other == nil {
		return false
	}
	lowercaseOther := strings.ToLower(fmt.Sprintf("%v", other))
	lowercaseSelf := strings.ToLower(string(p))

	return lowercaseOther == lowercaseSelf
}
