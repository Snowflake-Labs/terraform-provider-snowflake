package testhelpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFile(t *testing.T, filename string, data []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), filename)
	require.NoError(t, err)

	err = os.WriteFile(f.Name(), data, 0o600)
	require.NoError(t, err)
	return f.Name()
}
