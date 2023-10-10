package sdk_integration_tests

import (
	"context"
	"fmt"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalTables(t *testing.T) {
	client := sdk.testClient(t)
	ctx := context.Background()

	db, cleanupDB := sdk.createDatabase(t, client)
	t.Cleanup(cleanupDB)

	schema, _ := sdk.createSchema(t, client, db)

	err := client.Sessions.UseDatabase(ctx, db.ID())
	require.NoError(t, err)
	err = client.Sessions.UseSchema(ctx, schema.ID())
	require.NoError(t, err)

	stageID := sdk.NewAccountObjectIdentifier("EXTERNAL_TABLE_STAGE")
	stageLocation := "@external_table_stage"
	_, _ = sdk.createStageWithURL(t, client, stageID, "s3://snowflake-workshop-lab/weather-nyc")

	tag, _ := sdk.createTag(t, client, db, schema)

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

	minimalCreateExternalTableReq := func(id sdk.AccountObjectIdentifier) *sdk.CreateExternalTableRequest {
		return sdk.NewCreateExternalTableRequest(
			id,
			stageLocation,
			sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON),
		)
	}

	createExternalTableWithManualPartitioningReq := func(id sdk.AccountObjectIdentifier) *sdk.CreateWithManualPartitioningExternalTableRequest {
		return sdk.NewCreateWithManualPartitioningExternalTableRequest(
			id,
			stageLocation,
			sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON),
		).
			WithOrReplace(sdk.Bool(true)).
			WithColumns(columnsWithPartition).
			WithUserSpecifiedPartitionType(sdk.Bool(true)).
			WithPartitionBy([]string{"part_date"}).
			WithCopyGrants(sdk.Bool(true)).
			WithComment(sdk.String("some_comment")).
			WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")})
	}

	t.Run("Create: minimal", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: complete", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(
			ctx,
			sdk.NewCreateExternalTableRequest(
				externalTableID,
				stageLocation,
				sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeJSON),
			).
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

		externalTable, err := client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: infer schema", func(t *testing.T) {
		fileFormat, _ := sdk.createFileFormat(t, client, schema.ID())
		warehouse, warehouseCleanup := sdk.createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		err = client.Sessions.UseWarehouse(ctx, warehouse.ID())
		require.NoError(t, err)

		id := sdk.randomAccountObjectIdentifier(t)
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '%s', FILE_FORMAT=>'%s', ignore_case => true))`, stageLocation, fileFormat.ID().FullyQualifiedName())
		err = client.ExternalTables.CreateUsingTemplate(
			ctx,
			sdk.NewCreateExternalTableUsingTemplateRequest(
				id,
				stageLocation,
				sdk.NewExternalTableFileFormatRequest().WithName(sdk.String(fileFormat.ID().FullyQualifiedName())),
			).
				WithQuery(query).
				WithAutoRefresh(sdk.Bool(false)))
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(id))
		require.NoError(t, err)
	})

	t.Run("Create with manual partitioning: complete", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create delta lake: complete", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateDeltaLake(
			ctx,
			sdk.NewCreateDeltaLakeExternalTableRequest(
				externalTableID,
				stageLocation,
				sdk.NewExternalTableFileFormatRequest().WithFileFormatType(&sdk.ExternalTableFileFormatTypeParquet),
			).
				WithOrReplace(sdk.Bool(true)).
				WithColumns(columnsWithPartition).
				WithPartitionBy([]string{"filename"}).
				WithDeltaTableFormat(sdk.Bool(true)).
				WithAutoRefresh(sdk.Bool(false)).
				WithRefreshOnCreate(sdk.Bool(false)).
				WithCopyGrants(sdk.Bool(true)).
				WithComment(sdk.String("some_comment")).
				WithTag([]*sdk.TagAssociationRequest{sdk.NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Alter: refresh", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(externalTableID).
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(externalTableID).
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Drop(
			ctx,
			sdk.NewDropExternalTableRequest(externalTableID).
				WithIfExists(sdk.Bool(true)).
				WithDropOption(sdk.NewExternalTableDropOptionRequest().WithCascade(sdk.Bool(true))),
		)
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, sdk.NewShowExternalTableByIDRequest(externalTableID))
		require.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		et, err := client.ExternalTables.Show(
			ctx,
			sdk.NewShowExternalTableRequest().
				WithTerse(sdk.Bool(true)).
				WithLike(sdk.String(externalTableID.Name())).
				WithIn(sdk.NewShowExternalTableInRequest().WithDatabase(db.ID())).
				WithStartsWith(sdk.String(externalTableID.Name())).
				WithLimitFrom(sdk.NewLimitFromRequest().WithRows(sdk.Int(1))),
		)
		require.NoError(t, err)
		assert.Equal(t, 1, len(et))
		assert.Equal(t, externalTableID, et[0].ID())
	})

	t.Run("Describe: columns", func(t *testing.T) {
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		req := minimalCreateExternalTableReq(externalTableID)
		err := client.ExternalTables.Create(ctx, req)
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeColumns(ctx, sdk.NewDescribeExternalTableColumnsRequest(externalTableID))
		require.NoError(t, err)

		assert.Equal(t, len(req.columns)+1, len(d)) // +1 because there's underlying Value column
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
		externalTableID := sdk.randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
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
