package gen

import (
	"reflect"
)

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name           string
	ConcreteType   string
	UnderlyingType string
}

// TODO: test completely new struct with:
//   - basic type fields (string, int, float, bool)
//   - pointer to basic types fields
//   - time.Time, *time.Time
//   - enum (string and int)
//   - slice (string, enum)
//   - identifier (each type)
//   - slice (identifier)
//   - (?) slice of pointers
//   - (?) pointer to slice
//   - (?) struct
func ExtractStructDetails(s any) Struct {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	fields := make([]Field, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		currentField := v.Field(i)
		currentName := v.Type().Field(i).Name
		currentType := v.Type().Field(i).Type.String()
		//currentValue := currentField.Interface()

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
	return Struct{v.Type().String(), fields}
}
