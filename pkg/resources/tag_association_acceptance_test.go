package resources_test

import (
	"context"
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TagAssociation(t *testing.T) {
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name": config.StringVariable(tagId.Name()),
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
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "finance", string("DATABASE"))),
					resource.TestCheckResourceAttr(resourceName, "object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0", acc.TestClient().Ids.SchemaId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "finance"),
				),
			},
			{
				PreConfig: func() {
					acc.TestClient().Tag.Unset(t, "DATABASE", sdk.NewAccountObjectIdentifier(acc.TestDatabaseName), []sdk.ObjectIdentifier{tagId})
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ConfigVariables: m(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "finance", string("DATABASE"))),
					resource.TestCheckResourceAttr(resourceName, "object_type", "DATABASE"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0", acc.TestClient().Ids.SchemaId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "finance"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationSchema(t *testing.T) {
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	schemaId := acc.TestClient().Ids.SchemaId()
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":                    config.StringVariable(tagId.Name()),
			"database":                    config.StringVariable(acc.TestDatabaseName),
			"schema":                      config.StringVariable(acc.TestSchemaName),
			"schema_fully_qualified_name": config.StringVariable(schemaId.FullyQualifiedName()),
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
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeSchema))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifier.*", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationColumn(t *testing.T) {
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	columnId := sdk.NewTableColumnIdentifier(tableId.DatabaseName(), tableId.SchemaName(), tableId.Name(), "column")
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":                    config.StringVariable(tagId.Name()),
			"table_name":                  config.StringVariable(tableId.Name()),
			"database":                    config.StringVariable(acc.TestDatabaseName),
			"schema":                      config.StringVariable(acc.TestSchemaName),
			"column":                      config.StringVariable("column"),
			"column_fully_qualified_name": config.StringVariable(columnId.FullyQualifiedName()),
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
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifier.*", columnId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1202(t *testing.T) {
	tagName := acc.TestClient().Ids.Alpha()
	tableName := acc.TestClient().Ids.Alpha()
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
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tagName := tagId.Name()
	tableName := acc.TestClient().Ids.Alpha()
	tableName2 := acc.TestClient().Ids.Alpha()
	columnName := "test.column"
	resourceName := "snowflake_tag_association.test"
	objectID := sdk.NewTableColumnIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName, columnName)
	objectID2 := sdk.NewTableColumnIdentifier(acc.TestDatabaseName, acc.TestSchemaName, tableName2, columnName)
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
					testAccCheckTableColumnTagAssociation(tagId, objectID, "v1"),
					testAccCheckTableColumnTagAssociation(tagId, objectID2, "v1"),
				),
			},
		},
	})
}

func testAccCheckTableColumnTagAssociation(tagID sdk.SchemaObjectIdentifier, objectID sdk.ObjectIdentifier, tagValue string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acc.TestAccProvider.Meta().(*provider.Context).Client
		ctx := context.Background()
		tv, err := client.SystemFunctions.GetTag(ctx, tagID, objectID, sdk.ObjectTypeColumn)
		if err != nil {
			return err
		}
		if tv == nil {
			return fmt.Errorf("expected tag value %s, got nil", tagValue)
		}
		if tagValue != *tv {
			return fmt.Errorf("expected tag value %s, got %s", tagValue, *tv)
		}
		return nil
	}
}

func TestAcc_TagAssociationAccountIssues1910(t *testing.T) {
	// todo: use role with ORGADMIN in CI (SNOW-1165821)
	_ = testenvs.GetOrSkipTest(t, testenvs.TestAccountCreate)
	tagName := acc.TestClient().Ids.Alpha()
	accountName := acc.TestClient().Ids.Alpha()
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
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId1 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	columnId1 := sdk.NewTableColumnIdentifier(tableId1.DatabaseName(), tableId1.SchemaName(), tableId1.Name(), "init")
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":                    config.StringVariable(tagId.Name()),
			"table_name":                  config.StringVariable(tableId1.Name()),
			"column_name":                 config.StringVariable(columnId1.Name()),
			"column_fully_qualified_name": config.StringVariable(columnId1.FullyQualifiedName()),
			"database":                    config.StringVariable(acc.TestDatabaseName),
			"schema":                      config.StringVariable(acc.TestSchemaName),
		}
	}

	m2 := m()
	tableId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("table.test")
	columnId2 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "column")
	columnId3 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "column.test")
	m2["table_name"] = config.StringVariable(tableId2.Name())
	m2["column_name"] = config.StringVariable(columnId2.Name())
	m2["column_fully_qualified_name"] = config.StringVariable(columnId2.FullyQualifiedName())
	m3 := m()
	m3["table_name"] = config.StringVariable(tableId2.Name())
	m3["column_name"] = config.StringVariable(columnId3.Name())
	m3["column_fully_qualified_name"] = config.StringVariable(columnId3.FullyQualifiedName())
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
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifier.*", columnId1.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1926"),
				ConfigVariables: m2,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr(resourceName, "object_type", "COLUMN"),
					// resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					// resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.%", "3"),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.name", fmt.Sprintf("%s.%s", tableName2, columnName2)),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.database", acc.TestDatabaseName),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifier.*", columnId2.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1926"),
				ConfigVariables: m3,
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttr(resourceName, "object_type", "COLUMN"),
					// resource.TestCheckResourceAttr(resourceName, "tag_id", fmt.Sprintf("%s|%s|%s", acc.TestDatabaseName, acc.TestSchemaName, tagName)),
					// resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.%", "3"),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.name", fmt.Sprintf("%s.%s", tableName2, columnName3)),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.database", acc.TestDatabaseName),
					// resource.TestCheckResourceAttr(resourceName, "object_identifier.0.schema", acc.TestSchemaName),
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifier.*", columnId3.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}
