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
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
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
				Config: config.FromModel(t, baseModel()),
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
				Config:       config.FromModel(t, baseModel()),
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
				Config: config.FromModel(t, modelWithExtraFields),
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
				Config: config.FromModel(t, modelWithExtraFields),
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
				Config: config.FromModel(t, modelWithExtraFieldsModified),
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
				Config:       config.FromModel(t, modelWithExtraFieldsModified),
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
				Config:      config.FromModel(t, modelWithInvalidStageId),
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}
