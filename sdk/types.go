package sdk

import "strings"

// String returns a pointer to the given string.
func String(v string) *string {
	return &v
}

// StringSlice returns a pointer to the give strings.
func StringSlice(v []string) *[]string {
	return &v
}

// Bool returns a pointer to the given bool
func Bool(v bool) *bool {
	return &v
}

// Int returns a pointer to the given int.
func Int(v int) *int {
	return &v
}

// Int32 returns a pointer to the given int32.
func Int32(v int32) *int32 {
	return &v
}

// Int64 returns a pointer to the given int64.
func Int64(v int64) *int64 {
	return &v
}

// addQuote adds quotes for every string
func addQuote(raw []string) []string {
	quoted := []string{}
	for _, item := range raw {
		item = strings.Trim(item, "'")
		quoted = append(quoted, "'"+item+"'")
	}
	return quoted
}
