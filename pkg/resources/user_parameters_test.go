package resources

import (
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func Test_enrichWithReferenceToParameterDocs(t *testing.T) {
	t.Run("formats the docs reference correctly", func(t *testing.T) {
		description := random.Comment()

		enrichedDescription := enrichWithReferenceToParameterDocs(sdk.UserParameterAbortDetachedQuery, description)

		require.Equal(t, description+" "+"For more info, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query).", enrichedDescription)
	})
}

func Test_UserParametersSchema(t *testing.T) {
	t.Run("description references parameter docs correctly", func(t *testing.T) {
		require.True(t, strings.HasSuffix(UserParametersSchema["abort_detached_query"].Description, "For more info, check [ABORT_DETACHED_QUERY docs](https://docs.snowflake.com/en/sql-reference/parameters#abort-detached-query)."))
	})
}
