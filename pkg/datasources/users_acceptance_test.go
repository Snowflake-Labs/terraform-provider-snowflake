package datasources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Users_Complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.ConfigureClientOnce)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	pass := random.Password()
	key1, _ := random.GenerateRSAPublicKey(t)
	key2, _ := random.GenerateRSAPublicKey(t)

	userModelAllAttributes := model.User("u", id.Name()).
		WithPassword(pass).
		WithLoginName(id.Name() + "_login").
		WithDisplayName("Display Name").
		WithFirstName("Jan").
		WithMiddleName("Jakub").
		WithLastName("Testowski").
		WithEmail("fake@email.com").
		WithMustChangePassword("true").
		WithDisabled("false").
		WithDaysToExpiry(8).
		WithMinsToUnlock(9).
		WithDefaultWarehouse("some_warehouse").
		WithDefaultNamespace("some.namespace").
		WithDefaultRole("some_role").
		WithDefaultSecondaryRolesStringList("ALL").
		WithMinsToBypassMfa(10).
		WithRsaPublicKey(key1).
		WithRsaPublicKey2(key2).
		WithComment(comment).
		WithDisableMfa("true")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, userModelAllAttributes) + datasourceWithLike(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_users.test", "users.#", "1"),

					resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.show_output.0.name", id.Name()),
				),
			},
		},
	})
}

func datasourceWithLike() string {
	return `
	data "snowflake_users" "test" {
		like = snowflake_user.u.name
	}
`
}
