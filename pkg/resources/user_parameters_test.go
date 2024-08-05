package resources

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_UserParametersSchema(t *testing.T) {
	t.Run("description references parameter docs correctly", func(t *testing.T) {
		require.True(t, strings.HasSuffix(userParametersSchema["abort_detached_query"].Description, "For more information, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query)."))
	})
}
