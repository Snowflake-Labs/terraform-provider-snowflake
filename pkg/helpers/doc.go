package helpers

import (
	"fmt"
	"strings"
)

func PossibleValuesListed[T ~string | ~int](values []T) string {
	valuesWrapped := make([]string, len(values))
	for i, value := range values {
		valuesWrapped[i] = fmt.Sprintf("`%v`", value)
	}
	return strings.Join(valuesWrapped, " | ")
}
