package validation

import (
	"fmt"
)

const (
	ascii0 = 48
	ascii9 = 57
	asciiA = 65
	asciiZ = 90
	asciia = 97
	asciiz = 122
)

// ValidatePassword checks that you password meets the Snowflake Password Policy
//
// Must be at least 8 characters long.
// Must contain at least 1 digit.
// Must contain at least 1 uppercase letter and 1 lowercase letter.
func ValidatePassword(i interface{}, k string) (s []string, errs []error) {
	pass, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}

	if len(pass) < 8 {
		errs = append(errs, fmt.Errorf("Password must be at least 8 characters long"))
	}

	var digit, uppercase, lowercase bool
	for _, c := range pass {
		if c >= asciiA && c <= asciiZ {
			uppercase = true
		}
		if c >= asciia && c <= asciiz {
			lowercase = true
		}
		if c >= ascii0 && c <= ascii9 {
			digit = true
		}
	}

	if !uppercase {
		errs = append(errs, fmt.Errorf("Password must contain an uppercase character"))
	}

	if !lowercase {
		errs = append(errs, fmt.Errorf("Password must contain a lowercase character"))
	}

	if !digit {
		errs = append(errs, fmt.Errorf("Password must contain a digit"))
	}

	return
}
