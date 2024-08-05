package helpers

import (
	"strings"
)

const ResourceIdDelimiter = '|'

func ParseResourceIdentifier(identifier string) []string {
	if identifier == "" {
		return make([]string, 0)
	}
	return strings.Split(identifier, string(ResourceIdDelimiter))
}

func EncodeResourceIdentifier(parts ...string) string {
	return strings.Join(parts, string(ResourceIdDelimiter))
}
