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

func booleanStringToBool(value string) (bool, error) {
	switch value {
	case BooleanTrue:
		return true, nil
	case BooleanFalse:
		return false, nil
	default:
		return false, fmt.Errorf("cannot retrieve boolean value from %s", value)
	}
}

func booleanStringFieldDescription(description string) string {
	return fmt.Sprintf(`%s Available options are: "%s" or "%s". When the value is not set in the configuration the provider will put "%s" there which means to use the Snowflake default for this value.`, description, BooleanTrue, BooleanFalse, BooleanDefault)
}

func externalChangesNotDetectedFieldDescription(description string) string {
	return fmt.Sprintf(`%s External changes for this field won't be detected. In case you want to apply external changes, you can re-create the resource manually using "terraform taint".`, description)
}
