package datasources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	accConfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	secretWithClientCredentials      = "snowflake_secret_with_client_credentials"
	secretWithAuthorizationCodeGrant = "snowflake_secret_with_authorization_code_grant"
	secretWithBasicAuthentication    = "snowflake_secret_with_basic_authentication"
	secretWithGenericString          = "snowflake_secret_with_generic_string"
)

func TestAcc_Secrets_WithClientCredentials(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, dsName).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(sdk.SecretTypeOAuth2),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.show_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.0.show_output.0.oauth_scopes.*", "username")),
					assert.Check(resource.TestCheckTypeSetElemAttr(dsName, "secrets.0.show_output.0.oauth_scopes.*", "test_scope")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.secret_type", sdk.SecretTypeOAuth2)),
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

func TestAcc_Secrets_WithAuthorizationCodeGrant(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_oauth_client_id", "test_oauth_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "username"}, {Scope: "test_scope"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithAuthorizationCodeGrant("test", integrationId.Name(), id.DatabaseName(), id.SchemaName(), id.Name(), "test_token", time.Now().Add(24*time.Hour).Format(time.DateTime)).WithComment("test_comment")

	dataSecretsAuthorizationCode := accConfig.FromModel(t, secretModel) + secretsData(secretWithAuthorizationCodeGrant)

	dsName := "data.snowflake_secrets.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				Config: dataSecretsAuthorizationCode,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, dsName).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("test_comment").
						HasSecretType(sdk.SecretTypeOAuth2),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.show_output.0.oauth_scopes.#", "0")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.secret_type", sdk.SecretTypeOAuth2)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.comment", "test_comment")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.oauth_scopes.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "secrets.0.describe_output.0.oauth_refresh_token_expiry_time")),
				),
			},
		},
	})
}

func TestAcc_Secrets_WithBasicAuthentication(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	secretModel := model.SecretWithBasicAuthentication("test", id.DatabaseName(), id.Name(), "test_passwd", id.SchemaName(), "test_username")
	dataSecretsAuthorizationCode := accConfig.FromModel(t, secretModel) + secretsData(secretWithBasicAuthentication)

	dsName := "data.snowflake_secrets.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			{
				Config: dataSecretsAuthorizationCode,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, dsName).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(sdk.SecretTypePassword),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.show_output.0.oauth_scopes.#", "0")),

					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.secret_type", sdk.SecretTypePassword)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.username", "test_username")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.oauth_scopes.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Secrets_WithGenericString(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	secretModel := model.SecretWithGenericString("test", id.DatabaseName(), id.Name(), id.SchemaName(), "test_secret_string")

	dataSecretsAuthorizationCode := accConfig.FromModel(t, secretModel) + secretsData(secretWithGenericString)

	dsName := "data.snowflake_secrets.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithGenericString),
		Steps: []resource.TestStep{
			{
				Config: dataSecretsAuthorizationCode,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, dsName).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(sdk.SecretTypeGenericString),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.show_output.0.oauth_scopes.#", "0")),

					assert.Check(resource.TestCheckResourceAttrSet(dsName, "secrets.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.secret_type", sdk.SecretTypeGenericString)),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "secrets.0.describe_output.0.oauth_scopes.#", "0")),
				),
			},
		},
	})
}

func secretsData(secretResourceName string) string {
	return fmt.Sprintf(`
    data "snowflake_secrets" "test" {
        depends_on = [%s.test]
    }`, secretResourceName)
}

