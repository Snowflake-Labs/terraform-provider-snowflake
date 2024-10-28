package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
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

func TestAcc_StreamOnView_Basic(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	resourceName := "snowflake_stream_on_view.test"

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)
	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	baseModel := func() *model.StreamOnViewModel {
		return model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName())
	}

	modelWithExtraFields := baseModel().
		WithCopyGrants(false).
		WithComment("foo").
		WithAppendOnly(r.BooleanTrue).
		WithShowInitialRows(r.BooleanTrue).
		WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset": pluginconfig.StringVariable("0"),
		}))

	modelWithExtraFieldsDefaultMode := baseModel().
		WithCopyGrants(false).
		WithComment("foo").
		WithAppendOnly(r.BooleanFalse).
		WithShowInitialRows(r.BooleanTrue).
		WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset": pluginconfig.StringVariable("0"),
		}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnView),
		Steps: []resource.TestStep{
			// without optionals
			{
				Config: config.FromModel(t, baseModel()),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanDefault).
					HasViewString(view.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables([]sdk.SchemaObjectIdentifier{table.ID()}).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
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
					resourceassert.ImportedStreamOnViewResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAppendOnlyString(r.BooleanFalse).
						HasViewString(view.ID().FullyQualifiedName()),
				),
			},
			// set all fields
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFields),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanTrue).
					HasViewString(view.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables([]sdk.SchemaObjectIdentifier{table.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// external change
			{
				PreConfig: func() {
					acc.TestClient().Stream.Alter(t, sdk.NewAlterStreamRequest(id).WithSetComment("bar"))
				},
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFields),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanTrue).
					HasViewString(view.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables([]sdk.SchemaObjectIdentifier{table.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// update fields that recreate the object
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithExtraFieldsDefaultMode),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
					},
				},
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanFalse).
					HasViewString(view.ID().FullyQualifiedName()),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables([]sdk.SchemaObjectIdentifier{table.ID()}).
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
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			// import
			{
				Config:       config.FromModel(t, modelWithExtraFieldsDefaultMode),
				ResourceName: resourceName,
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedStreamOnViewResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAppendOnlyString(r.BooleanFalse).
						HasViewString(view.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_StreamOnView_CopyGrants(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_view.test"

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)
	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	var createdOn string

	model := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnView),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model.WithCopyGrants(false)),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			{
				Config: config.FromModel(t, model.WithCopyGrants(true)),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
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
				Config: config.FromModel(t, model.WithCopyGrants(false)),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
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

func TestAcc_StreamOnView_CheckGrantsAfterRecreation(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_view.test"

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)
	view2, cleanupView2 := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView2)

	role, cleanupRole := acc.TestClient().Role.CreateRole(t)
	t.Cleanup(cleanupRole)

	model1 := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName()).
		WithCopyGrants(true)
	model1WithoutCopyGrants := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName())
	model2 := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view2.ID().FullyQualifiedName()).
		WithCopyGrants(true)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnExternalTable),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model1) + grantStreamPrivilegesConfig(resourceName, role.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					// there should be more than one privilege, because we applied grant all privileges and initially there's always one which is ownership
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config: config.FromModel(t, model2) + grantStreamPrivilegesConfig(resourceName, role.ID()),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.#", "2"),
					resource.TestCheckResourceAttr("data.snowflake_grants.grants", "grants.1.privilege", "SELECT"),
				),
			},
			{
				Config:             config.FromModel(t, model1WithoutCopyGrants) + grantStreamPrivilegesConfig(resourceName, role.ID()),
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

func TestAcc_StreamOnView_PermadiffWhenIsStaleAndHasNoRetentionTime(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	resourceName := "snowflake_stream_on_view.test"

	schema, cleanupSchema := acc.TestClient().Schema.CreateSchemaWithOpts(t,
		acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(acc.TestClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(0),
			MaxDataExtensionTimeInDays: sdk.Pointer(0),
		},
	)
	t.Cleanup(cleanupSchema)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTrackingInSchema(t, schema.ID())
	t.Cleanup(cleanupTable)

	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	var createdOn string

	model := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnView),
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
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
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
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanTrue),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "true")),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						if value == createdOn {
							return fmt.Errorf("stream was not recreated")
						}
						return nil
					})),
				),
			},
		},
	})
}

