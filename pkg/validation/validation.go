package validation

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	ascii0 = 48
	ascii9 = 57
	asciiA = 65
	asciiZ = 90
	asciia = 97
	asciiz = 122
)

// ValidatePassword checks that your password meets the Snowflake Password Policy
//
// Must be at least 8 characters long.
// Must contain at least 1 digit.
// Must contain at least 1 uppercase letter and 1 lowercase letter.
//lintignore:V011
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

// ValidatePrivilege validates the privilege is in the authorized set.
// Will also check for the ALL privilege and hopefully provide a helpful error message.
func ValidatePrivilege(valid []string, ignoreCase bool) schema.SchemaValidateFunc {
	//lintignore:V013
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
			return warnings, errors
		}

		if v == "ALL" || (ignoreCase && strings.ToUpper(v) == "ALL") {
			errors = append(errors, fmt.Errorf("the ALL privilege is deprecated, see https://github.com/Snowflake-Labs/terraform-provider-snowflake/discussions/318"))
			return warnings, errors
		}

		return validation.StringInSlice(valid, ignoreCase)(i, k)
	}
}

// ValidateIsNotAccountLocator validates that the account value is not an account locator. Account locators have the
// following format: 8 characters where the first 3 characters are letters and the last 5 are digits. ex: ABC12345
// The desired format should be 'organization_name.account_name' ex: testOrgName.testAccName
func ValidateIsNotAccountLocator(i interface{}, k string) (s []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %s to be string", k))
		return
	}
	if !strings.Contains(v, ".") {
		errors = append(errors, fmt.Errorf("account locators are not allowed - please use 'organization_name.account_name"))
		return
	}
	if len(v) == 8 {
		isAccountLocator := true
		firstHalf := v[0:3]
		for _, r := range firstHalf {
			if !unicode.IsLetter(r) {
				isAccountLocator = false
			}
		}
		secondHalf := v[3:]
		for _, r := range secondHalf {
			if !unicode.IsDigit(r) {
				isAccountLocator = false
			}
		}
		if isAccountLocator {
			errors = append(errors, fmt.Errorf("account locators are not allowed - please use 'organization_name.account_name"))
		}
	}
	return
}

func ValidateFullyQualifiedTagPath(i interface{}, k string) (s []string, errors []error) {
	v, _ := i.(string)
	if !strings.Contains(v, "|") && !strings.Contains(v, ".") {
		errors = append(errors, fmt.Errorf("not a valid tag path. please use one of the following formats:"+
			"\n'dbName'.'schemaName'.'tagName' or dbName|schemaName|tagName or "))
	}
	tagArray := strings.Split(v, ".")
	if len(tagArray) != 3 {
		errors = append(errors, fmt.Errorf("not a valid tag path. please use one of the following formats:"+
			"\n'dbName'.'schemaName'.'tagName' or dbName|schemaName|tagName or "))
	}
	return
}
