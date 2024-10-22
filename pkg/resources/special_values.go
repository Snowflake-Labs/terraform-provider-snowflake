package resources

import (
	"fmt"
)

const (
	BooleanTrue    = "true"
	BooleanFalse   = "false"
	BooleanDefault = "default"

	IntDefault       = -1
	IntDefaultString = "-1"
)

var validateBooleanString = StringInSlice([]string{BooleanTrue, BooleanFalse}, false)

func booleanStringFromBool(value bool) string {
	if value {
		return BooleanTrue
	} else {
		return BooleanFalse
	}
}

func BooleanStringToBool(value string) (bool, error) {
	switch value {
	case BooleanTrue:
		return true, nil
	case BooleanFalse:
		return false, nil
	default:
		return false, fmt.Errorf("cannot retrieve boolean value from %s", value)
	}
}
