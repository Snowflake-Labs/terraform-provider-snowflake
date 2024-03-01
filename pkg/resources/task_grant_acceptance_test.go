package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TaskGrant(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(name, 8, normal, "OPERATE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "task_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// UPDATE MAX_CONCURRENCY_LEVEL
			{
				Config: taskGrantConfig(name, 10, normal, "OPERATE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "task_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "10"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// UPDATE PRIVILEGE
			{
				Config: taskGrantConfig(name, 10, normal, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "task_name", name),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_TaskGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(name, 8, onAll, "OPERATE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_task_grant.test", "task_name"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_TaskGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: taskGrantConfig(name, 8, onFuture, "OPERATE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_task_grant.test", "task_name"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OPERATE"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "max_concurrency_level", "8"),
					resource.TestCheckResourceAttr("snowflake_warehouse.test", "statement_timeout_in_seconds", "86400"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func taskGrantConfig(name string, concurrency int32, grantType grantType, privilege string, databaseName string, schemaName string) string {
	var taskNameConfig string
	switch grantType {
	case normal:
		taskNameConfig = "task_name \t= snowflake_task.test.name"
	case onFuture:
		taskNameConfig = "on_future = true"
	case onAll:
		taskNameConfig = "on_all = true"
	}

	s := `
resource "snowflake_role" "test" {
  	name = "%v"
}

resource "snowflake_warehouse" "test" {
	name                         = "%s"
	max_concurrency_level        = %d
	statement_timeout_in_seconds = 86400
	query_acceleration_max_scale_factor = 0
}

resource "snowflake_task" "test" {
	name     	    = "%s"
	database  		= "%s"
	schema   		= "%s"
	warehouse 		= snowflake_warehouse.test.name
	sql_statement = "SHOW FUNCTIONS"
	enabled  	  	= true
	schedule 	  	= "15 MINUTES"
	lifecycle {
		ignore_changes = [session_parameters]
	}
}

resource "snowflake_task_grant" "test" {
	%s
	database_name = "%s"
	roles         = [snowflake_role.test.name]
	schema_name   = "%s"
	privilege 	= "%s"
}
`
	return fmt.Sprintf(s, name, name, concurrency, name, databaseName, schemaName, taskNameConfig, databaseName, schemaName, privilege)
}

func TestAcc_TaskOwnershipGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	new_name := name + "_NEW"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// CREATE SCHEMA level FUTURE ownership grant to role <name>
			{
				Config: taskOwnershipGrantConfig(name, onFuture, "OWNERSHIP", name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OWNERSHIP"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "roles.0", name),
				),
			},
			// UPDATE SCHEMA level FUTURE OWNERSHIP grant to role <new_name>
			{
				Config: taskOwnershipGrantConfig(name, onFuture, "OWNERSHIP", new_name, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "privilege", "OWNERSHIP"),
					resource.TestCheckResourceAttr("snowflake_task_grant.test", "roles.0", new_name),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_task_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func taskOwnershipGrantConfig(name string, grantType grantType, privilege string, rolename string, databaseName string, schemaName string) string {
	var taskNameConfig string
	switch grantType {
	case normal:
		taskNameConfig = "task_name \t= snowflake_task.test.name"
	case onFuture:
		taskNameConfig = "on_future = true"
	case onAll:
		taskNameConfig = "on_all = true"
	}

	s := `
resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_role" "test_new" {
	name = "%v_NEW"
  }

resource "snowflake_task_grant" "test" {
  %s
  roles             = [ "%s" ]
  database_name 	= "%s"
  schema_name       = "%s"
  privilege 	    = "%s"
  with_grant_option = false
}
`
	return fmt.Sprintf(s, name, name, taskNameConfig, rolename, databaseName, schemaName, privilege)
}
