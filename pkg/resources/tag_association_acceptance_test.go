package resources_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/provider"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfconfig "github.com/hashicorp/terraform-plugin-testing/config"
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
	m := func(tagId sdk.SchemaObjectIdentifier, tagValue string) map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":                      tfconfig.StringVariable(tagId.Name()),
			"tag_value":                     tfconfig.StringVariable(tagValue),
			"database":                      tfconfig.StringVariable(databaseId.Name()),
			"schema":                        tfconfig.StringVariable(acc.TestSchemaName),
			"database_fully_qualified_name": tfconfig.StringVariable(databaseId.FullyQualifiedName()),
		}
	}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckResourceTagUnset(t),
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
			// external change - unset tag
			{
				PreConfig: func() {
					acc.TestClient().Tag.Unset(t, sdk.ObjectTypeDatabase, databaseId, []sdk.ObjectIdentifier{tagId})
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
			// external change - set a different value
			{
				PreConfig: func() {
					acc.TestClient().Tag.Set(t, sdk.ObjectTypeDatabase, databaseId, []sdk.TagAssociation{
						{
							Name:  tagId,
							Value: "external",
						},
					})
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
					acc.CheckTagUnset(t, tagId, acc.TestClient().Ids.DatabaseId(), sdk.ObjectTypeDatabase),
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

func TestAcc_TagAssociation_objectIdentifiers(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	dbRole1, dbRole1Cleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRole1Cleanup)
	dbRole2, dbRole2Cleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRole2Cleanup)
	dbRole3, dbRole3Cleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRole3Cleanup)

	model12 := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole1.ID(), dbRole2.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")
	model123 := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole1.ID(), dbRole2.ID(), dbRole3.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")
	model13 := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole1.ID(), dbRole3.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			acc.CheckResourceTagUnset(t),
		),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model12),
				Check: assert.AssertThat(t, resourceassert.TagAssociationResource(t, model12.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(2).
					HasTagValueString("foo"),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole1.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole2.ID().FullyQualifiedName())),
				),
			},
			{
				Config: config.FromModel(t, model123),
				Check: assert.AssertThat(t, resourceassert.TagAssociationResource(t, model12.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(3).
					HasTagValueString("foo"),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole1.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole2.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model12.ResourceReference(), "object_identifiers.*", dbRole3.ID().FullyQualifiedName())),
				),
			},
			{
				Config: config.FromModel(t, model13),
				Check: assert.AssertThat(t, resourceassert.TagAssociationResource(t, model13.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(2).
					HasTagValueString("foo"),
					assert.Check(resource.TestCheckTypeSetElemAttr(model13.ResourceReference(), "object_identifiers.*", dbRole1.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckTypeSetElemAttr(model13.ResourceReference(), "object_identifiers.*", dbRole3.ID().FullyQualifiedName())),
					assert.Check(acc.CheckTagUnset(t, tag.ID(), dbRole2.ID(), sdk.ObjectTypeDatabaseRole)),
				),
			},
		},
	})
}

