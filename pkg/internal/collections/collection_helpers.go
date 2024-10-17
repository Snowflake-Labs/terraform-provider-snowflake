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

func MapErr[T any, R any](collection []T, mapper func(T) (R, error)) ([]R, error) {
	result := make([]R, len(collection))
	errs := make([]error, 0)
	for i, elem := range collection {
		value, err := mapper(elem)
		if err != nil {
			errs = append(errs, err)
		}
		result[i] = value
	}
	return result, errors.Join(errs...)
}
