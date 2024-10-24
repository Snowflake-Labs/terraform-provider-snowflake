package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_UserAuthenticationPolicyAttachment(t *testing.T) {
	// TODO [SNOW-1423486]: unskip
	t.Skipf("Skip because error %s; will be fixed in SNOW-1423486", "Error: 000606 (57P03): No active warehouse selected in the current session.  Select an active warehouse with the 'use warehouse' command.")
	userId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	userName := userId.Name()
	newUserId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newUserName := newUserId.Name()
	authenticationPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	authenticationPolicyName := authenticationPolicyId.Name()
	newAuthenticationPolicyId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	newAuthenticationPolicyName := newAuthenticationPolicyId.Name()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             acc.CheckUserAuthenticationPolicyAttachmentDestroy(t),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userAuthenticationPolicyAttachmentConfig(userName, acc.TestDatabaseName, acc.TestSchemaName, authenticationPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "user_name", userName),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "authentication_policy_name", authenticationPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), authenticationPolicyId.FullyQualifiedName())),
				),
			},
			// UPDATE
			{
				Config: userAuthenticationPolicyAttachmentConfig(newUserName, acc.TestDatabaseName, acc.TestSchemaName, newAuthenticationPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "user_name", newUserName),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "authentication_policy_name", newAuthenticationPolicyId.FullyQualifiedName()),
					resource.TestCheckResourceAttr("snowflake_user_authentication_policy_attachment.ppa", "id", fmt.Sprintf("%s|%s", userId.FullyQualifiedName(), newAuthenticationPolicyId.FullyQualifiedName())),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_user_authentication_policy_attachment.ppa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func userAuthenticationPolicyAttachmentConfig(userName, databaseName, schemaName, authenticationPolicyName string) string {
	return fmt.Sprintf(`
resource "snowflake_user" "user" {
	name = "%s"
}

resource "snowflake_authentication_policy" "ap" {
	database   = "%s"
	schema     = "%s"
	name       = "%s"
}

resource "snowflake_user_authentication_policy_attachment" "apa" {
	authentication_policy_name = snowflake_authentication_policy.ap.fully_qualified_name
	user_name = snowflake_user.user.name
}
`, userName, databaseName, schemaName, authenticationPolicyName)
}
