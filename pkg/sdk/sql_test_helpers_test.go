package sdk

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func assertSqlEquals(t *testing.T, opts any, expected string) {
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
