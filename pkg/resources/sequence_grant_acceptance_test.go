package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_SequenceGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sequenceGrantConfig(name, onFuture, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_sequence_grant.test", "sequence_name"),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "privilege", "USAGE"),
				),
			},
			{
				ResourceName:      "snowflake_sequence_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_SequenceGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sequenceGrantConfig(name, onAll, "USAGE", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_sequence_grant.test", "sequence_name"),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_sequence_grant.test", "privilege", "USAGE"),
				),
			},
			{
				ResourceName:      "snowflake_sequence_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func sequenceGrantConfig(name string, grantType grantType, privilege string, databaseName string, schemaName string) string {
	var sequenceNameConfig string
	switch grantType {
	case onFuture:
		sequenceNameConfig = "on_future = true"
	case onAll:
		sequenceNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource "snowflake_role" "test" {
  name = "%s"
}

resource "snowflake_sequence_grant" "test" {
    database_name = "%s"
	roles         = [snowflake_role.test.name]
	schema_name   = "%s"
	%s
	privilege = "%s"
}
`, name, databaseName, schemaName, sequenceNameConfig, privilege)
}
