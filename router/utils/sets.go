package utils

// Set keeps a collection of unique items with efficient access for addition, removal and check existence of elements.
type Set[T comparable] map[T]bool

// Has checks if an element exist in the set.
// returns true if the element exist and false otherwise.
func (s Set[T]) Has(value T) bool { return s[value] }

// Add an element to the set.
// returns true if the element is a new element in the set and false if the element was already part of the set.
func (s Set[T]) Add(value T) bool {
	exist := s[value]
	s[value] = true
	return !exist
}

// Remove an element from the set.
// returns true if removed and false if the element was not on the set to begin with.
func (s Set[T]) Remove(value T) bool {
	exist := s[value]
	delete(s, value)
	return exist
}

// Clone creates a new copy of the set
func (s Set[T]) Clone() Set[T] {
	clone := make(Set[T], len(s))
	for key, value := range s {
		clone[key] = value
	}
	return clone
}

// NewSet creates a new Set and populate it with the provided elements. duplicate elements will be saved once
func NewSet[T comparable](elements ...T) Set[T] {
	set := make(Set[T], 0)
	for _, element := range elements {
		set.Add(element)
	}
	return set
}
