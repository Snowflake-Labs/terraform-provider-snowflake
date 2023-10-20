package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Tag(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_tag.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_tag.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_tag.test", "allowed_values.#", "2"),
					resource.TestCheckResourceAttr("snowflake_tag.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func tagConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test" {
	name = "%[1]v"
	database = "%[2]s"
	schema = "	%[3]s"
	allowed_values = ["alv1", "alv2"]
	comment = "Terraform acceptance test"
}
`, n, databaseName, schemaName)
}
