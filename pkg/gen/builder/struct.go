package builder

import (
	"fmt"
	"strings"
)

func QueryStruct(name string) *StructBuilder {
	return &StructBuilder{Name: name}
}

type FieldTransformer interface {
	transform(fb *FieldBuilder) *FieldBuilder
}

func (sb *StructBuilder) Field(fieldName string, typer Typer, ddlValues []string, sqlValues []string, transformer FieldTransformer) *StructBuilder {
	f := &FieldBuilder{Name: fieldName, Typer: typer}
	f.Tags = make(map[string][]string)
	f.Tags["ddl"] = make([]string, 0)
	f.Tags["ddl"] = append(f.Tags["ddl"], ddlValues...)
	f.Tags["sql"] = make([]string, 0)
	f.Tags["sql"] = append(f.Tags["sql"], sqlValues...)
	sb.Fields = append(sb.Fields, *f)
	if transformer != nil {
		f = transformer.transform(f)
	}
	return sb
}

// Static / SQL

func (sb *StructBuilder) SQL(sql string) *StructBuilder {
	return sb.Field(sqlToFieldName(sql, false), TypeBool, []string{"static"}, []string{sql}, nil)
}

func (sb *StructBuilder) Create() *StructBuilder {
	return sb.SQL("CREATE")
}

func (sb *StructBuilder) OrReplace() *StructBuilder {
	return sb.SQL("OR REPLACE")
}

func (sb *StructBuilder) IfNotExists() *StructBuilder {
	return sb.SQL("IF NOT EXISTS")
}

// Keyword / Value

// TODO: we can use varchar to skip if nothing was specified in the keywordOptions
func (sb *StructBuilder) Keyword(fieldName string, typer Typer, keywordOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(fieldName, true), typer, []string{"keyword"}, []string{}, keywordOptions)
}

func (sb *StructBuilder) Number(fieldName string, keywordOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(fieldName, true), TypeInt, []string{"keyword"}, []string{}, keywordOptions)
}

func (sb *StructBuilder) Text(fieldName string, keywordOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(fieldName, true), TypeString, []string{"keyword"}, []string{}, keywordOptions)
}

func (sb *StructBuilder) OptionalText(fieldName string, keywordOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(fieldName, true), TypeStringPtr, []string{"keyword"}, []string{}, keywordOptions)
}

func (sb *StructBuilder) Value(fieldName string, typer Typer, keywordOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(fieldName, true), typer, []string{"keyword"}, []string{}, keywordOptions)
}

func (sb *StructBuilder) OptionalValue(fieldName string, typer Typer, keywordOptions FieldTransformer) *StructBuilder {
	// TODO: Check if given typer.Type() returns a pointer ?
	return sb.Field(sqlToFieldName(fieldName, true), typer, []string{"keyword"}, []string{}, keywordOptions)
}

func (sb *StructBuilder) OptionalSQL(sql string, valueOptions FieldTransformer) *StructBuilder {
	return sb.Keyword(sqlToFieldName(sql, true), TypeBoolPtr, valueOptions)
}

// Parameter / Assignment

func (sb *StructBuilder) Parameter(sql string, typer Typer, parameterOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(sql, true), typer, []string{"parameter"}, []string{}, parameterOptions)
}

func (sb *StructBuilder) AssignNumber(sql string, parameterOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(sql, true), TypeInt, []string{"parameter"}, []string{}, parameterOptions)
}

func (sb *StructBuilder) AssignText(sql string, parameterOptions FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(sql, true), TypeString, []string{"parameter"}, []string{}, parameterOptions)
}

func (sb *StructBuilder) AssignValue(sql string, typer Typer, parameterOptions ...FieldTransformer) *StructBuilder {
	return sb.Field(sqlToFieldName(sql, true), typer, []string{"parameter"}, []string{}, parameterOptions[0])
}

// Identifier

func (sb *StructBuilder) Identifier(fieldName string, typer Typer) *StructBuilder {
	return sb.Field(fieldName, typer, []string{"identifier"}, []string{}, nil)
}

// List

func (sb *StructBuilder) List(fieldName string, typer Typer, listOptions ...FieldTransformer) *StructBuilder {
	return sb.Field(fieldName, typer, []string{"list"}, []string{}, listOptions[0])
}

// TODO: Refactor OneOf
func (sb *StructBuilder) OneOf(fields ...IntoFieldBuilder) *StructBuilder {
	for _, f := range fields {
		for _, fb := range f.IntoFieldBuilder() {
			sb.Fields = append(sb.Fields, fb)
		}
	}
	return sb
}

func (sb *StructBuilder) Type() string {
	return sb.Name
}

func (sb *StructBuilder) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{
		{
			Name:  sb.Name,
			Typer: TypeOfString(sb.Name),
			Tags:  make(map[string][]string),
		},
	}
}

// String temporary way to see what is being generated.
func (sb *StructBuilder) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("type %s struct {\n", sb.Name))
	for _, f := range sb.Fields {
		tags := make([]string, 0)
		for k, v := range f.Tags {
			if len(v) != 0 {
				tags = append(tags, fmt.Sprintf("%s:\"%s\"", k, strings.Join(v, ",")))
			}
		}
		s.WriteString(fmt.Sprintf("\t%-20s %-20s %-40s\n", f.Name, f.Typer.Type(), fmt.Sprintf("`%s`", strings.Join(tags, " "))))
	}
	s.WriteString("}\n")
	return s.String()
}
