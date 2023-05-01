package msg

import (
	"strings"

	"github.com/c0c0n3/resto/util/str"
)

// Header and trailer fields of HTTP messages.
// See:
// - https://www.rfc-editor.org/rfc/rfc9110.html#name-fields

// An HTTP field.
// Fields are name/(list of) values pairs for conveying extra info about
// the sender, message, content, or context (headers) or for communicating
// info obtained while sending the content (trailers).
// A field name is case insensitive and there may be multiple fields with
// the same name in a message. Each field value can be either a single string
// or a comma-separated list of strings. When a field is repeated, its values
// get combined, respecting message order, in a single list of values.
// Example from the spec
//
//     Example-Field: Foo, Bar
//     Example-Field: Baz
//
// The field value for "Example-Field" is the list "Foo, Bar, Baz".
type Field struct {
	name   str.CaseInsensitive
	values []string
}

// Create a new field with the given name and values.
func NewField(name string, values []string) *Field {
	f := &Field{}
	f.name = str.CaseInsensitive(name)
	f.values = make([]string, len(values))
	copy(f.values, values)

	return f
}

// The field name.
func (f *Field) Name() str.CaseInsensitive {
	return f.name
}

// Has this field any value?
func (f *Field) Empty() bool {
	return len(f.values) == 0
}

// How many values does this field have?
func (f *Field) Size() int {
	return len(f.values)
}

// The first value of this field or the empty string if this field has
// no values.
func (f *Field) First() string {
	return f.Get(0)
}

// The value at the specified index or the empty string if the index
// is out of bounds.
func (f *Field) Get(index int) string {
	lastIx := len(f.values) - 1
	if index < 0 || lastIx < index {
		return ""
	}
	return f.values[index]
}

// Lookup table of name/value pairs for conveying extra info about the
// sender, message, content, or context (headers) or for communicating
// info obtained while sending the content (trailers).
type Fields map[string]*Field

// NOTE. keys will always be lowercase b/c FieldsFromMap is the only ctor.

// List all the field names.
func (fs Fields) Names() []string {
	names := make([]string, len(fs))
	i := 0
	for k := range fs {
		names[i] = k
		i++
	}
	return names
}

// Look up the field with the given name using a case insensitive name
// match. If there's no match, the returned field will be nil and the
// ok variable set to false.
func (fs Fields) Lookup(name string) (f *Field, ok bool) {
	n := strings.ToLower(name)
	f, ok = fs[n]
	return f, ok
}

// Create a Fields map from a regular map of names to value lists.
// If multiple keys in the input map resolve to the same case-insensitive
// key, their respective values get merged in a single field, not necessarily
// preserving order though.
func FieldsFromMap(m map[string][]string) Fields {
	fs := make(map[string]*Field, len(m))
	for k, vs := range m {
		name := strings.ToLower(k)
		if f, ok := fs[name]; ok {
			f.values = append(f.values, vs...)
		} else {
			fs[name] = NewField(k, vs)
		}
	}
	return fs
}
