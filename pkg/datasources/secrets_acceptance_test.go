package datasources_test

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	accConfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	secretWithClientCredentials      = "snowflake_secret_with_client_credentials"
	secretWithAuthorizationCodeGrant = "snowflake_secret_with_authorization_code_grant"
	secretWithBasicAuthentication    = "snowflake_secret_with_basic_authentication"
	secretWithGenericString          = "snowflake_secret_with_generic_string"
)

func TestAcc_Secrets_MultipleTypes(t *testing.T) {
	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "A")
	idTwo := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "B")
	idThree := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "C")
	idFour := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "D")

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_oauth_client_id", "test_oauth_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "first_scope"}, {Scope: "second_scope"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	refreshTokenExpiryDateTime := time.Now().Add(24 * time.Hour).Format(time.DateTime)

	secretModelBasicAuth := model.SecretWithBasicAuthentication("test", idOne.DatabaseName(), idOne.Name(), "test_passwd", idOne.SchemaName(), "test_username")
	secretModelGenericString := model.SecretWithGenericString("test", idTwo.DatabaseName(), idTwo.Name(), idTwo.SchemaName(), "foo")
	secretModelClientCredentials := model.SecretWithClientCredentials("test", integrationId.Name(), idThree.DatabaseName(), idThree.SchemaName(), idThree.Name(), []string{"first_scope", "second_scope"})
	secretModelAuthorizationCodeGrant := model.SecretWithAuthorizationCodeGrant("test", integrationId.Name(), idFour.DatabaseName(), idFour.SchemaName(), idFour.Name(), "test_token", refreshTokenExpiryDateTime)

	multipleSecretModels := accConfig.FromModel(t, secretModelBasicAuth) +
		accConfig.FromModel(t, secretModelGenericString) +
		accConfig.FromModel(t, secretModelClientCredentials) +
		accConfig.FromModel(t, secretModelAuthorizationCodeGrant) +
		secretsDataAndVars()

	configVariables := config.Variables{
		"like": config.StringVariable(prefix + "%"),
	}

	dsName := "data.snowflake_secrets.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config:          multipleSecretModels,
				ConfigVariables: configVariables,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.#", "4")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.name", idOne.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.database_name", idOne.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.schema_name", idOne.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.secret_type", "PASSWORD")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.username", "test_username")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.1.describe_output.0.name", idTwo.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.1.describe_output.0.database_name", idTwo.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.1.describe_output.0.schema_name", idTwo.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.1.describe_output.0.secret_type", "GENERIC_STRING")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.2.describe_output.0.name", idThree.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.2.describe_output.0.database_name", idThree.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.2.describe_output.0.schema_name", idThree.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.2.describe_output.0.secret_type", "OAUTH2")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.2.describe_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.2.describe_output.0.oauth_scopes.*", "first_scope")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.2.describe_output.0.oauth_scopes.*", "second_scope")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.3.describe_output.0.name", idFour.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.3.describe_output.0.database_name", idFour.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.3.describe_output.0.schema_name", idFour.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.3.describe_output.0.secret_type", "OAUTH2")),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "secrets.3.describe_output.0.oauth_refresh_token_expiry_time")),
				),
			},
		},
	})
}

func secretsDataAndVars() string {
	return `
variable "like" {
  type = string
}

data "snowflake_secrets" "test" {
  depends_on = [snowflake_secret_with_basic_authentication.test, snowflake_secret_with_generic_string.test, snowflake_secret_with_client_credentials.test, snowflake_secret_with_authorization_code_grant.test]
  like = var.like
}

`
}

func secretsData(secretResourceName string) string {
	return fmt.Sprintf(`
data "snowflake_secrets" "test" {
  depends_on = [%s.test]
}`, secretResourceName)
}

func TestAcc_Secrets_WithClientCredentials(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_oauth_client_id", "test_oauth_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "username"}, {Scope: "test_scope"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithClientCredentials("test", integrationId.Name(), id.DatabaseName(), id.SchemaName(), id.Name(), []string{"username", "test_scope"})

	dataSecretsClientCredentials := accConfig.FromModel(t, secretModel) + secretsData(secretWithClientCredentials)

	dsName := "data.snowflake_secrets.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			{
				Config: dataSecretsClientCredentials,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.#", "1")),
				),
			},
		},
	})
}
