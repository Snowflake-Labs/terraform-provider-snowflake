package resources_test

import (
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

const (
	managedAccountComment = "Created by a Terraform acceptance test"
)

func TestAcc_ManagedAccount(t *testing.T) {
	// TODO [SNOW-1011985]: unskip the tests
	testenvs.SkipTestIfSet(t, testenvs.SkipManagedAccountTest, "error: 090337 (23001): Number of managed accounts allowed exceeded the limit. Please contact Snowflake support")

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	adminName := acc.TestClient().Ids.Alpha()
	adminPass := acc.TestClient().Ids.AlphaWithPrefix("A1")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.ManagedAccount),
		Steps: []resource.TestStep{
			{
				Config: managedAccountConfig(id.Name(), adminName, adminPass),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "admin_name", adminName),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "admin_password", adminPass),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "comment", managedAccountComment),
					resource.TestCheckResourceAttr("snowflake_managed_account.test", "type", "READER"),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_managed_account.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"admin_name", "admin_password"},
			},
		},
	})
}

func TestAcc_ManagedAccount_HandleShowOutputChanges_BCR_2024_08(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("w", userId.Name())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					acc.TestClient().BcrBundles.EnableBcrBundle(t, "2024_07")
					func() { acc.SetV097CompatibleConfigPathEnv(t) }()
				},
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.97.0",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: config.FromModels(t, userModel),
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasAllDefaults(userId, sdk.SecondaryRolesOptionDefault),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().BcrBundles.EnableBcrBundle(t, "2024_08")
					func() { acc.UnsetConfigPathEnv(t) }()
				},
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModels(t, userModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
				Check: assert.AssertThat(t,
					resourceassert.UserResource(t, userModel.ResourceReference()).
						HasAllDefaults(userId, sdk.SecondaryRolesOptionDefault),
				),
			},
		},
	})
}

func managedAccountConfig(accName, aName, aPass string) string {
	return fmt.Sprintf(`
resource "snowflake_managed_account" "test" {
	name           = "%v"
	admin_name     = "%v"
	admin_password = "%v"
	comment        = "%v"
}
`, accName, aName, aPass, managedAccountComment)
}
