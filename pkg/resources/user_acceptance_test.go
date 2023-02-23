package resources_test

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/Pallinder/go-randomdata"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/testhelpers"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"
)

func TestAcc_User(t *testing.T) {
	r := require.New(t)
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	prefix2 := "tst-terraform" + strings.ToUpper(randomdata.Email())
	sshkey1, err := testhelpers.Fixture("userkey1")
	r.NoError(err)
	sshkey2, err := testhelpers.Fixture("userkey2")
	r.NoError(err)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: uConfig(prefix, sshkey1, sshkey2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix))),
					resource.TestCheckResourceAttr("snowflake_user.w", "display_name", "Display Name"),
					resource.TestCheckResourceAttr("snowflake_user.w", "first_name", "Marcin"),
					resource.TestCheckResourceAttr("snowflake_user.w", "last_name", "Zukowski"),
					resource.TestCheckResourceAttr("snowflake_user.w", "email", "fake@email.com"),
					resource.TestCheckResourceAttr("snowflake_user.w", "disabled", "false"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_secondary_roles.0", "ALL"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "FOO"),
					resource.TestCheckResourceAttr("snowflake_user.w", "has_rsa_public_key", "true"),
					resource.TestCheckResourceAttr("snowflake_user.w", "must_change_password", "true"),
				),
			},
			// RENAME
			{
				Config: uConfig(prefix2, sshkey1, sshkey2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_user.w", "name", prefix2),
					resource.TestCheckResourceAttr("snowflake_user.w", "comment", "test comment"),
					resource.TestCheckResourceAttr("snowflake_user.w", "login_name", strings.ToUpper(fmt.Sprintf("%s_login", prefix2))),
					resource.TestCheckResourceAttr("snowflake_user.w", "display_name", "Display Name"),
					resource.TestCheckResourceAttr("snowflake_user.w", "first_name", "Marcin"),
					resource.TestCheckResourceAttr("snowflake_user.w", "last_name", "Zukowski"),
					resource.TestCheckResourceAttr("snowflake_user.w", "email", "fake@email.com"),
					resource.TestCheckResourceAttr("snowflake_user.w", "disabled", "false"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "foo"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_secondary_roles.0", "ALL"),
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
					resource.TestCheckResourceAttr("snowflake_user.w", "display_name", "New Name"),
					resource.TestCheckResourceAttr("snowflake_user.w", "first_name", "Benoit"),
					resource.TestCheckResourceAttr("snowflake_user.w", "last_name", "Dageville"),
					resource.TestCheckResourceAttr("snowflake_user.w", "email", "fake@email.net"),
					resource.TestCheckResourceAttr("snowflake_user.w", "disabled", "true"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_warehouse", "bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_role", "bar"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_secondary_roles.#", "0"),
					resource.TestCheckResourceAttr("snowflake_user.w", "default_namespace", "BAR"),
					resource.TestCheckResourceAttr("snowflake_user.w", "has_rsa_public_key", "false"),
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
	display_name = "Display Name"
	first_name = "Marcin"
	last_name = "Zukowski"
	email = "fake@email.com"
	disabled = false
	default_warehouse="foo"
	default_role="foo"
	default_secondary_roles=["ALL"]
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
	display_name = "New Name"
	first_name = "Benoit"
	last_name = "Dageville"
	email = "fake@email.net"
	disabled = true
	default_warehouse="bar"
	default_role="bar"
	default_secondary_roles=[]
	default_namespace="bar"
}
`
	log.Printf("[DEBUG] s2 %s", s)
	return fmt.Sprintf(s, prefix, prefix)
}
