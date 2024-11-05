package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/importchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SecondaryConnection_Basic(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := acc.SecondaryTestClient().Ids.RandomAccountObjectIdentifier()
	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	// create primary connection
	_, cleanup := acc.SecondaryTestClient().Connection.Create(t, id)
	t.Cleanup(cleanup)
	acc.SecondaryTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).WithEnableConnectionFailover(
		*sdk.NewEnableConnectionFailoverRequest([]sdk.AccountIdentifier{accountId})))

	secondartyConnectionModel := model.SecondaryConnection("t", primaryConnectionAsExternalId.FullyQualifiedName(), id.Name())
	secondartyConnectionModelWithComment := model.SecondaryConnection("t", primaryConnectionAsExternalId.FullyQualifiedName(), id.Name()).WithComment("secondary connection test comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.SecondaryConnection),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, secondartyConnectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondartyConnectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasAsReplicaOfIdentifier(primaryConnectionAsExternalId).
							HasIsPrimaryString("false").
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, secondartyConnectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(acc.SecondaryTestClient().Context.CurrentRegion(t)).
							HasAccountLocator(acc.SecondaryTestClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(false).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId).
							HasConnectionUrl(
								acc.SecondaryTestClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
							),
					),
				),
			},
			// set comment
			{
				Config: config.FromModel(t, secondartyConnectionModelWithComment),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondartyConnectionModelWithComment.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasCommentString("secondary connection test comment"),

						resourceshowoutputassert.ConnectionShowOutput(t, secondartyConnectionModelWithComment.ResourceReference()).
							HasComment("secondary connection test comment"),
					),
				),
			},
			// import
			{
				ResourceName:      secondartyConnectionModelWithComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", "secondary connection test comment"),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, secondartyConnectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondartyConnectionModel.ResourceReference()).
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, secondartyConnectionModel.ResourceReference()).
							HasComment(""),
					),
				),
			},
			// recreate when exteranlly promoted to primary
			{
				PreConfig: func() {
					acc.TestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).WithPrimary(true))
				},
				Config: config.FromModel(t, secondartyConnectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secondartyConnectionModel.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.SecondaryConnectionResource(t, secondartyConnectionModel.ResourceReference()).
							HasIsPrimaryString("false"),

						resourceshowoutputassert.ConnectionShowOutput(t, secondartyConnectionModel.ResourceReference()).
							HasIsPrimary(false),
					),
				),
			},
		},
	})
}
