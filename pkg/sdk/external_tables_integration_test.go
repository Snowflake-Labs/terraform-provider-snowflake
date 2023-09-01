package sdk

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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

	columns := []ExternalTableColumn{
		{
			Name:         "filename",
			Type:         DataTypeString,
			AsExpression: "metadata$filename::string",
		},
		{
			Name:         "city",
			Type:         DataTypeString,
			AsExpression: "value:city:findname::string",
		},
		{
			Name:         "time",
			Type:         DataTypeTimestamp,
			AsExpression: "to_timestamp(value:time::int)",
		},
		{
			Name:         "weather",
			Type:         DataTypeVariant,
			AsExpression: "value:weather::variant",
		},
	}

	columnsWithPartition := append(columns, []ExternalTableColumn{
		{
			Name:         "weather_date",
			Type:         DataTypeDate,
			AsExpression: "to_date(to_timestamp(value:time::int))",
		},
		{
			Name:         "part_date",
			Type:         DataTypeDate,
			AsExpression: "parse_json(metadata$external_table_partition):weather_date::date",
		},
	}...)

	minimalCreateExternalTableOpts := CreateExternalTableOpts{
		IfNotExists: Bool(true),
		Columns:     columns,
		Location:    stageLocation,
		FileFormat: []ExternalTableFileFormat{
			{
				Type: &ExternalTableFileFormatTypeJSON,
			},
		},
	}

	createExternalTableWithManualPartitioning := CreateWithManualPartitioningExternalTableOpts{
		OrReplace: Bool(true),
		Columns:   columnsWithPartition,
		// TODO Cloud provider params
		PartitionBy:                []string{"part_date"},
		Location:                   stageLocation,
		UserSpecifiedPartitionType: Bool(true),
		FileFormat: []ExternalTableFileFormat{
			{
				Type: &ExternalTableFileFormatTypeJSON,
			},
		},
		CopyGrants: Bool(true),
		Comment:    String("some_comment"),
		Tag: []TagAssociation{
			{
				Name:  tag.ID(),
				Value: "tag-value",
			},
		},
	}

	t.Run("Create: minimal", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create: complete", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.Create(ctx, externalTableID, &CreateExternalTableOpts{
			OrReplace:       Bool(true),
			name:            externalTableID,
			Columns:         columns,
			PartitionBy:     []string{"filename"},
			Location:        stageLocation,
			RefreshOnCreate: Bool(false),
			AutoRefresh:     Bool(false),
			Pattern:         String("weather-nyc/weather_2_3_0.json.gz"),
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJSON,
				},
			},
			CopyGrants: Bool(true),
			Comment:    String("some_comment"),
			Tag: []TagAssociation{
				{
					Name:  tag.ID(),
					Value: "tag-value",
				},
			},
		})
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
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
		err = client.ExternalTables.CreateUsingTemplate(ctx, id, &CreateExternalTableUsingTemplateOpts{
			name:     id,
			Query:    query,
			Location: stageLocation,
			FileFormat: []ExternalTableFileFormat{
				{
					Name: String(fileFormat.ID().FullyQualifiedName()),
				},
			},
			AutoRefresh: Bool(false),
		})
		require.NoError(t, err)

		_, err = client.ExternalTables.ShowByID(ctx, id)
		require.NoError(t, err)
	})

	t.Run("Create with manual partitioning: complete", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := createExternalTableWithManualPartitioning
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, externalTableID, &opts)
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Create delta lake: complete", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		err := client.ExternalTables.CreateDeltaLake(ctx, externalTableID, &CreateDeltaLakeExternalTableOpts{
			OrReplace: Bool(true),
			name:      externalTableID,
			Columns:   columnsWithPartition,
			// TODO Cloud provider params ?
			PartitionBy: []string{"filename"},
			Location:    stageLocation,
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeParquet,
				},
			},
			DeltaTableFormat: Bool(true),
			CopyGrants:       Bool(true),
			Comment:          String("some_comment"),
			Tag: []TagAssociation{
				{
					Name:  tag.ID(),
					Value: "tag-value",
				},
			},
		})
		require.NoError(t, err)

		externalTable, err := client.ExternalTables.ShowByID(ctx, externalTableID)
		require.NoError(t, err)
		assert.Equal(t, externalTableID.Name(), externalTable.Name)
	})

	t.Run("Alter: refresh", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			Refresh: &RefreshExternalTable{
				Path: "weather-nyc",
			},
		})
		require.NoError(t, err)
	})

	t.Run("Alter: add files", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		opts.Pattern = String("weather-nyc/weather_2_3_0.json.gz")
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			AddFiles: []ExternalTableFile{
				{
					Name: "weather-nyc/weather_0_0_0.json.gz",
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("Alter: remove files", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		opts.Pattern = String("weather-nyc/weather_2_3_0.json.gz")
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			AddFiles: []ExternalTableFile{
				{
					Name: "weather-nyc/weather_0_0_0.json.gz",
				},
			},
		})
		require.NoError(t, err)

		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			RemoveFiles: []ExternalTableFile{
				{
					Name: "weather-nyc/weather_0_0_0.json.gz",
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("Alter: set auto refresh", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			Set: &ExternalTableSet{
				AutoRefresh: Bool(true),
			},
		})
		require.NoError(t, err)
	})

	t.Run("Alter: set tags", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		tagValue := "tag-value"
		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			Set: &ExternalTableSet{
				Tag: []TagAssociation{
					{
						Name:  tag.ID(),
						Value: tagValue,
					},
				},
			},
		})
		require.NoError(t, err)

		tv, err := client.SystemFunctions.GetTag(ctx, tag.ID(), externalTableID, ObjectTypeExternalTable)
		// TODO: Add to the IntCreate tests
		require.NoError(t, err)
		assert.Equal(t, tagValue, tv)
	})

	t.Run("Alter: unset tags", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		opts.Tag = []TagAssociation{
			{
				Name:  tag.ID(),
				Value: "tag-value",
			},
		}
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.Alter(ctx, externalTableID, &AlterExternalTableOptions{
			IfExists: Bool(true),
			Unset: &ExternalTableUnset{
				Tag: []ObjectIdentifier{NewAccountObjectIdentifier(tag.ID().Name())},
			},
		})
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), externalTableID, ObjectTypeExternalTable)
		require.Error(t, err)
	})

	t.Run("Alter: add partitions", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := createExternalTableWithManualPartitioning
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(ctx, externalTableID, &AlterExternalTablePartitionOptions{
			IfExists: Bool(true),
			AddPartitions: []Partition{
				{
					ColumnName: "part_date",
					Value:      "2019-06-25",
				},
			},
			Location: "2019/06",
		})
		require.NoError(t, err)
	})

	t.Run("Alter: drop partitions", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := createExternalTableWithManualPartitioning
		err := client.ExternalTables.CreateWithManualPartitioning(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(ctx, externalTableID, &AlterExternalTablePartitionOptions{
			IfExists: Bool(true),
			AddPartitions: []Partition{
				{
					ColumnName: "part_date",
					Value:      "2019-06-25",
				},
			},
			Location: "2019/06",
		})
		require.NoError(t, err)

		err = client.ExternalTables.AlterPartitions(ctx, externalTableID, &AlterExternalTablePartitionOptions{
			IfExists:      Bool(true),
			DropPartition: Bool(true),
			Location:      "2019/06",
		})
		require.NoError(t, err)
	})

	t.Run("Drop", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		err = client.ExternalTables.Drop(ctx, externalTableID, &DropExternalTableOptions{
			IfExists: Bool(true),
			DropOption: &ExternalTableDropOption{
				Cascade: Bool(true),
			},
		})

		_, err = client.ExternalTables.ShowByID(ctx, externalTableID)
		require.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("Show", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		et, err := client.ExternalTables.Show(ctx, &ShowExternalTableOptions{
			Terse: Bool(true),
			Like: &Like{
				Pattern: String(externalTableID.Name()),
			},
			In: &In{
				Database: db.ID(),
			},
			StartsWith: String(externalTableID.Name()),
			LimitFrom: &LimitFrom{
				Rows: Int(1),
			},
		})
		require.NoError(t, err)
		assert.Equal(t, 1, len(et))
		assert.Equal(t, externalTableID, et[0].ID())
	})

	t.Run("Describe: columns", func(t *testing.T) {
		externalTableID := randomAccountObjectIdentifier(t)
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeColumns(ctx, externalTableID)
		require.NoError(t, err)

		assert.Equal(t, len(opts.Columns)+1, len(d)) // + 1 - because there's underlying Value column
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
		opts := minimalCreateExternalTableOpts
		err := client.ExternalTables.Create(ctx, externalTableID, &opts)
		require.NoError(t, err)

		d, err := client.ExternalTables.DescribeStage(ctx, externalTableID)
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
