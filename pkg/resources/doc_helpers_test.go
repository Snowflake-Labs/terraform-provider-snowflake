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

func Test_PossibleValuesListedInt(t *testing.T) {
	values := []int{42, 21}

	result := possibleValuesListedInt(values)

	assert.Equal(t, "`42` | `21`", result)
}

func Test_PossibleValuesListedInt_empty(t *testing.T) {
	var values []int

	result := possibleValuesListedInt(values)

	assert.Empty(t, result)
}
