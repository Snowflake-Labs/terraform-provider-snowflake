package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/stretchr/testify/require"
)

func TestAcc_StreamOnDirectoryTable_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	resourceName := "snowflake_stream_on_directory_table.test"

	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	baseModel := func() *model.StreamOnDirectoryTableModel {
		return model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName())
	}

	modelWithExtraFields := baseModel().
		WithCopyGrants(true).
		WithComment("foo")

	modelWithExtraFieldsModified := baseModel().
		WithCopyGrants(true).
		WithComment("bar")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: config.FromModels(t, baseModel()),
				Check: assert.AssertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// import without optionals
			{
				Config:       config.FromModels(t, baseModel()),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedStreamOnDirectoryTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageString(stage.ID().Name()),
				),
			},
			// set all fields
			{
				Config: config.FromModels(t, modelWithExtraFields),
				Check: assert.AssertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasComment("foo").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// external change
			{
				PreConfig: func() {
					acc.TestClient().Stream.Alter(t, sdk.NewAlterStreamRequest(id).WithSetComment("bar"))
				},
				Config: config.FromModels(t, modelWithExtraFields),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasComment("foo").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// update fields
			{
				Config: config.FromModels(t, modelWithExtraFieldsModified),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasStageString(stage.ID().Name()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(stage.ID().Name()).
						HasSourceType(sdk.StreamSourceTypeStage).
						HasBaseTablesPartiallyQualified(stage.ID().Name()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeDefault).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasComment("bar").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "bar")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// import
			{
				Config:       config.FromModels(t, modelWithExtraFieldsModified),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedStreamOnDirectoryTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasStageString(stage.ID().Name()).
						HasCommentString("bar"),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_CopyGrants(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_directory_table.test"

	var createdOn string

	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	model := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, model.WithCopyGrants(true)),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			{
				Config: config.FromModels(t, model.WithCopyGrants(false)),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("stream was recreated")
						}
						return nil
					})),
				),
			},
			{
				Config: config.FromModels(t, model.WithCopyGrants(true)),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("stream was recreated")
						}
						return nil
					})),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_CheckGrantsAfterRecreation(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_directory_table.test"

	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	stage2, cleanupStage2 := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage2)

	role, cleanupRole := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(cleanupRole)

	model1 := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName()).WithCopyGrants(true)
	model1WithoutCopyGrants := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName())
	model2 := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage2.ID().FullyQualifiedName()).WithCopyGrants(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, model1) + grantStreamPrivilegesConfig(resourceName, role.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: config.FromModels(t, model2) + grantStreamPrivilegesConfig(resourceName, role.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config:             config.FromModels(t, model1WithoutCopyGrants) + grantStreamPrivilegesConfig(resourceName, role.ID()),
				ExpectNonEmptyPlan: true,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_grant_privileges_to_account_role.grant", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "1"),
				),
			},
		},
	})
}

func grantStreamPrivilegesConfig(resourceName string, roleId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(`
resource "snowflake_grant_privileges_to_account_role" "grant" {
  privileges        = ["SELECT"]
  account_role_name = %[1]s
  on_schema_object {
    object_type = "STREAM"
    object_name = %[2]s.fully_qualified_name
  }
}

data "snowflake_grants" "grants" {
  depends_on = [snowflake_grant_privileges_to_account_role.grant, %[2]s]
  grants_on {
    object_type = "STREAM"
    object_name = %[2]s.fully_qualified_name
  }
}`, roleId.FullyQualifiedName(), resourceName)
}

// TODO (SNOW-1737932): Setting schema parameters related to retention time seems to have no affect on streams on directory tables.
// Adjust this test after this is fixed on Snowflake side.
func TestAcc_StreamOnDirectoryTable_RecreateWhenStale(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	resourceName := "snowflake_stream_on_directory_table.test"

	schema, cleanupSchema := acc.TestClient().Schema.CreateSchemaWithOpts(t,
		acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(acc.TestClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(0),
			MaxDataExtensionTimeInDays: sdk.Pointer(0),
		},
	)
	t.Cleanup(cleanupSchema)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	model := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, model),
				Check: assert.AssertThat(t, resourceassert.StreamOnDirectoryTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "false")),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_InvalidConfiguration(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	modelWithInvalidStageId := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), "invalid")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// invalid stage id
			{
				Config:      config.FromModels(t, modelWithInvalidStageId),
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable_ExternalStreamTypeChange(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	acc.TestAccPreCheck(t)
	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)
	model := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnDirectoryTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, model),
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.StreamOnDirectoryTableResource(t, model.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeStage)),
						resourceshowoutputassert.StreamShowOutput(t, model.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeStage),
					),
				),
			},
			// external change with a different type
			{
				PreConfig: func() {
					table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
					t.Cleanup(cleanupTable)
					acc.TestClient().Stream.DropFunc(t, id)()
					externalChangeStream, cleanup := acc.TestClient().Stream.CreateOnTableWithRequest(t, sdk.NewCreateOnTableStreamRequest(id, table.ID()))
					t.Cleanup(cleanup)
					require.Equal(t, sdk.StreamSourceTypeTable, *externalChangeStream.SourceType)
				},
				Config: config.FromModels(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(model.ResourceReference(), plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					assert.AssertThat(t,
						resourceassert.StreamOnDirectoryTableResource(t, model.ResourceReference()).
							HasStreamTypeString(string(sdk.StreamSourceTypeStage)),
						resourceshowoutputassert.StreamShowOutput(t, model.ResourceReference()).
							HasSourceType(sdk.StreamSourceTypeStage),
					),
				),
			},
		},
	})
}
