package validation

import (
	"fmt"
	"strings"
	"unicode"
)

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

func FormatFullyQualifiedObjectID(dbName, schemaName, objectName string) string {
	var n strings.Builder

	if dbName == "" {
		if schemaName == "" {
			if objectName == "" {
				return n.String()
			}
			n.WriteString(fmt.Sprintf(`"%v"`, objectName))
			return n.String()
		}
		n.WriteString(fmt.Sprintf(`"%v"`, schemaName))
		if objectName == "" {
			return n.String()
		}
		n.WriteString(fmt.Sprintf(`."%v"`, objectName))
		return n.String()
	} // dbName != ""
	n.WriteString(fmt.Sprintf(`"%v"`, dbName))
	if schemaName == "" {
		if objectName == "" {
			return n.String()
		}
		n.WriteString(fmt.Sprintf(`."%v"`, objectName))
		return n.String()
	} // schemaName != ""
	n.WriteString(fmt.Sprintf(`."%v"`, schemaName))
	if objectName == "" {
		return n.String()
	}
	n.WriteString(fmt.Sprintf(`."%v"`, objectName))
	return n.String()
}

func ParseFullyQualifiedObjectID(s string) (dbName, schemaName, objectName string) {
	parsedString := strings.ReplaceAll(s, "\"", "")

	var parts []string
	if strings.Contains(parsedString, "|") {
		parts = strings.Split(parsedString, "|")
	} else if strings.Contains(parsedString, ".") {
		parts = strings.Split(parsedString, ".")
	}
	for len(parts) < 3 {
		parts = append(parts, "")
	}
	return parts[0], parts[1], parts[2]
}
