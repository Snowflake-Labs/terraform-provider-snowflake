package resources_test

import (
	"fmt"
	"os"
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
		Providers:    acc.TestAccProviders(),
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
		Providers:    acc.TestAccProviders(),
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

func TestAcc_TagAssociationColumnIssues1926(t *testing.T) {
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	columnName := "test.column"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigColumnIssues1926(tagName, tableName, acc.TestDatabaseName, acc.TestSchemaName, columnName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_type", "COLUMN"),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_value", "v1"),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_identifier.0.%", "3"),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_identifier.0.name", fmt.Sprintf("%s.%s", tableName, columnName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_identifier.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_identifier.0.schema", acc.TestSchemaName),
				),
			},
		},
	})
}

func TestAcc_TagAssociationColumnIssues1909(t *testing.T) {
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	table1Name := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	table2Name := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	columnName := "test.column"

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigColumnIssues1909(tagName, table1Name, table2Name, acc.TestDatabaseName, acc.TestSchemaName, columnName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_type", "COLUMN"),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationTableIssues1202(t *testing.T) {
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigTableIssues1202(tagName, tableName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_type", "TABLE"),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationAccountIssues1910(t *testing.T) {
	// SNOWFLAKE_TEST_ACCOUNT_CREATE must be set to 1 to run this test
	if _, ok := os.LookupEnv("SNOWFLAKE_TEST_ACCOUNT_CREATE"); !ok {
		t.Skip("Skipping TestInt_AccountCreate")
	}

	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	accountName := "account_" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))

	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: tagAssociationConfigAccountIssues1910(tagName, accountName, acc.TestDatabaseName, acc.TestSchemaName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "object_type", "ACCOUNT"),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr("snowflake_tag_association.tag_column_association", "tag_value", "v1"),
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

func tagAssociationConfigColumnIssues1926(tagName, tableName string, databaseName string, schemaName string, columeName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test_tag" {
	name     = "%[1]v"
	database = "%[3]v"
	schema   = "%[4]v"
}
resource "snowflake_table" "test_table" {
	name                = "%[2]v"
	database            = "%[3]v"
	schema              = "%[4]v"
	column {
		name    = "%[5]v"
		type    = "VARIANT"
	}
}
resource "snowflake_tag_association" "tag_column_association" {
	object_identifier {
		database   = "%[3]v"
		schema     = "%[4]v"
		name       = "${snowflake_table.test_table.name}.${snowflake_table.test_table.column[0].name}"
	}
	object_type = "COLUMN"
	tag_id      = snowflake_tag.test_tag.id
	tag_value   = "v1"
}
`, tagName, tableName, databaseName, schemaName, columeName)
}

func tagAssociationConfigColumnIssues1909(tagName, table1Name string, table2Name string, databaseName string, schemaName string, columeName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test_tag" {
	name     = "%[1]v"
	database = "%[4]v"
	schema   = "%[5]v"
}
resource "snowflake_table" "test_table1" {
	name                = "%[2]v"
	database            = "%[4]v"
	schema              = "%[5]v"
	column {
		name    = "%[6]v"
		type    = "VARIANT"
	}
}
resource "snowflake_table" "test_table2" {
	name                = "%[3]v"
	database            = "%[4]v"
	schema              = "%[5]v"
	column {
		name    = "%[6]v"
		type    = "VARIANT"
	}
}
resource "snowflake_tag_association" "tag_column_association" {
	object_identifier {
		database   = "%[4]v"
		schema     = "%[5]v"
		name       = "${snowflake_table.test_table1.name}.${snowflake_table.test_table1.column[0].name}"
	}
	object_identifier {
		database   = "%[4]v"
		schema     = "%[5]v"
		name       = "${snowflake_table.test_table2.name}.${snowflake_table.test_table2.column[0].name}"
	}
	object_type = "COLUMN"
	tag_id      = snowflake_tag.test_tag.id
	tag_value   = "v1"
}
`, tagName, table1Name, table2Name, databaseName, schemaName, columeName)
}

func tagAssociationConfigTableIssues1202(tagName, tableName string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test_tag" {
	name     = "%[1]v"
	database = "%[3]v"
	schema   = "%[4]v"
}
resource "snowflake_table" "test_table" {
	name                = "%[2]v"
	database            = "%[3]v"
	schema              = "%[4]v"
	column {
		name    = "test_column"
		type    = "VARIANT"
	}
}
resource "snowflake_tag_association" "tag_column_association" {
	object_identifier {
		name = "${snowflake_table.test_table.name}"
	}
	object_type = "TABLE"
	tag_id      = snowflake_tag.test_tag.id
	tag_value   = "v1"
}
`, tagName, tableName, databaseName, schemaName)
}

func tagAssociationConfigAccountIssues1910(tagName, accountName string, databaseName string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_tag" "test_tag" {
	name     = "%[1]v"
	database = "%[3]v"
	schema   = "%[4]v"
}
resource "snowflake_account" "test_account" {
  	name = "%[2]v"
  	admin_name = "someadmin"
  	admin_password = "123456"
  	first_name = "Ad"
  	last_name = "Min"
  	email = "admin@example.com"
  	must_change_password = false
  	edition = "BUSINESS_CRITICAL"
  	grace_period_in_days = 4
}
resource "snowflake_tag_association" "tag_account_association" {
	object_identifier {
		name = "${snowflake_account.test_account.name}"
	}
	object_type = "ACCOUNT"
	tag_id      = snowflake_tag.test_tag.id
	tag_value   = "v1"
}
`, tagName, accountName, databaseName, schemaName)
}
