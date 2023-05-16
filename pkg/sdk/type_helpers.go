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
