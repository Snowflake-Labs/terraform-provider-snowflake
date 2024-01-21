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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		CheckDestroy:             testAccCheckYourResourceDestroy,
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: userPasswordPolicyAttachmentConfig("USER", acc.TestDatabaseName, acc.TestSchemaName, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_user_password_policy_attachment.ppa", "id"),
				),
				Destroy: false,
			},
			// UPDATE
			{
				Config: userPasswordPolicyAttachmentConfig(fmt.Sprintf("USER_%s", prefix), acc.TestDatabaseName, acc.TestSchemaName, prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("snowflake_user_password_policy_attachment.ppa", "id"),
					resource.TestCheckResourceAttr("snowflake_user_password_policy_attachment.ppa", "user_name", fmt.Sprintf("USER_%s", prefix)),
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

func testAccCheckYourResourceDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	ctx := context.Background()
	for _, rs := range s.RootModule().Resources {
		// Note: I leverage the fact that the state during the test is specific to the test case, so there should only be there resources created in this test
		if rs.Type != "snowflake_user_password_policy_attachment" {
			continue
		}
		userName := sdk.NewAccountObjectIdentifierFromFullyQualifiedName(rs.Primary.Attributes["user_name"])
		policyReferences, err := client.PolicyReferences.GetForEntity(ctx, &sdk.GetForEntityPolicyReferenceRequest{
			RefEntityName:   sdk.String(userName.FullyQualifiedName()),
			RefEntityDomain: sdk.String("user"),
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
