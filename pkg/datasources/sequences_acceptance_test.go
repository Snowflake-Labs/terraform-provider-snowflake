package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSequences(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	schemaName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	sequenceName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers: providers(),
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
		depends_on = [snowflake_sequence.t]
	}
	`, databaseName, schemaName, sequenceName)
}
