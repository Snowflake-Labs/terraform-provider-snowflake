package gen

import (
	"fmt"
	"regexp"
	"strings"
)

var splitOnTheWordsBeginnings = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
var splitRemainingWordBreaks = regexp.MustCompile("([a-z0-9])([A-Z]+)")

func ToSnakeCase(str string) string {
	wordsSplit := splitOnTheWordsBeginnings.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(splitRemainingWordBreaks.ReplaceAllString(wordsSplit, "${1}_${2}"))
}

// TODO: test and describe
func TabularOutput(columnWidth int, columns ...string) {
	var sb strings.Builder
	for i, column := range columns {
		d, rem := DivWithRemainder(columnWidth-len(column), 8)
		tabs := d
		if rem != 0 {
			tabs++
		}
		sb.WriteString(column)
		if i != len(columns) {
			sb.WriteString(strings.Repeat("\t", tabs))
		}
	}
	fmt.Println(sb.String())
}

// TODO: test and describe
func DivWithRemainder(numerator int, denominator int) (int, int) {
	return numerator / denominator, numerator % denominator
}
