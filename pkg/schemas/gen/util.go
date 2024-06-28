package gen

import (
	"regexp"
	"strings"
)

var (
	splitOnTheWordsBeginnings = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	splitRemainingWordBreaks  = regexp.MustCompile("([a-z0-9])([A-Z]+)")
)

// ToSnakeCase allows converting a CamelCase text to camel_case one (needed for schema attribute names). Examples:
// - CamelCase -> camel_case
// - ACamelCase -> a_camel_case
// - URLParser -> url_parser
// - Camel1Case -> camel1_case
// - camelCase -> camel_case
// - camelURL -> camel_url
// - camelURLSomething -> camel_url_something
func ToSnakeCase(str string) string {
	wordsSplit := splitOnTheWordsBeginnings.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(splitRemainingWordBreaks.ReplaceAllString(wordsSplit, "${1}_${2}"))
}

// ColumnOutput is a helper to align a tabular output with the given columnWidth, e.g. (for 20 spaces column width):
// Name                string              string
// State               sdk.WarehouseState  string
func ColumnOutput(columnWidth int, columns ...string) string {
	var sb strings.Builder
	for i, column := range columns {
		sb.WriteString(column)
		if i != len(columns)-1 {
			spaces := max(columnWidth-len(column), 1)
			sb.WriteString(strings.Repeat(" ", spaces))
		}
	}
	return sb.String()
}
