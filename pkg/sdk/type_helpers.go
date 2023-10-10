package sdk

import (
	"strconv"
)

// String returns a pointer to the given string.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the given bool.
func Bool(b bool) *bool {
	return &b
}

// toBool converts a string to a bool.
func toBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		panic(err)
	}
	return b
}

// Int returns a pointer to the given int.
func Int(i int) *int {
	return &i
}

// toInt converts a string to an int.
func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

// Float64 returns a pointer to the given float64.
func Float64(f float64) *float64 {
	return &f
}

// toFloat64 converts a string to a float64.
func toFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

// Pointer is a generic function that returns a pointer to a given value.
func Pointer[K any](v K) *K {
	return &v
}
