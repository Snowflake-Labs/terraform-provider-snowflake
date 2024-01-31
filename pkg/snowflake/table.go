package snowflake

import (
	"fmt"
)

func QuoteStringList(instrings []string) []string {
	clean := make([]string, 0, len(instrings))
	for _, word := range instrings {
		quoted := fmt.Sprintf(`"%s"`, word)
		clean = append(clean, quoted)
	}
	return clean
}
