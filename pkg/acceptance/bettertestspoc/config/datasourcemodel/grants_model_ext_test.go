package datasourcemodel_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/stretchr/testify/require"
)

func Test_GrantsModel(t *testing.T) {
	t.Run("on account", func(t *testing.T) {
		expected := `data "snowflake_grants" "test" {
  grants_on {
    account = true
  }
}
`

		result := config.FromModels(t, datasourcemodel.GrantsOnAccount("test"))

		require.Equal(t, expected, result)
	})
}
