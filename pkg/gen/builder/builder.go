package builder

import (
	"fmt"
	"reflect"
	"strings"
)

//// Builder

func QueryStruct(name string) *StructBuilder {
	return &StructBuilder{Name: name}
}

type Buildable interface {
	Build() interface{}
}

type StructBuilder struct {
	Name   string
	Fields []FieldBuilder
}

type Struct struct {
	name   string
	fields []FieldBuilder
}

// String temporary way to see what is being generated.
func (sb *StructBuilder) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("type %s struct {\n", sb.Name))
	for _, f := range sb.Fields {
		tags := make([]string, 0)
		for k, v := range f.Tags {
			tags = append(tags, fmt.Sprintf("%s:\"%s\"", k, strings.Join(v, ",")))
		}
		s.WriteString(fmt.Sprintf("\t%-20s %-20s %-40s\n", f.Name, f.Type, fmt.Sprintf("`%s`", strings.Join(tags, " "))))
	}
	s.WriteString("}\n")
	return s.String()
}

//// Fields

type Type string

func (k Type) Kind() string {
	return string(k)
}

const (
	TypeBoolPtr   Type = "*bool"
	TypeBool      Type = "bool"
	TypeStringPtr Type = "*string"
	TypeString    Type = "string"
	TypeIntPtr    Type = "*int"
	TypeInt       Type = "int"
)

type FieldBuilder struct {
	Name string
	Type string
	Tags map[string][]string
}

type IntoFieldBuilder interface {
	IntoFieldBuilder() []FieldBuilder
}

func (fb FieldBuilder) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{fb}
}

func (sb *StructBuilder) addField(f *FieldBuilder) {
	sb.Fields = append(sb.Fields, *f)
}

func (sb StructBuilder) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{
		{
			Name: sb.Name,
			Type: sb.Name,
			Tags: make(map[string][]string),
		},
	}
}

func (s Struct) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{
		{
			Name: s.name,
			Type: s.name,
			Tags: make(map[string][]string),
		},
	}
}

//// Enums

type EnumBuilder[T any] struct {
	name   string
	values []EnumValue[T]
}

type EnumValue[T any] struct {
	Name  string
	Value T
}

func EnumType[T any](name string) *EnumBuilder[T] {
	return &EnumBuilder[T]{
		name:   name,
		values: make([]EnumValue[T], 0),
	}
}

func (eb *EnumBuilder[T]) With(variableName string, value T) *EnumBuilder[T] {
	eb.values = append(eb.values, EnumValue[T]{
		Name:  variableName,
		Value: value,
	})
	return eb
}

func (eb *EnumBuilder[T]) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("type %s string\n\n", eb.name)) // TODO make generic to take in type or take in type when building EnumBuilder as string
	sb.WriteString("const (\n")
	for _, v := range eb.values {
		switch any(v.Value).(type) {
		case string:
			sb.WriteString(fmt.Sprintf("\t%-20s %s = %s\n", v.Name, eb.name, fmt.Sprintf(`"%s"`, v.Value)))
		default:
			sb.WriteString(fmt.Sprintf("\t%-20s %s = %v\n", v.Name, eb.name, v.Value))
		}
	}
	sb.WriteString(")\n")
	return sb.String()
}

func (eb *EnumBuilder[T]) IntoFieldBuilder() []FieldBuilder {
	return []FieldBuilder{
		{
			Name: eb.name,
			Type: eb.name,
			Tags: make(map[string][]string),
		},
	}
}

//// Options

type option func(sb *StructBuilder, f *FieldBuilder)

//type keywordOptions
//type staticOptions
//type parameterOptions

func WithTag(tagName string, tagValue string) func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		if _, ok := f.Tags[tagName]; !ok {
			f.Tags[tagName] = make([]string, 0)
		}
		f.Tags[tagName] = append(f.Tags[tagName], tagValue)
	}
}

func WithSQL(sql string) option {
	return WithTag("sql", sql)
}

