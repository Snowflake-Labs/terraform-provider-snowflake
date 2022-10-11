package validation

import (
	"fmt"
	"strings"
	"unicode"
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
// lintignore:V011
func ValidatePassword(i interface{}, k string) (s []string, errs []error) {
	pass, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be string", k)}
	}

	if len(pass) < 8 {
		errs = append(errs, fmt.Errorf("password must be at least 8 characters long"))
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
		errs = append(errs, fmt.Errorf("password must contain an uppercase character"))
	}

	if !lowercase {
		errs = append(errs, fmt.Errorf("password must contain a lowercase character"))
	}

	if !digit {
		errs = append(errs, fmt.Errorf("password must contain a digit"))
	}

	return
}

// ValidateIsNotAccountLocator validates that the account value is not an account locator. Account locators have the
// following format: 8 characters where the first 3 characters are letters and the last 5 are digits. ex: ABC12345
// The desired format should be 'organization_name.account_name' ex: testOrgName.testAccName.
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

func ValidateFullyQualifiedObjectID(i interface{}, k string) (s []string, errors []error) {
	v, _ := i.(string)
	if strings.Contains(v, ".") {
		tagArray := strings.Split(v, ".")
		if len(tagArray) != 3 {
			errors = append(errors, fmt.Errorf("%v, is not a valid id. If using period delimiter, three parts must be specified <db_name>.<schema_name>.<object_name>", v))
		}
	} else if strings.Contains(v, "|") {
		tagArray := strings.Split(v, "|")
		if len(tagArray) != 3 {
			errors = append(errors, fmt.Errorf("%v, is not a valid id. If using pipe delimiter, three parts must be specified <db_name>|<schema_name>|<object_name>", v))
		}
	} else {
		errors = append(errors, fmt.Errorf("%v, is not a valid id. please use one of the following formats:"+
			"\n'<db_name>'.'<schema_name>'.'<object_name>' or <db_name>|<schema_name>|<object_name>", v))
	}
	return
}

func FormatFullyQualifiedObjectID(dbName, schemaName, objectName string) string {
	var n strings.Builder

	if dbName != "" && schemaName != "" {
		n.WriteString(fmt.Sprintf(`"%v"."%v".`, dbName, schemaName))
	}

	if dbName != "" && schemaName == "" {
		n.WriteString(fmt.Sprintf(`"%v"..`, dbName))
	}

	if dbName == "" && schemaName != "" {
		n.WriteString(fmt.Sprintf(`"%v".`, schemaName))
	}

	n.WriteString(fmt.Sprintf(`"%v"`, objectName))

	return n.String()
}

func ParseAndFormatFullyQualifiedObectID(s string) string {
	dbName, schemaName, objectName := ParseFullyQualifiedObjectID(s)
	return FormatFullyQualifiedObjectID(dbName, schemaName, objectName)
}

func ParseFullyQualifiedObjectID(s string) (dbName, schemaName, objectName string) {
	parsedString := strings.Replace(s, "\"", "", -1)

	var parts []string
	if strings.Contains(parsedString, "|") {
		parts = strings.Split(parsedString, "|")
	} else if strings.Contains(parsedString, ".") {
		parts = strings.Split(parsedString, ".")
	}
	return parts[0], parts[1], parts[2]
}
