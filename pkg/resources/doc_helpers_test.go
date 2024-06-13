package resources

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PossibleValuesListed(t *testing.T) {
	values := []string{"abc", "DEF"}

	result := possibleValuesListed(values)

	assert.Equal(t, "`abc` | `DEF`", result)
}

func Test_PossibleValuesListed_empty(t *testing.T) {
	var values []string

	result := possibleValuesListed(values)

	assert.Empty(t, result)
}
