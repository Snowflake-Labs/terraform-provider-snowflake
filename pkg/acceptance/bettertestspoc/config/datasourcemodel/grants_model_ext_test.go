package datasourcemodel_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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

	t.Run("on account object", func(t *testing.T) {
		id := randomAccountObjectIdentifier()
		expected := fmt.Sprintf(`data "snowflake_grants" "test" {
  grants_on {
    object_name = "%s"
    object_type = "DATABASE"
  }
}
`, id.Name())

		result := config.FromModels(t, datasourcemodel.GrantsOnAccountObject("test", id, sdk.ObjectTypeDatabase))

		require.Equal(t, expected, result)
	})

	t.Run("on database object", func(t *testing.T) {
		id := randomDatabaseObjectIdentifier()
		expected := fmt.Sprintf(`data "snowflake_grants" "test" {
  grants_on {
    object_name = "\"%s\".\"%s\""
    object_type = "SCHEMA"
  }
}
`, id.DatabaseName(), id.Name())

		result := config.FromModels(t, datasourcemodel.GrantsOnDatabaseObject("test", id, sdk.ObjectTypeSchema))

		require.Equal(t, expected, result)
	})
}

func randomAccountObjectIdentifier() sdk.AccountObjectIdentifier {
	return sdk.NewAccountObjectIdentifier(random.StringN(12))
}

func randomDatabaseObjectIdentifier() sdk.DatabaseObjectIdentifier {
	return sdk.NewDatabaseObjectIdentifier(random.StringN(12), random.StringN(12))
}

func randomSchemaObjectIdentifier() sdk.SchemaObjectIdentifier {
	return sdk.NewSchemaObjectIdentifier(random.StringN(12), random.StringN(12), random.StringN(12))
}
