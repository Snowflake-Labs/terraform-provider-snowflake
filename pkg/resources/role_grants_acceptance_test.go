package resources_test

// import (
// 	"fmt"
// 	"strings"
// 	"testing"

// 	"github.com/hashicorp/terraform/helper/acctest"
// 	"github.com/hashicorp/terraform/helper/resource"
// )

// func TestAccGrantRole(t *testing.T) {
// 	prefix := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
// 	prefix2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

// 	resource.Test(t, resource.TestCase{
// 		Providers: providers(),
// 		Steps: []resource.TestStep{
// 			{
// 				Config: uConfig(prefix),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "name", strings.ToUpper(prefix)),
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "comment", "test comment"),
// 				),
// 			},
// 			// RENAME
// 			{
// 				Config: uConfig(prefix2),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "name", strings.ToUpper(prefix2)),
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "comment", "test comment"),
// 				),
// 			},
// 			// CHANGE PROPERTIES
// 			{
// 				Config: uConfig2(prefix2),
// 				Check: resource.ComposeTestCheckFunc(
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "name", strings.ToUpper(prefix2)),
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "comment", "test comment 2"),
// 					resource.TestCheckResourceAttr("snowflake_grant_role.w", "password", "best password"),
// 				),
// 			},
// 			// IMPORT
// 			{
// 				ResourceName:            "snowflake_grant_role.w",
// 				ImportState:             true,
// 				ImportStateVerify:       true,
// 				ImportStateVerifyIgnore: []string{"password"},
// 			},
// 		},
// 	})
// }

// func uConfig(prefix string) string {
// 	s := `
// resource "snowflake_grant_role" "w" {
// 	name = "%s"
// 	comment = "test comment"
// }
// `
// 	return fmt.Sprintf(s, prefix)
// }

// func uConfig2(prefix string) string {
// 	s := `
// resource "snowflake_grant_role" "w" {
// 	name = "%s"
// 	comment = "test comment 2"
// 	password = "best password"
// }
// `
// 	return fmt.Sprintf(s, prefix)
// }
