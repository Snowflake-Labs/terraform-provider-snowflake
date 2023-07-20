package sdk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func assertSqlEquals(t *testing.T, opts any, format string, args ...any) {
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	assert.Equal(t, fmt.Sprintf(format, args...), actual)
}