func WithSQLPrefix(sql string) option {
	return WithTag("sql", sql)
}

func WithDDL(ddl string) option {
	return WithTag("ddl", ddl)
}

func WithParen() option {
	return WithDDL("parentheses")
}

func WithNoQuotes() option {
	return WithDDL("no_quotes")
}

func NoEquals() option {
	return WithDDL("no_equals")
}

type Kinder interface {
	Kind() string
}

func (sb StructBuilder) Kind() string {
	return sb.Name
}

// func Optional(kinder Kinder) func(sb *StructBuilder, f *FieldBuilder) {
// 	return withType("*" + kinder.Kind())
// }

func Number() func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Type = TypeInt.Kind()
	}
}

func Text() func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Type = TypeInt.Kind()
	}
}

func SQLPrefix(s string) func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Type = TypeInt.Kind()
	}
}

func SingleQuotedText() func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Type = TypeInt.Kind()
	}
}

func WithType(kinder Kinder) func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Type = kinder.Kind()
	}
}

func withType(kind string) func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		f.Type = kind
	}
}

func WithTypeT[T any]() func(sb *StructBuilder, f *FieldBuilder) {
	return func(sb *StructBuilder, f *FieldBuilder) {
		tType := reflect.TypeOf(new(T))
		fmt.Printf("LUL: %v", new(T))
		f.Type = tType.Name()
	}
}

func WithPointerType(kinder Kinder) func(sb *StructBuilder, f *FieldBuilder) {
	return withType("*" + kinder.Kind())
}

func ListOf(kinder Kinder) func(sb *StructBuilder, f *FieldBuilder) {
	return withType("[]" + kinder.Kind())
}

func WithTypeBoolPtr() option {
	return WithType(TypeBoolPtr)
}

//// Building functions

func (sb *StructBuilder) Field(ddlValue string, fieldName string, options ...option) *StructBuilder {
	f := &FieldBuilder{Name: fieldName}
	f.Tags = make(map[string][]string)
	f.Tags["ddl"] = make([]string, 0)
	f.Tags["ddl"] = append(f.Tags["ddl"], ddlValue)
	for _, opt := range options {
		opt(sb, f)
	}
	sb.addField(f)
	return sb
}

// Static
func (sb *StructBuilder) Static(fieldName string, sql string, options ...option) *StructBuilder {
	return sb.Field("static", fieldName, options...)
}

func (sb *StructBuilder) SQL(sql string, options ...option) *StructBuilder {
	return sb.Field("static", fieldName, options...)
}

func (sb *StructBuilder) Create() *StructBuilder {
	return sb.Static("create", WithType(TypeBool), WithSQL("CREATE"))
}

// Keyword
func (sb *StructBuilder) Tag() *StructBuilder {
	return sb
}

func (sb *StructBuilder) Keyword(fieldName string, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("keyword", fieldName, options...)
}

func (sb *StructBuilder) OptionalAssignment(fieldName string, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("keyword", fieldName, options...)
}

func (sb *StructBuilder) Assignment(fieldName string, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("keyword", fieldName, options...)
}

func (sb *StructBuilder) OptionalValue(fieldName string, value any, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("keyword", fieldName, options...)
}

func (sb *StructBuilder) Value(fieldName string, value any, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("keyword", fieldName, options...)
}

func (sb *StructBuilder) OptionalSQL(sql string, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("keyword", sql, options...)
}

func Keyword(fieldName string, options ...option) FieldBuilder {
	fb := FieldBuilder{Name: fieldName, Type: TypeBoolPtr.Kind(), Tags: map[string][]string{"ddl": {"keyword"}}}
	for _, opt := range options {
		opt(nil, &fb)
	}
	return fb
}

func Value(fieldName string, options ...option) FieldBuilder {
	fb := FieldBuilder{Name: fieldName, Type: TypeBoolPtr.Kind(), Tags: map[string][]string{"ddl": {"keyword"}}}
	for _, opt := range options {
		opt(nil, &fb)
	}
	return fb
}

