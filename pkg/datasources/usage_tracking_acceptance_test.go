package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/tracking"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/datasources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_CompleteUsageTracking(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	schemaModel := model.Schema("test", id.DatabaseName(), id.Name())

	assertQueryMetadataExists := func(t *testing.T, query string) resource.TestCheckFunc {
		t.Helper()
		return func(state *terraform.State) error {
			queryHistory := acc.TestClient().InformationSchema.GetQueryHistory(t, 100)
			expectedMetadata := tracking.NewVersionedDatasourceMetadata(datasources.Schemas)
			if _, err := collections.FindFirst(queryHistory, func(history helpers.QueryHistory) bool {
				metadata, err := tracking.ParseMetadata(history.QueryText)
				return err == nil &&
					expectedMetadata == metadata &&
					strings.Contains(history.QueryText, query)
			}); err != nil {
				return fmt.Errorf("query history does not contain query metadata: %v with query containing: %s", expectedMetadata, query)
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, schemaModel) + schemaDatasourceConfigWithDependency(schemaModel.ResourceReference(), id),
				Check: assert.AssertThat(t,
					resourceassert.SchemaResource(t, schemaModel.ResourceReference()).
						HasNameString(id.Name()),
					assert.Check(assertQueryMetadataExists(t, fmt.Sprintf(`SHOW SCHEMAS LIKE '%s' IN DATABASE "%s"`, id.Name(), id.DatabaseName()))),
					// SHOW PARAMETERS IN SCHEMA "acc_test_db_AT_1AB7E1DE_1A10_89C3_C13C_899754A250B6"."FPGDHEAT_1AB7E1DE_1A10_89C3_C13C_899754A250B6" --terraform_provider_usage_tracking {"json_schema_version":"1","version":"v0.99.0","datasource":"snowflake_schemas","operation":"read"}
					assert.Check(assertQueryMetadataExists(t, fmt.Sprintf(`SHOW PARAMETERS IN SCHEMA %s`, id.FullyQualifiedName()))),
					assert.Check(assertQueryMetadataExists(t, fmt.Sprintf(`DESCRIBE SCHEMA %s`, id.FullyQualifiedName()))),
				),
			},
		},
	})
}

func schemaDatasourceConfigWithDependency(schemaResourceReference string, id sdk.DatabaseObjectIdentifier) string {
	return fmt.Sprintf(`
data "snowflake_schemas" "test" {
	depends_on = [ %[1]s ]
	in {
		database = "%[2]s"
	}
	like = "%[3]s"
}
`, schemaResourceReference, id.DatabaseName(), id.Name())
}
