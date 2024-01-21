package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAcc_UserPasswordPolicyAttachment(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		// resource.Test(t, resource.TestCase{
		Providers: acc.TestAccProviders(),
		// ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck: func() { acc.TestAccPreCheck(t) },
		// CheckDestroy:             testAccCheckYourResourceDestroy,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// CREATE
			{
				// TODO: handle the case where the user is in lowercase
				Config: userPasswordPolicyAttachmentConfig("USER", acc.TestDatabaseName, acc.TestSchemaName, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_user_password_policy_attachment.ppa", "id"),
					// printState(prefix),
				),
				Destroy: false,
			},
			// UPDATE
			{
				Config: userPasswordPolicyAttachmentConfig("USER2", acc.TestDatabaseName, acc.TestSchemaName, prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_user_password_policy_attachment.ppa", "id"),
					// TODO: change the USER
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", "USER2"),
					// printState(prefix2),
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

// func printState(prefix2 string) resource.TestCheckFunc {
// 	return func(s *terraform.State) error {
// 		for _, module := range s.Modules {
// 			for key, rs := range module.Resources {
// 				fmt.Printf("Resource Address: %s\n", key)
// 				fmt.Printf("Resource Type: %s\n", rs.Type)
// 				fmt.Printf("Resource ID: %s\n", rs.Primary.ID)
// 				for attrKey, attrValue := range rs.Primary.Attributes {
// 					fmt.Printf("  %s = %s\n", attrKey, attrValue)
// 				}
// 			}
// 		}
// 		fmt.Printf("ANDTHERENDERsnowflake_user_password_policy_attachment.password_policy_attachment_acceptance_test_user_%s\n", prefix2)
// 		return nil
// 	}
// }

func testAccCheckYourResourceDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		// Note: I leverage the fact that the state during the test is specific to the test case, so there should only be there resources created in this test
		if rs.Type != "snowflake_user_password_policy_attachment" {
			continue
		}
		user_name := rs.Primary.Attributes["user_name"]
		policyReferences, err := client.PolicyReferences.GetForEntity(ctx, &sdk.GetForEntityPolicyReferenceRequest{
			RefEntityName:   user_name,
			RefEntityDomain: "user",
		})
		if err != nil {
			if strings.Contains(err.Error(), "does not exist or not authorized") {
				// Note: this can happen if the Policy Reference or the User have been deleted as well; in this case, just ignore the error
				continue
			}
			return err
		}
		if len(policyReferences) > 0 {
			return fmt.Errorf("User Password Policy attachment %v still exists", policyReferences[0].PolicyName)
		}
	}
	return nil
}

// func userPasswordPolicyAttachmentConfig(userName, databaseName, schemaName, prefix string) string {
// 	s := `
// resource "snowflake_user" "user_password_policy_attachment_acceptance_test_user_%s" {
// 	name = "%s"
// }
// resource "snowflake_password_policy" "password_policy_password_policy_attachment_acceptance_test_user_%s" {
// 	database   = "%s"
// 	schema     = "%s"
// 	name       = "%v"
// }

// resource "snowflake_user_password_policy_attachment" "password_policy_attachment_acceptance_test_user_%s" {
// 	password_policy = snowflake_password_policy.password_policy_password_policy_attachment_acceptance_test_user_%s.qualified_name
// 	user_name = snowflake_user.user_password_policy_attachment_acceptance_test_user_%s.name
// }
// `
// 	return fmt.Sprintf(s, prefix, userName, prefix, databaseName, schemaName, prefix, prefix, prefix, prefix)
// }

// TODO: the USER needs to be suffixed, but for this we need to fix the lower - upper case problem
func userPasswordPolicyAttachmentConfig(userName, databaseName, schemaName, prefix string) string {
	s := `
resource "snowflake_user" "user" {
	name = "%s"
}
resource "snowflake_password_policy" "pp" {
	database   = "%s"
	schema     = "%s"
	name       = "pp_%v"
}

resource "snowflake_user_password_policy_attachment" "ppa" {
	password_policy_database = snowflake_password_policy.pp.database
	password_policy_schema = snowflake_password_policy.pp.schema
	password_policy_name = snowflake_password_policy.pp.name
	user_name = snowflake_user.user.name
}
`
	return fmt.Sprintf(s, userName, databaseName, schemaName, prefix)
}
