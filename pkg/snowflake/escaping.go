package snowflake

import "strings"

// EscapeString will escape only the ' character. Would prefer a more robust OSS solution, but this should
// prevent some dumb errors for now.
func EscapeString(in string) string {
	out := strings.Replace(in, `\`, `\\`, -1)
	out = strings.Replace(out, `'`, `\'`, -1)
	return out
}

// UnescapeString reverses EscapeString
func UnescapeString(in string) string {
	out := strings.Replace(in, `\\`, `\`, -1)
	out = strings.Replace(out, `\'`, `'`, -1)
	return out
}
