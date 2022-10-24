package helpers

import (
	"fmt"
	"regexp"
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
