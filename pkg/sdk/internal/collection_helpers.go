package internal

import (
	"errors"
)

var ErrObjectNotFound = errors.New("object does not exist")

func FindOne[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		if condition(o) {
			return &o, nil
		}
	}
	return nil, ErrObjectNotFound
}
