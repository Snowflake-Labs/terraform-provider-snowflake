package sdk

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

// couple of helper functions
func parentheses(s string) string {
	return fmt.Sprintf("(%s)", s)
}

type quoteType string

const (
	NoQuotes     quoteType = "no_quotes"
	DoubleQuotes quoteType = "double_quotes"
	SingleQuotes quoteType = "single_quotes"
)

func (qt quoteType) Quote(v interface{}) string {
	s := fmt.Sprintf("%v", v)
	switch qt {
	case NoQuotes:
		return s
	case DoubleQuotes:
		escapedString := strings.ReplaceAll(s, qt.String(), qt.String()+qt.String())
		return fmt.Sprintf(`%v%v%v`, qt.String(), escapedString, qt.String())
	case SingleQuotes:
		escapedString := strings.Trim(s, qt.String())
		return fmt.Sprintf(`%v%v%v`, qt.String(), escapedString, qt.String())
	default:
		return s
	}
}

func (qt quoteType) String() string {
	switch qt {
	case NoQuotes:
		return ""
	case DoubleQuotes:
		return "\""
	case SingleQuotes:
		return "'"
	}
	return ""
}

// getQuoteTypeFromTag returns the quote type from a struct tag.
func getQuoteTypeFromTag(tag reflect.StructTag, tagName string) quoteType {
	t := strings.ToLower(tag.Get(tagName))
	if t == "" {
		return NoQuotes
	}

	parts := strings.Split(t, ",")
	for _, part := range parts {
		if strings.Contains(part, "quotes") {
			return quoteType(strings.TrimSpace(part))
		}
	}
	return NoQuotes
}

func getUseParenthesesFromTag(tag reflect.StructTag, tagName string, defaultParentheses bool) bool {
	t := strings.ToLower(tag.Get(tagName))
	if t == "" {
		return defaultParentheses
	}
	parts := strings.Split(t, ",")
	for _, part := range parts {
		switch strings.TrimSpace(part) {
		case "parentheses":
			return true
		case "no_parentheses":
			return false
		}
	}
	return defaultParentheses
}

type sqlBuilder struct{}

// sql builds a SQL statement from sqlClauses.
func (b *sqlBuilder) sql(clauses ...sqlClause) string {
	// remove nil and empty strings
	sList := make([]string, 0)
	for _, c := range clauses {
		if c != nil && c.String() != "" {
			sList = append(sList, c.String())
		}
	}

	return strings.Trim(strings.Join(sList, " "), " ")
}

// parseStruct parses a struct and returns a slice of sqlClauses.
func (b *sqlBuilder) parseStruct(s interface{}) ([]sqlClause, error) {
	clauses := make([]sqlClause, 0)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct, got %s", v.Kind())
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// unexported fields need to be handled separately.
		if !value.CanInterface() {
			fieldClauses, err := b.parseUnexportedField(field, value)
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, fieldClauses...)
			continue
		}

		// skip nil pointers for attributes, since they are not set. Otherwise dereference them
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}

		if value.Kind() == reflect.Slice {
			// check if there is any keyword
			ddlTag := field.Tag.Get("ddl")
			if ddlTag != "" {
				ddlTagParts := strings.Split(ddlTag, ",")
				ddlType := ddlTagParts[0]
				switch ddlType {
				case "keyword":
					clauses = append(clauses, sqlKeywordClause{
						value: field.Tag.Get("db"),
						qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
					})
				case "list":
					listClauses := make([]sqlClause, 0)
					// loop through the slice call parseStruct on each element (since the elements could be structs)
					for i := 0; i < value.Len(); i++ {
						v := value.Index(i).Interface()
						// test if v is an ObjectIdentifier. If it is it needs to be handled separately
						objectIdentifer, ok := v.(ObjectIdentifier)
						if ok {
							listClauses = append(listClauses, sqlIdentifierClause{
								value: objectIdentifer,
							})
							continue
						}
						structClauses, err := b.parseStruct(value.Index(i).Interface())
						if err != nil {
							return nil, err
						}
						// each element of the slice needs to be pre-rendered before the commas are added
						renderedStructClauses := b.sql(structClauses...)
						sClause := sqlStaticClause(renderedStructClauses)
						listClauses = append(listClauses, sClause)
					}
					if len(listClauses) < 1 {
						continue
					}
					clauses = append(clauses, sqlListClause{
						clauses:        listClauses,
						sep:            ",",
						useParentheses: getUseParenthesesFromTag(field.Tag, "ddl", true),
						keyword:        field.Tag.Get("db"),
					})
				}
			}
		}

		if value.Kind() == reflect.Struct {
			// check if there is any keyword on the struct
			// if there is, then we need to add it to the clause
			// if there is not, then we need to recurse into the struct
			// and get the clauses from there
			ddlTag := field.Tag.Get("ddl")
			if ddlTag != "" {
				ddlTagParts := strings.Split(ddlTag, ",")
				ddlType := ddlTagParts[0]
				switch ddlType {
				case "keyword":
					clauses = append(clauses, sqlKeywordClause{
						value: field.Tag.Get("db"),
						qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
					})
				case "identifier":
					if value.Interface().(ObjectIdentifier).FullyQualifiedName() == "" {
						continue
					}
					clauses = append(clauses, sqlIdentifierClause{
						key:   field.Tag.Get("db"),
						value: value.Interface().(ObjectIdentifier),
					})
				}
			}
			structClauses, err := b.parseStruct(value.Interface())
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, structClauses...)
			continue
		}

		// default case, if not a struct then it is a field
		fieldClauses, err := b.parseField(field, value)
		if err != nil {
			return nil, err
		}
		clauses = append(clauses, fieldClauses...)
	}
	return clauses, nil
}

