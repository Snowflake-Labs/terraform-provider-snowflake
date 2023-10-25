// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TableGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name, onAll, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "SELECT"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.0", name),

					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}),
				),
			},
			{
				ResourceName:      "snowflake_table_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_TableGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name, onFuture, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "SELECT"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "roles.0", name),

					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_table_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_TableGrant_defaults(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableGrantConfig(name, normal, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "table_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "SELECT"),
					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: tableGrantConfig(name, normal, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "table_name", name),
					resource.TestCheckResourceAttr("snowflake_table_grant.g", "privilege", "ALL PRIVILEGES"),
					testRolesAndShares(t, "snowflake_table_grant.g", []string{name}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_table_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func tableGrantConfig(name string, grantType grantType, privilege string, databaseName string, schemaName string) string {
	var tableNameConfig string
	switch grantType {
	case normal:
		tableNameConfig = "table_name = snowflake_table.t.name"
	case onFuture:
		tableNameConfig = "on_future = true"
	case onAll:
		tableNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource snowflake_role r {
  name = "%s"
}

resource snowflake_table t {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "id"
		type = "NUMBER(38,0)"
	}
}

resource snowflake_table_grant g {
	database_name = "%s"
	schema_name   = "%s"
	%s
	privilege = "%s"
	roles = [
		snowflake_role.r.name
	]
}

`, name, name, databaseName, schemaName, databaseName, schemaName, tableNameConfig, privilege)
}
