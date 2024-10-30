package datasources_test

import (
	"fmt"
	"regexp"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	testconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/internal/snowflakeroles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_Streams(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	streamOnTable := model.StreamOnTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), table.ID().FullyQualifiedName()).
		WithAppendOnly("true").
		WithComment("foo")

	dsName := "data.snowflake_streams.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streams/optionals_set"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, streamOnTable),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.#", "1")),

					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.table_name", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streams/optionals_unset"),
				ConfigVariables: tfconfig.ConfigVariablesFromModel(t, streamOnTable),

				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.#", "1")),

					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.#", "0")),
				),
			},
		},
	})
}

func TestAcc_StreamOnTable(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	streamOnTable := model.StreamOnTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), table.ID().FullyQualifiedName()).
		WithAppendOnly("true").
		WithComment("foo")

	resourceName := "snowflake_stream_on_table.test"
	dsName := "data.snowflake_streams.test"
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, streamOnTable) + streamsDatasource(id.Name(), resourceName),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(table.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeTable).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.table_name", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeTable))),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_StreamOnExternalTable(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_external_table.test"
	dsName := "data.snowflake_streams.test"

	stageID := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := acc.TestClient().Stage.CreateStageWithURL(t, stageID)
	t.Cleanup(stageCleanup)

	externalTable, externalTableCleanup := acc.TestClient().ExternalTable.CreateWithLocation(t, stageLocation)
	t.Cleanup(externalTableCleanup)

	model := model.StreamOnExternalTableBase("test", id, externalTable.ID()).
		WithCopyGrants(true).
		WithComment("foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model) + streamsDatasource(id.Name(), resourceName),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(externalTable.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeExternalTable).
						HasBaseTables(externalTable.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeInsertOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.table_name", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeExternalTable))),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.0", externalTable.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.mode", string(sdk.StreamModeInsertOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_StreamOnDirectoryTable(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_directory_table.test"
	dsName := "data.snowflake_streams.test"

	stage, cleanupStage := acc.TestClient().Stage.CreateStageWithDirectory(t)
	t.Cleanup(cleanupStage)

	model := model.StreamOnDirectoryTable("test", id.DatabaseName(), id.Name(), id.SchemaName(), stage.ID().FullyQualifiedName()).
		WithComment("foo")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model) + streamsDatasource(id.Name(), resourceName),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
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
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.table_name", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeStage))),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.0", stage.ID().Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.mode", string(sdk.StreamModeDefault))),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func TestAcc_StreamOnView(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)
	id := acc.TestClient().Ids.RandomSchemaObjectIdentifier()
	resourceName := "snowflake_stream_on_view.test"
	dsName := "data.snowflake_streams.test"

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)
	statement := fmt.Sprintf("SELECT * FROM %s", table.ID().FullyQualifiedName())
	view, cleanupView := acc.TestClient().View.CreateView(t, statement)
	t.Cleanup(cleanupView)

	model := model.StreamOnView("test", id.DatabaseName(), id.Name(), id.SchemaName(), view.ID().FullyQualifiedName()).
		WithComment("foo").
		WithAppendOnly("true")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model) + streamsDatasource(id.Name(), resourceName),
				Check: assert.AssertThat(t,
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.#", "1")),
					resourceshowoutputassert.StreamsDatasourceShowOutput(t, "snowflake_streams.test").
						HasCreatedOnNotEmpty().
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasOwner(snowflakeroles.Accountadmin.Name()).
						HasTableName(view.ID().FullyQualifiedName()).
						HasSourceType(sdk.StreamSourceTypeView).
						HasBaseTables(table.ID()).
						HasType("DELTA").
						HasStale("false").
						HasMode(sdk.StreamModeAppendOnly).
						HasStaleAfterNotEmpty().
						HasInvalidReason("N/A").
						HasOwnerRoleType("ROLE"),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.created_on")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.name", id.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.database_name", id.DatabaseName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.schema_name", id.SchemaName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner", snowflakeroles.Accountadmin.Name())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.comment", "foo")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.table_name", view.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.source_type", string(sdk.StreamSourceTypeView))),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.#", "1")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.base_tables.0", table.ID().FullyQualifiedName())),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.type", "DELTA")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.stale", "false")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.mode", string(sdk.StreamModeAppendOnly))),
					assert.Check(resource.TestCheckResourceAttrSet(dsName, "streams.0.describe_output.0.stale_after")),
					assert.Check(resource.TestCheckResourceAttr(dsName, "streams.0.describe_output.0.owner_role_type", "ROLE")),
				),
			},
		},
	})
}

func streamsDatasource(like, resourceName string) string {
	return fmt.Sprintf(`
data "snowflake_streams" "test" {
	depends_on = [%s]

	like = "%s"
}
`, resourceName, like)
}

func TestAcc_Streams_Filtering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	prefix := random.AlphaN(4)
	id1 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id2 := acc.TestClient().Ids.RandomSchemaObjectIdentifierWithPrefix(prefix)
	id3 := acc.TestClient().Ids.RandomSchemaObjectIdentifier()

	table, cleanupTable := acc.TestClient().Table.CreateWithChangeTracking(t)
	t.Cleanup(cleanupTable)

	model1 := model.StreamOnTable("test_1", id1.DatabaseName(), id1.Name(), id1.SchemaName(), table.ID().FullyQualifiedName())
	model2 := model.StreamOnTable("test_2", id2.DatabaseName(), id2.Name(), id2.SchemaName(), table.ID().FullyQualifiedName())
	model3 := model.StreamOnTable("test_3", id3.DatabaseName(), id3.Name(), id3.SchemaName(), table.ID().FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck: func() { acc.TestAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testconfig.FromModel(t, model1) + testconfig.FromModel(t, model2) + testconfig.FromModel(t, model3) + streamsDatasourceLike(id1.Name()),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streams.test", "streams.#", "1"),
				),
			},
			{
				Config: testconfig.FromModel(t, model1) + testconfig.FromModel(t, model2) + testconfig.FromModel(t, model3) + streamsDatasourceLike(prefix+"%"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_streams.test", "streams.#", "2"),
				),
			},
		},
	})
}

func streamsDatasourceLike(like string) string {
	return fmt.Sprintf(`
data "snowflake_streams" "test" {
	depends_on = [snowflake_stream_on_table.test_1, snowflake_stream_on_table.test_2, snowflake_stream_on_table.test_3]

	like = "%s"
}
`, like)
}

func TestAcc_Streams_emptyIn(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config:      streamsDatasourceEmptyIn(),
				ExpectError: regexp.MustCompile("Invalid combination of arguments"),
			},
		},
	})
}

func streamsDatasourceEmptyIn() string {
	return `
data "snowflake_streams" "test" {
  in {
  }
}
`
}

func TestAcc_Streams_NotFound_WithPostConditions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_Streams/non_existing"),
				ExpectError:     regexp.MustCompile("there should be at least one stream"),
			},
		},
	})
}
