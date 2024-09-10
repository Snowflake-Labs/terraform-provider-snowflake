package util

import "strings"

// TrimAllPrefixes removes all prefixes from the input. Order matters.
func TrimAllPrefixes(text string, prefixes ...string) string {
	result := text
	for _, prefix := range prefixes {
		result = strings.TrimPrefix(result, prefix)
	}
	return result
}