func TestAcc_TagAssociation_objectType(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	role, roleCleanup := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(roleCleanup)
	dbRole, dbRoleCleanup := acc.TestClient().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(dbRoleCleanup)

	baseModel := model.TagAssociation("test", []sdk.ObjectIdentifier{role.ID()}, string(sdk.ObjectTypeRole), tag.ID().FullyQualifiedName(), "foo")
	modelWithDifferentObjectType := model.TagAssociation("test", []sdk.ObjectIdentifier{dbRole.ID()}, string(sdk.ObjectTypeDatabaseRole), tag.ID().FullyQualifiedName(), "foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(
			acc.CheckResourceTagUnset(t),
		),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, baseModel),
				Check: assert.AssertThat(t, resourceassert.TagAssociationResource(t, baseModel.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(1).
					HasTagValueString("foo"),
				),
			},
			{
				Config: config.FromModel(t, modelWithDifferentObjectType),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(modelWithDifferentObjectType.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.TagAssociationResource(t, baseModel.ResourceReference()).
					HasObjectTypeString(string(sdk.ObjectTypeDatabaseRole)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(1).
					HasTagValueString("foo"),
					assert.Check(acc.CheckTagUnset(t, tag.ID(), role.ID(), sdk.ObjectTypeRole)),
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
	m := func() map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":                    tfconfig.StringVariable(tagId.Name()),
			"database":                    tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":                      tfconfig.StringVariable(acc.TestSchemaName),
			"schema_fully_qualified_name": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
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

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/3235 is fixed
func TestAcc_TagAssociation_lowercaseObjectType(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tag, tagCleanup := acc.TestClient().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)
	objectType := strings.ToLower(string(sdk.ObjectTypeSchema))
	objectId := acc.TestClient().Ids.SchemaId()

	model := model.TagAssociation("test", []sdk.ObjectIdentifier{objectId}, objectType, tag.ID().FullyQualifiedName(), "foo")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.TagAssociationResource(t, model.ResourceReference()).
					HasIdString(helpers.EncodeSnowflakeID(tag.ID().FullyQualifiedName(), "foo", string(sdk.ObjectTypeSchema))).
					HasObjectTypeString(string(sdk.ObjectTypeSchema)).
					HasTagIdString(tag.ID().FullyQualifiedName()).
					HasObjectIdentifiersLength(1).
					HasTagValueString("foo"),
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
	m := func() map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":                    tfconfig.StringVariable(tagId.Name()),
			"table_name":                  tfconfig.StringVariable(tableId.Name()),
			"database":                    tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":                      tfconfig.StringVariable(acc.TestSchemaName),
			"column":                      tfconfig.StringVariable("column"),
			"column_fully_qualified_name": tfconfig.StringVariable(columnId.FullyQualifiedName()),
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
	m := func() map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":   tfconfig.StringVariable(tagId.Name()),
			"table_name": tfconfig.StringVariable(tableName),
			"database":   tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":     tfconfig.StringVariable(acc.TestSchemaName),
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
	m := func() map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":                     tfconfig.StringVariable(tagId.Name()),
			"table_name":                   tfconfig.StringVariable(tableId1.Name()),
			"table_name2":                  tfconfig.StringVariable(tableId2.Name()),
			"column_name":                  tfconfig.StringVariable("test.column"),
			"column_fully_qualified_name":  tfconfig.StringVariable(columnId1.FullyQualifiedName()),
			"column2_fully_qualified_name": tfconfig.StringVariable(columnId2.FullyQualifiedName()),
			"database":                     tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":                       tfconfig.StringVariable(acc.TestSchemaName),
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
	accountId := acc.TestClient().Context.CurrentAccountIdentifier(t)
	resourceName := "snowflake_tag_association.test"
	m := func() map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":                     tfconfig.StringVariable(tagId.Name()),
			"account_fully_qualified_name": tfconfig.StringVariable(accountId.FullyQualifiedName()),
			"database":                     tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":                       tfconfig.StringVariable(acc.TestSchemaName),
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
	m := func() map[string]tfconfig.Variable {
		return map[string]tfconfig.Variable{
			"tag_name":                    tfconfig.StringVariable(tagId.Name()),
			"table_name":                  tfconfig.StringVariable(tableId1.Name()),
			"column_name":                 tfconfig.StringVariable(columnId1.Name()),
			"column_fully_qualified_name": tfconfig.StringVariable(columnId1.FullyQualifiedName()),
			"database":                    tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":                      tfconfig.StringVariable(acc.TestSchemaName),
		}
	}

	m2 := m()
	tableId2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix("table.test")
	columnId2 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "column")
	columnId3 := sdk.NewTableColumnIdentifier(tableId2.DatabaseName(), tableId2.SchemaName(), tableId2.Name(), "column.test")
	m2["table_name"] = tfconfig.StringVariable(tableId2.Name())
	m2["column_name"] = tfconfig.StringVariable(columnId2.Name())
	m2["column_fully_qualified_name"] = tfconfig.StringVariable(columnId2.FullyQualifiedName())
	m3 := m()
	m3["table_name"] = tfconfig.StringVariable(tableId2.Name())
	m3["column_name"] = tfconfig.StringVariable(columnId3.Name())
	m3["column_fully_qualified_name"] = tfconfig.StringVariable(columnId3.FullyQualifiedName())
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

func TestAcc_TagAssociation_migrateFromVersion_0_98_0(t *testing.T) {
	t.Setenv(string(testenvs.ConfigureClientOnce), "")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	tagId := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_tag_association.test"
	schemaId := acc.TestClient().Ids.SchemaId()

	m := func() tfconfig.Variables {
		return tfconfig.Variables{
			"tag_name":                    tfconfig.StringVariable(tagId.Name()),
			"database":                    tfconfig.StringVariable(acc.TestDatabaseName),
			"schema":                      tfconfig.StringVariable(acc.TestSchemaName),
			"schema_fully_qualified_name": tfconfig.StringVariable(schemaId.FullyQualifiedName()),
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
