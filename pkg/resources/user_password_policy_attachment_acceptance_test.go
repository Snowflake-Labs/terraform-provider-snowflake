package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TODO [SNOW-1423486]: Fix using warehouse; remove unsetting testenvs.ConfigureClientOnce
func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	t.Setenv(string(testenvs.ConfigureClientOnce), "")

	user1, user1Cleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(user1Cleanup)

	user2, user2Cleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(user2Cleanup)

	userId := user1.ID()
	newUserId := user2.ID()
	passwordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newPasswordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             acc.CheckUserPasswordPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userPasswordPolicyAttachmentConfig(userId, passwordPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", userId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", passwordPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), passwordPolicyId.FullyQualifiedName())),
				),
			},
			// UPDATE
			{
				Config: userPasswordPolicyAttachmentConfig(newUserId, newPasswordPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", newUserId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", newPasswordPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", newUserId.FullyQualifiedName(), newPasswordPolicyId.FullyQualifiedName())),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_password_policy_attachment.ppa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// Adding this test to check if it will fail sometimes. It should, based on:
// - https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3005
// - https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2627
// but haven't (at least during manual runs).
// The behavior was fixed in https://github.com/Snowflake-Labs/terraform-provider-snowflake/pull/2627
// so the problem should not occur in the newest provider versions.
func TestAcc_UserPasswordPolicyAttachment_gh3005(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	user, userCleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	userId := user.ID()
	passwordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	resource.Test(t, resource.TestCase{
		ExternalProviders: acc.ExternalProviderWithExactVersion("0.87.0"),
		PreCheck:          func() { acc.TestAccPreCheck(t) },
		CheckDestroy:      acc.CheckUserPasswordPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				PreConfig: func() { acc.SetV097CompatibleConfigPathEnv(t) },
				Config:    userPasswordPolicyAttachmentConfigV087(userId, passwordPolicyId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", userId.Name()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", passwordPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), passwordPolicyId.FullyQualifiedName())),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_password_policy_attachment.ppa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userPasswordPolicyAttachmentConfig(userId sdk.AccountObjectIdentifier, passwordPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_password_policy" "pp" {
	database   = "%[2]s"
	schema     = "%[3]s"
	name       = "%[4]s"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
	password_policy_name = snowflake_password_policy.pp.fully_qualified_name
	user_name = "%[1]s"
}
`, userId.Name(), passwordPolicyId.DatabaseName(), passwordPolicyId.SchemaName(), passwordPolicyId.Name())
}

func userPasswordPolicyAttachmentConfigV087(userId sdk.AccountObjectIdentifier, passwordPolicyId sdk.SchemaObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_password_policy" "pp" {
	database   = "%[2]s"
	schema     = "%[3]s"
	name       = "%[4]s"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
	depends_on = [snowflake_password_policy.pp]
	password_policy_name = "\"%[2]s\".\"%[3]s\".\"%[4]s\""
	user_name =  "%[1]s"
}
`, userId.Name(), passwordPolicyId.DatabaseName(), passwordPolicyId.SchemaName(), passwordPolicyId.Name())
}
