package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidDataType(t *testing.T) {
	t.Run("with valid data type", func(t *testing.T) {
		ok := IsValidDataType("VARCHAR")
		assert.Equal(t, ok, true)
	})

	t.Run("with invalid data type", func(t *testing.T) {
		ok := IsValidDataType("foo")
		assert.Equal(t, ok, false)
	})
}
