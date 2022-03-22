package utils

func Keys[K comparable, V any, M ~map[K]V](m M) []K {
	keys := make([]K, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	return keys
}

func Values[K comparable, V any, M ~map[K]V](m M) []V {
	values := make([]V, len(m))
	i := 0
	for _, value := range m {
		values[i] = value
		i++
	}
	return values
}

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

func Entries[K comparable, V any, M ~map[K]V](m M) []Entry[K, V] {
	values := make([]Entry[K, V], len(m))
	i := 0
	for key, value := range m {
		values[i] = Entry[K, V]{Key: key, Value: value}
		i++
	}
	return values
}

func FromEntries[K comparable, V any](entries []Entry[K, V]) map[K]V {
	m := make(map[K]V, len(entries))
	for _, entry := range entries {
		m[entry.Key] = entry.Value
	}
	return m
}

func Clone[K comparable, V any, M ~map[K]V](m M) M {
	clone := make(M, len(m))
	for key, value := range m {
		clone[key] = value
	}
	return m
}
