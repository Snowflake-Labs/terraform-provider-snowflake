package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_GrantPrivilegesToRole_onAccount(t *testing.T) {
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
	func TestAcc_GrantPrivilegesToRole_onAccountAllPrivileges(t *testing.T) {
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

func TestAcc_GrantPrivilegesToRole_onAccountObject(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountObjectConfig(name, []string{"CREATE DATABASE ROLE"}, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_type", "DATABASE"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "CREATE DATABASE ROLE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onAccountObjectConfig(name, []string{"MONITOR", "CREATE SCHEMA"}, acc.TestDatabaseName),
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

func TestAcc_GrantPrivilegesToRole_onAccountObjectAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onAccountObjectConfigAllPrivileges(name, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_account_object.0.object_name", acc.TestDatabaseName),
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

func grantPrivilegesToRole_onAccountObjectConfig(name string, privileges []string, databaseName string) string {
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
			object_name = "%s"
		}
	}
	`, name, privilegesString, databaseName)
}

func grantPrivilegesToRole_onAccountObjectConfigAllPrivileges(name string, databaseName string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		all_privileges = true
		role_name  = snowflake_role.r.name
		on_account_object {
			object_type = "DATABASE"
			object_name = "%s"
		}
	}
	`, name, databaseName)
}

func TestAcc_GrantPrivilegesToRole_onSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaConfig(name, []string{"MONITOR", "USAGE"}, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.schema_name", fmt.Sprintf("\"%v\".\"%v\"", acc.TestDatabaseName, acc.TestSchemaName)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "USAGE"),
				),
			},
			// ADD PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaConfig(name, []string{"MONITOR"}, acc.TestDatabaseName, acc.TestSchemaName),
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

func grantPrivilegesToRole_onSchemaConfig(name string, privileges []string, databaseName string, schemaName string) string {
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
		  schema_name = "\"%s\".\"%s\""
		}
	}
	`, name, privilegesString, databaseName, schemaName)
}

func grantPrivilegesToRole_onSchemaConfigAllPrivileges(name string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "r" {
		name = "%v"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		role_name = snowflake_role.r.name
		all_privileges = true
		on_schema {
			schema_name = "\"%s\".\"%s\""
		}
	}
	`, name, databaseName, schemaName)
}

func TestAcc_GrantPrivilegesToRole_onSchemaConfigAllPrivileges(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaConfigAllPrivileges(name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.schema_name", fmt.Sprintf("\"%v\".\"%v\"", acc.TestDatabaseName, acc.TestSchemaName)),
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

func TestAcc_GrantPrivilegesToRole_onSchema_allSchemasInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name, []string{"MONITOR", "USAGE"}, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.all_schemas_in_database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "MONITOR"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "USAGE"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name, []string{"MONITOR"}, acc.TestDatabaseName),
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

func TestAcc_GrantPrivilegesToRole_onSchema_futureSchemasInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchema_futureSchemasInDatabaseConfig(name, []string{"MONITOR", "USAGE"}, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema.0.future_schemas_in_database", acc.TestDatabaseName),
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

func grantPrivilegesToRole_onSchema_allSchemasInDatabaseConfig(name string, privileges []string, databaseName string) string {
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
			all_schemas_in_database = "%s"

		}
	}
	`, name, privilegesString, databaseName)
}

func grantPrivilegesToRole_onSchema_futureSchemasInDatabaseConfig(name string, privileges []string, databaseName string) string {
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
			future_schemas_in_database = "%s"

		}
	}
	`, name, privilegesString, databaseName)
}

func TestAcc_GrantPrivilegesToRole_onSchemaObject_objectType(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_objectType(name, []string{"SELECT", "REFERENCES"}, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.object_type", "VIEW"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.object_name", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_objectType(name, []string{"SELECT"}, acc.TestDatabaseName, acc.TestSchemaName),
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

func grantPrivilegesToRole_onSchemaObject_objectType(name string, privileges []string, databaseName string, schemaName string) string {
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
		database    = "%s"
		schema      = "%s"
		is_secure   = true
		statement   = "SELECT ROLE_NAME, ROLE_OWNER FROM INFORMATION_SCHEMA.APPLICABLE_ROLES"
	}

	resource "snowflake_grant_privileges_to_role" "g" {
		depends_on = [ snowflake_view.v]
		role_name = snowflake_role.r.name
		privileges = [%s]
		on_schema_object {
			object_type = "VIEW"
			object_name = "\"%s\".\"%s\".\"%s\""
		}
	}
	`, name, name, databaseName, schemaName, privilegesString, databaseName, schemaName, name)
}

func TestAcc_GrantPrivilegesToRole_onSchemaObject_allInSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInSchema(name, []string{"SELECT", "REFERENCES"}, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.in_schema", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInSchema(name, []string{"SELECT"}, acc.TestDatabaseName, acc.TestSchemaName),
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

func grantPrivilegesToRole_onSchemaObject_allInSchema(name string, privileges []string, databaseName string, schemaName string) string {
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
		on_schema_object {
			all {
				object_type_plural = "TABLES"
				in_schema = "\"%s\".\"%s\""
			}
		}
	}
	`, name, privilegesString, databaseName, schemaName)
}

