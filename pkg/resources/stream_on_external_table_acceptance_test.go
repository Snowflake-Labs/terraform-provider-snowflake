package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	r "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	tfjson "github.com/hashicorp/terraform-json"
	pluginconfig "github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_StreamOnExternalTable_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	resourceName := "snowflake_stream_on_external_table.test"

	stageID := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t, stageID)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateWithLocation(t, stageLocation)
	t.Cleanup(externalTableCleanup)

	var createdOn string

	baseModel := model.StreamOnExternalTableBase("test", id, externalTable.ID())

	modelWithExtraFields := model.StreamOnExternalTableBase("test", id, externalTable.ID()).
		WithCopyGrants(true).
		WithComment("foo").
		WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset": pluginconfig.StringVariable("0"),
		}))

	modelWithExtraFieldsModified := model.StreamOnExternalTableBase("test", id, externalTable.ID()).
		WithCopyGrants(true).
		WithComment("bar").
		WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset": pluginconfig.StringVariable("0"),
		}))

	modelWithExtraFieldsModifiedCauseRecreation := model.StreamOnExternalTableBase("test", id, externalTable.ID()).
		WithCopyGrants(true).
		WithComment("bar").
		WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset": pluginconfig.StringVariable("0"),
		}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: config.FromModel(t, baseModel),
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			// import without optionals
			{
				Config:       config.FromModel(t, baseModel),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedStreamOnExternalTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasInsertOnlyString(r.BooleanTrue).
						HasExternalTableString(externalTable.ID().FullyQualifiedName()),
				),
			},
			// set all fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFields),
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// external change
			{
				PreConfig: func() {
					acc.TestClient().Stream.Alter(t, sdk.NewAlterStreamRequest(id).WithSetComment("bar"))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFields),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// update fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFieldsModified),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// update fields to force recreation
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFieldsModifiedCauseRecreation),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("view was not recreated")
						}
						return nil
					})),
				),
			},
			// import
			{
				Config:       config.FromModel(t, modelWithExtraFieldsModified),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedStreamOnExternalTableResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasInsertOnlyString(r.BooleanTrue).
						HasExternalTableString(externalTable.ID().FullyQualifiedName()).
						HasCommentString("bar"),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_CopyGrants(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_external_table.test"

	var createdOn string

	stageID := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t, stageID)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateWithLocation(t, stageLocation)
	t.Cleanup(externalTableCleanup)

	model := model.StreamOnExternalTable("test", id.DatabaseName(), externalTable.ID().FullyQualifiedName(), id.Name(), id.SchemaName()).WithInsertOnly(r.BooleanTrue)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model.WithCopyGrants(true)),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			{
				Config: config.FromModel(t, model.WithCopyGrants(false)),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					})),
				),
			},
			{
				Config: config.FromModel(t, model.WithCopyGrants(true)),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value != createdOn {
							return fmt.Errorf("view was recreated")
						}
						return nil
					})),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable_RecreateWhenStale(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	resourceName := "snowflake_stream_on_external_table.test"

	schema, cleanupSchema := acc.TestClient().Schema.CreateSchemaWithOpts(t,
		acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(acc.TestClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(0),
			MaxDataExtensionTimeInDays: sdk.Pointer(0),
		},
	)
	t.Cleanup(cleanupSchema)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	stageID := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t, stageID)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateInSchemaWithLocation(t, stageLocation, schema.ID())
	t.Cleanup(externalTableCleanup)

	var createdOn string

	model := model.StreamOnExternalTable("test", id.DatabaseName(), externalTable.ID().FullyQualifiedName(), id.Name(), id.SchemaName()).WithInsertOnly(r.BooleanTrue)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			// check that stale state is marked properly and forces an update
			{
				Config: config.FromModel(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceName, "stale", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanTrue),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "true")),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			// check that the resource was recreated
			// note that it is stale again because we still have schema parameters set to 0
			{
				Config: config.FromModel(t, model),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceName, "stale", tfjson.ActionUpdate, sdk.String(r.BooleanTrue), sdk.String(r.BooleanFalse)),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanTrue),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "true")),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("view was not recreated")
						}
						return nil
					})),
				),
			},
			// set schema parameters to bigger values ensuring that the stream is not stale
			{
				PreConfig: func() {
					acc.TestClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
						Set: &sdk.SchemaSet{
							DataRetentionTimeInDays:    sdk.Int(1),
							MaxDataExtensionTimeInDays: sdk.Int(1),
						},
					})
				},
				RefreshState: true,
				RefreshPlanChecks: resource.RefreshPlanChecks{
					PostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionNoop),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "false")),
				),
			},
		},
	})
}

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnExternalTable_At(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_external_table.test"

	stageID := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t, stageID)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateWithLocation(t, stageLocation)
	t.Cleanup(externalTableCleanup)

	commonModel := func() *model.StreamOnExternalTableModel {
		return model.StreamOnExternalTableBase("test", id, externalTable.ID()).
			WithComment("foo").
			WithInsertOnly(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := commonModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps and statements
		},
	})
}

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnExternalTable_Before(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_external_table.test"

	stageID := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t, stageID)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateWithLocation(t, stageLocation)
	t.Cleanup(externalTableCleanup)

	commonModel := func() *model.StreamOnExternalTableModel {
		return model.StreamOnExternalTableBase("test", id, externalTable.ID()).
			WithComment("foo").
			WithInsertOnly(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := commonModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assert.AssertThat(t, resourceassert.StreamOnExternalTableResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasExternalTableString(externalTable.ID().FullyQualifiedName()).
					HasInsertOnlyString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables([]sdk.SchemaObjectIdentifier{externalTable.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assert.AssertThat(t, resourceassert.StreamOnTableResource(t, resourceName).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps and statements
		},
	})
}

func TestAcc_StreamOnExternalTable_InvalidConfiguration(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	modelWithInvalidExternalTableId := model.StreamOnExternalTable("test", id.DatabaseName(), "invalid", id.Name(), id.SchemaName())

	modelWithBefore := model.StreamOnExternalTable("test", id.DatabaseName(), "foo.bar.hoge", id.Name(), id.SchemaName()).
		WithComment("foo").
		WithCopyGrants(false).
		WithInsertOnly(r.BooleanTrue).
		WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset":    pluginconfig.StringVariable("0"),
			"timestamp": pluginconfig.StringVariable("0"),
			"statement": pluginconfig.StringVariable("0"),
			"stream":    pluginconfig.StringVariable("0"),
		}))

	modelWithAt := model.StreamOnExternalTable("test", id.DatabaseName(), "foo.bar.hoge", id.Name(), id.SchemaName()).
		WithComment("foo").
		WithCopyGrants(false).
		WithInsertOnly(r.BooleanTrue).
		WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset":    pluginconfig.StringVariable("0"),
			"timestamp": pluginconfig.StringVariable("0"),
			"statement": pluginconfig.StringVariable("0"),
			"stream":    pluginconfig.StringVariable("0"),
		}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			// multiple excluding options - before
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithBefore),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// multiple excluding options - at
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnExternalTable/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithAt),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// invalid external table id
			{
				Config:      config.FromModel(t, modelWithInvalidExternalTableId),
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}
