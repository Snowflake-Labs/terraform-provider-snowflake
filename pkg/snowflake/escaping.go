package snowflake

import (
	"fmt"
	"strings"
)

// EscapeString will escape only the ' character. Would prefer a more robust OSS solution, but this should
// prevent some dumb errors for now.
func EscapeString(in string) string {
	out := strings.ReplaceAll(in, `\`, `\\`)
	out = strings.ReplaceAll(out, `'`, `\'`)
	return out
}

// EscapeSnowflakeString will escape single quotes with the SQL native double single quote.
func EscapeSnowflakeString(in string) string {
	out := strings.ReplaceAll(in, `'`, `''`)
	return fmt.Sprintf(`'%v'`, out)
}

// UnescapeSnowflakeString reverses EscapeSnowflakeString.
func UnescapeSnowflakeString(in string) string {
	out := strings.TrimPrefix(in, `'`)
	out = strings.TrimSuffix(out, `'`)
	out = strings.ReplaceAll(out, `''`, `'`)
	return out
}
