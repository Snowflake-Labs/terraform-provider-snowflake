package snowflake

import (
	"fmt"
	"regexp"
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

// AddressEscape wraps a name inside double quotes only if required by Snowflake
func AddressEscape(in ...string) string {
	quoteCheck := regexp.MustCompile(`[^A-Z0-9_]`)
	address := make([]string, len(in))

	for i, n := range in {
		if quoteCheck.MatchString(n) {
			address[i] = fmt.Sprintf(`"%s"`, strings.Replace(n, `"`, `\"`, -1))
		} else {
			address[i] = n
		}
	}

	return strings.Join(address, ".")
}
