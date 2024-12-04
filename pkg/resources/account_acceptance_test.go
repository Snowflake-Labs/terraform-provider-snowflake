package resources_test

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Account_complete(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)

	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	firstName := acc.TestClient().Ids.Alpha()
	lastName := acc.TestClient().Ids.Alpha()
	email := random.Email()
	name := random.AdminName()
	key, _ := random.GenerateRSAPublicKey(t)
	region := acc.TestClient().Context.CurrentRegion(t)
	comment := random.Comment()

	configModel := model.Account("test", id.Name(), string(sdk.EditionStandard), email, name).
		// TODO: WithAdminUserType()
		WithAdminRsaPublicKey(key).
		WithFirstName(firstName).
		WithLastName(lastName).
		WithMustChangePassword(true).
		WithRegionGroup("PUBLIC").
		WithRegion(region).
		WithComment(comment).
		WithIsOrgAdmin(true).
		WithGracePeriodInDays(3)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Account),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, configModel),
				Check: assert.AssertThat(t,
					resourceassert.AccountResource(t, configModel.ResourceReference()).
						HasNameString(id.Name()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAdminNameString(name).
						HasAdminRsaPublicKeyString(key).
						HasEmailString(email).
						HasFirstNameString(firstName).
						HasLastNameString(lastName).
						HasMustChangePasswordString(r.BooleanTrue).
						HasRegionGroupString("PUBLIC").
						HasRegionString(region).
						HasCommentString(comment).
						HasIsOrgAdminString(r.BooleanTrue).
						HasGracePeriodInDaysString("3"),
					// TODO: Show output
				),
			},
		},
	})
}

//func accountConfig(name string, password string, comment string, gracePeriodInDays int) string {
//	return fmt.Sprintf(`
//data "snowflake_current_account" "current" {}
//
//resource "snowflake_account" "test" {
//  name = "%s"
//  admin_name = "someadmin"
//  admin_password = "%s"
//  first_name = "Ad"
//  last_name = "Min"
//  email = "admin@example.com"
//  must_change_password = false
//  edition = "BUSINESS_CRITICAL"
//  comment = "%s"
//  region = data.snowflake_current_account.current.region
//  grace_period_in_days = %d
//}
//`, name, password, comment, gracePeriodInDays)
//}
