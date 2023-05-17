package sdk

import (
	"reflect"
)

func IsValidDataType(v string) bool {
	dt := DataTypeFromString(v)
	return dt != DataTypeUnknown
}

func validObjectidentifier(objectIdentifier ObjectIdentifier) bool {
	// https://docs.snowflake.com/en/sql-reference/identifiers-syntax#double-quoted-identifiers
	l := len(objectIdentifier.FullyQualifiedName())
	if l == 0 || l > 255 {
		return false
	}
	return true
}

func anyValueSet(values ...interface{}) bool {
	for _, v := range values {
		if valueSet(v) {
			return true
		}
	}
	return false
}

func exactlyOneValueSet(values ...interface{}) bool {
	var count int
	for _, v := range values {
		if valueSet(v) {
			count++
		}
	}
	return count == 1
}

func everyValueSet(values ...interface{}) bool {
	for _, v := range values {
		if !valueSet(v) {
			return false
		}
	}
	return true
}

func everyValueNil(values ...interface{}) bool {
	for _, v := range values {
		if valueSet(v) {
			return false
		}
	}
	return true
}

func valueSet(value interface{}) bool {
	if value == nil {
		return false
	}
	reflectedValue := reflect.ValueOf(value)
	if reflectedValue.Kind() == reflect.Ptr {
		reflectedValue = reflectedValue.Elem()
	}
	if reflectedValue.Kind() == reflect.Slice {
		return reflectedValue.Len() > 0
	}
	if reflectedValue.Kind() == reflect.Invalid || reflectedValue.IsZero() {
		return false
	}
	if reflectedValue.CanInterface() {
		if _, ok := reflectedValue.Interface().(ObjectIdentifier); ok {
			return validObjectidentifier(reflectedValue.Interface().(ObjectIdentifier))
		}
	}
	if reflectedValue.Kind() != reflect.Struct && reflectedValue.Kind() != reflect.Bool {
		return !reflectedValue.IsNil()
	}

	return true
}

func validateIntInRange(value int, min int, max int) bool {
	if value < min || value > max {
		return false
	}
	return true
}

func validateIntGreaterThanOrEqual(value int, min int) bool {
	return value >= min
}
