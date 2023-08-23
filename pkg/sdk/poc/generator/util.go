package generator

import (
	"fmt"
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
