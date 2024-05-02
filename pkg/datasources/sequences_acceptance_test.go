package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Sequences(t *testing.T) {
	databaseName := acc.TestClient().Ids.Alpha()
	schemaName := acc.TestClient().Ids.Alpha()
	sequenceName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: sequences(databaseName, schemaName, sequenceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "database", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "schema", schemaName),
					resource.TestCheckResourceAttrSet("data.snowflake_sequences.t", "sequences.#"),
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "sequences.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_sequences.t", "sequences.0.name", sequenceName),
				),
			},
		},
	})
}

func sequences(databaseName string, schemaName string, sequenceName string) string {
	return fmt.Sprintf(`

	resource snowflake_database "d" {
		name = "%v"
	}

	resource snowflake_schema "s"{
		name 	 = "%v"
		database = snowflake_database.d.name
	}

	resource snowflake_sequence "t"{
		name 	 = "%v"
		database = snowflake_schema.s.database
		schema 	 = snowflake_schema.s.name
	}

	data snowflake_sequences "t" {
		database = snowflake_sequence.t.database
		schema = snowflake_sequence.t.schema
	}
	`, databaseName, schemaName, sequenceName)
}
