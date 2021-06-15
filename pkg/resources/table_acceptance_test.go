package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Table(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	table2Name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: tableConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARIANT"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "VARCHAR(16)"),
				),
			},
			{
				Config: tableConfig2(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.name", "column2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.0.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.name", "column3"),
					resource.TestCheckResourceAttr("snowflake_table.test_table", "column.1.type", "FLOAT"),
					resource.TestCheckNoResourceAttr("snowflake_table.test_table", "cluster_by"),
				),
			},
			{
				Config: tableConfig3(accName, table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.0", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.type", "FLOAT"),
				),
			},
			{
				Config: tableConfig4(accName, table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.1", "\"col2\""),
				),
			},
			{
				Config: tableConfig5(accName, table2Name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "name", table2Name),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "database", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.name", "COL1"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.0.type", "VARCHAR(16777216)"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "column.1.name", "col2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.#", "2"),
					resource.TestCheckResourceAttr("snowflake_table.test_table2", "cluster_by.0", "\"col2\""),
				),
			},
		},
	})
}

func tableConfig(name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16)"
	}
}
`
	return fmt.Sprintf(s, name, name, name)
}

func tableConfig2(name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "column3"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, name, name, name)
}

func tableConfig3(name string, table2Name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table2" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	cluster_by = ["COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, name, name, table2Name)
}

func tableConfig4(name string, table2Name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table2" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	cluster_by = ["COL1","\"col2\""]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, name, name, table2Name)
}

func tableConfig5(name string, table2Name string) string {
	s := `
resource "snowflake_database" "test_database" {
	name    = "%s"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test_schema" {
	name     = "%s"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}

resource "snowflake_table" "test_table2" {
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	cluster_by = ["\"col2\"","COL1"]
	column {
		name = "COL1"
		type = "VARCHAR(16777216)"
	}
	column {
		name = "col2"
		type = "FLOAT"
	}
}
`
	return fmt.Sprintf(s, name, name, table2Name)
}
