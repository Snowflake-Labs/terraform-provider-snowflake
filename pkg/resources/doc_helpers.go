package resources

import (
	"fmt"
	"strings"
)

func possibleValuesListed(values []string) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%s`", value)
	}
	return strings.Join(valuesWrapped, " | ")
}
