package sdk

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"
)

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
		if c != nil {
			sList[i] = c.String()
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
					clauses = append(clauses, sqlClauseKeyword{
						value: field.Tag.Get("db"),
						qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
					})
				case "identifier":
					if value.Interface().(ObjectIdentifier).FullyQualifiedName() == "" {
						continue
					}
					clauses = append(clauses, sqlClauseIdentifier{
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
		clauses = append(clauses, sqlClauseStatic(dbTag))
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
				clause = sqlClauseKeyword{
					value: dbTag,
					qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
				}
			} else {
				return nil, nil
			}
		} else {
			clause = sqlClauseKeyword{
				value: value.Interface().(string),
				qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
			}
		}
	case "command":
		clause = sqlClauseCommand{
			key:   dbTag,
			value: value.Interface(),
			qt:    getQuoteTypeFromTag(field.Tag, "ddl"),
		}

	case "parameter":
		clause = sqlClauseParameter{
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
	case "identifier":
		id := b.getUnexportedField(value).(ObjectIdentifier)
		if id.FullyQualifiedName() != "" {
			clause = sqlClauseIdentifier{
				key:   dbTag,
				value: id,
			}
		}
	case "static":
		clause = sqlClauseStatic(dbTag)
	}
	return append(clauses, clause), nil
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
	return v.qt.Quote(v.value)
}

type sqlClauseIdentifier struct {
	key   string
	value ObjectIdentifier
}

func (v sqlClauseIdentifier) String() string {
	if v.key != "" {
		return fmt.Sprintf("%s %s", v.key, v.value.FullyQualifiedName())
	}
	return v.value.FullyQualifiedName()
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
		result += v.qt.Quote(v.value.(string))
	} else {
		result += fmt.Sprintf("%v", v.value)
	}

	return result
}

type sqlClauseCommand struct {
	key   string
	value interface{}
	qt    quoteType
}

func (v sqlClauseCommand) String() string {
	return fmt.Sprintf("%s %s", v.key, v.qt.Quote(v.value))
}
