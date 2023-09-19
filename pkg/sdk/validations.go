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

func moreThanOneValueSet(values ...interface{}) bool {
	var count int
	for _, v := range values {
		if valueSet(v) {
			count++
		}
	}
	return count > 1
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

// TODO This have to be changed or new function should be created with more options or better defaults
//
//		because there are cases where validation is incorrect because of this function,
//		e.g. you want to alter some resource and you want to provide
//		empty array to unset all of the values, sometimes you cannot do it because in the validation you used this function
//		to see if anything was set... well it was set, but it was empty array which is considered not set, but in this case
//	 	it was valid option (and in some cases it may not be!).
func valueSet(value interface{}) bool {
	if value == nil {
		return false
	}
	reflectedValue := reflect.ValueOf(value)
	if reflectedValue.Kind() == reflect.Ptr {
		reflectedValue = reflectedValue.Elem()
	}
	switch reflectedValue.Kind() {
	case reflect.Slice, reflect.String:
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
