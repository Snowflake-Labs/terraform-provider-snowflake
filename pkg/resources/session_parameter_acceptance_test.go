package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_SessionParameterWithUser(t *testing.T) {
	prefix := "TEST_USER_" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sessionParameterWithUser(prefix, "BINARY_OUTPUT_FORMAT", "BASE64"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "key", "BINARY_OUTPUT_FORMAT"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "value", "BASE64"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "user", prefix),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func TestAcc_SessionParameterOnAccount(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sessionParameterOnAccount("AUTOCOMMIT", "false"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "key", "AUTOCOMMIT"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "value", "false"),
					resource.TestCheckResourceAttr("snowflake_session_parameter.p", "on_account", "true"),
				),
			},
		},
	})
}

func sessionParameterWithUser(user, key, value string) string {
	s := `
resource "snowflake_user" "u" {
	name = "%s"
}

resource "snowflake_session_parameter" "p" {
	key = "%s"
	value = "%s"
	user = snowflake_user.u.name
}
`
	return fmt.Sprintf(s, user, key, value)
}

func sessionParameterOnAccount(key, value string) string {
	s := `
resource "snowflake_session_parameter" "p" {
	key = "%s"
	value = "%s"
	on_account = true
}
`
	return fmt.Sprintf(s, key, value)
}
