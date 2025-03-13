package oswrapper_test

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/unix"
)

func TestReadFileSafeFailsForFileThatIsTooBig(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file sizes on Windows is currently done in manual tests package")
	}
	c := make([]byte, 11*1024*1024)
	configPath := testhelpers.TestFile(t, "config", c)

	_, err := oswrapper.ReadFileSafe(configPath)
	require.ErrorContains(t, err, fmt.Sprintf("config file %s is too big - maximum allowed size is 10MB", configPath))
}

func TestReadFileSafeFailsForFileThatDoesNotExist(t *testing.T) {
	configPath := "non-existing"
	_, err := oswrapper.ReadFileSafe(configPath)
	require.ErrorContains(t, err, fmt.Sprintf("reading information about the config file: stat %s: no such file or directory", configPath))
}

func TestReadFileSafeFailsForFileWithTooWidePermissions(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	tests := []struct {
		permissions fs.FileMode
	}{
		{permissions: 0o707},
		{permissions: 0o706},
		{permissions: 0o705},
		{permissions: 0o704},
		{permissions: 0o703},
		{permissions: 0o702},
		{permissions: 0o701},

		{permissions: 0o770},
		{permissions: 0o760},
		{permissions: 0o750},
		{permissions: 0o740},
		{permissions: 0o730},
		{permissions: 0o720},
		{permissions: 0o710},
	}
	oldMask := unix.Umask(0o000)
	defer unix.Umask(oldMask)
	for _, tt := range tests {
		t.Run(fmt.Sprintf("reading file with too wide permisssions %#o", tt.permissions), func(t *testing.T) {
			c := random.String()
			path := testhelpers.CreateTestFileWithPermissions(t, "config", []byte(c), tt.permissions)
			_, err := oswrapper.ReadFileSafe(path)
			require.ErrorContains(t, err, fmt.Sprintf("config file %s has unsafe permissions", path))
		})
	}
}

func TestReadFileSafeFailsForFileWithTooRestrictivePermissions(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	tests := []struct {
		permissions fs.FileMode
	}{
		{permissions: 0o300},
		{permissions: 0o200},
		{permissions: 0o100},
	}
	oldMask := unix.Umask(0o000)
	defer unix.Umask(oldMask)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("reading file with too restrictive permisssions %#o", tt.permissions), func(t *testing.T) {
			c := random.String()
			path := testhelpers.CreateTestFileWithPermissions(t, "config", []byte(c), tt.permissions)
			_, err := oswrapper.ReadFileSafe(path)
			require.ErrorContains(t, err, fmt.Sprintf("open %s: permission denied", path))
		})
	}
}

func TestReadFileSafeReadsFileWithCorrectPermissions(t *testing.T) {
	if oswrapper.IsRunningOnWindows() {
		t.Skip("checking file permissions on Windows is currently done in manual tests package")
	}
	tests := []struct {
		permissions fs.FileMode
	}{
		{permissions: 0o700},
		{permissions: 0o600},
		{permissions: 0o500},
		{permissions: 0o400},
		// {permissions: 0o300},
		// {permissions: 0o200},
		// {permissions: 0o100},
	}
	oldMask := unix.Umask(0o000)
	defer unix.Umask(oldMask)

	for _, tt := range tests {
		t.Run(fmt.Sprintf("reading file with correct permisssions %#o", tt.permissions), func(t *testing.T) {
			c := random.String()
			path := testhelpers.CreateTestFileWithPermissions(t, "config", []byte(c), tt.permissions)
			_, err := oswrapper.ReadFileSafe(path)
			require.NoError(t, err)
		})
	}
}
