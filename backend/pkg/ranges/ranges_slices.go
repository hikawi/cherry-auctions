// Package ranges provides some set of functions to allow the usage
// of functional transformations in Go.
package ranges

// Each takes a slice, and applies a transformation function
// on each element, and aggregates into an array.
func Each[T any, R any](input []T, mapper func(T) R) []R {
	var result []R
	for _, val := range input {
		result = append(result, mapper(val))
	}
	return result
}
