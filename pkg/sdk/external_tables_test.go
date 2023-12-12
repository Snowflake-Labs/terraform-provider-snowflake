package sdk

import (
	"testing"
)

func TestExternalTablesCreate(t *testing.T) {
	t.Run("basic options", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			IfNotExists: Bool(true),
			name:        NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorageIntegration: String("123"),
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJSON,
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE IF NOT EXISTS "db"."schema"."external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`)
	})

	t.Run("every optional field", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			OrReplace: Bool(true),
			name:      NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorageIntegration: String("123"),
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJSON,
				},
			},
			AwsSnsTopic: String("aws_sns_topic"),
			CopyGrants:  Bool(true),
			RowAccessPolicy: &TableRowAccessPolicy{
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE "db"."schema"."external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON) AWS_SNS_TOPIC = 'aws_sns_topic' COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
			name:        NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errOneOf("CreateExternalTableOptions", "OrReplace", "IfNotExists"),
			ErrInvalidObjectIdentifier,
			errNotSet("CreateExternalTableOptions", "Location"),
			errExactlyOneOf("CreateExternalTableOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location:      "@s1/logs/",
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE "db"."schema"."external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`)
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location: "@s1/logs/",
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateExternalTableOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTablesCreateWithManualPartitioning(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			OrReplace: Bool(true),
			name:      NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			CloudProviderParams: &CloudProviderParams{
				GoogleCloudStorageIntegration: String("123"),
			},
			Location: "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Type: &ExternalTableFileFormatTypeJSON,
				},
			},
			CopyGrants: Bool(true),
			RowAccessPolicy: &TableRowAccessPolicy{
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE "db"."schema"."external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) INTEGRATION = '123' LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON) COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
			name:        NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errOneOf("CreateWithManualPartitioningExternalTableOptions", "OrReplace", "IfNotExists"),
			ErrInvalidObjectIdentifier,
			errNotSet("CreateWithManualPartitioningExternalTableOptions", "Location"),
			errExactlyOneOf("CreateWithManualPartitioningExternalTableOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location:      "@s1/logs/",
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE "db"."schema"."external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`)
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateWithManualPartitioningExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location: "@s1/logs/",
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateWithManualPartitioningExternalTableOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTablesCreateDeltaLake(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			OrReplace: Bool(true),
			name:      NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:             "column",
					Type:             "varchar",
					AsExpression:     []string{"value::column::varchar"},
					InlineConstraint: nil,
				},
			},
			CloudProviderParams: &CloudProviderParams{
				MicrosoftAzureIntegration: String("123"),
			},
			PartitionBy: []string{"column"},
			Location:    "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Name: String("JSON"),
				},
			},
			DeltaTableFormat: Bool(true),
			CopyGrants:       Bool(true),
			RowAccessPolicy: &TableRowAccessPolicy{
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
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE "db"."schema"."external_table" (column varchar AS (value::column::varchar)) INTEGRATION = '123' PARTITION BY (column) LOCATION = @s1/logs/ FILE_FORMAT = (FORMAT_NAME = 'JSON') TABLE_FORMAT = DELTA COPY GRANTS COMMENT = 'some_comment' ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
			name:        NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			errOneOf("CreateDeltaLakeExternalTableOptions", "OrReplace", "IfNotExists"),
			ErrInvalidObjectIdentifier,
			errNotSet("CreateDeltaLakeExternalTableOptions", "Location"),
			errExactlyOneOf("CreateDeltaLakeExternalTableOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location:      "@s1/logs/",
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE "db"."schema"."external_table" (column varchar AS (value::column::varchar) NOT NULL CONSTRAINT my_constraint UNIQUE) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`)
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateDeltaLakeExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Columns: []ExternalTableColumn{
				{
					Name:         "column",
					Type:         "varchar",
					AsExpression: []string{"value::column::varchar"},
					NotNull:      Bool(true),
					InlineConstraint: &ColumnInlineConstraint{
						Name: String("my_constraint"),
						Type: ColumnConstraintTypeUnique,
					},
				},
			},
			Location: "@s1/logs/",
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateDeltaLakeExternalTableOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTableUsingTemplateOpts(t *testing.T) {
	t.Run("valid options", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			OrReplace:  Bool(true),
			name:       NewSchemaObjectIdentifier("db", "schema", "external_table"),
			CopyGrants: Bool(true),
			Query:      []string{"query statement"},
			CloudProviderParams: &CloudProviderParams{
				MicrosoftAzureIntegration: String("123"),
			},
			PartitionBy: []string{"column"},
			Location:    "@s1/logs/",
			FileFormat: []ExternalTableFileFormat{
				{
					Name: String("JSON"),
				},
			},
			Comment: String("some_comment"),
			RowAccessPolicy: &TableRowAccessPolicy{
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
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE OR REPLACE EXTERNAL TABLE "db"."schema"."external_table" COPY GRANTS USING TEMPLATE (query statement) INTEGRATION = '123' PARTITION BY (column) LOCATION = @s1/logs/ FILE_FORMAT = (FORMAT_NAME = 'JSON') COMMENT = 'some_comment' ROW ACCESS POLICY "db"."schema"."row_access_policy" ON (value1, value2) TAG ("tag1" = 'value1', "tag2" = 'value2')`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			name: NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errNotSet("CreateExternalTableUsingTemplateOptions", "Query"),
			errNotSet("CreateExternalTableUsingTemplateOptions", "Location"),
			errExactlyOneOf("CreateExternalTableUsingTemplateOptions", "RawFileFormat", "FileFormat"),
		)
	})

	t.Run("raw file format", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Location: "@s1/logs/",
			Query: []string{
				"query statement",
			},
			RawFileFormat: &RawFileFormat{Format: "TYPE = JSON"},
		}
		assertOptsValidAndSQLEquals(t, opts, `CREATE EXTERNAL TABLE "db"."schema"."external_table" USING TEMPLATE (query statement) LOCATION = @s1/logs/ FILE_FORMAT = (TYPE = JSON)`)
	})

	t.Run("validation: neither raw file format is set, nor file format", func(t *testing.T) {
		opts := &CreateExternalTableUsingTemplateOptions{
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Location: "@s1/logs/",
			Query: []string{
				"query statement",
			},
		}
		assertOptsInvalid(t, opts, errExactlyOneOf("CreateExternalTableUsingTemplateOptions", "RawFileFormat", "FileFormat"))
	})
}

