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

func TestAcc_RowAccessPolicy(t *testing.T) {
	if _, ok := os.LookupEnv("SKIP_ROW_ACCESS_POLICY_TESTS"); ok {
		t.Skip("Skipping TestAccRowAccessPolicy")
	}

	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: rowAccessPolicyConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "row_access_expression", "case when current_role() in ('ANALYST') then true else false end"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.N", "VARCHAR"),
					resource.TestCheckResourceAttr("snowflake_row_access_policy.test", "signature.V", "VARCHAR"),
				),
			},
		},
	})
}

func rowAccessPolicyConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_row_access_policy" "test" {
	name = "%v"
	database = "%s"
	schema = "%s"
	signature = {
		N = "VARCHAR"
		V = "VARCHAR",
	}
	row_access_expression = "case when current_role() in ('ANALYST') then true else false end"
	comment = "Terraform acceptance test"
}
`, n, databaseName, schemaName)
}
