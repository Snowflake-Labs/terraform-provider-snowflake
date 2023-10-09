package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGrantPrivilegesToRole_onAccount(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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

/*
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
*/
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
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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

	resource "snowflake_grant_privileges_to_role" "g" {
		privileges = [%v]
		role_name  = snowflake_role.r.name
		on_account_object {
			object_type = "DATABASE"
			object_name = "terraform_test_database"
		}
	}
	`, name, privilegesString)
}

func grantPrivilegesToRole_onAccountObjectConfigAllPrivileges(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		all_privileges = true
		role_name  = snowflake_role.r.name
		on_account_object {
			object_type = "DATABASE"
			object_name = "terraform_test_database"
		}
	}
	`, name)
}

func TestAccGrantPrivilegesToRole_onSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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

	resource "snowflake_grant_privileges_to_role" "g" {
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema {
		  schema_name = "\"terraform_test_database\".\"terraform_test_schema\""
		}
	}
	`, name, privilegesString)
}

func grantPrivilegesToRole_onSchemaConfigAllPrivileges(name string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		all_privileges = true
		on_schema {
			schema_name = "\"terraform_test_database\".\"terraform_test_schema\""
		}
	}
	`, name)
}

func TestAccGrantPrivilegesToRole_onSchemaConfigAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
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
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name, []string{"MONITOR", "USAGE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.all_schemas_in_database", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "USAGE"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name, []string{"MONITOR"}),
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

func TestAccGrantPrivilegesToRole_onSchema_futureSchemasInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchema_futureSchemasInDatabaseConfig(name, []string{"MONITOR", "USAGE"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.future_schemas_in_database", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "USAGE"),
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

func grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name string, privileges []string) string {
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
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema {
			all_schemas_in_database = "terraform_test_database"

		}
	}
	`, name, privilegesString)
}

func grantPrivilegesToRole_onSchema_futureSchemasInDatabaseConfig(name string, privileges []string) string {
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
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema {
			future_schemas_in_database = "terraform_test_database"

		}
	}
	`, name, privilegesString)
}

func TestAccGrantPrivilegesToRole_onSchemaObject_objectType(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_objectType(name, []string{"SELECT", "REFERENCES"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.object_type", "VIEW"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.object_name", fmt.Sprintf(`"%v"."%v"."%v"`, name, name, name)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_objectType(name, []string{"SELECT"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "SELECT"),
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

func grantPrivilegesToRole_onSchemaObject_objectType(name string, privileges []string) string {
	doubleQuotePrivileges := make([]string, len(privileges))
	for i, p := range privileges {
		doubleQuotePrivileges[i] = fmt.Sprintf(`"%v"`, p)
	}
	privilegesString := strings.Join(doubleQuotePrivileges, ",")
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_view" "v" {
		name        = "%v"
		database    = snowflake_database.d.name
		schema      = snowflake_schema.s.name
		is_secure   = true
		statement   = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		depends_on = [ snowflake_view.v]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema_object {
			object_type = "VIEW"
			object_name = "\"terraform_test_database\".\"terraform_test_schema\".\"%s\""
		}
	}
	`, name, name, privilegesString, name)
}

func TestAccGrantPrivilegesToRole_onSchemaObject_allInSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInSchema(name, []string{"SELECT", "REFERENCES"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.in_schema", fmt.Sprintf(`"%v"."%v"`, name, name)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInSchema(name, []string{"SELECT"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "SELECT"),
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

func grantPrivilegesToRole_onSchemaObject_allInSchema(name string, privileges []string) string {
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
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema_object {
			all {
				object_type_plural = "TABLES"
				in_schema = "\"terraform_test_database\".\"terraform_test_schema\""
			}
		}
	}
	`, name, privilegesString)
}

func TestAccGrantPrivilegesToRole_onSchemaObject_allInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInDatabase(name, []string{"SELECT", "REFERENCES"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.in_database", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInDatabase(name, []string{"SELECT"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "SELECT"),
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

func grantPrivilegesToRole_onSchemaObject_allInDatabase(name string, privileges []string) string {
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
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema_object {
			all {
				object_type_plural = "TABLES"
				in_database = "terraform_test_database"
			}
		}
	}
	`, name, privilegesString)
}

func TestAccGrantPrivilegesToRole_onSchemaObject_futureInSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInSchema(name, []string{"SELECT", "REFERENCES"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.in_schema", fmt.Sprintf(`"%v"."%v"`, name, name)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInSchema(name, []string{"SELECT"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "SELECT"),
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

func grantPrivilegesToRole_onSchemaObject_futureInSchema(name string, privileges []string) string {
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
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema_object {
			future {
				object_type_plural = "TABLES"
				in_schema = "\"terraform_test_database\".\"terraform_test_schema\""
			}
		}
	}
	`, name, privilegesString)
}

func TestAccGrantPrivilegesToRole_onSchemaObject_futureInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	objectType := "TABLES"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT", "REFERENCES"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.in_database", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "SELECT"),
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

func grantPrivilegesToRole_onSchemaObject_futureInDatabase(name string, objectType string, privileges []string) string {
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
		depends_on = [ snowflake_schema.s ]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema_object {
			future {
				object_type_plural = "%s"
				in_database = "terraform_test_database"
			}
		}
	}
	`, name, privilegesString, objectType)
}

func TestAccGrantPrivilegesToRole_multipleResources(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_multipleResources(name, []string{"CREATE ACCOUNT", "CREATE ROLE"}, []string{"IMPORT SHARE", "MANAGE GRANTS"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g1", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g1", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g1", "privileges.0", "CREATE ACCOUNT"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g1", "privileges.1", "CREATE ROLE"),

					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g2", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g2", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g2", "privileges.0", "IMPORT SHARE"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g2", "privileges.1", "MANAGE GRANTS"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_grant_privileges_to_role.g1",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				ResourceName:      "snowflake_grant_privileges_to_role.g2",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func grantPrivilegesToRole_multipleResources(name string, privileges1, privileges2 []string) string {
	doubleQuotePrivileges1 := make([]string, len(privileges1))
	for i, p := range privileges1 {
		doubleQuotePrivileges1[i] = fmt.Sprintf(`"%v"`, p)
	}
	privilegesString1 := strings.Join(doubleQuotePrivileges1, ",")

	doubleQuotePrivileges2 := make([]string, len(privileges2))
	for i, p := range privileges2 {
		doubleQuotePrivileges2[i] = fmt.Sprintf(`"%v"`, p)
	}
	privilegesString2 := strings.Join(doubleQuotePrivileges2, ",")

	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g1" {
		role_name  = snowflake_role.r.name
		privileges = [%s]
		on_account = true
	}

	resource "snowflake_grant_privileges_to_role" "g2" {
		role_name  = snowflake_role.r.name
		privileges = [%s]
		on_account = true
	}
	`, name, privilegesString1, privilegesString2)
}

func TestAccGrantPrivilegesToRole_onSchemaObject_futureInDatabase_externalTable(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	objectType := "EXTERNAL TABLES"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT", "REFERENCES"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.object_type_plural", "EXTERNAL TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.in_database", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT"}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "SELECT"),
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
