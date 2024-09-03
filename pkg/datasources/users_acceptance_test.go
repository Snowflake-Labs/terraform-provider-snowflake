package datasources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceparametersassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Users_Complete(t *testing.T) {
	id := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	comment := random.Comment()
	pass := random.Password()
	key1, key1Fp := random.GenerateRSAPublicKey(t)
	key2, key2Fp := random.GenerateRSAPublicKey(t)

	userModelNoAttributes := model.User("u", id.Name())
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
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.#", "1")),
					resourceshowoutputassert.UsersDatasourceShowOutput(t, "snowflake_users.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasLoginName(fmt.Sprintf("%s_LOGIN", id.Name())).
						HasDisplayName("Display Name").
						HasFirstName("Jan").
						HasLastName("Testowski").
						HasEmail("fake@email.com").
						HasMustChangePassword(true).
						HasDisabled(false).
						HasSnowflakeLock(false).
						HasDaysToExpiryNotEmpty().
						HasMinsToUnlockNotEmpty().
						HasDefaultWarehouse("some_warehouse").
						HasDefaultNamespace("some.namespace").
						HasDefaultRole("some_role").
						HasDefaultSecondaryRoles(`["ALL"]`).
						HasMinsToBypassMfaNotEmpty().
						HasHasRsaPublicKey(true).
						HasComment(comment),
					resourceparametersassert.UsersDatasourceParameters(t, "snowflake_users.test").
						HasAllDefaults(),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.comment", comment)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.display_name", "Display Name")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.login_name", fmt.Sprintf("%s_LOGIN", id.Name()))),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.first_name", "Jan")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.middle_name", "Jakub")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.last_name", "Testowski")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.email", "fake@email.com")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.password", "********")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.must_change_password", "true")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.disabled", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.snowflake_lock", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.snowflake_support", "false")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.days_to_expiry")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.mins_to_unlock")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_warehouse", "some_warehouse")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_namespace", "some.namespace")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_role", "some_role")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_secondary_roles", `["ALL"]`)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.ext_authn_duo", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.ext_authn_uid", "")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.mins_to_bypass_mfa")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.mins_to_bypass_network_policy", "0")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key", key1)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key_fp", "SHA256:"+key1Fp)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key2", key2)),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key2_fp", "SHA256:"+key2Fp)),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.password_last_set_time")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.custom_landing_page_url", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.custom_landing_page_url_flush_next_ui_load", "false")),
				),
			},
			{
				Config: config.FromModel(t, userModelNoAttributes) + datasourceWithLike(),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.#", "1")),
					resourceshowoutputassert.UsersDatasourceShowOutput(t, "snowflake_users.test").
						HasName(id.Name()).
						HasCreatedOnNotEmpty().
						HasLoginName(strings.ToUpper(id.Name())).
						HasDisplayName("").
						HasFirstName("").
						HasLastName("").
						HasEmail("").
						HasMustChangePassword(false).
						HasDisabled(false).
						HasSnowflakeLock(false).
						HasDaysToExpiry("").
						HasMinsToUnlock("").
						HasDefaultWarehouse("").
						HasDefaultNamespace("").
						HasDefaultRole("").
						HasDefaultSecondaryRoles("").
						HasMinsToBypassMfa("").
						HasHasRsaPublicKey(false).
						HasComment(""),
					resourceparametersassert.UsersDatasourceParameters(t, "snowflake_users.test").
						HasAllDefaults(),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.display_name", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.login_name", strings.ToUpper(id.Name()))),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.first_name", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.middle_name", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.last_name", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.email", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.password", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.must_change_password", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.disabled", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.snowflake_lock", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.snowflake_support", "false")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.days_to_expiry")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.mins_to_unlock")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_warehouse", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_namespace", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_role", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.default_secondary_roles", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.ext_authn_duo", "false")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.ext_authn_uid", "")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.mins_to_bypass_mfa")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.mins_to_bypass_network_policy", "0")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key_fp", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key2", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.rsa_public_key2_fp", "")),
					assert.Check(resource.TestCheckResourceAttrSet("data.snowflake_users.test", "users.0.describe_output.0.password_last_set_time")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.custom_landing_page_url", "")),
					assert.Check(resource.TestCheckResourceAttr("data.snowflake_users.test", "users.0.describe_output.0.custom_landing_page_url_flush_next_ui_load", "false")),
				),
			},
		},
	})
}

func TestAcc_Users_DifferentFiltering(t *testing.T) {
	prefix := random.AlphaN(4)
	id := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id2 := acc.TestClient().Ids.RandomAccountObjectIdentifierWithPrefix(prefix)
	id3 := acc.TestClient().Ids.RandomAccountObjectIdentifier()

	userModel := model.User("u", id.Name())
	user2Model := model.User("u2", id2.Name())
	user3Model := model.User("u3", id3.Name())

	allUsersConfig := config.FromModel(t, userModel) + config.FromModel(t, user2Model) + config.FromModel(t, user3Model)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.User),
		Steps: []resource.TestStep{
			{
				Config: allUsersConfig + datasourceWithLikeMultipleUsers(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_users.test", "users.#", "1"),
				),
			},
			{
				Config: allUsersConfig + datasourceWithStartsWithMultipleUsers(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_users.test", "users.#", "2"),
				),
			},
			{
				Config: allUsersConfig + datasourceWithLimitMultipleUsers(1, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_users.test", "users.#", "1"),
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

func datasourceWithLikeMultipleUsers() string {
	return `
	data "snowflake_users" "test" {
		depends_on = [snowflake_user.u, snowflake_user.u2, snowflake_user.u3]
		like = snowflake_user.u.name
	}
`
}

func datasourceWithStartsWithMultipleUsers(startsWith string) string {
	return fmt.Sprintf(`
	data "snowflake_users" "test" {
		depends_on = [snowflake_user.u, snowflake_user.u2, snowflake_user.u3]
		starts_with = "%s"
	}
`, startsWith)
}

func datasourceWithLimitMultipleUsers(rows int, from string) string {
	return fmt.Sprintf(`
	data "snowflake_users" "test" {
		depends_on = [snowflake_user.u, snowflake_user.u2, snowflake_user.u3]
		limit {
    		rows = %d
			from = "%s"
  		}
	}
`, rows, from)
}

func TestAcc_Users_UserNotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Users/without_user"),
				ExpectError:     regexp.MustCompile("there should be at least one user"),
			},
		},
	})
}
