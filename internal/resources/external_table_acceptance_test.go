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

func TestAcc_ExternalTable(t *testing.T) {
	env := os.Getenv("SKIP_EXTERNAL_TABLE_TEST")
	if env != "" {
		t.Skip("Skipping TestAcc_ExternalTable")
	}
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	bucketURL := os.Getenv("AWS_EXTERNAL_BUCKET_URL")
	if bucketURL == "" {
		t.Skip("Skipping TestAcc_ExternalTable")
	}
	roleName := os.Getenv("AWS_EXTERNAL_ROLE_NAME")
	if roleName == "" {
		t.Skip("Skipping TestAcc_ExternalTable")
	}
	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalTableConfig(accName, bucketURL, roleName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "name", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "database", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "schema", accName),
					resource.TestCheckResourceAttr("snowflake_external_table.test_table", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func externalTableConfig(name string, bucketURL string, roleName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_storage_integration" "i" {
	name = "%v"
	storage_allowed_locations = ["%s"]
	storage_provider = "S3"
	storage_aws_role_arn = "%s"
}

resource "snowflake_stage" "test" {
	name = "%v"
	url = "%s"
	database = "%s"
	schema = "%s"
	storage_integration = snowflake_storage_integration.i.name
}

resource "snowflake_external_table" "test_table" {
	name     = "%s"
	database = "%s"
	schema = "%s"
	comment  = "Terraform acceptance test"
	column {
		name = "column1"
		type = "STRING"
		as   = "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"
	}
	column {
		name = "column2"
		type = "TIMESTAMP_NTZ(9)"
		as   = "($1:\"CreatedDate\"::timestamp)"
	}
  file_format = "TYPE = CSV"
  location = "@\"%s\".\"%s\".\"${snowflake_stage.test.name}\""
}
`
	return fmt.Sprintf(s, name, bucketURL, roleName, name, bucketURL, databaseName, schemaName, name, databaseName, schemaName, databaseName, schemaName)
}
