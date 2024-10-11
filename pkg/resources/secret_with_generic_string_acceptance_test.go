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

func TestAcc_SecretWithGenericString_BasicFlow(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()

	secretModel := model.SecretWithGenericString("s", id.DatabaseName(), name, id.SchemaName(), "foo")
	secretModelEmptySecretString := model.SecretWithGenericString("s", id.DatabaseName(), name, id.SchemaName(), "")

	secretName := secretModel.ResourceReference()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecretWithGenericString),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, secretModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithGenericStringResource(t, secretModel.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasSecretStringString("foo").
							HasCommentString(""),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasName(name).
							HasDatabaseName(id.DatabaseName()).
							HasSecretType(sdk.SecretTypeGenericString).
							HasSchemaName(id.SchemaName()).
							HasComment(""),
					),

					resource.TestCheckResourceAttr(secretName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttrSet(secretName, "describe_output.0.created_on"),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.name", name),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.database_name", id.DatabaseName()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.schema_name", id.SchemaName()),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.secret_type", sdk.SecretTypeGenericString),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.username", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.comment", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_access_token_expiry_time", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_refresh_token_expiry_time", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.integration_name", ""),
					resource.TestCheckResourceAttr(secretName, "describe_output.0.oauth_scopes.#", "0"),
				),
			},
			// set secret_string and comment
			{
				Config: config.FromModel(t, secretModel.
					WithSecretString("bar").
					WithComment(comment),
				),

				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithGenericStringResource(t, secretName).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasSecretStringString("bar").
							HasCommentString(comment),

						resourceshowoutputassert.SecretShowOutput(t, secretModel.ResourceReference()).
							HasSecretType(sdk.SecretTypeGenericString).
							HasComment(comment),
					),

					resource.TestCheckResourceAttr(secretName, "describe_output.0.comment", comment),
				),
			},
			// set comment externally, external changes for secret_string are not being detected
			{
				PreConfig: func() {
					acc.TestClient().Secret.Alter(t, sdk.NewAlterSecretRequest(id).
						WithSet(*sdk.NewSecretSetRequest().
							WithComment("test_comment"),
						),
					)
				},
				Config: config.FromModel(t, secretModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretName, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(secretName, "comment", sdk.String(comment), sdk.String("test_comment")),
						planchecks.ExpectChange(secretName, "comment", tfjson.ActionUpdate, sdk.String("test_comment"), sdk.String(comment)),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.SecretWithGenericStringResource(t, secretName).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasSecretStringString("bar").
						HasCommentString(comment),
				),
			},
			// import
			{
				ResourceName:            secretName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_string"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, secretModelEmptySecretString),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secretModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(secretModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModelEmptySecretString.ResourceReference()).
						HasCommentString(""),
				),
			},
			// import with no fields set
			{
				ResourceName:            secretModel.ResourceReference(),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret_string"},
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "database", id.DatabaseId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "schema", id.SchemaId().Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", ""),
				),
			},
			// destroy
			{
				Config:  config.FromModel(t, secretModel),
				Destroy: true,
			},
			// create with empty secret_string
			{
				Config: config.FromModel(t, secretModelEmptySecretString),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecretWithGenericStringResource(t, secretModelEmptySecretString.ResourceReference()).
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasSecretStringString("").
							HasCommentString(""),
					),
				),
			},
		},
	})
}
