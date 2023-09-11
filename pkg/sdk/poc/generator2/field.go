package generator2

import (
	"fmt"
	"golang.org/x/exp/slices"
	"strings"
)

// Field defines properties of a single field or struct (by defining Fields)
type Field struct {
	// Name is how field is called in parent struct
	Name string
	// Kind is an interface to get type with helper functions TypeOf(string), TypeOfT[type]()
	Kind Kind
	// Tags should contain ddl and sql tags used for SQL generation
	Tags map[string][]string
	// Required is used to mark fields which are essential (it's used e.g. for DTO builders generation)
	Required bool
}

func NewField(name string, kind Kind, tags map[string][]string) *Field {
	return &Field{
		Name: name,
		Kind: kind,
		Tags: tags,
	}
}

func (f *Field) WithRequired(required bool) *Field {
	f.Required = required
	return f
}

// ShouldBeInDto checks if field is not some static SQL field which should not be interacted with by SDK user
// TODO: this is a very naive implementation, consider fixing it with DSL builder connection
func (f *Field) ShouldBeInDto() bool {
	return !slices.Contains(f.Tags["ddl"], "static")
}

// TagsPrintable defines how tags are printed in options structs, it ensures the same order of tags for every field
func (f *Field) TagsPrintable() string {
	var tagNames = []string{"ddl", "sql"}
	var tagParts []string
	for _, tagName := range tagNames {
		var v, ok = f.Tags[tagName]
		if ok {
			tagParts = append(tagParts, fmt.Sprintf(`%s:"%s"`, tagName, strings.Join(v, ",")))
		}
	}
	return fmt.Sprintf("`%s`", strings.Join(tagParts, " "))
}

func queryTags(ddlTags []string, sqlTags []string) map[string][]string {
	tags := make(map[string][]string)
	if len(ddlTags) > 0 {
		tags["ddl"] = ddlTags
	}
	if len(sqlTags) > 0 {
		tags["sql"] = sqlTags
	}
	return tags
}

func StructField(s *Struct, ddlTags []string, sqlTags []string) *Field {
	return nil
}

// Static / SQL

func SQL(sql string) *Field {
	return NewField(sqlToFieldName(sql, false), KindOf("bool"), queryTags([]string{"static"}, []string{sql}))
}

func Create() *Field {
	return SQL("CREATE")
}

func Alter() *Field {
	return SQL("ALTER")
}

func Drop() *Field {
	return SQL("DROP")
}

func Show() *Field {
	return SQL("SHOW")
}

func Describe() *Field {
	return SQL("DESCRIBE")
}

// Keyword / Value

func OptionalSQL(sql string) *Field {
	return NewField(sqlToFieldName(sql, true), KindOf("bool"), queryTags([]string{"keyword"}, []string{sql}))
}

func OrReplace() *Field {
	return OptionalSQL("OR REPLACE")
}

func IfNotExists() *Field {
	return OptionalSQL("IF NOT EXISTS")
}

func IfExists() *Field {
	return OptionalSQL("IF EXISTS")
}

// Parameters

type parameterOptions struct {
	singleQuotes bool
}

func ParameterOptions() *parameterOptions {
	return &parameterOptions{}
}

func (v *parameterOptions) SingleQuotes(value bool) *parameterOptions {
	v.singleQuotes = value
	return v
}

func (v *parameterOptions) toOptions() []string {
	opts := make([]string, 0)
	if v.singleQuotes {
		opts = append(opts, "single_quotes")
	}
	return opts
}

func OptionalTextAssignment(sqlPrefix string, paramOptions *parameterOptions) *Field {
	if paramOptions != nil {
		return NewField(sqlToFieldName(sqlPrefix, true), KindOf("*string"), queryTags(append([]string{"parameter"}, paramOptions.toOptions()...), []string{sqlPrefix}))
	}
	return NewField(sqlToFieldName(sqlPrefix, true), KindOf("*string"), queryTags([]string{"parameter"}, []string{sqlPrefix}))
}

// Identifier

func DatabaseObjectIdentifier(fieldName string) *Field {
	return NewField(fieldName, KindOf("DatabaseObjectIdentifier"), queryTags([]string{"identifier"}, nil)).WithRequired(true)
}
