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

func TestAcc_Stages(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	stageName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: stages(databaseName, schemaName, stageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_stages.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_stages.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_stages.t", "stages.#"),
					resource.TestCheckResourceAttr("data.snowflake_stages.t", "stages.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_stages.t", "stages.0.name", stageName),
				),
			},
		},
	})
}

func stages(databaseName string, schemaName string, stageName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	resource snowflake_stage "t"{
		name 	 = "%v"
		database = snowflake_schema.s.database
		schema 	 = snowflake_schema.s.name
	}

	data snowflake_stages "t" {
		database = snowflake_stage.t.database
		schema = snowflake_stage.t.schema
		depends_on = [snowflake_stage.t]
	}
	`, databaseName, schemaName, stageName)
}
