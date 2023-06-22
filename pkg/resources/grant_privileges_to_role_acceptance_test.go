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
				Config: grantPrivilegesToRole_onAccountObjectConfig(name, []string{"CREATE DATABASE ROLE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "CREATE DATABASE ROLE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onAccountObjectConfig(name, []string{"MONITOR", "CREATE SCHEMA"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "CREATE SCHEMA"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "MONITOR"),
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

func TestAccGrantPrivilegesToRole_onAccountObjectAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountObjectConfigAllPrivileges(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_type", "DATABASE"),
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

func grantPrivilegesToRole_onAccountObjectConfig(name string, privileges []string) string {
	doubleQuotePrivileges := make([]string, len(privileges))
	for i, p := range privileges {
		doubleQuotePrivileges[i] = fmt.Sprintf(`"%v"`, p)
	}
	privilegesString := strings.Join(doubleQuotePrivileges, ",")
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_database" "d" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		privileges = [%v]
		role_name  = snowflake_role.r.name
		on_account_object {
			object_type = "DATABASE"
			object_name = snowflake_database.d.name
		}
	}
	`, name, name, privilegesString)
}

func grantPrivilegesToRole_onAccountObjectConfigAllPrivileges(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_database" "d" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		all_privileges = true
		role_name  = snowflake_role.r.name
		on_account_object {
			object_type = "DATABASE"
			object_name = snowflake_database.d.name
		}
	}
	`, name, name)
}

func TestAccGrantPrivilegesToRole_onSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaConfig(name, []string{"MONITOR", "USAGE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.schema_name", fmt.Sprintf("\"%v\".\"%v\"", name, name)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "USAGE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaConfig(name, []string{"MONITOR"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
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

func TestAccGrantPrivilegesToRole_onSchemaConfigAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaConfigAllPrivileges(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.schema_name", fmt.Sprintf("\"%v\".\"%v\"", name, name)),
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

func TestAccGrantPrivilegesToRole_onSchema_allSchemasInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name, []string{"MONITOR", "USAGE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.scma_name", fmt.Sprintf("\"%v\".\"%v\"", name, name)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "USAGE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaConfig(name, []string{"MONITOR"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
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

func grantPrivilegesToRole_onSchemaConfig(name string, privileges []string) string {
	doubleQuotePrivileges := make([]string, len(privileges))
	for i, p := range privileges {
		doubleQuotePrivileges[i] = fmt.Sprintf(`"%v"`, p)
	}
	privilegesString := strings.Join(doubleQuotePrivileges, ",")
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_database" "d" {
		name = "%v"
	}

	resource "snowflake_schema" "s" {
		name = "%v"
		database = snowflake_database.d.name
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema {
		  schema_name = "\"%s\".\"%s\""
		}
	}
	`, name, name, name, privilegesString, name, name)
}

func grantPrivilegesToRole_onSchemaConfigAllPrivileges(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_database" "d" {
		name = "%v"
	}

	resource "snowflake_schema" "s" {
		name = "%v"
		database = snowflake_database.d.name
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		all_privileges = true
		on_schema {
		  schema_name = "\"%s\".\"%s\""
		}
	}
	`, name, name, name, name, name)
}
