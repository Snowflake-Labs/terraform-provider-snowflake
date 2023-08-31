package sdk

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExternalTablesCreate(t *testing.T) {
	t.Run("basic options", func(t *testing.T) {
		opts := &CreateExternalTableOpts{
			IfNotExists: Bool(true),
			name:        NewAccountObjectIdentifier("external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: "value::column::varchar",
					InlineConstraint: &ColumnInlineConstraint{
						Name:    String("my_constraint"),
						NotNull: Bool(true),
						Type:    &ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorage: &GoogleCloudStorageParams{
					Integration: String("123"),
				},
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJSON,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE EXTERNAL TABLE IF NOT EXISTS "external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`
		assert.Equal(t, expected, actual)
	})

	t.Run("every optional field", func(t *testing.T) {
		opts := &CreateExternalTableOpts{
			OrReplace: Bool(true),
			name:      NewAccountObjectIdentifier("external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: "value::column::varchar",
					InlineConstraint: &ColumnInlineConstraint{
						Name:    String("my_constraint"),
						NotNull: Bool(true),
						Type:    &ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorage: &GoogleCloudStorageParams{
					Integration: String("123"),
				},
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJSON,
				},
			},
			AwsSnsTopic: String("aws_sns_topic"),
			CopyGrants:  Bool(true),
			RowAccessPolicy: &RowAccessPolicy{
				Name: NewSchemaObjectIdentifier("db", "schema", "row_access_policy"),
				On:   []string{"value1", "value2"},
			},
			Tag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "value2",
				},
			},
			Comment: String("some_comment"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `CREATE OR REPLACE EXTERNAL TABLE "external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON) AWS_SNS_TOPIC = 'aws_sns_topic' COPY GRANTS ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2') COMMENT = 'some_comment'`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesCreateWithManualPartitioning(t *testing.T) {
	opts := &CreateWithManualPartitioningExternalTableOpts{
		OrReplace: Bool(true),
		name:      NewAccountObjectIdentifier("external_table"),
		Columns: []ExternalTableColumn{
			{
				Name:         "column",
				Type:         "varchar",
				AsExpression: "value::column::varchar",
				InlineConstraint: &ColumnInlineConstraint{
					Name:    String("my_constraint"),
					NotNull: Bool(true),
					Type:    &ColumnConstraintTypeUnique,
				},
			},
		},
		CloudProviderParams: &CloudProviderParams{
			GoogleCloudStorage: &GoogleCloudStorageParams{
				Integration: String("123"),
			},
		},
		Location:                   "@s1/logs/",
		UserSpecifiedPartitionType: Bool(true),
		FileFormat: []ExternalTableFileFormat{
			{
				Type: &ExternalTableFileFormatTypeJSON,
			},
		},
		CopyGrants: Bool(true),
		RowAccessPolicy: &RowAccessPolicy{
			Name: NewSchemaObjectIdentifier("db", "schema", "row_access_policy"),
			On:   []string{"value1", "value2"},
		},
		Tag: []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag2"),
				Value: "value2",
			},
		},
		Comment: String("some_comment"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	expected := `CREATE OR REPLACE EXTERNAL TABLE "external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (TYPE = JSON) COPY GRANTS ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2') COMMENT = 'some_comment'`
	assert.Equal(t, expected, actual)
}

func TestExternalTablesCreateDeltaLake(t *testing.T) {
	opts := &CreateDeltaLakeExternalTableOpts{
		OrReplace: Bool(true),
		name:      NewAccountObjectIdentifier("external_table"),
		Columns: []ExternalTableColumn{
			{
				Name:             "column",
				Type:             "varchar",
				AsExpression:     "value::column::varchar",
				InlineConstraint: nil,
			},
		},
		CloudProviderParams: &CloudProviderParams{
			MicrosoftAzure: &MicrosoftAzureParams{
				Integration: String("123"),
			},
		},
		PartitionBy:                []string{"column"},
		Location:                   "@s1/logs/",
		UserSpecifiedPartitionType: Bool(true),
		FileFormat: []ExternalTableFileFormat{
			{
				Name: String("JSON"),
			},
		},
		DeltaTableFormat: Bool(true),
		CopyGrants:       Bool(true),
		RowAccessPolicy: &RowAccessPolicy{
			Name: NewSchemaObjectIdentifier("db", "schema", "row_access_policy"),
			On:   []string{"value1", "value2"},
		},
		Tag: []TagAssociation{
			{
				Name:  NewAccountObjectIdentifier("tag1"),
				Value: "value1",
			},
			{
				Name:  NewAccountObjectIdentifier("tag2"),
				Value: "value2",
			},
		},
		Comment: String("some_comment"),
	}
	actual, err := structToSQL(opts)
	require.NoError(t, err)
	// TODO with line break to read / edit / review - write util to trim newline characters and double spaces
	expected := `CREATE OR REPLACE EXTERNAL TABLE "external_table" (column varchar AS (value::column::varchar)) INTEGRATION = '123' PARTITION BY (column) LOCATION = @s1/logs/ PARTITION_TYPE = USER_SPECIFIED FILE_FORMAT = (FORMAT_NAME = 'JSON') TABLE_FORMAT = DELTA COPY GRANTS ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2') COMMENT = 'some_comment'`
	assert.Equal(t, expected, actual)
}

func TestExternalTablesAlter(t *testing.T) {
	t.Run("refresh without path", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			Refresh:  &RefreshExternalTable{},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE IF EXISTS "external_table" REFRESH`
		assert.Equal(t, expected, actual)
	})

	t.Run("refresh with path", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			Refresh: &RefreshExternalTable{
				Path: "some/path",
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE IF EXISTS "external_table" REFRESH 'some/path'`
		assert.Equal(t, expected, actual)
	})

	t.Run("add files", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewAccountObjectIdentifier("external_table"),
			AddFiles: []ExternalTableFile{
				{Name: "one/file.txt"},
				{Name: "second/file.txt"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" ADD FILES ('one/file.txt', 'second/file.txt')`
		assert.Equal(t, expected, actual)
	})

	t.Run("remove files", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewAccountObjectIdentifier("external_table"),
			RemoveFiles: []ExternalTableFile{
				{Name: "one/file.txt"},
				{Name: "second/file.txt"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" REMOVE FILES ('one/file.txt', 'second/file.txt')`
		assert.Equal(t, expected, actual)
	})

	t.Run("set", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewAccountObjectIdentifier("external_table"),
			Set: &ExternalTableSet{
				AutoRefresh: Bool(true),
				Tag: []TagAssociation{
					{
						Name:  NewAccountObjectIdentifier("tag1"),
						Value: "tag_value1",
					},
					{
						Name:  NewAccountObjectIdentifier("tag2"),
						Value: "tag_value2",
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" SET AUTO_REFRESH = TRUE TAG ("tag1" = 'tag_value1', "tag2" = 'tag_value2')`
		assert.Equal(t, expected, actual)
	})

	t.Run("unset", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewAccountObjectIdentifier("external_table"),
			Unset: &ExternalTableUnset{
				Tag: []ObjectIdentifier{
					NewAccountObjectIdentifier("tag1"),
					NewAccountObjectIdentifier("tag2"),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" UNSET TAG "tag1", "tag2"`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesAlterPartitions(t *testing.T) {
	t.Run("add partition", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:     NewAccountObjectIdentifier("external_table"),
			IfExists: Bool(true),
			AddPartitions: []Partition{
				{
					ColumnName: "one",
					Value:      "one_value",
				},
				{
					ColumnName: "two",
					Value:      "two_value",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" IF EXISTS ADD PARTITION ("one" = 'one_value', "two" = 'two_value')`
		assert.Equal(t, expected, actual)
	})

	t.Run("remove partition", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:          NewAccountObjectIdentifier("external_table"),
			IfExists:      Bool(true),
			DropPartition: Bool(true),
			Location:      "partition_location",
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `ALTER EXTERNAL TABLE "external_table" IF EXISTS DROP PARTITION LOCATION 'partition_location'`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesDrop(t *testing.T) {
	t.Run("restrict", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP EXTERNAL TABLE IF EXISTS "external_table" RESTRICT`
		assert.Equal(t, expected, actual)
	})

	t.Run("cascade", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			IfExists: Bool(true),
			name:     NewAccountObjectIdentifier("external_table"),
			DropOption: &ExternalTableDropOption{
				Cascade: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DROP EXTERNAL TABLE IF EXISTS "external_table" CASCADE`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesShow(t *testing.T) {
	t.Run("all options", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			Like: &Like{
				Pattern: String("some_pattern"),
			},
			In: &In{
				Account: Bool(true),
			},
			StartsWith: String("some_external_table"),
			LimitFrom: &LimitFrom{
				Rows: Int(123),
				From: String("some_string"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE EXTERNAL TABLES LIKE 'some_pattern' IN ACCOUNT STARTS WITH 'some_external_table' LIMIT 123 FROM 'some_string'`
		assert.Equal(t, expected, actual)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Database: NewAccountObjectIdentifier("database_name"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE EXTERNAL TABLES IN DATABASE "database_name"`
		assert.Equal(t, expected, actual)
	})

	t.Run("in schema", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Schema: NewDatabaseObjectIdentifier("database_name", "schema_name"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `SHOW TERSE EXTERNAL TABLES IN SCHEMA "database_name"."schema_name"`
		assert.Equal(t, expected, actual)
	})
}

func TestExternalTablesDescribe(t *testing.T) {
	t.Run("type columns", func(t *testing.T) {
		opts := &describeExternalTableColumns{
			name: NewAccountObjectIdentifier("external_table"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE EXTERNAL TABLE "external_table" TYPE = COLUMNS`
		assert.Equal(t, expected, actual)
	})

	t.Run("type stage", func(t *testing.T) {
		opts := &describeExternalTableStage{
			name: NewAccountObjectIdentifier("external_table"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := `DESCRIBE EXTERNAL TABLE "external_table" TYPE = STAGE`
		assert.Equal(t, expected, actual)
	})
}