func TestExternalTablesAlter(t *testing.T) {
	t.Run("refresh without path", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Refresh:  &RefreshExternalTable{},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS "db"."schema"."external_table" REFRESH ''`)
	})

	t.Run("refresh with path", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			IfExists: Bool(true),
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
			Refresh: &RefreshExternalTable{
				Path: "some/path",
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS "db"."schema"."external_table" REFRESH 'some/path'`)
	})

	t.Run("add files", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			AddFiles: []ExternalTableFile{
				{Name: "one/file.txt"},
				{Name: "second/file.txt"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE "db"."schema"."external_table" ADD FILES ('one/file.txt', 'second/file.txt')`)
	})

	t.Run("remove files", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			RemoveFiles: []ExternalTableFile{
				{Name: "one/file.txt"},
				{Name: "second/file.txt"},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE "db"."schema"."external_table" REMOVE FILES ('one/file.txt', 'second/file.txt')`)
	})

	t.Run("set auto refresh", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name:        NewSchemaObjectIdentifier("db", "schema", "external_table"),
			AutoRefresh: Bool(true),
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE "db"."schema"."external_table" SET AUTO_REFRESH = true`)
	})

	t.Run("set tag", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			SetTag: []TagAssociation{
				{
					Name:  NewAccountObjectIdentifier("tag1"),
					Value: "tag_value1",
				},
				{
					Name:  NewAccountObjectIdentifier("tag2"),
					Value: "tag_value2",
				},
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE "db"."schema"."external_table" SET TAG "tag1" = 'tag_value1', "tag2" = 'tag_value2'`)
	})

	t.Run("unset tag", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
			UnsetTag: []ObjectIdentifier{
				NewAccountObjectIdentifier("tag1"),
				NewAccountObjectIdentifier("tag2"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE "db"."schema"."external_table" UNSET TAG "tag1", "tag2"`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &AlterExternalTableOptions{
			name:        NewSchemaObjectIdentifier("", "", ""),
			AddFiles:    []ExternalTableFile{{Name: "some file"}},
			RemoveFiles: []ExternalTableFile{{Name: "some other file"}},
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errExactlyOneOf("AlterExternalTableOptions", "Refresh", "AddFiles", "RemoveFiles", "AutoRefresh", "SetTag", "UnsetTag"),
		)
	})
}

func TestExternalTablesAlterPartitions(t *testing.T) {
	t.Run("add partition", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
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
			Location: "123",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS "db"."schema"."external_table" ADD PARTITION (one = 'one_value', two = 'two_value') LOCATION '123'`)
	})

	t.Run("remove partition", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:          NewSchemaObjectIdentifier("db", "schema", "external_table"),
			IfExists:      Bool(true),
			DropPartition: Bool(true),
			Location:      "partition_location",
		}
		assertOptsValidAndSQLEquals(t, opts, `ALTER EXTERNAL TABLE IF EXISTS "db"."schema"."external_table" DROP PARTITION LOCATION 'partition_location'`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &AlterExternalTablePartitionOptions{
			name:          NewSchemaObjectIdentifier("", "", ""),
			AddPartitions: []Partition{{ColumnName: "colName", Value: "value"}},
			DropPartition: Bool(true),
		}
		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errOneOf("AlterExternalTablePartitionOptions", "AddPartitions", "DropPartition"),
		)
	})
}

