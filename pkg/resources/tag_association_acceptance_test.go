package resources_test

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TagAssociation(t *testing.T) {
	tagName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name": config.StringVariable(tagName),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "finance"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationSchema(t *testing.T) {
	tagName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name": config.StringVariable(tagName),
			"database": config.StringVariable(acc.TestDatabaseName),
			"schema":   config.StringVariable(acc.TestSchemaName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/schema"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "SCHEMA"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationColumn(t *testing.T) {
	tagName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	tableName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":   config.StringVariable(tagName),
			"table_name": config.StringVariable(tableName),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/column"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "COLUMN"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.name", fmt.Sprintf("%s.column_name", tableName)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.schema", acc.TestSchemaName)),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1202(t *testing.T) {
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":   config.StringVariable(tagName),
			"table_name": config.StringVariable(tableName),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1202"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "TABLE"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1909(t *testing.T) {
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName2 := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	columnName := "test.column"
	resourceName := "snowflake_tag_association.test"
	objectID := sdk.NewTableColumnIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName, columnName)
	objectID2 := sdk.NewTableColumnIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName2, columnName)
	tagID := sdk.NewSchemaObjectIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tagName)
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":    config.StringVariable(tagName),
			"table_name":  config.StringVariable(tableName),
			"table_name2": config.StringVariable(tableName2),
			"column_name": config.StringVariable("test.column"),
			"database":    config.StringVariable(acc.TestDatabaseName),
			"schema":      config.StringVariable(acc.TestSchemaName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1909"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "COLUMN"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
					testAccCheckTableColumnTagAssociation(tagID, objectID, "v1"),
					testAccCheckTableColumnTagAssociation(tagID, objectID2, "v1"),
				),
			},
		},
	})
}

func testAccCheckTableColumnTagAssociation(tagID sdk.SchemaObjectIdentifier, objectID sdk.ObjectIdentifier, tagValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		db := acc.TestAccProvider.Meta().(*sql.DB)
		client := sdk.NewClientFromDB(db)
		ctx := context.Background()
		tv, err := client.SystemFunctions.GetTag(ctx, tagID, objectID, sdk.ObjectTypeColumn)
		if err != nil {
			return err
		}
		if tagValue != tv {
			return fmt.Errorf("expected tag value %s, got %s", tagValue, tv)
		}
		return nil
	}
}

func TestAcc_TagAssociationAccountIssues1910(t *testing.T) {
	// todo: use role with ORGADMIN in CI (SNOW-1165821)
	// SNOWFLAKE_TEST_ACCOUNT_CREATE must be set to 1 to run this test
	if _, ok := os.LookupEnv("SNOWFLAKE_TEST_ACCOUNT_CREATE"); !ok {
		t.Skip("Skipping TestInt_AccountCreate")
	}
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	accountName := "account_" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":     config.StringVariable(tagName),
			"account_name": config.StringVariable(accountName),
			"database":     config.StringVariable(acc.TestDatabaseName),
			"schema":       config.StringVariable(acc.TestSchemaName),
		}
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1910"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "ACCOUNT"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1926(t *testing.T) {
	tagName := "tag-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	tableName := "table-" + strings.ToUpper(acctest.RandStringFromCharSet(4, acctest.CharSetAlpha))
	resourceName := "snowflake_tag_association.test"
	columnName := "test.column"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":    config.StringVariable(tagName),
			"table_name":  config.StringVariable(tableName),
			"column_name": config.StringVariable(columnName),
			"database":    config.StringVariable(acc.TestDatabaseName),
			"schema":      config.StringVariable(acc.TestSchemaName),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1926"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "COLUMN"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.name", fmt.Sprintf("%s.%s", tableName, columnName)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.schema", acc.TestSchemaName),
				),
			},
		},
	})
}
