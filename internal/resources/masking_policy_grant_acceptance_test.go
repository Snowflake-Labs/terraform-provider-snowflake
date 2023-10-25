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

func TestAcc_MaskingPolicyGrant(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: maskingPolicyGrantConfig(accName, "APPLY", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "masking_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "privilege", "APPLY"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: maskingPolicyGrantConfig(accName, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "masking_policy_name", accName),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_masking_policy_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_masking_policy_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func maskingPolicyGrantConfig(name string, privilege string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
	resource "snowflake_role" "test" {
		name = "%v"
	}

	resource "snowflake_masking_policy" "test" {
		name = "%v"
		database = "%s"
		schema = "%s"
		signature {
			column {
				name = "val"
				type = "VARCHAR"
			}
		}
		masking_expression = "case when current_role() in ('ANALYST') then val else sha2(val, 512) end"
		return_data_type = "VARCHAR"
		comment = "Terraform acceptance test"
	}

	resource "snowflake_masking_policy_grant" "test" {
		masking_policy_name = snowflake_masking_policy.test.name
		database_name = "%s"
		roles         = [snowflake_role.test.name]
		schema_name   = "%s"
		privilege = "%s"
	}
	`, name, name, databaseName, schemaName, databaseName, schemaName, privilege)
}
