package gencommons

import (
	"reflect"
	"strings"
)

type StructDetails struct {
	Name   string
	Fields []Field
}

func (s StructDetails) ObjectName() string {
	return s.Name
}

type Field struct {
	Name           string
	ConcreteType   string
	UnderlyingType string
}

func (f *Field) IsPointer() bool {
	return strings.HasPrefix(f.ConcreteType, "*")
}

func (f *Field) IsSlice() bool {
	return strings.HasPrefix(f.ConcreteType, "[]")
}

func (f *Field) GetImportedType() (string, bool) {
	parts := strings.Split(f.ConcreteType, ".")
	return parts[0], len(parts) > 1
}

func ExtractStructDetails(s any) StructDetails {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	fields := make([]Field, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		currentField := v.Field(i)
		currentName := v.Type().Field(i).Name
		currentType := v.Type().Field(i).Type.String()

		var kind reflect.Kind
		var isPtr bool

		if currentField.Kind() == reflect.Pointer {
			isPtr = true
			kind = currentField.Type().Elem().Kind()
		} else {
			kind = currentField.Kind()
		}

		var underlyingType string
		if isPtr {
			underlyingType = "*"
		}
		underlyingType += kind.String()

		fields[i] = Field{currentName, currentType, underlyingType}
	}
	return StructDetails{v.Type().String(), fields}
}