func TestAcc_Secrets_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	idFour := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_oauth_client_id", "test_oauth_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "first_scope"}, {Scope: "second_scope"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	// ERROR Insufficient privileges to operate on account ... for role BASIC_PRIVILEGES
	/*
		appPkg, appPkgCleanup := acc.TestClient().ApplicationPackage.CreateApplicationPackage(t)
		t.Cleanup(appPkgCleanup)

		_, appCleanup := acc.TestClient().Application.CreateApplication(t, appPkg.ID(), "1")
		t.Cleanup(appCleanup)
	*/

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchemaInDatabase(t, acc.TestClient().Ids.DatabaseId())
	t.Cleanup(schemaCleanup)

	idFive := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	secretModelBasicAuth := model.SecretWithBasicAuthentication("s", idOne.DatabaseName(), idOne.Name(), "test_passwd", idOne.SchemaName(), "test_username")
	secretModelGenericString := model.SecretWithGenericString("s2", idTwo.DatabaseName(), idTwo.Name(), idTwo.SchemaName(), "foo")
	secretModelClientCredentials := model.SecretWithClientCredentials("s3", integrationId.Name(), idThree.DatabaseName(), idThree.SchemaName(), idThree.Name(), []string{"first_scope", "second_scope"})
	secretModelAuthorizationCodeGrant := model.SecretWithAuthorizationCodeGrant("s4", integrationId.Name(), idFour.DatabaseName(), idFour.SchemaName(), idFour.Name(), "test_token", time.Now().Add(24*time.Hour).Format(time.DateTime))
	secretModelInDifferentSchema := model.SecretWithBasicAuthentication("s5", idFive.DatabaseName(), idFive.Name(), "test_passwd", idFive.SchemaName(), "test_username")

	multipleSecretModels := accConfig.FromModel(t, secretModelBasicAuth) +
		accConfig.FromModel(t, secretModelGenericString) +
		accConfig.FromModel(t, secretModelClientCredentials) +
		accConfig.FromModel(t, secretModelAuthorizationCodeGrant) +
		accConfig.FromModel(t, secretModelInDifferentSchema)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.ComposeCheckDestroy(t,
			resources.SecretWithClientCredentials,
			resources.SecretWithAuthorizationCodeGrant,
			resources.SecretWithBasicAuthentication,
			resources.SecretWithGenericString,
		),
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// like with one type
			{
				Config: multipleSecretModels + datasourceWithLikeMultipleSecretTypes("snowflake_secret_with_basic_authentication.s.name"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "1"),
				),
			},
			// like with prefix
			{
				Config: multipleSecretModels + datasourceWithLikeMultipleSecretTypes("\""+prefix+"%"+"\""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "3"),
				),
			},
			// In schema
			{
				Config: multipleSecretModels + secretDatasourceWithIn("schema", idFive.SchemaId().FullyQualifiedName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "1"),
				),
			},
			// In Database
			{
				Config: multipleSecretModels + secretDatasourceWithIn("database", idFive.DatabaseId().FullyQualifiedName()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "5"),
				),
			},
			/*
				// In Application Package
				// ERROR Insufficient privileges to operate on account ... for role BASIC_PRIVILEGES
				{
					Config: multipleSecretModels + secretDatasourceWithIn("application_package", idFive.DatabaseId().FullyQualifiedName()),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "5"),
					),
				},
				// In Application
				// ERROR Insufficient privileges to operate on account ... for role BASIC_PRIVILEGES
				{
					Config: multipleSecretModels + secretDatasourceWithIn("application", idFive.DatabaseId().FullyQualifiedName()),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "5"),
					),
				},
				// In Account
				// ERROR Insufficient privileges to operate on 'SYSTEM' for role BASIC_PRIVILEGES
					{
						Config: multipleSecretModels + secretDatasourceWithIn("account", acc.TestClient().Account.GetAccountIdentifier(t).FullyQualifiedName()),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr("data.snowflake_secrets.test", "secrets.#", "5"),
						),
					},
			*/
		},
	})
}

func datasourceWithLikeMultipleSecretTypes(like string) string {
	return fmt.Sprintf(`
    data "snowflake_secrets" "test" {
        depends_on = [snowflake_secret_with_basic_authentication.s, snowflake_secret_with_generic_string.s2, snowflake_secret_with_client_credentials.s3, snowflake_secret_with_authorization_code_grant.s4]
        like = %s
    }
`, like)
}

func secretDatasourceWithIn(objectName, objectFullyQualifiedName string) string {
	return fmt.Sprintf(`
    data "snowflake_secrets" "test" {
        in {
            %s = "%s"
        }
    }
`, objectName, strings.ReplaceAll(objectFullyQualifiedName, `"`, ""))
}

func TestAcc_Secrets_EmptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      secretDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func secretDatasourceEmptyIn() string {
	return `
    data "snowflake_secrets" "test" {
        in {
        }
    }
`
}
