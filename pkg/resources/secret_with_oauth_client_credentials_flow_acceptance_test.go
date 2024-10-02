package resources_test

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecretWithClientCredentials_BasicFlow(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := "aaa"
	newComment := random.Comment()

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithClientCredentials("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, []string{"foo", "bar"}).WithComment(comment)
	secretModelWithoutComment := model.SecretWithClientCredentials("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, []string{"foo", "bar"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasCommentString(comment),
					),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "oauth_scopes.*", "bar"),
				),
			},
			// set oauth_scopes and comment in config
			{
				Config: config.FromModel(t, secretModel.
					WithOauthScopes([]string{"foo"}).
					WithComment(newComment)),
				Check: assert.AssertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, "snowflake_secret_with_client_credentials.s").
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(integrationId.Name()).
						HasCommentString(newComment),
					assert.Check(resource.TestCheckResourceAttr("snowflake_secret_with_client_credentials.s", "oauth_scopes.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr("snowflake_secret_with_client_credentials.s", "oauth_scopes.*", "foo")),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, secretModelWithoutComment.WithOauthScopes([]string{"foo"})),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(secretModelWithoutComment.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(newComment), nil),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModelWithoutComment.ResourceReference()).
						HasCommentString(""),
				),
			},
			// destroy
			{
				Config:  config.FromModel(t, secretModelWithoutComment.WithOauthScopes([]string{"foo"})),
				Destroy: true,
			},
			// create without comment
			{
				Config: config.FromModel(t, secretModelWithoutComment.WithOauthScopes([]string{"foo", "bar"})),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, "snowflake_secret_with_client_credentials.s").
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasCommentString(""),
					),
					resource.TestCheckResourceAttr("snowflake_secret_with_client_credentials.s", "oauth_scopes.#", "2"),
					resource.TestCheckTypeSetElemAttr("snowflake_secret_with_client_credentials.s", "oauth_scopes.*", "foo"),
					resource.TestCheckTypeSetElemAttr("snowflake_secret_with_client_credentials.s", "oauth_scopes.*", "bar"),
				),
			},
			// import
			{
				ResourceName:      secretModelWithoutComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "api_authentication", integrationId.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "oauth_scopes.#", "2"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// destroy
			{
				Config:  config.FromModel(t, secretModel),
				Destroy: true,
			},
		},
	})
}

func TestAcc_SecretWithClientCredentials_EmptyScopesList(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()

	integrationId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	_, apiIntegrationCleanup := acc.TestClient().SecurityIntegration.CreateApiAuthenticationClientCredentialsWithRequest(t,
		sdk.NewCreateApiAuthenticationWithClientCredentialsFlowSecurityIntegrationRequest(integrationId, true, "foo", "foo").
			WithOauthAllowedScopes([]sdk.AllowedScope{{Scope: "foo"}, {Scope: "bar"}}),
	)
	t.Cleanup(apiIntegrationCleanup)

	secretModel := model.SecretWithClientCredentials("s", integrationId.Name(), id.DatabaseName(), id.SchemaName(), name, []string{})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithClientCredentials),
		Steps: []resource.TestStep{
			// create secret without providing oauth_scopes value
			{
				Config: config.FromModel(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasApiAuthenticationString(integrationId.Name()).
							HasCommentString(""),
					),
					resource.TestCheckResourceAttr(secretModel.ResourceReference(), "oauth_scopes.#", "0"),
				),
			},
			// Set oauth_scopes
			{
				Config: config.FromModel(t, secretModel.
					WithOauthScopes([]string{"foo"}),
				),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(secretModel.ResourceReference(), "oauth_scopes", tfjson.ActionUpdate, sdk.String("[]"), sdk.String("[foo]")),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasApiAuthenticationString(integrationId.Name()),
					assert.Check(resource.TestCheckResourceAttr(secretModel.ResourceReference(), "oauth_scopes.#", "1")),
					assert.Check(resource.TestCheckTypeSetElemAttr(secretModel.ResourceReference(), "oauth_scopes.*", "foo")),
				),
			},
		},
	})
}
