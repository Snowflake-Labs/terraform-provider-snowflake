package sdk

import (
	"fmt"
	"reflect"
	"strings"
)

type quoteType string

const (
	NoQuotes     quoteType = "no_quotes"
	DoubleQuotes quoteType = "double_quotes"
	SingleQuotes quoteType = "single_quotes"
)

func (v quoteType) String() string {
	switch v {
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
			return quoteType(part)
		}
	}
	return NoQuotes
}

type sqlBuilder struct{}

// sql builds a SQL statement from sqlClauses.
func (b *sqlBuilder) sql(clauses ...sqlClause) string {
	sList := make([]string, len(clauses))
	for i, c := range clauses {
		sList[i] = c.String()
	}
	return strings.Join(sList, " ")
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
			clause := b.parseUnexportedField(field, value)
			if clause != nil {
				clauses = append(clauses, clause)
			}
			continue
		}

		// skip nil pointers for attributes, since they are not set. Otherwise dereference them
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue
			}
			value = value.Elem()
		}

		if value.Kind() == reflect.Struct {
			// check if there is any keyword on the struct
			// if there is, then we need to add it to the clause
			// if there is not, then we need to recurse into the struct
			// and get the clauses from there
			ddlTag := field.Tag.Get("ddl")
			if ddlTag != "" {
				ddlTagParts := strings.Split(ddlTag, ",")
				if ddlTagParts[0] == "keyword" {
					clauses = append(clauses, sqlClauseKeyword{
						value: field.Tag.Get("db"),
						qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
					})
				}
			}
			innerClauses, err := b.parseStruct(value.Interface())
			if err != nil {
				return nil, err
			}
			clauses = append(clauses, innerClauses...)
			continue
		}

		// default case, if not a struct then it is a field
		clause := b.parseField(field, value)
		if clause != nil {
			clauses = append(clauses, clause)
		}
	}
	return clauses, nil
}

// parseField parses an exported struct field and returns a sqlClause.
func (b *sqlBuilder) parseField(field reflect.StructField, value reflect.Value) sqlClause {
	if field.Tag.Get("ddl") == "" {
		return nil
	}

	ddlTag := strings.Split(field.Tag.Get("ddl"), ",")[0]
	dbTag := field.Tag.Get("db")

	switch ddlTag {
	case "static":
		return sqlClauseStatic(dbTag)
	case "keyword":
		if value.Kind() == reflect.Bool {
			useKeyword := value.Interface().(bool)
			if useKeyword {
				return sqlClauseKeyword{
					value: dbTag,
					qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
				}
			}
			return nil
		}
		return sqlClauseKeyword{
			value: value.Interface().(string),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	case "command":
		return sqlClauseCommand{
			key:   dbTag,
			value: value.Interface().(string),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	case "parameter":
		return sqlClauseParameter{
			key:   dbTag,
			value: value.Interface(),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	default:
		return nil
	}
}

// parseUnexportedField parses an unexported struct field and returns a sqlClause.
func (b *sqlBuilder) parseUnexportedField(field reflect.StructField, value reflect.Value) sqlClause {
	if field.Tag.Get("ddl") == "" {
		return nil
	}
	tagParts := strings.Split(field.Tag.Get("ddl"), ",")
	ddlType := tagParts[0]
	dbTag := field.Tag.Get("db")
	switch ddlType {
	case "static":
		return sqlClauseStatic(dbTag)
	case "name":
		return sqlClauseName(value.String())
	case "keyword":
		if value.Kind() == reflect.Bool {
			useKeyword := value.Bool()
			if !useKeyword {
				return nil
			}
		}
		return sqlClauseKeyword{
			value: value.String(),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}
	}

	return nil
}

type sqlClause interface {
	String() string
}

type sqlClauseStatic string

func (v sqlClauseStatic) String() string {
	return string(v)
}

type sqlClauseKeyword struct {
	value string
	qt    quoteType
}

func (v sqlClauseKeyword) String() string {
	// make sure we dont double quote the keyword
	trimmedValue := strings.Trim(v.value, v.qt.String())
	return fmt.Sprintf("%s%s%s", v.qt.String(), trimmedValue, v.qt.String())
}

type sqlClauseName string

func (v sqlClauseName) String() string {
	return string(v)
}

type sqlClauseParameter struct {
	key   string
	value interface{} // string list, string, string literal, bool, int
	qt    quoteType
}

func (v sqlClauseParameter) String() string {
	vType := reflect.TypeOf(v.value)
	var result string
	if v.key != "" {
		result = fmt.Sprintf("%s = ", v.key)
	}
	if vType.Kind() == reflect.String {
		trimmedValue := strings.Trim(v.value.(string), v.qt.String())
		result += fmt.Sprintf("%s%s%s", v.qt.String(), trimmedValue, v.qt.String())
	} else {
		result += fmt.Sprintf("%v", v.value)
	}

	return result
}

type sqlClauseCommand struct {
	key   string
	value string
	qt    quoteType
}

func (v sqlClauseCommand) String() string {
	trimmedValue := strings.Trim(v.value, v.qt.String())
	return fmt.Sprintf("%s %s%s%s", v.key, v.qt.String(), trimmedValue, v.qt.String())
}
