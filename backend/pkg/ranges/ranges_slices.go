// Package ranges provides some set of functions to allow the usage
// of functional transformations in Go.
package ranges

// Each takes a slice, and applies a transformation function
// on each element, and aggregates into an array.
func Each[T any, R any](input []T, mapper func(T) R) []R {
	result := make([]R, 0)
	for _, val := range input {
		result = append(result, mapper(val))
	}
	return result
}

func EachAddress[T any, R any](input []T, mapper func(*T) R) []R {
	result := make([]R, 0)
	for _, val := range input {
		result = append(result, mapper(&val))
	}
	return result
}

func Filter[T any](input []T, predicate func(T) bool) []T {
	var result []T
	for _, val := range input {
		if predicate(val) {
			result = append(result, val)
		}
	}
	return result
}
