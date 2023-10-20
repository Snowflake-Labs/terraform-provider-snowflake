package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ExternalStage(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: externalStageConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_stage.test", "name", accName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr("snowflake_stage.test", "comment", "Terraform acceptance test"),
				),
			},
		},
	})
}

func externalStageConfig(n, databaseName, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_stage" "test" {
	name = "%v"
	url = "s3://com.example.bucket/prefix"
	database = "%s"
	schema = "%s"
	comment = "Terraform acceptance test"
}
`, n, databaseName, schemaName)
}
