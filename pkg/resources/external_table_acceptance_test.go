package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ExternalTable_basic(t *testing.T) {
	env := os.Getenv("SKIP_EXTERNAL_TABLE_TEST")
	if env != "" {
		t.Skip("Skipping TestAcc_ExternalTable")
	}
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	bucketURL := os.Getenv("AWS_EXTERNAL_BUCKET_URL")
	if bucketURL == "" {
		t.Skip("Skipping TestAcc_ExternalTable")
	}
	roleName := os.Getenv("AWS_EXTERNAL_ROLE_NAME")
	if roleName == "" {
		t.Skip("Skipping TestAcc_ExternalTable")
	}
	resourceName := "snowflake_external_table.test_table"

	configVariables := map[string]config.Variable{
		"name":     config.StringVariable(name),
		"location": config.StringVariable(bucketURL),
		"aws_arn":  config.StringVariable(roleName),
		"database": config.StringVariable(acc.TestDatabaseName),
		"schema":   config.StringVariable(acc.TestSchemaName),
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: testAccCheckExternalTableDestroy,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: configVariables,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "location", fmt.Sprintf(`@"%s"."%s"."%s"`, acc.TestDatabaseName, acc.TestSchemaName, name)),
					resource.TestCheckResourceAttr(resourceName, "file_format", "TYPE = CSV"),
					resource.TestCheckResourceAttr(resourceName, "comment", "Terraform acceptance test"),
					resource.TestCheckResourceAttr(resourceName, "column.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "column[0].name", "column1"),
					resource.TestCheckResourceAttr(resourceName, "column[0].type", "STRING"),
					resource.TestCheckResourceAttr(resourceName, "column[0].as", "TO_VARCHAR(TO_TIMESTAMP_NTZ(value:unix_timestamp_property::NUMBER, 3), 'yyyy-mm-dd-hh')"),
					resource.TestCheckResourceAttr(resourceName, "column[1].name", "column2"),
					resource.TestCheckResourceAttr(resourceName, "column[1].type", "TIMESTAMP_NTZ(9)"),
					resource.TestCheckResourceAttr(resourceName, "column[1].as", "($1:\"CreatedDate\"::timestamp)"),
				),
			},
		},
	})
}

func testAccCheckExternalTableDestroy(s *terraform.State) error {
	db := acc.TestAccProvider.Meta().(*sql.DB)
	client := sdk.NewClientFromDB(db)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "snowflake_external_table" {
			continue
		}
		ctx := context.Background()
		id := sdk.NewSchemaObjectIdentifier(rs.Primary.Attributes["database"], rs.Primary.Attributes["schema"], rs.Primary.Attributes["name"])
		dynamicTable, err := client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(id))
		if err == nil {
			return fmt.Errorf("external table %v still exists", dynamicTable.Name)
		}
	}
	return nil
}
