package sdk

import (
	"reflect"
	"strconv"
)

// String returns a pointer to the given string.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the given bool.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to the given int.
func Int(i int) *int {
	return &i
}

// toInt converts a string to an int.
func toInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func copyFields(src, dst interface{}) {
	if src == nil || dst == nil {
		return
	}
	srcVal := reflect.ValueOf(src).Elem()
	if !srcVal.IsValid() {
		return
	}
	dstVal := reflect.ValueOf(dst).Elem()
	if dstVal.IsZero() {
		dst = reflect.New(reflect.TypeOf(dst).Elem())
	}
	for i := 0; i < srcVal.NumField(); i++ {
		fieldName := srcVal.Type().Field(i).Name
		// find the field in dst and set it to the value in src
		dstField := dstVal.FieldByName(fieldName)
		if dstField.CanSet() {
			dstField.Set(srcVal.Field(i))
		}
	}
}
