package snowflake

import (
	"fmt"
	"strings"
)

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

// EscapeSnowflakeString will escape single quotes with the SQL native double single quote
func EscapeSnowflakeString(in string) string {
	out := strings.Replace(in, `'`, `''`, -1)
	return fmt.Sprintf(`'%v'`, out)
}

// UnescapeSnowflakeString reverses EscapeSnowflakeString
func UnescapeSnowflakeString(in string) string {
	out := strings.TrimPrefix(in, `'`)
	out = strings.TrimSuffix(out, `'`)
	out = strings.Replace(out, `''`, `'`, -1)
	return out
}
