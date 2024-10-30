package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accConfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func connectionsData(secretResourceName string) string {
	return fmt.Sprintf(`
    data "snowflake_connections" "test" {
        depends_on = [%s.test]
    }`, secretResourceName)
}

func TestAcc_Connections(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	connectionModel := model.Connection("test", id.Name())

	dataConnections := accConfig.FromModel(t, connectionModel) //+ secretsData(secretWithClientCredentials)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Connection),
		Steps: []resource.TestStep{
			{
				Config: dataSecretsClientCredentials,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, "snowflake_secrets.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(string(sdk.SecretTypeOAuth2)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.show_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.0.show_output.0.oauth_scopes.*", "username")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.0.show_output.0.oauth_scopes.*", "test_scope")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.0.describe_output.0.oauth_scopes.*", "username")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.0.describe_output.0.oauth_scopes.*", "test_scope")),
				),
			},
		},
	})
}