func TestAcc_StreamOnView_StaleWithExternalChanges(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	resourceName := "snowflake_stream_on_view.test"

	schema, cleanupSchema := acc.TestClient().Schema.CreateSchemaWithOpts(t,
		acc.TestClient().Ids.RandomDatabaseObjectIdentifierInDatabase(acc.TestClient().Ids.DatabaseId()),
		&sdk.CreateSchemaOptions{
			DataRetentionTimeInDays:    sdk.Pointer(1),
			MaxDataExtensionTimeInDays: sdk.Pointer(1),
		},
	)
	t.Cleanup(cleanupSchema)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifierInSchema(schema.ID())

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTrackingInSchema(t, schema.ID())
	t.Cleanup(cleanupTable)

	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateViewInSchema(t, statement, schema.ID())
	t.Cleanup(cleanupView)

	var createdOn string

	model := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName())
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnView),
		Steps: []resource.TestStep{
			// initial creation does not lead to stale stream
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttrWith(resourceName, "show_output.0.created_on", func(value string) error {
						createdOn = value
						return nil
					})),
				),
			},
			// changing the value externally on schema
			{
				PreConfig: func() {
					acc.TestClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
						Set: &sdk.SchemaSet{
							DataRetentionTimeInDays:    sdk.Int(0),
							MaxDataExtensionTimeInDays: sdk.Int(0),
						},
					})
					assert.AssertThatObject(t, objectassert.Stream(t, id).
						HasName(id.Name()).
						HasStale(true),
					)

					acc.TestClient().Schema.Alter(t, schema.ID(), &sdk.AlterSchemaOptions{
						Set: &sdk.SchemaSet{
							DataRetentionTimeInDays:    sdk.Int(1),
							MaxDataExtensionTimeInDays: sdk.Int(1),
						},
					})
					assert.AssertThatObject(t, objectassert.Stream(t, id).
						HasName(id.Name()).
						HasStale(false),
					)
				},
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasStaleString(r.BooleanFalse),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "show_output.0.stale", "false")),
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

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnView_At(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceId := helpers.EncodeResourceIdentifier(id)
	resourceName := "snowflake_stream_on_view.test"

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)
	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	acc.TestClient().Table.InsertInt(t, table.ID())
	lastQueryId := acc.TestClient().Context.LastQueryId(t)

	baseModel := func() *model.StreamOnViewModel {
		return model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName()).
			WithComment("foo").
			WithAppendOnly(r.BooleanTrue).
			WithShowInitialRows(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := baseModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := baseModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))
	modelWithStatement := baseModel().WithAtValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"statement": pluginconfig.StringVariable(lastQueryId),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnView),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasViewString(view.ID().FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanTrue).
					HasShowInitialRowsString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables([]sdk.SchemaObjectIdentifier{table.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", "APPEND_ONLY")),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStatement),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				ResourceName:    resourceName,
				ImportState:     true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedStreamOnViewResource(t, resourceId).
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasSchemaString(id.SchemaName()).
						HasFullyQualifiedNameString(id.FullyQualifiedName()).
						HasAppendOnlyString(r.BooleanTrue).
						HasViewString(view.ID().FullyQualifiedName()),
				),
			},
		},
	})
}

// There is no way to check at/before fields in show and describe. That's why we try creating with these values, but do not assert them.
func TestAcc_StreamOnView_Before(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_view.test"

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)
	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	acc.TestClient().Table.InsertInt(t, table.ID())
	lastQueryId := acc.TestClient().Context.LastQueryId(t)

	baseModel := func() *model.StreamOnViewModel {
		return model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName()).
			WithComment("foo").
			WithAppendOnly(r.BooleanTrue).
			WithShowInitialRows(r.BooleanTrue).
			WithCopyGrants(false)
	}

	modelWithOffset := baseModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"offset": pluginconfig.StringVariable("0"),
	}))
	modelWithStream := baseModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"stream": pluginconfig.StringVariable(id.FullyQualifiedName()),
	}))
	modelWithStatement := baseModel().WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
		"statement": pluginconfig.StringVariable(lastQueryId),
	}))

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.StreamOnView),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithOffset),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()).
					HasDatabaseString(id.DatabaseName()).
					HasSchemaString(id.SchemaName()).
					HasFullyQualifiedNameString(id.FullyQualifiedName()).
					HasViewString(view.ID().FullyQualifiedName()).
					HasAppendOnlyString(r.BooleanTrue).
					HasShowInitialRowsString(r.BooleanTrue).
					HasCommentString("foo"),
					resourceshowoutputassert.StreamShowOutput(t, resourceName).
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasComment("foo").
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables([]sdk.SchemaObjectIdentifier{table.ID()}).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.mode", "APPEND_ONLY")),
					assert.Check(resource.TestCheckResourceAttrSet(resourceName, "describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(resourceName, "describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStream),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithStatement),
				Check: assert.AssertThat(t, resourceassert.StreamOnViewResource(t, resourceName).
					HasNameString(id.Name()),
				),
			},
			// TODO(SNOW-1689111): test timestamps
		},
	})
}

func TestAcc_StreamOnView_InvalidConfiguration(t *testing.T) {
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	modelWithInvalidTableId := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), "invalid")

	modelWithBefore := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), "foo.bar.hoge").
		WithComment("foo").
		WithCopyGrants(false).
		WithAppendOnly(r.BooleanFalse).
		WithShowInitialRows(r.BooleanFalse).
		WithBeforeValue(pluginconfig.MapVariable(map[string]pluginconfig.Variable{
			"offset":    pluginconfig.StringVariable("0"),
			"timestamp": pluginconfig.StringVariable("0"),
			"statement": pluginconfig.StringVariable("0"),
			"stream":    pluginconfig.StringVariable("0"),
		}))

	modelWithAt := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), "foo.bar.hoge").
		WithComment("foo").
		WithCopyGrants(false).
		WithAppendOnly(r.BooleanFalse).
		WithShowInitialRows(r.BooleanFalse).
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
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/before"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithBefore),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// multiple excluding options - at
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_StreamOnView/at"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, modelWithAt),
				ExpectError:     regexp.MustCompile("Error: Invalid combination of arguments"),
			},
			// invalid view id
			{
				Config:      config.FromModel(t, modelWithInvalidTableId),
				ExpectError: regexp.MustCompile("Error: Invalid identifier type"),
			},
		},
	})
}
