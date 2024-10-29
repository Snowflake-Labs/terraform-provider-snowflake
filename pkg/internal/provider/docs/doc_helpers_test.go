package docs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PossibleValuesListedStrings(t *testing.T) {
	values := []string{"abc", "DEF"}

	result := PossibleValuesListed(values)

	assert.Equal(t, "`abc` | `DEF`", result)
}

func Test_PossibleValuesListedInts(t *testing.T) {
	values := []int{42, 21}

	result := PossibleValuesListed(values)

	assert.Equal(t, "`42` | `21`", result)
}

func Test_PossibleValuesListed_empty(t *testing.T) {
	var values []string

	result := PossibleValuesListed(values)

	assert.Empty(t, result)
}
