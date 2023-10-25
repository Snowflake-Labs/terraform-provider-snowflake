package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TableConstraint_fk(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintFKConfig(name, acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "type", "FOREIGN KEY"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "enforced", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "deferrable", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.fk", "comment", "hello fk"),
				),
			},
		},
	})
}

func tableConstraintFKConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table" "fk_t" {
	name     = "fk_%s"
	database = "%s"
	schema   = "%s"
	column {
		name     = "fk_col1"
		type     = "text"
		nullable = false
	  }
}

resource "snowflake_table_constraint" "fk" {
	name="%s"
	type= "FOREIGN KEY"
	table_id = snowflake_table.t.id
	columns = ["col1"]
	foreign_key_properties {
	  references {
		table_id = snowflake_table.fk_t.id
		columns = ["fk_col1"]
	  }
	}
	enforced = false
	deferrable = false
	initially = "IMMEDIATE"
	comment = "hello fk"
}

`, n, databaseName, schemaName, n, databaseName, schemaName, n)
}

func TestAcc_TableConstraint_unique(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintUniqueConfig(name, acc.TestDatabaseName, acc.TestSchemaName),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "type", "UNIQUE"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "enforced", "true"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "deferrable", "false"),
					resource.TestCheckResourceAttr("snowflake_table_constraint.unique", "comment", "hello unique"),
				),
			},
		},
	})
}

func tableConstraintUniqueConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_table" "t" {
	name     = "%s"
	database = "%s"
	schema   = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table" "unique_t" {
	name     = "unique_%s"
	database = "%s"
	schema   = "%s"

	column {
		name     = "unique_col1"
		type     = "text"
		nullable = false
	  }
}

resource "snowflake_table_constraint" "unique" {
	name="%s"
	type= "UNIQUE"
	table_id = snowflake_table.t.id
	columns = ["col1"]
	enforced = true
	deferrable = false
	initially = "IMMEDIATE"
	comment = "hello unique"
}

`, n, databaseName, schemaName, n, databaseName, schemaName, n)
}
