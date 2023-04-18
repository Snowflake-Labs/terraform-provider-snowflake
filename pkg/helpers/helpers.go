package helpers

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ToDo: We can merge these two functions together and also add more functions here with similar functionality

// This function converts list of string into snowflake formated string like 'ele1', 'ele2'.
func ListToSnowflakeString(list []string) string {
	for index, element := range list {
		list[index] = fmt.Sprintf(`'%v'`, strings.ReplaceAll(element, "'", "\\'"))
	}

	return fmt.Sprintf("%v", strings.Join(list, ", "))
}

// IPListToString formats a list of IPs into a Snowflake-DDL friendly string, e.g. ('192.168.1.0', '192.168.1.100').
func IPListToSnowflakeString(ips []string) string {
	for index, element := range ips {
		ips[index] = fmt.Sprintf(`'%v'`, element)
	}

	return fmt.Sprintf("(%v)", strings.Join(ips, ", "))
}

// ListContentToString strips list elements of double quotes or brackets.
func ListContentToString(listString string) string {
	re := regexp.MustCompile(`[\"\[\]]`)
	return re.ReplaceAllString(listString, "")
}

// StringListToList splits a string into a slice of strings, separated by a separator. It also removes empty strings and trims whitespace.
func StringListToList(s string) []string {
	var v []string
	for _, elem := range strings.Split(s, ",") {
		if strings.TrimSpace(elem) != "" {
			v = append(v, strings.TrimSpace(elem))
		}
	}
	return v
}

// StringToBool converts a string to a bool.
func StringToBool(s string) bool {
	return strings.ToLower(s) == "true"
}

// SnowflakeID generates a unique ID for a resource.
func SnowflakeID(attributes ...interface{}) string {
	var parts []string
	for i, attr := range attributes {
		if attr == nil {
			attributes[i] = ""
		}
		switch reflect.TypeOf(attr).Kind() {
		case reflect.String:
			parts = append(parts, attr.(string))
		case reflect.Bool:
			parts = append(parts, strconv.FormatBool(attr.(bool)))
		case reflect.Slice:
			parts = append(parts, strings.Join(attr.([]string), ","))
		}
	}
	return strings.Join(parts, "|")
}

const IDDelimiter = "|"
