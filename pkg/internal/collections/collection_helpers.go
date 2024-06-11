package collections

import (
	"errors"
)

var ErrObjectNotFound = errors.New("object does not exist")

// TODO [SNOW-1473414]: move collection helpers fully with a separate PR
func FindOne[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		if condition(o) {
			return &o, nil
		}
	}
	return nil, ErrObjectNotFound
}
