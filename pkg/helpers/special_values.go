package helpers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	BooleanTrue    = "true"
	BooleanFalse   = "false"
	BooleanDefault = "default"

	IntDefault       = -1
	IntDefaultString = "-1"
)

// StringInSlice has the same implementation as validation.StringInSlice, but adapted to schema.SchemaValidateDiagFunc
func StringInSlice(valid []string, ignoreCase bool) schema.SchemaValidateDiagFunc {
	return func(i interface{}, path cty.Path) diag.Diagnostics {
		v, ok := i.(string)
		if !ok {
			return diag.Errorf("expected type of %v to be string", path)
		}

		for _, str := range valid {
			if v == str || (ignoreCase && strings.EqualFold(v, str)) {
				return nil
			}
		}

		return diag.Errorf("expected %v to be one of %q, got %s", path, valid, v)
	}
}

var ValidateBooleanString = StringInSlice([]string{BooleanTrue, BooleanFalse}, false)

var ValidateBooleanStringWithDefault = StringInSlice([]string{BooleanTrue, BooleanFalse, BooleanDefault}, false)

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
