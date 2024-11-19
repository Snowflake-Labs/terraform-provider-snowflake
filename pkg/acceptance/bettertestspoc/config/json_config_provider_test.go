package config_test

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_JsonConfigProvider(t *testing.T) {

	t.Run("test resource json config", func(t *testing.T) {
		model := Some("some_name", "abc")
		expectedResult := `{
    "resource": {
        "snowflake_share": {
            "some_name": {
                "name": "abc"
            }
        }
    }
}`

		result, err := config.DefaultJsonConfigProvider.ResourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result[:]))
	})

	t.Run("test datasource json config", func(t *testing.T) {
		model := datasourcemodel.Databases("some_name")
		expectedResult := `{
    "data": {
        "snowflake_databases": {
            "some_name": {}
        }
    }
}`

		result, err := config.DefaultJsonConfigProvider.DatasourceJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result[:]))
	})

	t.Run("test provider json config", func(t *testing.T) {
		model := providermodel.SnowflakeProvider()
		expectedResult := `{
    "provider": {
        "snowflake": {}
    }
}`

		result, err := config.DefaultJsonConfigProvider.ProviderJsonFromModel(model)
		require.NoError(t, err)
		assert.Equal(t, expectedResult, string(result[:]))
	})
}
