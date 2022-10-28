package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTableConstraint_fk(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tableConstraintFKConfig(name),

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

func tableConstraintFKConfig(n string) string {
	return fmt.Sprintf(`

resource "snowflake_database" "d" {
	name = "%s"
}

resource "snowflake_schema" "s" {
	name = "%s"
	database = snowflake_database.d.name
}

resource "snowflake_table" "t" {
	database = snowflake_database.d.name
	schema   = snowflake_schema.s.name
	name     = "%s"

	column {
		name = "col1"
		type = "NUMBER(38,0)"
	}
}

resource "snowflake_table" "fk_t" {
	database = snowflake_database.d.name
	schema   = snowflake_schema.s.name
	name     = "fk_%s"

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

`, n, n, n, n, n)
}
