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
		if v != nil && !reflect.ValueOf(v).IsNil() {
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
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		return !v.IsNil()
	}
	if v.CanInterface() {
		// if the value is an identifier, check if it is valid
		if _, ok := v.Interface().(ObjectIdentifier); ok {
			return validObjectidentifier(v.Interface().(ObjectIdentifier))
		}
	}
	return false
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
