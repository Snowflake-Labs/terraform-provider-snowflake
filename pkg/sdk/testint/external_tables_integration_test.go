package testint

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalTables(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	stageID := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, "EXTERNAL_TABLE_STAGE")
	stageLocation := fmt.Sprintf("@%s", stageID.FullyQualifiedName())
	_, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t, stageID, nycWeatherDataURL)
	t.Cleanup(stageCleanup)

	tag, tagCleanup := testClientHelper().Tag.CreateTag(t)
	t.Cleanup(tagCleanup)

	defaultColumns := func() []*sdk.ExternalTableColumnRequest {
		return []*sdk.ExternalTableColumnRequest{
			sdk.NewExternalTableColumnRequest("filename", sdk.DataTypeString, "metadata$filename::string"),
			sdk.NewExternalTableColumnRequest("city", sdk.DataTypeString, "value:city:findname::string"),
			sdk.NewExternalTableColumnRequest("time", sdk.DataTypeTimestamp, "to_timestamp(value:time::int)"),
			sdk.NewExternalTableColumnRequest("weather", sdk.DataTypeVariant, "value:weather::variant"),
		}
	}

	columns := defaultColumns()
	columnsWithPartition := append(defaultColumns(), []*sdk.ExternalTableColumnRequest{
		sdk.NewExternalTableColumnRequest("weather_date", sdk.DataTypeDate, "to_date(to_timestamp(value:time::int))"),
		sdk.NewExternalTableColumnRequest("part_date", sdk.DataTypeDate, "parse_json(metadata$external_table_partition):weather_date::date"),
	}...)

	minimalCreateExternalTableReq := func(name string) *sdk.CreateExternalTableRequest {
		return sdk.NewCreateExternalTableRequest(
			sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name),
			stageLocation,
		).WithFileFormat(sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON))
	}

	createExternalTableWithManualPartitioningReq := func(name string) *sdk.CreateWithManualPartitioningExternalTableRequest {
		return sdk.NewCreateWithManualPartitioningExternalTableRequest(
			sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name),
			stageLocation,
		).
			WithFileFormat(sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON)).
			WithOrReplace(sdk.Bool(true)).
			WithColumns(columnsWithPartition).
			WithPartitionBy([]string{"part_date"}).
			WithCopyGrants(sdk.Bool(true)).
			WithComment(sdk.String("some_comment")).
			WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")})
	}

	t.Run("Create: minimal", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Create: with raw file format", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name).
			WithFileFormat(nil).
			WithRawFileFormat(sdk.String("TYPE = JSON")),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Create: complete", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(
			ctx,
			sdk.NewCreateExternalTableRequest(
				externalTableID,
				stageLocation,
			).
				WithFileFormat(sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON)).
				WithOrReplace(sdk.Bool(true)).
				WithColumns(columns).
				WithPartitionBy([]string{"filename"}).
				WithRefreshOnCreate(sdk.Bool(false)).
				WithAutoRefresh(sdk.Bool(false)).
				WithPattern(sdk.String("weather-nyc/weather_2_3_0.json.gz")).
				WithCopyGrants(sdk.Bool(true)).
				WithComment(sdk.String("some_comment")).
				WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Create: infer schema", func(t *testing.T) {
		fileFormat, _ := createFileFormat(t, client, testSchema(t).ID())

		err := client.Sessions.UseWarehouse(ctx, testWarehouse(t).ID())
		require.NoError(t, err)

		name := random.AlphanumericN(32)
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '%s', FILE_FORMAT=>'%s', ignore_case => true))`, stageLocation, fileFormat.ID().FullyQualifiedName())
		err = client.ExternalTables.CreateUsingTemplate(
			ctx,
			sdk.NewCreateExternalTableUsingTemplateRequest(
				id,
				stageLocation,
			).
				WithFileFormat(sdk.NewExternalTableFileFormatRequest().WithName(sdk.String(fileFormat.ID().FullyQualifiedName()))).
				WithQuery(query).
				WithAutoRefresh(sdk.Bool(false)))
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Create with manual partitioning: complete", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(name))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Create delta lake: complete", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.CreateDeltaLake(
			ctx,
			sdk.NewCreateDeltaLakeExternalTableRequest(
				externalTableID,
				stageLocation,
			).
				WithFileFormat(sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeParquet)).
				WithOrReplace(sdk.Bool(true)).
				WithColumns(columnsWithPartition).
				WithPartitionBy([]string{"filename"}).
				WithAutoRefresh(sdk.Bool(false)).
				WithRefreshOnCreate(sdk.Bool(false)).
				WithCopyGrants(sdk.Bool(true)).
				WithComment(sdk.String("some_comment")).
				WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, name, externalTable.Name)
	})

	t.Run("Alter: refresh", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name))
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithRefresh(sdk.NewRefreshExternalTableRequest("weather-nyc")),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: add files", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(name).
				WithPattern(sdk.String("weather-nyc/weather_2_3_0.json.gz")),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithAddFiles([]*sdk.ExternalTableFileRequest{sdk.NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: remove files", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(name).
				WithPattern(sdk.String("weather-nyc/weather_2_3_0.json.gz")),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithAddFiles([]*sdk.ExternalTableFileRequest{sdk.NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithRemoveFiles([]*sdk.ExternalTableFileRequest{sdk.NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: set auto refresh", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name))
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			sdk.NewAlterExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithAutoRefresh(sdk.Bool(true)),
		)
		require.NoError(t, err)
	})

	// TODO: (SNOW-919981) Uncomment when the problem with alter external table set / unset tags is solved
	// t.Run("Alter: set tags", func(t *testing.T) {
	//	externalTableID := randomAccountObjectIdentifier(t)
	//	err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
	//	require.NoError(t, err)
	//
	//	tagValue := "tag-value"
	//	err = client.ExternalTables.Alter(
	//		ctx,
	//		NewAlterExternalTableRequest(externalTableID).
	//			WithIfExists(Bool(true)).
	//			WithSetTag([]*TagAssociationRequest{NewTagAssociationRequest(tag.ID(), tagValue)}))
	//	require.NoError(t, err)
	//
	//	tv, err := client.SystemFunctions.GetTag(ctx, tag.ID(), externalTableID, ObjectTypeExternalTable)
	//	require.NoError(t, err)
	//	assert.Equal(t, tagValue, tv)
	// })
	//
	// t.Run("Alter: unset tags", func(t *testing.T) {
	//	externalTableID := randomAccountObjectIdentifier(t)
	//	err := client.ExternalTables.Create(
	//		ctx,
	//		minimalCreateExternalTableReq(externalTableID).
	//			WithTag([]*TagAssociationRequest{NewTagAssociationRequest(tag.ID(), "tag-value")}),
	//	)
	//	require.NoError(t, err)
	//	tv, err := client.SystemFunctions.GetTag(ctx, tag.ID(), externalTableID, ObjectTypeExternalTable)
	//	require.NoError(t, err)
	//	assert.Equal(t, "tag-value", tv)
	//
	//	err = client.ExternalTables.Alter(
	//		ctx,
	//		NewAlterExternalTableRequest(externalTableID).
	//			WithIfExists(Bool(true)).
	//			WithUnsetTag([]ObjectIdentifier{
	//				NewAccountObjectIdentifier(tag.ID().Name()),
	//			}),
	//	)
	//	require.NoError(t, err)
	//
	//	_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), externalTableID, ObjectTypeExternalTable)
	//	require.Error(t, err)
	// })

	t.Run("Alter: add partitions", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(name))
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			sdk.NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithAddPartitions([]*sdk.PartitionRequest{sdk.NewPartitionRequest("part_date", "2019-06-25")}).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: drop partitions", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(name))
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			sdk.NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithAddPartitions([]*sdk.PartitionRequest{sdk.NewPartitionRequest("part_date", "2019-06-25")}).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			sdk.NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithDropPartition(sdk.Bool(true)).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name))
		require.NoError(t, err)

		err = client.ExternalTables.Drop(
			ctx,
			sdk.NewDropExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithDropOption(sdk.NewExternalTableDropOptionRequest().WithCascade(sdk.Bool(true))),
		)
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, externalTableID)
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("Show", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name))
		require.NoError(t, err)

		et, err := client.ExternalTables.Show(
			ctx,
			sdk.NewShowExternalTableRequest().
				WithTerse(sdk.Bool(true)).
				WithLike(sdk.String(name)).
				WithIn(sdk.NewShowExternalTableInRequest().WithDatabase(testDb(t).ID())).
				WithStartsWith(sdk.String(name)).
				WithLimitFrom(sdk.NewLimitFromRequest().WithRows(sdk.Int(1))),
		)
		require.NoError(t, err)
		assert.Equal(t, 1, len(et))
		assert.Equal(t, externalTableID, et[0].ID())
	})

	t.Run("Describe: columns", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		req := minimalCreateExternalTableReq(name)
		err := client.ExternalTables.Create(ctx, req)
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeColumns(ctx, sdk.NewDescribeExternalTableColumnsRequest(externalTableID))
		require.NoError(t, err)

		assert.Equal(t, len(req.GetColumns())+1, len(d)) // +1 because there's underlying Value column
		assert.Contains(t, d, sdk.ExternalTableColumnDetails{
			Name:       "VALUE",
			Type:       "VARIANT",
			Kind:       "COLUMN",
			IsNullable: true,
			Default:    nil,
			IsPrimary:  false,
			IsUnique:   false,
			Check:      nil,
			Expression: nil,
			Comment:    sdk.String("The value of this row"),
			PolicyName: nil,
		})
	})

	t.Run("Describe: stage", func(t *testing.T) {
		name := random.AlphanumericN(32)
		externalTableID := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(name))
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeStage(ctx, sdk.NewDescribeExternalTableStageRequest(externalTableID))
		require.NoError(t, err)

		assert.Contains(t, d, sdk.ExternalTableStageDetails{
			ParentProperty:  "STAGE_FILE_FORMAT",
			Property:        "TIME_FORMAT",
			PropertyType:    "String",
			PropertyValue:   "AUTO",
			PropertyDefault: "AUTO",
		})
	})
}

func TestInt_ExternalTablesShowByID(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	databaseTest, schemaTest := testDb(t), testSchema(t)
	stage := sdk.NewSchemaObjectIdentifier(TestDatabaseName, TestSchemaName, random.AlphaN(6))
	_, stageCleanup := testClientHelper().Stage.CreateStageWithURL(t, stage, nycWeatherDataURL)
	t.Cleanup(stageCleanup)

	stageLocation := fmt.Sprintf("@%s", stage.FullyQualifiedName())
	cleanupExternalTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) func() {
		t.Helper()
		return func() {
			err := client.ExternalTables.Drop(ctx, sdk.NewDropExternalTableRequest(id))
			if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
				return
			}
			require.NoError(t, err)
		}
	}

	createExternalTableHandle := func(t *testing.T, id sdk.SchemaObjectIdentifier) {
		t.Helper()

		request := sdk.NewCreateExternalTableRequest(id, stageLocation).WithFileFormat(sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON))
		err := client.ExternalTables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupExternalTableHandle(t, id))
	}

	t.Run("show by id - same name in different schemas", func(t *testing.T) {
		schema, schemaCleanup := testClientHelper().Schema.CreateSchema(t)
		t.Cleanup(schemaCleanup)

		name := random.AlphaN(4)
		id1 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schemaTest.Name, name)
		id2 := sdk.NewSchemaObjectIdentifier(databaseTest.Name, schema.Name, name)

		createExternalTableHandle(t, id1)
		createExternalTableHandle(t, id2)

		e1, err := client.ExternalTables.ShowByID(ctx, id1)
		require.NoError(t, err)
		require.Equal(t, id1, e1.ID())

		e2, err := client.ExternalTables.ShowByID(ctx, id2)
		require.NoError(t, err)
		require.Equal(t, id2, e2.ID())
	})
}
