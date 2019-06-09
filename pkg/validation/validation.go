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

// ValidateAccount checks that your account name is valid
//
// the identifier must start with an alphabetic character and cannot contain spaces or special characters unless
// the entire identifier string is enclosed in double quotes (e.g. "My object").
func ValidateAccount(i interface{}, k string) ([]string, []error) {
	name, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}

	if firstChar := name[0]; firstChar < asciiA || firstChar > asciiz || (firstChar > asciiZ && firstChar < asciia) {
		return nil, []error{fmt.Errorf("%v is not a valid starting character for the identifier (must be alphabetic)", firstChar)}
	}

	return nil, nil
}

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
