package datasources_test

import (
	"fmt"
	"testing"

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

func connectionsData() string {
	return `
    data "snowflake_connections" "test" {
        depends_on = [snowflake_connection.test]
    }`
}

func TestAcc_Connections(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	//_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	accountId := acc.TestClient().Account.GetAccountIdentifier(t)
	prefix := random.AlphaN(4)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	connectionModel := model.Connection("test", id.Name())

	primaryConnectionAsExternalId := sdk.NewExternalObjectIdentifier(accountId, id)

	dataConnections := accConfig.FromModel(t, connectionModel) + connectionsData()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Connection),
		Steps: []resource.TestStep{
			{
				Config: dataConnections,
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "1")),
					resourceshowoutputassert.ConnectionShowOutput(t, "snowflake_connections.test").
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
							acc.TestClient().Connection.GetConnectionUrl(accountId.OrganizationName(), id.Name()),
						),
				),
			},
		},
	})
}

func TestAcc_Connections_Filtering(t *testing.T) {
	// TODO: [SNOW-1002023]: Unskip; Business Critical Snowflake Edition needed
	_ = testenvs.GetOrSkipTest(t, testenvs.TestFailoverGroups)

	prefix := random.AlphaN(4)
	idOne := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idTwo := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	idThree := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)

	connectionModelOne := model.Connection("c1", idOne.Name())
	connectionModelTwo := model.Connection("c2", idTwo.Name())
	connectionModelThree := model.Connection("c3", idThree.Name())

	configWithLike := accConfig.FromModel(t, connectionModelOne) +
		accConfig.FromModel(t, connectionModelTwo) +
		accConfig.FromModel(t, connectionModelThree)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Connection),
		Steps: []resource.TestStep{
			// with like
			{
				Config: configWithLike + connectionDatasourceWithLike(prefix+"%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_connections.test", "connections.#", "3"),
				),
			},
		},
	})
}

func connectionDatasourceWithLike(like string) string {
	return fmt.Sprintf(`
    data "snowflake_connections" "test" {
        depends_on = ["snowflake_connection.c1"]
        like = "%s"
    }
`, like)
}
