package resources_test

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/chanzuckerberg/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
)

func checkBool(path, attr string, value bool) func(*terraform.State) error {
	return func(state *terraform.State) error {
		is := state.RootModule().Resources[path].Primary
		d := is.Attributes[attr]
		b, err := strconv.ParseBool(d)
		if err != nil {
			return err
		}
		if b != value {
			return fmt.Errorf("at %s expected %t but got %t", path, value, b)
		}
		return nil
	}
}

func TestAccUser(t *testing.T) {
	r := require.New(t)
	prefix := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix2 := randomdata.Email()
	sshkey1, err := testhelpers.Fixture("userkey1")
	r.NoError(err)
	sshkey2, err := testhelpers.Fixture("userkey2")
	r.NoError(err)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: uConfig(prefix, sshkey1, sshkey2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix))),
					checkBool("snowflake_user.w", "disabled", false),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "FOO"),
					checkBool("snowflake_user.w", "has_rsa_public_key", true),
					checkBool("snowflake_user.w", "must_change_password", true),
				),
			},
			// RENAME
			{
				Config: uConfig(prefix2, sshkey1, sshkey2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix2))),
					checkBool("snowflake_user.w", "disabled", false),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "FOO"),
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
					checkBool("snowflake_user.w", "disabled", true),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "BAR"),
					checkBool("snowflake_user.w", "has_rsa_public_key", false),
				),
			},
			// IMPORT
			{
				ResourceName:            "snowflake_user.w",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "rsa_public_key", "rsa_public_key_2", "must_change_password"},
			},
		},
	})
}

func uConfig(prefix, key1, key2 string) string {
	s := `
resource "snowflake_user" "w" {
	name = "%s"
	comment = "test comment"
	login_name = "%s_login"
	disabled = false
	default_warehouse="foo"
	default_role="foo"
	default_namespace="foo"
	rsa_public_key = <<KEY
%s
KEY
	rsa_public_key_2 = <<KEY
%s
KEY
	must_change_password = true
}
`
	s = fmt.Sprintf(s, prefix, prefix, key1, key2)
	log.Printf("[DEBUG] s %s", s)
	return s
}

func uConfig2(prefix string) string {
	s := `
resource "snowflake_user" "w" {
	name = "%s"
	comment = "test comment 2"
	password = "best password"
	login_name = "%s_login"
	disabled = true
	default_warehouse="bar"
	default_role="bar"
	default_namespace="bar"
}
`
	log.Printf("[DEBUG] s2 %s", s)
	return fmt.Sprintf(s, prefix, prefix)
}
