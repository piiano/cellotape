package utils

func Find[T any](elements []T, findFn func(T) bool) (T, bool) {
	for _, el := range elements {
		if findFn(el) {
			return el, true
		}
	}
	var empty T
	return empty, false
}

func IndexOf[T any](elements []T, findFn func(T) bool) int {
	for i, el := range elements {
		if findFn(el) {
			return i
		}
	}
	return -1
}

func LastIndexOf[T any](elements []T, findFn func(T) bool) int {
	for i := len(elements) - 1; i >= 0; i-- {
		if findFn(elements[i]) {
			return i
		}
	}
	return -1
}

func Filter[T any](elements []T, filterFn func(T) bool) []T {
	filtered := make([]T, 0)
	for _, element := range elements {
		if filterFn(element) {
			filtered = append(filtered, element)
		}
	}
	return filtered
}

func Map[T, R any](elements []T, mapFn func(T) R) []R {
	mappedElements := make([]R, len(elements))
	for i, element := range elements {
		mappedElements[i] = mapFn(element)
	}
	return mappedElements
}

func ConcatSlices[T any](slices ...[]T) []T {
	target := make([]T, 0)
	for _, s := range slices {
		target = append(target, s...)
	}
	return target
}

func Ptr[T any](value T) *T {
	return &value
}
