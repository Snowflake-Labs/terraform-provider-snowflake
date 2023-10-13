package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_TagAssociation(t *testing.T) {
	accName := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfig(accName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "object_type", "DATABASE"),
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "tag_id", fmt.Sprintf("%s|%s|%s", accName, accName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.test", "tag_value", "finance"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationSchema(t *testing.T) {
	accName := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigSchema(accName),
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
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigColumn(accName, accName2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_type", "COLUMN"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "tag_id", fmt.Sprintf("%s|%s|%s", accName, accName, accName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "tag_value", "TAG_VALUE"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.%", "3"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.name", "test_table.column_name"),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.database", accName2),
					resource.TestCheckResourceAttr("snowflake_tag_association.columnTag", "object_identifier.0.schema", accName2),
				),
			},
		},
	})
}

func tagAssociationConfig(n string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "test" {
	name = "%[1]v"
	comment = "Terraform acceptance test"
}

resource "snowflake_schema" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	comment = "Terraform acceptance test"
}

resource "snowflake_tag" "test" {
	name = "%[1]v"
	database = snowflake_database.test.name
	schema = snowflake_schema.test.name
	allowed_values = ["finance", "hr"]
	comment = "Terraform acceptance test"
}

resource "snowflake_tag_association" "test" {
	object_identifier {
		name = snowflake_database.test.name
	  }
	object_type = "DATABASE"
	tag_id = snowflake_tag.test.id
	tag_value = "finance"
}
`, n)
}

func tagAssociationConfigSchema(n string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "db" {
	name = "test_db"
}

resource "snowflake_schema" "sch" {
	database = snowflake_database.db.name
	name = "test_sch"
	comment = "%v"
}

resource "snowflake_tag" "tag1" {
 database = snowflake_database.db.name
 name     = "EXAMPLE_TAG"
 schema   = "PUBLIC"

 allowed_values = []
}

resource "snowflake_tag_association" "schema" {
  object_identifier {
    database = snowflake_database.db.name
    name     = snowflake_schema.sch.name
  }

  object_type = "SCHEMA"
  tag_id      = snowflake_tag.tag1.id
  tag_value   = "TAG_VALUE"
}
`, n)
}

func tagAssociationConfigColumn(n1 string, n2 string) string {
	return fmt.Sprintf(`
resource "snowflake_database" "tag_db" {
	name = "%[1]v"
}

resource "snowflake_schema" "tag_sch" {
	database = snowflake_database.tag_db.name
	name = "%[1]v"
}

resource "snowflake_database" "table_db" {
	name = "%[2]v"
}

resource "snowflake_schema" "table_sch" {
	database = snowflake_database.table_db.name
	name = "%[2]v"
}

resource "snowflake_tag" "tag1" {
 database = snowflake_database.tag_db.name
 name     = "%[1]v"
 schema   = snowflake_schema.tag_sch.name
}

resource "snowflake_table" "test_table" {
 database            = snowflake_database.table_db.name
 schema              = snowflake_schema.table_sch.name
 name                = "test_table"

 column {
   name    = "column_name"
   type    = "VARIANT"
 }
}

resource "snowflake_tag_association" "columnTag" {
  object_identifier {
    database = snowflake_database.table_db.name
	schema   = snowflake_schema.table_sch.name
    name     = "${snowflake_table.test_table.name}.${snowflake_table.test_table.column[0].name}"
  }

  object_type = "COLUMN"
  tag_id      = snowflake_tag.tag1.id
  tag_value   = "TAG_VALUE"
}
`, n1, n2)
}
