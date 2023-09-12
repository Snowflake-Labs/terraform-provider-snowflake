package generator2

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

type TagBuilder struct {
	ddl []string
	sql []string
}

func Tags() *TagBuilder {
	return &TagBuilder{
		ddl: make([]string, 0),
		sql: make([]string, 0),
	}
}

func (v *TagBuilder) Static() *TagBuilder {
	v.ddl = append(v.ddl, "static")
	return v
}

func (v *TagBuilder) Keyword() *TagBuilder {
	v.ddl = append(v.ddl, "keyword")
	return v
}

func (v *TagBuilder) Parameter() *TagBuilder {
	v.ddl = append(v.ddl, "parameter")
	return v
}

func (v *TagBuilder) Identifier() *TagBuilder {
	v.ddl = append(v.ddl, "identifier")
	return v
}

func (v *TagBuilder) List() *TagBuilder {
	v.ddl = append(v.ddl, "list")
	return v
}

func (v *TagBuilder) NoParentheses() *TagBuilder {
	v.ddl = append(v.ddl, "no_parentheses")
	return v
}

func (v *TagBuilder) DDL(ddl ...string) *TagBuilder {
	v.ddl = append(v.ddl, ddl...)
	return v
}

func (v *TagBuilder) SQL(sql ...string) *TagBuilder {
	v.sql = append(v.sql, sql...)
	return v
}

func (v *TagBuilder) Build() map[string][]string {
	return map[string][]string{
		"ddl": v.ddl,
		"sql": v.sql,
	}
}
