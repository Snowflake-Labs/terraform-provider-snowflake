package collections

import (
	"errors"
	"fmt"
)

var ErrObjectNotFound = errors.New("object does not exist")

func FindOne[T any](collection []T, condition func(T) bool) (*T, error) {
	for _, o := range collection {
		fmt.Printf("%v\n", o)
		if condition(o) {
			return &o, nil
		}
	}
	return nil, ErrObjectNotFound
}