func TestExternalTablesDrop(t *testing.T) {
	t.Run("restrict", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			IfExists: Bool(true),
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL TABLE IF EXISTS "db"."schema"."external_table" RESTRICT`)
	})

	t.Run("cascade", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			IfExists: Bool(true),
			name:     NewSchemaObjectIdentifier("db", "schema", "external_table"),
			DropOption: &ExternalTableDropOption{
				Cascade: Bool(true),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `DROP EXTERNAL TABLE IF EXISTS "db"."schema"."external_table" CASCADE`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			name: NewSchemaObjectIdentifier("", "", ""),
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
				Cascade:  Bool(true),
			},
		}

		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errOneOf("ExternalTableDropOption", "Restrict", "Cascade"),
		)
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
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EXTERNAL TABLES LIKE 'some_pattern' IN ACCOUNT STARTS WITH 'some_external_table' LIMIT 123 FROM 'some_string'`)
	})

	t.Run("in database", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Database: NewAccountObjectIdentifier("database_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EXTERNAL TABLES IN DATABASE "database_name"`)
	})

	t.Run("in schema", func(t *testing.T) {
		opts := &ShowExternalTableOptions{
			Terse: Bool(true),
			In: &In{
				Schema: NewDatabaseObjectIdentifier("database_name", "schema_name"),
			},
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TERSE EXTERNAL TABLES IN SCHEMA "database_name"."schema_name"`)
	})

	t.Run("invalid options", func(t *testing.T) {
		opts := &DropExternalTableOptions{
			name: NewSchemaObjectIdentifier("", "", ""),
			DropOption: &ExternalTableDropOption{
				Restrict: Bool(true),
				Cascade:  Bool(true),
			},
		}

		assertOptsInvalidJoinedErrors(
			t, opts,
			ErrInvalidObjectIdentifier,
			errOneOf("ExternalTableDropOption", "Restrict", "Cascade"),
		)
	})
}

func TestExternalTablesDescribe(t *testing.T) {
	t.Run("type columns", func(t *testing.T) {
		opts := &describeExternalTableColumnsOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL TABLE "db"."schema"."external_table" TYPE = COLUMNS`)
	})

	t.Run("type stage", func(t *testing.T) {
		opts := &describeExternalTableStageOptions{
			name: NewSchemaObjectIdentifier("db", "schema", "external_table"),
		}
		assertOptsValidAndSQLEquals(t, opts, `DESCRIBE EXTERNAL TABLE "db"."schema"."external_table" TYPE = STAGE`)
	})
}
