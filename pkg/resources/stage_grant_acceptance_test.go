package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_StageGrant_defaults(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: stageGrantConfig(name, normal, "READ", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "privilege", "READ"),
					testRolesAndShares(t, "snowflake_stage_grant.g", []string{name}),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: stageGrantConfig(name, normal, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_role.r", "name", name),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "privilege", "ALL PRIVILEGES"),
					testRolesAndShares(t, "snowflake_stage_grant.g", []string{name}),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_stage_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func stageGrantConfig(name string, grantType grantType, privilege string, databaseName string, schemaName string) string {
	var stageNameConfig string
	switch grantType {
	case normal:
		stageNameConfig = "stage_name = snowflake_stage.s.name"
	case onFuture:
		stageNameConfig = "on_future = true"
	case onAll:
		stageNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
	resource snowflake_stage s {
		name = "%s"
		database = "%s"
		schema = "%s"
		comment = "Terraform acceptance test"
	}

	resource snowflake_role r {
		name = "%s"
	}

	resource snowflake_stage_grant g {
		database_name = "%s"
		schema_name = "%s"
		%s

		privilege = "%s"

		roles = [
			snowflake_role.r.name
		]
	}
`, name, databaseName, schemaName, name, databaseName, schemaName, stageNameConfig, privilege)
}

func TestAcc_StageFutureGrant(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: stageGrantConfig(name, onFuture, "READ", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_stage_grant.g", "stage_name"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "privilege", "READ"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_stage_grant.g",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
					"on_all",                 // not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_StageGrantOnAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: stageGrantConfig(name, onAll, "READ", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_stage_grant.g", "stage_name"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_stage_grant.g", "privilege", "READ"),
				),
			},
		},
	})
}
