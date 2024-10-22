package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecretWithBasicAuthentication_BasicFlow(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()

	secretModel := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "foo", id.SchemaName(), "foo")
	secretModelDifferentCredentialsWithComment := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "bar", id.SchemaName(), "bar").WithComment(comment)
	secretModelWithoutComment := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "bar", id.SchemaName(), "bar")
	secretModelEmptyCredentials := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "", id.SchemaName(), "")

	secretName := secretModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretName).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("foo").
							HasPasswordString("foo").
							HasCommentString(""),

						resourceshowoutputassert.SecretShowOutput(t, secretName).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(string(sdk.SecretTypePassword)).
							HasSchemaName(id.SchemaName()).
							HasComment(""),
					),

					resource.TestCheckResourceAttr(secretName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(secretName, "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.name", name),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.secret_type", string(sdk.SecretTypePassword)),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.username", "foo"),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.comment", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_refresh_token_expiry_time", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.integration_name", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// set username, password and comment
			{
				Config: config.FromModel(t, secretModelDifferentCredentialsWithComment),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,

						resourceassert.SecretWithBasicAuthenticationResource(t, secretName).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("bar").
							HasPasswordString("bar").
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretName).
							HasSecretType(string(sdk.SecretTypePassword)).
							HasComment(comment),
					),

					resource.TestCheckResourceAttr(secretName, "describe_output.0.username", "bar"),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.comment", comment),
				),
			},
			// set username and comment externally
			{
				PreConfig: func() {
					acc.TestClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).
						WithSet(*sdk.NewSecretSetRequest().
							WithComment("test_comment").
							WithSetForFlow(*sdk.NewSetForFlowRequest().
								WithSetForBasicAuthentication(*sdk.NewSetForBasicAuthenticationRequest().
									WithUsername("test_username"),
								),
							),
						),
					)
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretName, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(secretName, "comment", sdk.String(comment), sdk.String("test_comment")),
						planchecks.ExpectDrift(secretName, "username", sdk.String("bar"), sdk.String("test_username")),

						planchecks.ExpectChange(secretName, "comment", tfjson.ActionUpdate, sdk.String("test_comment"), sdk.String(comment)),
						planchecks.ExpectChange(secretName, "username", tfjson.ActionUpdate, sdk.String("test_username"), sdk.String("bar")),
					},
				},
				Config: config.FromModel(t, secretModelDifferentCredentialsWithComment),
				Check: assert.AssertThat(t,
					resourceassert.SecretWithBasicAuthenticationResource(t, secretName).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUsernameString("bar").
						HasPasswordString("bar").
						HasCommentString(comment),
				),
			},
			// import
			{
				ResourceName:            secretName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "username", "bar"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, secretModelWithoutComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretName, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretName, "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretName).
						HasCommentString(""),
				),
			},
			// import with no fields set
			{
				ResourceName:            secretName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "username", "bar"),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// set empty username and password
			{
				Config: config.FromModel(t, secretModelEmptyCredentials),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretName).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("").
							HasPasswordString("").
							HasCommentString(""),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithBasicAuthentication_CreateWithEmptyCredentials(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	secretModelEmptyCredentials := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "", id.SchemaName(), "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, secretModelEmptyCredentials),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretModelEmptyCredentials.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("").
							HasPasswordString("").
							HasCommentString(""),
					),
				),
			},
		},
	})
}

func TestAcc_SecretWithBasicAuthentication_ExternalSecretTypeChange(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	secretModel := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "test_pswd", id.SchemaName(), "test_usr")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithBasicAuthentication),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypePassword)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypePassword)),
					),
				),
			},
			// create or replace with different secret type
			{
				PreConfig: func() {
					acc.TestClient().Secret.DropFunc(t, id)()
					_, cleanup := acc.TestClient().Secret.CreateWithGenericString(t, id, "test_secret_string")
					t.Cleanup(cleanup)
				},
				Config: config.FromModel(t, secretModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithBasicAuthenticationResource(t, secretModel.ResourceReference()).
							HasSecretTypeString(string(sdk.SecretTypePassword)),
						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(string(sdk.SecretTypePassword)),
					),
				),
			},
		},
	})
}
