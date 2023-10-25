// Copyright (c) Snowflake, Inc.
// SPDX-License-Identifier: MIT

package resources_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Pipe(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_PIPE_TESTS"); ok {
		t.Skip("Skipping TestAccPipe")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: pipeConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_pipe.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "auto_ingest", "false"),
					resource.TestCheckResourceAttr("snowflake_pipe.test", "notification_channel", ""),
				),
			},
		},
	})
}

func pipeConfig(name string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_table" "test" {
	database = "%s"
  	schema   = "%s"
	name     = "%s"

	  column {
			name = "id"
			type = "NUMBER(5,0)"
	  }

	  column {
		name = "data"
		type = "VARCHAR(16)"
	  }
}

resource "snowflake_stage" "test" {
	name = "%s"
	database = "%s"
	schema = "%s"
	comment = "Terraform acceptance test"
}


resource "snowflake_pipe" "test" {
  database       = "%s"
  schema         = "%s"
  name           = "%s"
  comment        = "Terraform acceptance test"
  copy_statement = <<CMD
COPY INTO "${snowflake_table.test.database}"."${snowflake_table.test.schema}"."${snowflake_table.test.name}"
  FROM @"${snowflake_stage.test.database}"."${snowflake_stage.test.schema}"."${snowflake_stage.test.name}"
  FILE_FORMAT = (TYPE = CSV)
CMD
  auto_ingest    = false
}
`
	return fmt.Sprintf(s, databaseName, schemaName, name, name, databaseName, schemaName, databaseName, schemaName, name)
}
