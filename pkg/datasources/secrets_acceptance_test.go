//go:build !account_level_tests

package datasources_test

import (
	"regexp"
	"testing"
	"time"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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
	secretsModel := datasourcemodel.Secrets("test").
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(secretModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, secretModel, secretsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, "snowflake_secrets.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(string(sdk.SecretTypeOAuth2)),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.show_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretsModel.DatasourceReference(), "secrets.0.show_output.0.oauth_scopes.*", "username")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretsModel.DatasourceReference(), "secrets.0.show_output.0.oauth_scopes.*", "test_scope")),

					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "2")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.*", "username")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.*", "test_scope")),
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
	secretsModel := datasourcemodel.Secrets("test").
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(secretModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithAuthorizationCodeGrant),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, secretModel, secretsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, "snowflake_secrets.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("test_comment").
						HasSecretType(string(sdk.SecretTypeOAuth2)),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.show_output.0.oauth_scopes.#", "0")),

					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeOAuth2))),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.comment", "test_comment")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "0")),
					assert.Check(resource.TestCheckResourceAttrSet(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_refresh_token_expiry_time")),
				),
			},
		},
	})
}

func TestAcc_Secrets_WithBasicAuthentication(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	secretModel := model.SecretWithBasicAuthentication("test", id.DatabaseName(), id.Name(), "test_passwd", id.SchemaName(), "test_username")
	secretsModel := datasourcemodel.Secrets("test").
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(secretModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, secretModel, secretsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, "snowflake_secrets.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(string(sdk.SecretTypePassword)),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.show_output.0.oauth_scopes.#", "0")),

					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypePassword))),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.username", "test_username")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Secrets_WithGenericString(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	secretModel := model.SecretWithGenericString("test", id.DatabaseName(), id.Name(), id.SchemaName(), "test_secret_string")
	secretsModel := datasourcemodel.Secrets("test").
		WithInDatabase(id.DatabaseId()).
		WithDependsOn(secretModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithGenericString),
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, secretModel, secretsModel),
				Check: assertThat(t,
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.#", "1")),
					resourceshowoutputassert.SecretsDatasourceShowOutput(t, "snowflake_secrets.test").
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasComment("").
						HasSecretType(string(sdk.SecretTypeGenericString)),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.show_output.0.oauth_scopes.#", "0")),

					assert.Check(resource.TestCheckResourceAttrSet(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.secret_type", string(sdk.SecretTypeGenericString))),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.username", "")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(secretsModel.DatasourceReference(), "secrets.0.describe_output.0.oauth_scopes.#", "0")),
				),
			},
		},
	})
}

func TestAcc_Secrets_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "test_oauth_client_id", "test_oauth_client_secret").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "first_scope"}, {Scope: "second_scope"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	schema, schemaCleanup := acc.TestClient().Schema.CreateSchema(t)
	t.Cleanup(schemaCleanup)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "1")
	idTwo := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "2")
	idThree := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix + "3")
	idFour := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	idFive := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	pass := random.Password()

	secretModelBasicAuth := model.SecretWithBasicAuthentication("s", idOne.DatabaseName(), idOne.Name(), pass, idOne.SchemaName(), "test_username")
	secretModelGenericString := model.SecretWithGenericString("s2", idTwo.DatabaseName(), idTwo.Name(), idTwo.SchemaName(), pass)
	secretModelClientCredentials := model.SecretWithClientCredentials("s3", integrationId.Name(), idThree.DatabaseName(), idThree.SchemaName(), idThree.Name(), []string{"first_scope", "second_scope"})
	secretModelAuthorizationCodeGrant := model.SecretWithAuthorizationCodeGrant("s4", integrationId.Name(), idFour.DatabaseName(), idFour.SchemaName(), idFour.Name(), pass, time.Now().Add(24*time.Hour).Format(time.DateTime))
	secretModelInDifferentSchema := model.SecretWithBasicAuthentication("s5", idFive.DatabaseName(), idFive.Name(), pass, idFive.SchemaName(), "test_username")
	allSecretModels := []accconfig.ResourceModel{secretModelBasicAuth, secretModelGenericString, secretModelClientCredentials, secretModelAuthorizationCodeGrant, secretModelInDifferentSchema}
	allReferences := collections.Map(allSecretModels, func(resourceModel accconfig.ResourceModel) string { return resourceModel.ResourceReference() })

	secretsModelWithLike := datasourcemodel.Secrets("test").
		WithLike(idOne.Name()).
		WithInDatabase(idOne.DatabaseId()).
		WithDependsOn(allReferences...)
	secretsModelWithLikePrefix := datasourcemodel.Secrets("test").
		WithLike(prefix + "%").
		WithInDatabase(idOne.DatabaseId()).
		WithDependsOn(allReferences...)
	secretsModelInSchema := datasourcemodel.Secrets("test").
		WithInSchema(idFive.SchemaId()).
		WithDependsOn(allReferences...)
	secretsModelInDatabase := datasourcemodel.Secrets("test").
		WithInDatabase(idFive.DatabaseId()).
		WithDependsOn(allReferences...)
	secretsModelInAccount := datasourcemodel.Secrets("test").
		WithInAccount().
		WithLike(prefix + "%").
		WithDependsOn(allReferences...)

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
				Config: accconfig.FromModels(t, secretModelBasicAuth, secretModelGenericString, secretModelClientCredentials, secretModelAuthorizationCodeGrant, secretModelInDifferentSchema, secretsModelWithLike),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secretsModelWithLike.DatasourceReference(), "secrets.#", "1"),
				),
			},
			// like with prefix
			{
				Config: accconfig.FromModels(t, secretModelBasicAuth, secretModelGenericString, secretModelClientCredentials, secretModelAuthorizationCodeGrant, secretModelInDifferentSchema, secretsModelWithLikePrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secretsModelWithLikePrefix.DatasourceReference(), "secrets.#", "3"),
				),
			},
			// In schema
			{
				Config: accconfig.FromModels(t, secretModelBasicAuth, secretModelGenericString, secretModelClientCredentials, secretModelAuthorizationCodeGrant, secretModelInDifferentSchema, secretsModelInSchema),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secretsModelInSchema.DatasourceReference(), "secrets.#", "1"),
				),
			},
			// In Database
			{
				Config: accconfig.FromModels(t, secretModelBasicAuth, secretModelGenericString, secretModelClientCredentials, secretModelAuthorizationCodeGrant, secretModelInDifferentSchema, secretsModelInDatabase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secretsModelInDatabase.DatasourceReference(), "secrets.#", "5"),
				),
			},
			// In Account
			{
				Config: accconfig.FromModels(t, secretModelBasicAuth, secretModelGenericString, secretModelClientCredentials, secretModelAuthorizationCodeGrant, secretModelInDifferentSchema, secretsModelInAccount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(secretsModelInAccount.DatasourceReference(), "secrets.#", "3"),
				),
			},
		},
	})
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
