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

func TestAcc_StreamGrant_basic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfig(name, streamName, normal, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "stream_name", streamName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_all", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "privilege", "SELECT"),
				),
			},
			// UPDATE ALL PRIVILEGES
			{
				Config: streamGrantConfig(name, streamName, normal, "ALL PRIVILEGES", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "stream_name", streamName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "privilege", "ALL PRIVILEGES"),
				),
			},
			{
				ResourceName:      "snowflake_stream_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_StreamGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfig(name, streamName, onAll, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_stream_grant.test", "stream_name"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_all", "true"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "privilege", "SELECT"),
				),
			},
			{
				ResourceName:      "snowflake_stream_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func TestAcc_StreamGrant_onFuture(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfig(name, streamName, onFuture, "SELECT", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", acc.TestSchemaName),
					resource.TestCheckNoResourceAttr("snowflake_stream_grant.test", "stream_name"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_future", "true"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "privilege", "SELECT"),
				),
			},
			{
				ResourceName:      "snowflake_stream_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"enable_multiple_grants", // feature flag attribute not defined in Snowflake, can't be imported
				},
			},
		},
	})
}

func streamGrantConfig(name string, streamName string, grantType grantType, privilege string, databaseName string, schemaName string) string {
	var streamNameConfig string
	switch grantType {
	case normal:
		streamNameConfig = "stream_name = snowflake_stream.test.name"
	case onFuture:
		streamNameConfig = "on_future = true"
	case onAll:
		streamNameConfig = "on_all = true"
	}

	return fmt.Sprintf(`
resource "snowflake_role" "test" {
    name = "%s"
}
resource "snowflake_table" "test" {
	name            = "%s"
	database        = "%s"
	schema          = "%s"
	change_tracking = true
	comment         = "Terraform acceptance test"

	column {
		name = "column1"
		type = "VARIANT"
	}
	column {
		name = "column2"
		type = "VARCHAR(16777216)"
	}
}

resource "snowflake_stream" "test" {
	name     = "%s"
	database = "%s"
	schema   = "%s"
	comment  = "Terraform acceptance test"
	on_table = "\"${snowflake_table.test.database}\".\"${snowflake_table.test.schema}\".${snowflake_table.test.name}"
}

resource "snowflake_stream_grant" "test" {
    database_name = "%s"
	roles         = [snowflake_role.test.name]
	schema_name   = "%s"
	%s
    privilege = "%s"
}
`, name, name, databaseName, schemaName, streamName, databaseName, schemaName, databaseName, schemaName, streamNameConfig, privilege)
}
