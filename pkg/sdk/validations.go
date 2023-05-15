package sdk

import (
	"fmt"
	"reflect"
)

func IsValidDataType(v string) bool {
	dt := DataTypeFromString(v)
	return dt != DataTypeUnknown
}

func exactlyOneValueSet(values ...interface{}) error {
	var count int
	for _, v := range values {
		if !reflect.ValueOf(v).IsNil() {
			count++
		}
	}
	if count != 1 {
		return fmt.Errorf("exactly one of the following values must be non-nil: %v", values)
	}
	return nil
}

func everyValueSet(values ...interface{}) bool {
	for _, v := range values {
		if !valueSet(v) {
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
