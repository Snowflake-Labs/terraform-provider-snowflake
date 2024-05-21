package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	// TODO [SNOW-1423486]: unskip
	t.Skipf("Skip because error %s; will be fixed in SNOW-1423486", "Error: 000606 (57P03): No active warehouse selected in the current session.  Select an active warehouse with the 'use warehouse' command.")
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
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), newPasswordPolicyId.FullyQualifiedName())),
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
	password_policy_name = "\"${snowflake_password_policy.pp.database}\".\"${snowflake_password_policy.pp.schema}\".\"${snowflake_password_policy.pp.name}\""
	user_name = snowflake_user.user.name
}
`, userName, databaseName, schemaName, passwordPolicyName)
}
