package resources_test

import (
	// "context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	// "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform2" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: acc.TestAccProviders(),
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		// CheckDestroy: testAccCheckYourResourceDestroy,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// CREATE
			{
				// TODO: handle the case where the user is in lowercase
				Config: userPasswordPolicyAttachmentConfig("USER_PASSWORD_POLICY_ATTACHMENT_ACCEPTANCE_TEST_USER", acc.TestDatabaseName, acc.TestSchemaName, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fmt.Sprintf("snowflake_user_password_policy_attachment.password_policy_attachment_acceptance_test_user_%s", prefix), "id"),
				),
			},
			// UPDATE
			{
				Config: userPasswordPolicyAttachmentConfig("USER_PASSWORD_POLICY_ATTACHMENT_ACCEPTANCE_TEST_USER2", acc.TestDatabaseName, acc.TestSchemaName, prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(fmt.Sprintf("snowflake_user_password_policy_attachment.password_policy_attachment_acceptance_test_user_%s", prefix2), "id"),
					resource.TestCheckResourceAttr(fmt.Sprintf("snowflake_user_password_policy_attachment.password_policy_attachment_acceptance_test_user_%s", prefix2), "user_name", "USER_PASSWORD_POLICY_ATTACHMENT_ACCEPTANCE_TEST_USER2"),
				),
			},

			// TODO: importer
		},
	})
}

func testAccCheckYourResourceDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	// client := sdk.NewClientFromDB(db)
	// ctx := context.Background()
	// for _, rs := range s.RootModule().Resources {
	// 	// Note: I leverage the fact that the state during the test is specific to the test case, so there should only be there resources created in this test
	// 	if rs.Type != "snowflake_user_password_policy_attachment" {
	// 		continue
	// 	}
	// 	user_name := rs.Primary.Attributes["user_name"]
	// 	policyReferences, err := client.PolicyReferences.GetForEntity(ctx, &sdk.GetForEntityPolicyReferenceRequest{
	// 		RefEntityName:   user_name,
	// 		RefEntityDomain: "user",
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if len(policyReferences) > 0 {
	// 		return fmt.Errorf("User Password Policy attachment %v still exists", policyReferences[0].PolicyName)
	// 	}
	// }
	if db != nil {
		return nil
	}
	return nil
}

func userPasswordPolicyAttachmentConfig(userName, databaseName, schemaName, prefix string) string {
	s := `
resource "snowflake_user" "user_password_policy_attachment_acceptance_test_user_%s" {
	name = "%s"
}
resource "snowflake_password_policy" "password_policy_password_policy_attachment_acceptance_test_user_%s" {
	database   = "%s"
	schema     = "%s"
	name       = "%v"
}

resource "snowflake_user_password_policy_attachment" "password_policy_attachment_acceptance_test_user_%s" {
	password_policy = snowflake_password_policy.password_policy_password_policy_attachment_acceptance_test_user_%s.qualified_name
	user_name = snowflake_user.user_password_policy_attachment_acceptance_test_user_%s.name
}
`
	return fmt.Sprintf(s, prefix, userName, prefix, databaseName, schemaName, prefix, prefix, prefix, prefix)
}
