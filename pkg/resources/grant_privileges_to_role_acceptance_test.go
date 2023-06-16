package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGrantPrivilegesToRole_onAccount(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountConfig(name, []string{"MONITOR USAGE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account", "true"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR USAGE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onAccountConfig(name, []string{"MONITOR USAGE", "MANAGE GRANTS"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account", "true"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MANAGE GRANTS"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "MONITOR USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_grant_privileges_to_role.g",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGrantPrivilegesToRole_onAccountAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountConfigAllPrivileges(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account", "true"),
					resource.TestCheckNoResourceAttr("snowflake_grant_privileges_to_role.g", "privileges"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "all_privileges", "true"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_grant_privileges_to_role.g",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func grantPrivilegesToRole_onAccountConfig(name string, privileges []string) string {
	doubleQuotePrivileges := make([]string, len(privileges))
	for i, p := range privileges {
		doubleQuotePrivileges[i] = fmt.Sprintf(`"%v"`, p)
	}
	privilegesString := strings.Join(doubleQuotePrivileges, ",")
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		privileges = [%v]
		role_name  = snowflake_role.r.name
		on_account = true
	  }
	`, name, privilegesString)
}

func grantPrivilegesToRole_onAccountConfigAllPrivileges(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		all_privileges = true
		role_name  = snowflake_role.r.name
		on_account = true
	  }
	`, name)
}

func TestAccGrantPrivilegesToRole_onAccountObject(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountConfig(name, []string{"MONITOR USAGE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account", "true"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR USAGE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onAccountConfig(name, []string{"MONITOR USAGE", "MANAGE GRANTS"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account", "true"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MANAGE GRANTS"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "MONITOR USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_grant_privileges_to_role.g",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGrantPrivilegesToRole_onAccountAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountConfigAllPrivileges(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account", "true"),
					resource.TestCheckNoResourceAttr("snowflake_grant_privileges_to_role.g", "privileges"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "all_privileges", "true"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_grant_privileges_to_role.g",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