// parseField parses an exported struct field and returns all nested sqlClauses.
func (b *sqlBuilder) parseField(field reflect.StructField, value reflect.Value) ([]sqlClause, error) {
	// recurse into structs
	if field.Type.Kind() == reflect.Struct {
		return b.parseStruct(value.Interface())
	}
	if field.Tag.Get("ddl") == "" {
		return nil, nil
	}

	// dereference any pointers
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	ddlTag := strings.Split(field.Tag.Get("ddl"), ",")[0]
	dbTag := field.Tag.Get("db")
	clauses := make([]sqlClause, 0)
	var clause sqlClause

	// static must be applied no matter what
	if ddlTag == "static" {
		clauses = append(clauses, sqlStaticClause(dbTag))
		return clauses, nil
	}

	if value.Kind() == 0 {
		return nil, nil
	}

	switch ddlTag {
	case "keyword":
		if value.Kind() == reflect.Bool {
			useKeyword := value.Interface().(bool)
			if useKeyword {
				clause = sqlKeywordClause{
					value: dbTag,
					qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
				}
			} else {
				return nil, nil
			}
		} else {
			clause = sqlKeywordClause{
				value: value.Interface(),
				qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
			}
		}
	case "command":
		clause = sqlCommandClause{
			key:   dbTag,
			value: value.Interface(),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	case "identifier":
		clause = sqlIdentifierClause{
			key:   dbTag,
			value: value.Interface().(ObjectIdentifier),
		}
	case "parameter":
		clause = sqlParameterClause{
			key:   dbTag,
			value: value.Interface(),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	default:
		return nil, nil
	}
	return append(clauses, clause), nil
}

// getUnexportedField returns the value of an unexported field.
func (b *sqlBuilder) getUnexportedField(field reflect.Value) interface{} {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().Interface()
}

// parseUnexportedField parses an unexported struct field and returns a sqlClause.
func (b *sqlBuilder) parseUnexportedField(field reflect.StructField, value reflect.Value) ([]sqlClause, error) {
	clauses := make([]sqlClause, 0)
	if field.Tag.Get("ddl") == "" {
		return clauses, nil
	}
	tagParts := strings.Split(field.Tag.Get("ddl"), ",")
	ddlType := tagParts[0]
	dbTag := field.Tag.Get("db")
	var clause sqlClause
	switch ddlType {
	case "list":
		// if it is a list just get the type and go back to parseStruct
		f := b.getUnexportedField(value)
		if f == nil {
			return nil, nil
		}

		listClauses := make([]sqlClause, 0)
		// loop through the slice call parseStruct on each element (since the elements could be structs)
		for i := 0; i < value.Len(); i++ {
			u := b.getUnexportedField(value.Index(i))
			structClauses, err := b.parseStruct(u)
			if err != nil {
				return nil, err
			}
			// each element of the slice needs to be pre-rendered before the commas are added
			renderedStructClauses := b.sql(structClauses...)
			sClause := sqlStaticClause(renderedStructClauses)
			listClauses = append(listClauses, sClause)
		}
		clauses = append(clauses, sqlListClause{
			clauses:        listClauses,
			sep:            ",",
			keyword:        field.Tag.Get("db"),
			useParentheses: getUseParenthesesFromTag(field.Tag, "ddl", true),
		})
		return clauses, nil
	case "identifier":
		id := b.getUnexportedField(value).(ObjectIdentifier)
		if id.FullyQualifiedName() != "" {
			clause = sqlIdentifierClause{
				key:   dbTag,
				value: id,
			}
		}
	case "keyword":
		clause = sqlKeywordClause{
			value: b.getUnexportedField(value),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	case "command":
		clause = sqlCommandClause{
			key:   dbTag,
			value: b.getUnexportedField(value),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	case "static":
		clause = sqlStaticClause(dbTag)
	}
	return append(clauses, clause), nil
}

type sqlListClause struct {
	keyword        string
	clauses        []sqlClause
	sep            string
	useParentheses bool
}

func (v sqlListClause) String() string {
	var s string
	// unclear if we should return parentheses at all.
	if len(v.clauses) == 0 {
		return s
	}
	clauseStrings := make([]string, len(v.clauses))
	for i, clause := range v.clauses {
		clauseStrings[i] = clause.String()
	}
	s = strings.Join(clauseStrings, v.sep)
	if v.useParentheses {
		s = parentheses(s)
	}
	if v.keyword != "" {
		s = fmt.Sprintf("%s %s", v.keyword, s)
	}
	return s
}

type sqlClause interface {
	String() string
}

type sqlStaticClause string

func (v sqlStaticClause) String() string {
	return string(v)
}

type sqlKeywordClause struct {
	value interface{}
	qt    quoteType
}

func (v sqlKeywordClause) String() string {
	return v.qt.Quote(v.value)
}

type sqlIdentifierClause struct {
	key   string
	value ObjectIdentifier
}

func (v sqlIdentifierClause) String() string {
	if v.key != "" {
		return fmt.Sprintf("%s %s", v.key, v.value.FullyQualifiedName())
	}
	return v.value.FullyQualifiedName()
}

type sqlParameterClause struct {
	key   string
	value interface{} // string list, string, string literal, bool, int
	qt    quoteType
}

func (v sqlParameterClause) String() string {
	vType := reflect.TypeOf(v.value)
	var result string
	if v.key != "" {
		result = fmt.Sprintf("%s = ", v.key)
	}
	if vType.Kind() == reflect.String {
		result += v.qt.Quote(v.value.(string))
	} else {
		result += fmt.Sprintf("%v", v.value)
	}

	return result
}

type sqlCommandClause struct {
	key   string
	value interface{}
	qt    quoteType
}

func (v sqlCommandClause) String() string {
	return fmt.Sprintf("%s %s", v.key, v.qt.Quote(v.value))
}
