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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Connection_Basic(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	// comment := random.Comment()

	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	// secondaryAccountId := acc.SecondaryTestClient().Account.GetAccountIdentifier(t)
	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	getConnectionUrl := func(organizationName, objectName string) string {
		return strings.ToLower(fmt.Sprintf("%s-%s.snowflakecomputing.com", organizationName, objectName))
	}

	connectionModel := model.Connection("t", id.Name())

	/*
		connectionModelWithFailover := model.Connection("test_connection_with_failover", id.Name()).
			WithEnableFailover([]sdk.AccountIdentifier{secondaryAccountId})

		connectionModelWithComment := model.Connection("test_connection", id.Name()).
			WithComment(comment)

		replicatedConnectionModel := model.Connection("replicated_test_connection", id.Name()).
			WithAsReplicaOf(primaryConnectionAsExternalId.FullyQualifiedName())

		replicatedConnectionModelWithPrimary := model.Connection("replicated_test_connection", id.Name()).
			WithAsReplicaOf(primaryConnectionAsExternalId.FullyQualifiedName()).
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
							HasEnableFailover([]sdk.AccountIdentifier{accountId}).
							HasIsPrimaryString("true").
							HasNoComment(),

						resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
							HasName(id.Name()).
							HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
							HasAccountLocator(acc.TestClient().GetAccountLocator()).
							HasAccountName(accountId.AccountName()).
							HasOrganizationName(accountId.OrganizationName()).
							HasComment("").
							HasIsPrimary(true).
							HasPrimaryIdentifier(primaryConnectionAsExternalId).
							HasFailoverAllowedToAccounts([]sdk.AccountIdentifier{accountId}).
							HasConnectionUrl(
								getConnectionUrl(accountId.OrganizationName(), id.Name()),
							),
					),
				),
			},
			/*
				// TODO: [SNOW-1763442]
				// enable failover to secondary account
				{
					Config: config.FromModel(t, connectionModelWithFailover),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, connectionModelWithFailover.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasNoAsReplicaOf().
								HasIsPrimaryString("true").
								HasEnableFailover([]sdk.AccountIdentifier{
									accountId,
									secondaryAccountId,
								}),

							resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithFailover.ResourceReference()).
								HasIsPrimary(true).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts([]sdk.AccountIdentifier{
									accountId,
									secondaryAccountId,
								}),
						),
					),
				},
				// TODO: [SNOW-1763442]
				// create as replica of
				{
					Config: config.FromModel(t, replicatedConnectionModel),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, replicatedConnectionModel.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasAsReplicaOfString(primaryConnectionAsExternalId.FullyQualifiedName()).
								HasEnableFailover([]sdk.AccountIdentifier{accountId, secondaryAccountId}).
								HasIsPrimaryString("false").
								HasNoComment(),

							resourceshowoutputassert.ConnectionShowOutput(t, replicatedConnectionModel.ResourceReference()).
								HasName(id.Name()).
								HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
								HasAccountLocator(acc.TestClient().GetAccountLocator()).
								HasAccountName(accountId.AccountName()).
								HasOrganizationName(accountId.OrganizationName()).
								HasComment("").
								HasIsPrimary(false).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts([]sdk.AccountIdentifier{accountId, secondaryAccountId}).
								HasConnectionUrl(
									getConnectionUrl(secondaryAccountId.OrganizationName(), id.Name()),
								),
						),
					),
				},
				// TODO: [SNOW-1763442]
				// promote to primary
				{
					Config: config.FromModel(t, replicatedConnectionModel),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, replicatedConnectionModel.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasAsReplicaOfString(primaryConnectionAsExternalId.FullyQualifiedName()).
								HasEnableFailover([]sdk.AccountIdentifier{
									accountId,
									secondaryAccountId,
								}).
								HasIsPrimaryString("false").
								HasNoComment(),

							resourceshowoutputassert.ConnectionShowOutput(t, replicatedConnectionModel.ResourceReference()).
								HasName(id.Name()).
								HasSnowflakeRegion(acc.TestClient().Context.CurrentRegion(t)).
								HasAccountLocator(acc.TestClient().GetAccountLocator()).
								HasAccountName(accountId.AccountName()).
								HasOrganizationName(accountId.OrganizationName()).
								HasComment("").
								HasIsPrimary(false).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts([]sdk.AccountIdentifier{
									accountId,
									secondaryAccountId,
								}).
								HasConnectionUrl(
									getConnectionUrl(secondaryAccountId.OrganizationName(), id.Name()),
								),
						),
					),
				},
				// TODO: [SNOW-1763442]
				// disable failover to accounts
				{
					Config: config.FromModel(t, replicatedConnectionModelWithPrimary),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, replicatedConnectionModelWithPrimary.ResourceReference()).
								HasIsPrimaryString("true").
								HasEnableFailover([]sdk.AccountIdentifier{
									accountId,
									secondaryAccountId,
								}),

							resourceshowoutputassert.ConnectionShowOutput(t, replicatedConnectionModelWithPrimary.ResourceReference()).
								HasIsPrimary(true).
								HasPrimaryIdentifier(primaryConnectionAsExternalId).
								HasFailoverAllowedToAccounts([]sdk.AccountIdentifier{
									accountId,
									secondaryAccountId,
								}),
						),
					),
				},
				// set comment
				{
					Config: config.FromModel(t, connectionModelWithComment),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, connectionModelWithComment.ResourceReference()).
								HasNameString(id.Name()).
								HasFullyQualifiedNameString(id.FullyQualifiedName()).
								HasNoAsReplicaOf().
								HasEnableFailover([]sdk.AccountIdentifier{accountId}).
								HasIsPrimaryString("true").
								HasCommentString(comment),

							resourceshowoutputassert.ConnectionShowOutput(t, connectionModelWithComment.ResourceReference()).
								HasComment(comment),
						),
					),
				},
				// unset comment
				{
					Config: config.FromModel(t, connectionModel),
					Check: resource.ComposeTestCheckFunc(
						assert.AssertThat(t,
							resourceassert.ConnectionResource(t, connectionModel.ResourceReference()).
								HasNoComment(),

							resourceshowoutputassert.ConnectionShowOutput(t, connectionModel.ResourceReference()).
								HasComment(""),
						),
					),
				},
			*/
		},
	})
}
