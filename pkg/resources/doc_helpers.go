package resources

import (
	"fmt"
	"strings"
)

func possibleValuesListed[T ~string](values []T) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%s`", value)
	}
	return strings.Join(valuesWrapped, " | ")
}

func characterList(values []rune) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%c`", value)
	}
	return strings.Join(valuesWrapped, ", ")
}
