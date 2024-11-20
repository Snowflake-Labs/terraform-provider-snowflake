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
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TagAssociation(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tag2Id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tagValue := "foo"
	tagValue2 := "bar"
	databaseId := acc.TestClient().Ids.DatabaseId()
	resourceName := "snowflake_tag_association.test"
	m := func(tagId sdk.SchemaObjectIdentifier, tagValue string) map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":                      config.StringVariable(tagId.Name()),
			"tag_value":                     config.StringVariable(tagValue),
			"database":                      config.StringVariable(databaseId.Name()),
			"schema":                        config.StringVariable(acc.TestSchemaName),
			"database_fully_qualified_name": config.StringVariable(databaseId.FullyQualifiedName()),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckTagValueEmpty(t),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ConfigVariables: m(tagId, tagValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", tagValue),
				),
			},
			// external change
			{
				PreConfig: func() {
					acc.TestClient().Tag.Unset(t, sdk.ObjectTypeDatabase, sdk.NewAccountObjectIdentifier(acc.TestDatabaseName), []sdk.ObjectIdentifier{tagId})
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ConfigVariables: m(tagId, tagValue),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", tagValue),
				),
			},
			// change tag value
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ConfigVariables: m(tagId, tagValue2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), tagValue2, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", tagValue2),
				),
			},
			// change tag id
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ConfigVariables: m(tag2Id, tagValue2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tag2Id.FullyQualifiedName(), tagValue2, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tag2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", tagValue2),
				),
			},
			{
				ConfigVariables:   m(tag2Id, tagValue2),
				ConfigDirectory:   acc.ConfigurationDirectory("TestAcc_TagAssociation/basic"),
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// object_identifiers does not get set because during the import, the configuration is considered as empty
				ImportStateVerifyIgnore: []string{"skip_validation", "object_identifiers.#", "object_identifiers.0"},
			},
			// after refreshing the state, object_identifiers is correct
			{
				RefreshState: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tag2Id.FullyQualifiedName(), tagValue2, string(sdk.ObjectTypeDatabase))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeDatabase)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", acc.TestClient().Ids.DatabaseId().FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tag2Id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", tagValue2),
				),
			},
		},
	})
}

func TestAcc_TagAssociationSchema(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeSchema))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationColumn(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", columnId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1202(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableName := acc.TestClient().Ids.Alpha()
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":   config.StringVariable(tagId.Name()),
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", "TABLE"),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1909(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId1 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	tableId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	columnId1 := sdk.NewTableColumnIdentifier(tableId1.DatabaseName(), tableId1.SchemaName(), tableId1.Name(), "test.column")
	columnId2 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "test.column")
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":                     config.StringVariable(tagId.Name()),
			"table_name":                   config.StringVariable(tableId1.Name()),
			"table_name2":                  config.StringVariable(tableId2.Name()),
			"column_name":                  config.StringVariable("test.column"),
			"column_fully_qualified_name":  config.StringVariable(columnId1.FullyQualifiedName()),
			"column2_fully_qualified_name": config.StringVariable(columnId2.FullyQualifiedName()),
			"database":                     config.StringVariable(acc.TestDatabaseName),
			"schema":                       config.StringVariable(acc.TestSchemaName),
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
					testAccCheckTableColumnTagAssociation(tagId, columnId1, "v1"),
					testAccCheckTableColumnTagAssociation(tagId, columnId2, "v1"),
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

// TODO(SNOW-1165821): use a separate account with ORGADMIN in CI

func TestAcc_TagAssociationAccountIssues1910(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	accountName := acc.TestClient().Context.CurrentAccountName(t)
	orgName := acc.TestClient().Context.CurrentOrganizationName(t)
	accountId := sdk.NewAccountIdentifier(orgName, accountName)
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]config.Variable {
		return map[string]config.Variable{
			"tag_name":                     config.StringVariable(tagId.Name()),
			"account_fully_qualified_name": config.StringVariable(accountId.FullyQualifiedName()),
			"database":                     config.StringVariable(acc.TestDatabaseName),
			"schema":                       config.StringVariable(acc.TestSchemaName),
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeAccount)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", accountId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "v1"),
				),
			},
		},
	})
}

func TestAcc_TagAssociationIssue1926(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", columnId1.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1926"),
				ConfigVariables: m2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", columnId2.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_TagAssociation/issue1926"),
				ConfigVariables: m3,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeColumn))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeColumn)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", columnId3.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func TestAcc_Tag_migrateFromVersion_0_98_0(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	// tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	// t.Cleanup(tagCleanuop)
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_tag_association.test"
	schemaId := acc.TestClient().Ids.SchemaId()

	m := func() config.Variables {
		return config.Variables{
			"tag_name":                    config.StringVariable(tagId.Name()),
			"database":                    config.StringVariable(acc.TestDatabaseName),
			"schema":                      config.StringVariable(acc.TestSchemaName),
			"schema_fully_qualified_name": config.StringVariable(schemaId.FullyQualifiedName()),
		}
	}

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},

		Steps: []resource.TestStep{
			{
				ExternalProviders: acc.ExternalProviderWithExactVersion("0.98.0"),
				Config:            tagAssociation_v_0_98_0(tagId, "TAG_VALUE", sdk.ObjectTypeSchema, schemaId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.DatabaseName(), tagId.SchemaName(), tagId.Name())),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.name", schemaId.Name()),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.database", schemaId.DatabaseName()),
					resource.TestCheckResourceAttr(resourceName, "object_identifier.0.schema", ""),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				ConfigDirectory:          acc.ConfigurationDirectory("TestAcc_TagAssociation/schema"),
				ConfigVariables:          m(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "id", helpers.EncodeSnowflakeID(tagId.FullyQualifiedName(), "TAG_VALUE", string(sdk.ObjectTypeSchema))),
					resource.TestCheckResourceAttr(resourceName, "object_type", string(sdk.ObjectTypeSchema)),
					resource.TestCheckResourceAttr(resourceName, "object_identifiers.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "object_identifiers.*", schemaId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_id", tagId.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "tag_value", "TAG_VALUE"),
				),
			},
		},
	})
}

func tagAssociation_v_0_98_0(tagId sdk.SchemaObjectIdentifier, tagValue string, objectType sdk.ObjectType, objectId sdk.DatabaseObjectIdentifier) string {
	s := `
resource "snowflake_tag_association" "test" {
	tag_id					= snowflake_tag.test.fully_qualified_name
	tag_value				= "%[1]s"
	object_type				= "%[2]s"
	object_identifier {
		name = "%[3]s"
		database = "%[4]s"
	}
}

resource "snowflake_tag" "test" {
  name           = "%[5]s"
  database       = "%[6]s"
  schema         = "%[7]s"
}
`
	return fmt.Sprintf(s, tagValue, objectType, objectId.Name(), objectId.DatabaseName(), tagId.Name(), tagId.DatabaseName(), tagId.SchemaName())
}
