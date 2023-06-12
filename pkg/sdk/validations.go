package sdk

import (
	"reflect"
)

func IsValidDataType(v string) bool {
	_, err := ToDataType(v)
	return err == nil
}

func IsValidWarehouseSize(v string) bool {
	_, err := ToWarehouseSize(v)
	return err == nil
}

func validObjectidentifier(objectIdentifier ObjectIdentifier) bool {
	// https://docs.snowflake.com/en/sql-reference/identifiers-syntax#double-quoted-identifiers
	l := len(objectIdentifier.Name())
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
	switch reflectedValue.Kind() {
	case reflect.Slice:
		return reflectedValue.Len() > 0
	case reflect.Invalid:
		return false
	case reflect.Struct:
		if _, ok := reflectedValue.Interface().(ObjectIdentifier); ok {
			return validObjectidentifier(reflectedValue.Interface().(ObjectIdentifier))
		}
		return reflectedValue.Interface() != nil
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
