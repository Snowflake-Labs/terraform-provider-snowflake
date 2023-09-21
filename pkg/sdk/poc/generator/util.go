package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// IsNil is used for special cases where x != nil might not work (e.g. passing nil instead of interface implementation)
func IsNil(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return false
}

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

var (
	// Split by empty space or underscore
	splitSQLPattern   = regexp.MustCompile(`\s+|_`)
	englishLowerCaser = cases.Lower(language.English)
	englishTitleCaser = cases.Title(language.English)
)

func sqlToFieldName(sql string, shouldExport bool) string {
	sqlWords := splitSQLPattern.Split(sql, -1)
	for i, s := range sqlWords {
		if !shouldExport && i == 0 {
			sqlWords[i] = englishLowerCaser.String(s)
			continue
		}
		sqlWords[i] = englishTitleCaser.String(s)
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

func KindOfTPointer[T any]() string {
	return KindOfPointer(KindOfT[T]())
}

func KindOfTSlice[T any]() string {
	return KindOfSlice(KindOfT[T]())
}

func KindOfPointer(kind string) string {
	return "*" + kind
}

func KindOfSlice(kind string) string {
	return "[]" + kind
}

type TagBuilder struct {
	db  []string
	ddl []string
	sql []string
}

func Tags() *TagBuilder {
	return &TagBuilder{
		db:  make([]string, 0),
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

func (v *TagBuilder) DB(db ...string) *TagBuilder {
	v.db = append(v.db, db...)
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
	if len(v.db) > 0 {
		res["db"] = v.db
	}
	if len(v.ddl) > 0 {
		res["ddl"] = v.ddl
	}
	if len(v.sql) > 0 {
		res["sql"] = v.sql
	}
	return res
}
