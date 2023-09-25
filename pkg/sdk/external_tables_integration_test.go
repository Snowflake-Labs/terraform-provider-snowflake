package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_ExternalTables(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	db, cleanupDB := createDatabase(t, client)
	t.Cleanup(cleanupDB)

	schema, _ := createSchema(t, client, db)

	err := client.Sessions.UseDatabase(ctx, db.ID())
	require.NoError(t, err)
	err = client.Sessions.UseSchema(ctx, schema.ID())
	require.NoError(t, err)

	stageID := NewAccountObjectIdentifier("EXTERNAL_TABLE_STAGE")
	stageLocation := "@external_table_stage"
	_, _ = createStageWithURL(t, client, stageID, "s3://snowflake-workshop-lab/weather-nyc")

	tag, _ := createTag(t, client, db, schema)

	columns := []*ExternalTableColumnRequest{
		NewExternalTableColumnRequest("filename", DataTypeString, "metadata$filename::string"),
		NewExternalTableColumnRequest("city", DataTypeString, "value:city:findname::string"),
		NewExternalTableColumnRequest("time", DataTypeTimestamp, "to_timestamp(value:time::int)"),
		NewExternalTableColumnRequest("weather", DataTypeVariant, "value:weather::variant"),
	}

	var columnsWithPartition []*ExternalTableColumnRequest
	copy(columnsWithPartition, columns)
	columnsWithPartition = append(columnsWithPartition, []*ExternalTableColumnRequest{
		NewExternalTableColumnRequest("weather_date", DataTypeDate, "to_date(to_timestamp(value:time::int))"),
		NewExternalTableColumnRequest("part_date", DataTypeDate, "parse_json(metadata$external_table_partition):weather_date::date"),
	}...)

	minimalCreateExternalTableReq := func(id AccountObjectIdentifier) *CreateExternalTableRequest {
		return NewCreateExternalTableRequest(
			id,
			stageLocation,
			NewExternalTableFileFormatRequest().WithFileFormatType(&ExternalTableFileFormatTypeJSON),
		)
	}

	createExternalTableWithManualPartitioningReq := func(id AccountObjectIdentifier) *CreateWithManualPartitioningExternalTableRequest {
		return NewCreateWithManualPartitioningExternalTableRequest(
			id,
			stageLocation,
			NewExternalTableFileFormatRequest().WithFileFormatType(&ExternalTableFileFormatTypeJSON),
		).
			WithOrReplace(Bool(true)).
			WithColumns(columnsWithPartition).
			WithUserSpecifiedPartitionType(Bool(true)).
			WithPartitionBy([]string{"part_date"}).
			WithCopyGrants(Bool(true)).
			WithComment(String("some_comment")).
			WithTag([]*TagAssociationRequest{NewTagAssociationRequest(tag.ID(), "tag-value")})
	}

	t.Run("Create: minimal", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: complete", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(
			ctx,
			NewCreateExternalTableRequest(
				externalTableID,
				stageLocation,
				NewExternalTableFileFormatRequest().WithFileFormatType(&ExternalTableFileFormatTypeJSON),
			).
				WithOrReplace(Bool(true)).
				WithColumns(columns).
				WithPartitionBy([]string{"filename"}).
				WithRefreshOnCreate(Bool(false)).
				WithAutoRefresh(Bool(false)).
				WithPattern(String("weather-nyc/weather_2_3_0.json.gz")).
				WithCopyGrants(Bool(true)).
				WithComment(String("some_comment")).
				WithTag([]*TagAssociationRequest{NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: infer schema", func(t *testing.T) {
		fileFormat, _ := createFileFormat(t, client, schema.ID())
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		err = client.Sessions.UseWarehouse(ctx, warehouse.ID())
		require.NoError(t, err)

		id := randomAccountObjectIdentifier(t)
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '%s', FILE_FORMAT=>'%s', ignore_case => true))`, stageLocation, fileFormat.ID().FullyQualifiedName())
		err = client.ExternalTables.CreateUsingTemplate(
			ctx,
			NewCreateExternalTableUsingTemplateRequest(
				id,
				stageLocation,
				NewExternalTableFileFormatRequest().WithName(String(fileFormat.ID().FullyQualifiedName())),
			).
				WithQuery(query).
				WithAutoRefresh(Bool(false)))
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, NewShowExternalTableByIDRequest(id))
		require.NoError(t, err)
	})

	t.Run("Create with manual partitioning: complete", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create delta lake: complete", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateDeltaLake(
			ctx,
			NewCreateDeltaLakeExternalTableRequest(
				externalTableID,
				stageLocation,
				NewExternalTableFileFormatRequest().WithFileFormatType(&ExternalTableFileFormatTypeParquet),
			).
				WithOrReplace(Bool(true)).
				WithColumns(columnsWithPartition).
				WithPartitionBy([]string{"filename"}).
				WithDeltaTableFormat(Bool(true)).
				WithAutoRefresh(Bool(false)).
				WithRefreshOnCreate(Bool(false)).
				WithCopyGrants(Bool(true)).
				WithComment(String("some_comment")).
				WithTag([]*TagAssociationRequest{NewTagAssociationRequest(tag.ID(), "tag-value")}),
		)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, NewShowExternalTableByIDRequest(externalTableID))
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Alter: refresh", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			NewAlterExternalTableRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithRefresh(NewRefreshExternalTableRequest("weather-nyc")),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: add files", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(externalTableID).
				WithPattern(String("weather-nyc/weather_2_3_0.json.gz")),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			NewAlterExternalTableRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithAddFiles([]*ExternalTableFileRequest{NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: remove files", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(
			ctx,
			minimalCreateExternalTableReq(externalTableID).
				WithPattern(String("weather-nyc/weather_2_3_0.json.gz")),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			NewAlterExternalTableRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithAddFiles([]*ExternalTableFileRequest{NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			NewAlterExternalTableRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithRemoveFiles([]*ExternalTableFileRequest{NewExternalTableFileRequest("weather-nyc/weather_0_0_0.json.gz")}),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: set auto refresh", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Alter(
			ctx,
			NewAlterExternalTableRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithAutoRefresh(Bool(true)),
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
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithAddPartitions([]*PartitionRequest{NewPartitionRequest("part_date", "2019-06-25")}).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)
	})

	t.Run("Alter: drop partitions", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, createExternalTableWithManualPartitioningReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithAddPartitions([]*PartitionRequest{NewPartitionRequest("part_date", "2019-06-25")}).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(
			ctx,
			NewAlterExternalTablePartitionRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithDropPartition(Bool(true)).
				WithLocation("2019/06"),
		)
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		err = client.ExternalTables.Drop(
			ctx,
			NewDropExternalTableRequest(externalTableID).
				WithIfExists(Bool(true)).
				WithDropOption(NewExternalTableDropOptionRequest().WithCascade(Bool(true))),
		)
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, NewShowExternalTableByIDRequest(externalTableID))
		require.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		et, err := client.ExternalTables.Show(
			ctx,
			NewShowExternalTableRequest().
				WithTerse(Bool(true)).
				WithLike(String(externalTableID.Name())).
				WithIn(NewShowExternalTableInRequest().WithDatabase(db.ID())).
				WithStartsWith(String(externalTableID.Name())).
				WithLimitFrom(NewLimitFromRequest().WithRows(Int(1))),
		)
		require.NoError(t, err)
		assert.Equal(t, 1, len(et))
		assert.Equal(t, externalTableID, et[0].ID())
	})

	t.Run("Describe: columns", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		req := minimalCreateExternalTableReq(externalTableID)
		err := client.ExternalTables.Create(ctx, req)
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeColumns(ctx, NewDescribeExternalTableColumnsRequest(externalTableID))
		require.NoError(t, err)

		assert.Equal(t, len(req.columns)+1, len(d)) // +1 because there's underlying Value column
		assert.Contains(t, d, ExternalTableColumnDetails{
			Name:       "VALUE",
			Type:       "VARIANT",
			Kind:       "COLUMN",
			IsNullable: true,
			Default:    nil,
			IsPrimary:  false,
			IsUnique:   false,
			Check:      nil,
			Expression: nil,
			Comment:    String("The value of this row"),
			PolicyName: nil,
		})
	})

	t.Run("Describe: stage", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, minimalCreateExternalTableReq(externalTableID))
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeStage(ctx, NewDescribeExternalTableStageRequest(externalTableID))
		require.NoError(t, err)

		assert.Contains(t, d, ExternalTableStageDetails{
			ParentProperty:  "STAGE_FILE_FORMAT",
			Property:        "TIME_FORMAT",
			PropertyType:    "String",
			PropertyValue:   "AUTO",
			PropertyDefault: "AUTO",
		})
	})
}
