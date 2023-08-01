package builder

import (
	"fmt"
	"reflect"
)

type Typer interface {
	Type() string
}

type Type string

func (k Type) Type() string {
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

type Buildable interface {
	Build() interface{}
}

type IntoFieldBuilder interface {
	IntoFieldBuilder() []FieldBuilder
}

type Struct struct {
	name   string
	fields []Field
}

type StructBuilder struct {
	Name   string
	Fields []FieldBuilder
}

type Field struct {
	Name  string
	Typer Typer
	Tags  map[string][]string
}

type FieldBuilder struct {
	Name  string
	Typer Typer
	Tags  map[string][]string
}

type EnumBuilder[T any] struct {
	name   string
	values []EnumValue[T]
}

type EnumValue[T any] struct {
	Name  string
	Value T
}

type stringTyper struct {
	TypeName string
}

func newStringTyper(typeName string) Typer {
	return &stringTyper{
		TypeName: typeName,
	}
}

func (s *stringTyper) Type() string {
	return s.TypeName
}

func TypeOfString(typeName string) Typer {
	return newStringTyper(typeName)
}

func TypeOfQueryStruct(sb *StructBuilder) Typer {
	return sb
}

func TypeT[T any]() Typer {
	return newStringTyper(getTypeName[T]())
}

func getTypeName[T any]() string {
	t := reflect.TypeOf((*T)(nil)).Elem()
	return fmt.Sprint(t)
}
