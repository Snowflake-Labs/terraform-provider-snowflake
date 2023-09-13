package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"os"
	"path/filepath"
	"reflect"
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

func WriteCodeToFile(buffer *bytes.Buffer, fileName string) {
	wd, errWd := os.Getwd()
	if errWd != nil {
		log.Panicln(errWd)
	}
	outputPath := filepath.Join(wd, fileName)
	src, errSrcFormat := format.Source(buffer.Bytes())
	if errSrcFormat != nil {
		log.Panicln(errSrcFormat)
	}
	if err := os.WriteFile(outputPath, src, 0o600); err != nil {
		log.Panicln(err)
	}
}

func KindOfT[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	return t.Name()
}

func KindOfPointer(kind string) string {
	return "*" + kind
}

func KindOfSlice(kind string) string {
	return "[]" + kind
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
	res := make(map[string][]string)
	if len(v.ddl) > 0 {
		res["ddl"] = v.ddl
	}
	if len(v.sql) > 0 {
		res["sql"] = v.sql
	}
	return res
}
