package resources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_StreamGrant_basic(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfig(name, streamName, normal),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", name),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "stream_name", streamName),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "with_grant_option", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_future", "false"),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "on_all", "false"),
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

func TestAcc_StreamGrant_onAll(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	streamName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfig(name, streamName, onAll),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", name),
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

	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: streamGrantConfig(name, streamName, onFuture),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "database_name", name),
					resource.TestCheckResourceAttr("snowflake_stream_grant.test", "schema_name", name),
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

func streamGrantConfig(name string, streamName string, grantType grantType) string {
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
resource "snowflake_database" "test" {
    name = "%s"
}

resource "snowflake_schema" "test" {
	name = "%s"
	database = snowflake_database.test.name
}

resource "snowflake_role" "test" {
    name = "%s"
}
resource "snowflake_table" "test" {
	database        = snowflake_database.test.name
	schema          = snowflake_schema.test.name
	name            = "%s"
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
	database = snowflake_database.test.name
	schema   = snowflake_schema.test.name
	name     = "%s"
	comment  = "Terraform acceptance test"
	on_table = "${snowflake_database.test.name}.${snowflake_schema.test.name}.${snowflake_table.test.name}"
}

resource "snowflake_stream_grant" "test" {
    database_name = snowflake_database.test.name
	roles         = [snowflake_role.test.name]
	schema_name   = snowflake_schema.test.name
	%s
    privilege = "SELECT"
}
`, name, name, name, name, streamName, streamNameConfig)
}
