package resources

import (
	"fmt"
	"strings"
)

// TODO [SNOW-999049]: replace during identifiers rework
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

// TODO [SNOW-999049]: replace during identifiers rework
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
