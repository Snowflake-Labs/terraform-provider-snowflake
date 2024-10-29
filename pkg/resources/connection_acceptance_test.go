package resources_test

import (
	"fmt"
	"strings"
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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Connection_Basic(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	comment := random.Comment()

	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	getConnectionUrl := func(organizationName, objectName string) string {
		return strings.ToLower(fmt.Sprintf("%s-%s.snowflakecomputing.com", organizationName, objectName))
	}

	connectionModel := model.Connection("t", id.Name())
	connectionModelWithComment := model.Connection("t", id.Name()).WithComment(comment)

	// TODO: [SNOW-1763442]
	/*
	   secondaryAccountId := acc.SecondaryTestClient().Account.GetAccountIdentifier(t)

	   connectionModelWithFailover := model.Connection("t", id.Name()).WithEnableFailover(secondaryAccountId)

	   replicatedConnectionModel := model.Connection("replicated_connection", id.Name()).
	       WithAsReplicaOfIdentifier(primaryConnectionAsExternalId)

	   replicatedConnectionModelWithPrimary := model.Connection("replicated_connection", id.Name()).
	       WithAsReplicaOfIdentifier(primaryConnectionAsExternalId).
	       WithIsPrimary(true)
	*/

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Connection),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasNoAsReplicaOf().
							HasNoEnableFailoverToAccounts().
							HasNoIsPrimary().
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
							HasAccountLocator(acc.TestClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(true).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId).
							HasConnectionUrl(
								getConnectionUrl(accountId.OrganizationName(), id.Name()),
							),
					),
				),
			},
			// TODO: [SNOW-1763442]
			// enable failover to secondary account
			/*
				{
					Config: config.FromModel(t, connectionModelWithFailover),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, connectionModelWithFailover.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasNoAsReplicaOf().
								HasNoIsPrimary().
								HasEnableFailoverToAccounts(secondaryAccountId),

							resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithFailover.ResourceReference()).
								HasIsPrimary(true).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts(
									accountId,
									secondaryAccountId,
								),
						),
					),
				},
			*/
			// TODO: [SNOW-1763442]
			// create as replica of
			/*
				{
					Config: config.FromModel(t, replicatedConnectionModel),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, replicatedConnectionModel.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasAsReplicaOfIdentifier(primaryConnectionAsExternalId).
								HasNoEnableFailoverToAccounts().
								HasIsPrimaryString("false").
								HasNoComment(),

							resourceshowoutputassert.ConnectionShowOutput(t, replicatedConnectionModel.ResourceReference()).
								HasName(id.Name()).
								HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
								HasAccountLocator(acc.TestClient().GetAccountLocator()).
								HasAccountName(secondaryAccountId.AccountName()).
								HasOrganizationName(secondaryAccountId.OrganizationName()).
								HasComment("").
								HasIsPrimary(false).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts(secondaryAccountId).
								HasConnectionUrl(
									getConnectionUrl(secondaryAccountId.OrganizationName(), id.Name()),
								),
						),
					),
				},
			*/
			// TODO: [SNOW-1763442]
			// promote to primary
			/*
				{
					Config: config.FromModel(t, replicatedConnectionModelWithPrimary),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, replicatedConnectionModelWithPrimary.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasAsReplicaOfString(primaryConnectionAsExternalId.FullyQualifiedName()).
								HasNoEnableFailoverToAccounts().
								HasIsPrimaryString("true").
								HasNoComment(),

							resourceshowoutputassert.ConnectionShowOutput(t, replicatedConnectionModelWithPrimary.ResourceReference()).
								HasName(id.Name()).
								HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
								HasAccountLocator(acc.TestClient().GetAccountLocator()).
								HasAccountName(secondaryAccountId.AccountName()).
								HasOrganizationName(secondaryAccountId.OrganizationName()).
								HasComment("").
								HasIsPrimary(true).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts(
									secondaryAccountId,
								).
								HasConnectionUrl(
									getConnectionUrl(secondaryAccountId.OrganizationName(), id.Name()),
								),
						),
					),
				},
			*/
			// TODO: [SNOW-1763442]
			// disable failover to accounts
			/*
				{
					Config: config.FromModel(t, connectionModel),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
								HasIsPrimaryString("true").
								HasNoEnableFailoverToAccounts(),

							resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
								HasIsPrimary(true).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts(
									accountId,
								),
						),
					),
				},
			*/
			// set comment
			{
				Config: config.FromModel(t, connectionModelWithComment),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, connectionModelWithComment.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasNoAsReplicaOf().
							HasNoEnableFailoverToAccounts().
							HasNoIsPrimary().
							HasCommentString(comment),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithComment.ResourceReference()).
							HasComment(comment),
					),
				),
			},
			// import
			{
				ResourceName:      connectionModelWithComment.ResourceReference(),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateCheck: importchecks.ComposeImportStateCheck(
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "name", id.Name()),
					importchecks.TestCheckResourceAttrInstanceState(helpers.EncodeResourceIdentifier(id), "comment", comment),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasComment(""),
					),
				),
			},
		},
	})
}

