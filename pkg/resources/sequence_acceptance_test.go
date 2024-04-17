package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Sequence(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accRename := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: acc.CheckDestroy(t, resources.Sequence),
		Steps: []resource.TestStep{
			// CREATE
			{
				Config: sequenceConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "ordering", "ORDER"),
				),
			},
			// Set comment and rename
			{
				Config: sequenceConfigWithComment(accRename, "look at me I am a comment", acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accRename),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", "look at me I am a comment"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, accRename)),
				),
			},
			// Unset comment and set increment
			{
				Config: sequenceConfigWithIncrement(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "name", accName),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "comment", ""),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "next_value", "1"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "increment", "32"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "ordering", "NOORDER"),
					resource.TestCheckResourceAttr("snowflake_sequence.test_sequence", "fully_qualified_name", fmt.Sprintf(`"%v"."%v"."%v"`, acc.TestDatabaseName, acc.TestSchemaName, accName)),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_sequence.test_sequence",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func sequenceConfigWithIncrement(sequenceName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_sequence" {
	name       = "%s"
	database   = "%s"
	schema     = "%s"
    increment = 32
	ordering = "NOORDER"
}
`
	return fmt.Sprintf(s, sequenceName, databaseName, schemaName)
}

func sequenceConfig(sequenceName string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_sequence" {
	name     = "%s"
	database   = "%s"
	schema     = "%s"
}
`
	return fmt.Sprintf(s, sequenceName, databaseName, schemaName)
}

func sequenceConfigWithComment(sequenceName, comment string, databaseName string, schemaName string) string {
	s := `
resource "snowflake_sequence" "test_sequence" {
	name     = "%s"
	database   = "%s"
	schema     = "%s"
    comment  = "%s"
}
`
	return fmt.Sprintf(s, sequenceName, databaseName, schemaName, comment)
}