func OptionalSQL(sql string, options ...option) FieldBuilder {
	fb := FieldBuilder{Name: sql, Type: TypeBoolPtr.Kind(), Tags: map[string][]string{"ddl": {"keyword"}}}
	for _, opt := range options {
		opt(nil, &fb)
	}
	return fb
}

func OptionalText(sql string, options ...option) FieldBuilder {
	fb := FieldBuilder{Name: sql, Type: TypeBoolPtr.Kind(), Tags: map[string][]string{"ddl": {"keyword"}}}
	for _, opt := range options {
		opt(nil, &fb)
	}
	return fb
}

func OptionalValue(sql string, typeObj any, options ...option) FieldBuilder {
	fb := FieldBuilder{Name: sql, Type: TypeBoolPtr.Kind(), Tags: map[string][]string{"ddl": {"keyword"}}}
	for _, opt := range options {
		opt(nil, &fb)
	}
	return fb
}

func (sb *StructBuilder) OrReplace() *StructBuilder {
	return sb.Keyword("OrReplace", WithTypeBoolPtr())
}

func (sb *StructBuilder) IfExists() *StructBuilder {
	return sb.Keyword("IfExists", WithTypeBoolPtr())
}

func (sb *StructBuilder) IfNotExists() *StructBuilder {
	return sb.Keyword("IfNotExists", WithTypeBoolPtr())
}

func (sb *StructBuilder) Transient() *StructBuilder {
	return sb.Keyword("Transient", WithTypeBoolPtr())
}

func (sb *StructBuilder) Number(fieldName string, options ...option) *StructBuilder {
	return sb.Field("identifier", fieldName, options...)
}

func (sb *StructBuilder) Text(fieldName string, options ...option) *StructBuilder {
	return sb.Field("identifier", fieldName, options...)
}

func (sb *StructBuilder) OptionalText(fieldName string, options ...option) *StructBuilder {
	return sb.Field("identifier", fieldName, options...)
}

// Identifier
func (sb *StructBuilder) Identifier(fieldName string, typeObj any, options ...option) *StructBuilder {
	return sb.Field("identifier", fieldName, options...)
}

type oneof struct {
	fields []IntoFieldBuilder
}

func OneOf(fields ...IntoFieldBuilder) oneof {
	return oneof{
		fields: fields,
	}
}

func (o oneof) IntoFieldBuilder() []FieldBuilder {
	fbs := make([]FieldBuilder, 0)
	for _, f := range o.fields {
		for _, fb := range f.IntoFieldBuilder() {
			fbs = append(fbs, fb)
		}
	}
	return fbs
}

func (sb *StructBuilder) OneOf(fields ...IntoFieldBuilder) *StructBuilder {
	for _, f := range fields {
		for _, fb := range f.IntoFieldBuilder() {
			sb.addField(&fb)
		}
	}
	return sb
}

// Parameter
func (sb *StructBuilder) Parameter(fieldName string, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("parameter", fieldName, options...)
}

// Parameter
func (sb *StructBuilder) List(fieldName string, typeObj any, options ...option) *StructBuilder {
	// TODO add sql tag ?
	return sb.Field("list", fieldName, options...)
}

type Group struct {
	IntoFieldBuilders []IntoFieldBuilder
}

func (g Group) IntoFieldBuilder() []FieldBuilder {
	fbs := make([]FieldBuilder, 0)
	for _, f := range g.IntoFieldBuilders {
		fbs = append(fbs, f.IntoFieldBuilder()...)
	}
	return fbs
}

func GroupFields(fields ...IntoFieldBuilder) Group {
	return Group{
		IntoFieldBuilders: fields,
	}
}

func (sb *StructBuilder) Build() Struct {
	return Struct{}
}

//// API

type API struct{}

func BuildAPI() API {
	return API{}
}
