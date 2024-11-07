package collections

import (
	"errors"
)

var ErrObjectNotFound = errors.New("object does not exist")

func FindFirst[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		if condition(o) {
			return &o, nil
		}
	}
	return nil, ErrObjectNotFound
}

func Map[T any, R any](collection []T, mapper func(T) R) []R {
	result := make([]R, len(collection))
	for i, elem := range collection {
		result[i] = mapper(elem)
	}
	return result
}

// TODO(SNOW-1479870): Test
// MergeMaps takes any number of maps (of the same type) and concatenates them.
// In case of key collision, the value will be selected from the map that is provided
// later in the src function parameter.
func MergeMaps[M ~map[K]V, K comparable, V any](src ...M) M {
	merged := make(M)
	for _, m := range src {
		for k, v := range m {
			merged[k] = v
		}
	}
	return merged
}