func TestAcc_Connection_ExternalChanges(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	connectionModel := model.Connection("t", id.Name()).WithComment("config comment")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Connection),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasNoAsReplicaOf().
							HasNoEnableFailoverToAccounts().
							HasNoIsPrimary().
							HasCommentString("config comment"),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
							HasAccountLocator(acc.TestClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("config comment").
							HasIsPrimary(true).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId),
					),
				),
			},
			// change comment externally
			{
				PreConfig: func() {
					acc.TestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).
						WithSet(*sdk.NewSetRequest().
							WithComment("external comment")))
				},
				Config: config.FromModel(t, connectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(connectionModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(connectionModel.ResourceReference(), "comment", tfjson.ActionUpdate, sdk.String("external comment"), sdk.String("config comment")),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
							HasCommentString("config comment"),
						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasComment("config comment"),
					),
				),
			},
		},
	})
}

// TODO: [SNOW-1763442]
/*
func TestAcc_Connection_ExternalChangeToIsPrimary(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	secondaryAccountId := acc.SecondaryTestClient().Account.GetAccountIdentifier(t)

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	connectionModel := model.Connection("t", id.Name()).WithEnableFailover(secondaryAccountId)
	secondaryConnectionModel := model.Connection("replicated", id.Name()).WithAsReplicaOfIdentifier(primaryConnectionAsExternalId)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Connection),
		Steps: []resource.TestStep{
			// create
			{
				Config: config.FromModel(t, connectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasNoAsReplicaOf().
							HasEnableFailoverToAccounts(secondaryAccountId).
							HasNoIsPrimary().
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
							HasAccountLocator(acc.TestClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(true).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(accountId, secondaryAccountId),
					),
				),
			},
			// create replication
			{
				Config: config.FromModel(t, secondaryConnectionModel),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, secondaryConnectionModel.ResourceReference()).
							HasNameString(id.Name()).
							HasFullyQualifiedNameString(id.FullyQualifiedName()).
							HasAsReplicaOfIdentifier(primaryConnectionAsExternalId).
							HasNoEnableFailoverToAccounts().
							HasIsPrimaryString("false").
							HasCommentString(""),

						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(acc.SecondaryTestClient().Context.CurrentRegion(t)).
							HasAccountLocator(acc.SecondaryTestClient().GetAccountLocator()).
							HasAccountName(secondaryAccountId.AccountName()).
							HasOrganizationName(secondaryAccountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(false).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts(secondaryAccountId),
					),
				),
			},
			// change to primary externally
			{
				PreConfig: func() {
					acc.SecondaryTestClient().Connection.Alter(t, sdk.NewAlterConnectionRequest(id).WithPrimary(true))
				},
				Config: config.FromModel(t, secondaryConnectionModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(secondaryConnectionModel.ResourceReference(), plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.ConnectionResource(t, secondaryConnectionModel.ResourceReference()).
							HasIsPrimaryString("false"),
						resourceshowoutputassert.ConnectionShowOutput(t, secondaryConnectionModel.ResourceReference()).
							HasIsPrimary(false),
					),
				),
			},
		},
	})
}
*/
