package testhelpers

import (
	"io/fs"
	"os"
	"testing"

	"golang.org/x/sys/unix"

	"github.com/stretchr/testify/require"
)

// TestFile creates a temporary file with the given filename and data with the default permissions.
// The directory is automatically removed when the test and all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory.
func TestFile(t *testing.T, filename string, data []byte) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), filename)
	require.NoError(t, err)

	err = os.WriteFile(f.Name(), data, 0o600)
	require.NoError(t, err)
	return f.Name()
}

// TestFileWithPermissions creates a temporary file with the given filename and permissions.
// The directory is automatically removed when the test and all its subtests complete.
// Each subsequent call to t.TempDir returns a unique directory.
func CreateTestFileWithPermissions(t *testing.T, filename string, perms fs.FileMode) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), filename)
	require.NoError(t, err)

	oldMask := unix.Umask(0o000)
	defer unix.Umask(oldMask)

	err = os.Chmod(f.Name(), perms)
	require.NoError(t, err)

	return f.Name()
}
