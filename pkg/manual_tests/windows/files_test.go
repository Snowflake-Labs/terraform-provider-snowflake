package windows_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/oswrapper"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFileThatIsTooBig(t *testing.T) {
	if !oswrapper.IsRunningOnWindows {
		t.Skip("checking file sizes on other platforms is currently done in the sdk package")
	}
	c := make([]byte, 11*1024*1024)
	configPath := testhelpers.TestFile(t, "config", c)

	_, err := sdk.LoadConfigFile(configPath)
	require.ErrorContains(t, err, fmt.Sprintf("config file %s is too big - maximum allowed size is 10MB", configPath))
}