func TestAcc_GrantPrivilegesToRole_onSchemaObject_allInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInDatabase(name, []string{"SELECT", "REFERENCES"}, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.all.0.in_database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_allInDatabase(name, []string{"SELECT"}, acc.TestDatabaseName),
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

func grantPrivilegesToRole_onSchemaObject_allInDatabase(name string, privileges []string, databaseName string) string {
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
		on_schema_object {
			all {
				object_type_plural = "TABLES"
				in_database = "%s"
			}
		}
	}
	`, name, privilegesString, databaseName)
}

func TestAcc_GrantPrivilegesToRole_onSchemaObject_futureInSchema(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInSchema(name, []string{"SELECT", "REFERENCES"}, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.in_schema", fmt.Sprintf(`"%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName)),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInSchema(name, []string{"SELECT"}, acc.TestDatabaseName, acc.TestSchemaName),
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

func grantPrivilegesToRole_onSchemaObject_futureInSchema(name string, privileges []string, databaseName string, schemaName string) string {
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
		on_schema_object {
			future {
				object_type_plural = "TABLES"
				in_schema = "\"%s\".\"%s\""
			}
		}
	}
	`, name, privilegesString, databaseName, schemaName)
}

func TestAcc_GrantPrivilegesToRole_onSchemaObject_futureInDatabase(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	objectType := "TABLES"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT", "REFERENCES"}, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.object_type_plural", "TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.in_database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT"}, acc.TestDatabaseName),
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

func grantPrivilegesToRole_onSchemaObject_futureInDatabase(name string, objectType string, privileges []string, databaseName string) string {
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
		on_schema_object {
			future {
				object_type_plural = "%s"
				in_database = "%s"
			}
		}
	}
	`, name, privilegesString, objectType, databaseName)
}

func TestAcc_GrantPrivilegesToRole_multipleResources(t *testing.T) {
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

func TestAcc_GrantPrivilegesToRole_onSchemaObject_futureInDatabase_externalTable(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	objectType := "EXTERNAL TABLES"
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT", "REFERENCES"}, acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "role_name", name),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.object_type_plural", "EXTERNAL TABLES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "on_schema_object.0.future.0.in_database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.#", "2"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.0", "REFERENCES"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.g", "privileges.1", "SELECT"),
				),
			},
			// REMOVE PRIVILEGE
			{
				Config: grantPrivilegesToRole_onSchemaObject_futureInDatabase(name, objectType, []string{"SELECT"}, acc.TestDatabaseName),
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

func TestAcc_GrantPrivilegesToRole_onSchemaObject_futureIcebergTables(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "snowflake_role" "role" {
  name = "TEST_ROLE_123"
}

resource "snowflake_grant_privileges_to_role" "grant" {
  role_name  = snowflake_role.role.name
  privileges = ["SELECT"]
  on_schema_object {
    future {
      object_type_plural = "ICEBERG TABLES"
      in_schema          = "\"%s\".\"%s\""
    }
  }
}
`, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.grant", "on_schema_object.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.grant", "on_schema_object.0.future.#", "1"),
					resource.TestCheckResourceAttr("snowflake_grant_privileges_to_role.grant", "on_schema_object.0.future.0.object_type_plural", string(sdk.PluralObjectTypeIcebergTables)),
				),
			},
		},
	})
}

func TestAcc_GrantPrivilegesToRole_ValidatedIdentifiers(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "snowflake_role" "role" {
  name = "TEST_ROLE_123"
}

resource "snowflake_grant_privileges_to_role" "test_invalidation" {
  role_name  = snowflake_role.role.name
  privileges = ["SELECT"]
  on_schema_object {
    future {
      object_type_plural = "ICEBERG TABLES"
      in_schema          = "%s"
    }
  }
}`, acc.TestSchemaName),
				ExpectError: regexp.MustCompile(".*Expected DatabaseObjectIdentifier identifier type.*"),
			},
		},
	})
}
