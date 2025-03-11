package os_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/os"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFileThatIsTooBig(t *testing.T) {
	if os.IsRunningOnWindows {
		t.Skip("checking file sizes on Windows is currently done in manual tests package")
	}
	c := make([]byte, 11*1024*1024)
	configPath := testhelpers.TestFile(t, "config", c)

	_, err := os.ReadFileSafe(configPath)
	require.ErrorContains(t, err, fmt.Sprintf("config file %s is too big - maximum allowed size is 10MB", configPath))
}

func TestLoadConfigFileThatDoesNotExist(t *testing.T) {
	configPath := "non-existing"
	_, err := os.ReadFileSafe(configPath)
	require.ErrorContains(t, err, fmt.Sprintf("reading information about the config file: stat %s: no such file or directory", configPath))
}
