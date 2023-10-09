package datasources_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Database(t *testing.T) {
	databaseName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	comment := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    providers(),
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: database(databaseName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_database.t", "name", databaseName),
					resource.TestCheckResourceAttr("data.snowflake_database.t", "comment", comment),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "created_on"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "owner"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "retention_time"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "is_current"),
					resource.TestCheckResourceAttrSet("data.snowflake_database.t", "is_default"),
				),
			},
		},
	})
}

func database(databaseName, comment string) string {
	return fmt.Sprintf(`
		resource snowflake_database "test_database" {
			name = "%v"
			comment = "%v"
		}
		data snowflake_database "t" {
			depends_on = [snowflake_database.test_database]
			name = "%v"
		}
	`, databaseName, comment, databaseName)
}
