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

func TestAcc_Tag(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
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
	name = "%s"
	database = "%s"
	schema = "%s"
	allowed_values = ["alv1", "alv2"]
	comment = "Terraform acceptance test"
}
`, n, databaseName, schemaName)
}
