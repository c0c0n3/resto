package set

// Simple immutable set.
type Set[T comparable] interface {
	// Is the given value in this set?
	Member(v T) bool
	// Count how many elements this set holds.
	Size() int
	// Put the set elements in a list.
	List() []T
}

type hashSet[T comparable] struct {
	elements map[T]struct{}
}

// Create a map-backed Set from the given elements.
func FromList[T comparable](vs []T) Set[T] {
	set := &hashSet[T]{
		elements: make(map[T]struct{}, len(vs)),
	}
	for _, v := range vs {
		set.elements[v] = struct{}{}
	}
	return set
}

// Create a map-backed Set from the given elements.
func From[T comparable](vs ...T) Set[T] {
	return FromList(vs)
}

func (s *hashSet[T]) Member(v T) bool {
	_, hasKey := s.elements[v]
	return hasKey
}

func (s *hashSet[T]) Size() int {
	return len(s.elements)
}

func (s *hashSet[T]) List() []T {
	vs := make([]T, 0, len(s.elements))
	for v := range s.elements {
		vs = append(vs, v)
	}
	return vs
}
