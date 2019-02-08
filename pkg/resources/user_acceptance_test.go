package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccUser(t *testing.T) {
	t.Parallel()
	prefix := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix2 := randomdata.Email()

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: uConfig(prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix))),
				),
			},
			// RENAME
			{
				Config: uConfig(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix2))),
				),
			},
			// CHANGE PROPERTIES
			{
				Config: uConfig2(prefix2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment 2"),
					resource.TestCheckResourceAttr("snowflake_user.w", "password", "best password"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix2))),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_user.w",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func uConfig(prefix string) string {
	s := `
resource "snowflake_user" "w" {
	name = "%s"
	comment = "test comment"
	login_name = "%s_login"
}
`
	return fmt.Sprintf(s, prefix, prefix)
}

func uConfig2(prefix string) string {
	s := `
resource "snowflake_user" "w" {
	name = "%s"
	comment = "test comment 2"
	password = "best password"
	login_name = "%s_login"
}
`
	return fmt.Sprintf(s, prefix, prefix)
}
