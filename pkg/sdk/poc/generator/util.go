package generator

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode/utf8"
)

func startingWithLowerCase(s string) string {
	firstLetter, _ := utf8.DecodeRuneInString(s)
	return strings.ToLower(string(firstLetter)) + s[1:]
}

func startingWithUpperCase(s string) string {
	firstLetter, _ := utf8.DecodeRuneInString(s)
	return strings.ToUpper(string(firstLetter)) + s[1:]
}

func wrapWith(s string, with string) string {
	return fmt.Sprintf("%s%s%s", with, s, with)
}

func sqlToFieldName(sql string, shouldExport bool) string {
	sqlWords := strings.Split(sql, " ")
	for i, s := range sqlWords {
		if !shouldExport && i == 0 {
			sqlWords[i] = cases.Lower(language.English).String(s)
			continue
		}
		sqlWords[i] = cases.Title(language.English).String(s)
	}
	return strings.Join(sqlWords, "")
}
