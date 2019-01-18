package snowflake

import "strings"

// QuoteString will escape only the ' character. Would prefer a more robust OSS solution, but this should
// prevent some dumb errors for now.
func EscapeString(in string) string {
	return strings.Replace(in, `'`, `\'`, -1)

}
