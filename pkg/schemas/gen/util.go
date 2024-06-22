package gen

import (
	"regexp"
	"strings"
)

var splitOnTheWordsBeginnings = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
var splitRemainingWordBreaks = regexp.MustCompile("([a-z0-9])([A-Z]+)")

// TODO: describe
func ToSnakeCase(str string) string {
	wordsSplit := splitOnTheWordsBeginnings.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(splitRemainingWordBreaks.ReplaceAllString(wordsSplit, "${1}_${2}"))
}

// TODO: describe
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
