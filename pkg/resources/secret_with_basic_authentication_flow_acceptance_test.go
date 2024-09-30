package resources_test

import (
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"testing"
)

func TestAcc_SecretWithBasicAuthentication_BasicFlow(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	name := id.Name()
	comment := random.Comment()

	secretModel := model.SecretWithBasicAuthentication("s", id.DatabaseName(), name, "foo", id.SchemaName(), "foo")

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
							HasNameString(name).
							HasDatabaseString(id.DatabaseName()).
							HasSchemaString(id.SchemaName()).
							HasUsernameString("foo").
							HasPasswordString("foo").
							HasNoComment(),
						//HasCommentString(""),
					),
				),
			},
			// set username, password and comment
			{
				Config: config.FromModel(t, secretModel.
					WithUsername("bar").
					WithPassword("bar").
					WithComment(comment),
				),
				Check: assert.AssertThat(t,
					resourceassert.SecretWithBasicAuthenticationResource(t, secretModel.ResourceReference()).
						HasNameString(name).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasUsernameString("bar").
						HasPasswordString("bar").
						HasCommentString(comment),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, secretModel.WithComment("")),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planchecks.ExpectChange(secretModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String(comment), nil),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.SecretWithClientCredentialsResource(t, secretModel.ResourceReference()).
						HasCommentString(""),
				),
			},
			// destroy
			{
				Config:  config.FromModel(t, secretModel),
				Destroy: true,
			},
			/*
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
			*/
		},
	})
}
