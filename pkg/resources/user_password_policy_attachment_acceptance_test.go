package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	userName := userId.Name()
	newUserId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newUserName := newUserId.Name()
	passwordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	passwordPolicyName := passwordPolicyId.Name()
	newPasswordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newPasswordPolicyName := newPasswordPolicyId.Name()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             acc.CheckUserPasswordPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userPasswordPolicyAttachmentConfig(userName, acc.TestDatabaseName, acc.TestSchemaName, passwordPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", userName),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "password_policy_name", passwordPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), passwordPolicyId.FullyQualifiedName())),
				),
			},
			// UPDATE
			{
				Config: userPasswordPolicyAttachmentConfig(newUserName, acc.TestDatabaseName, acc.TestSchemaName, newPasswordPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", newUserName),
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
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	userName := userId.Name()
	passwordPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	passwordPolicyName := passwordPolicyId.Name()

	resource.Test(t, resource.TestCase{
		ExternalProviders: map[string]resource.ExternalProvider{
			"snowflake": {
				VersionConstraint: "=0.87.0",
				Source:            "Snowflake-Labs/snowflake",
			},
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckUserPasswordPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userPasswordPolicyAttachmentConfigV087(userName, acc.TestDatabaseName, acc.TestSchemaName, passwordPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", userName),
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

func userPasswordPolicyAttachmentConfig(userName, databaseName, schemaName, passwordPolicyName string) string {
	return fmt.Sprintf(`
resource "snowflake_user" "user" {
	name = "%s"
}

resource "snowflake_password_policy" "pp" {
	database   = "%s"
	schema     = "%s"
	name       = "%s"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
	password_policy_name = snowflake_password_policy.pp.fully_qualified_name
	user_name = snowflake_user.user.name
}
`, userName, databaseName, schemaName, passwordPolicyName)
}

func userPasswordPolicyAttachmentConfigV087(userName, databaseName, schemaName, passwordPolicyName string) string {
	return fmt.Sprintf(`
resource "snowflake_user" "user" {
	name = "%[1]s"
}

resource "snowflake_password_policy" "pp" {
	database   = "%[2]s"
	schema     = "%[3]s"
	name       = "%[4]s"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
	depends_on = [snowflake_password_policy.pp]
	password_policy_name = "\"%[2]s\".\"%[3]s\".\"%[4]s\""
	user_name = snowflake_user.user.name
}
`, userName, databaseName, schemaName, passwordPolicyName)
}
