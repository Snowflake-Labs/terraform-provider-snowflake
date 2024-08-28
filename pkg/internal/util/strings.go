package util

import "strings"

func TrimAllPrefixes(text string, prefixes ...string) string {
	result := text
	for _, prefix := range prefixes {
		result = strings.TrimPrefix(result, prefix)
	}
	return result
}
