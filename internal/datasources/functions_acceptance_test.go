// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Functions(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	functionName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: functions(databaseName, schemaName, functionName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_functions.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_functions.t", "schema", schemaName),
				),
			},
		},
	})
}

func functions(databaseName string, schemaName string, functionName string) string {
	s := `
resource "snowflake_database" "test_database" {
	name 	  = "%v"
	comment = "Terraform acceptance test"
}
resource "snowflake_schema" "test_schema" {
	name 	   = "%v"
	database = snowflake_database.test_database.name
	comment  = "Terraform acceptance test"
}
resource "snowflake_function" "test_funct_simple" {
	name = "%s"
	database = snowflake_database.test_database.name
	schema   = snowflake_schema.test_schema.name
	return_type = "float"
	statement = "3.141592654::FLOAT"
}

data snowflake_functions "t" {
	database = snowflake_database.test_database.name
	schema = snowflake_schema.test_schema.name
	depends_on = [snowflake_function.test_funct_simple]
}
`
	return fmt.Sprintf(s, databaseName, schemaName, functionName)
}
