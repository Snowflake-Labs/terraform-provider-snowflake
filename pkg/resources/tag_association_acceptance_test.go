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

func TestAcc_TagAssociation(t *testing.T) {
	accName := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfig(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "object_type", "DATABASE"),
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "tag_value", "finance"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationSchema(t *testing.T) {
	accName := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigSchema(accName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.schema", "object_type", "SCHEMA"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationColumn(t *testing.T) {
	accName := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	accName2 := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigColumn(accName, accName2, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_type", "COLUMN"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "tag_value", "TAG_VALUE"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.%", "3"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.name", fmt.Sprintf("%s.column_name", accName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.schema", acc.TestSchemaName),
				),
			},
		},
	})
}

func tagAssociationConfig(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test" {
	name = "%[1]v"
	database = "%[2]s"
	schema = "%[3]s"
	allowed_values = ["finance", "hr"]
	comment = "Terraform acceptance test"
}

resource "snowflake_tag_association" "test" {
	object_identifier {
		name = "%[2]s"
	  }
	object_type = "DATABASE"
	tag_id = snowflake_tag.test.id
	tag_value = "finance"
}
`, n, databaseName, schemaName)
}

func tagAssociationConfigSchema(n string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "tag1" {
 name     = "%s"
 database = "%s"
 schema   = "%s"

 allowed_values = []
}

resource "snowflake_tag_association" "schema" {
  object_identifier {
    database = "%s"
    name     = "%s"
  }

  object_type = "SCHEMA"
  tag_id      = snowflake_tag.tag1.id
  tag_value   = "TAG_VALUE"
}
`, n, databaseName, schemaName, databaseName, schemaName)
}

func tagAssociationConfigColumn(n1, n2 string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "tag1" {
	name     = "%[1]v"
	database = "%[3]v"
	schema   = "%[4]v"
}

resource "snowflake_table" "test_table" {
	name                = "%[1]v"
	database            = "%[3]v"
	schema              = "%[4]v"

	column {
		name    = "column_name"
		type    = "VARIANT"
	}
}

resource "snowflake_tag_association" "columnTag" {
	object_identifier {
		database            = "%[3]v"
		schema              = "%[4]v"
		name     = "${snowflake_table.test_table.name}.${snowflake_table.test_table.column[0].name}"
	}

	object_type = "COLUMN"
	tag_id      = snowflake_tag.tag1.id
	tag_value   = "TAG_VALUE"
}
`, n1, n2, databaseName, schemaName)
}
