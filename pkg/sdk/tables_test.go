package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {

	t.Run("empty options", func(t *testing.T) {
		opts := &TableCreateOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE TABLE ( )"
		assert.Equal(t, expected, actual)
	})
	t.Run("with complete options", func(t *testing.T) {
		columnComment := randomString(t)
		tableComment := randomString(t)
		collation := "de"
		id := randomSchemaObjectIdentifier(t)
		columnName := "FIRST_COLUMN"
		columnType, err := ToDataType("VARCHAR")
		maskingPolicy := ColumnMaskingPolicy{
			Name:  randomSchemaObjectIdentifier(t),
			Using: []string{"FOO", "BAR"},
		}
		columnTags := []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "column_tag1"),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "column_tag2"),
				Value: "v2",
			},
		}

		tableTags := []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "table_tag1"),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier("db", "schema", "table_tag2"),
				Value: "v2",
			},
		}
		inlineConstraint := ColumnInlineConstraint{
			Name: "INLINE_CONSTRAINT",
			Type: ColumnConstraintTypePrimaryKey,
		}
		require.NoError(t, err)
		outOfLineConstraint := OutOfLineConstraint{
			Name:    "OUT_OF_LINE_CONSTRAINT",
			Type:    ColumnConstraintTypeForeignKey,
			Columns: []string{"COLUMN_1", "COLUMN_2"},
			ForeignKey: &OutOfLineForeignKey{
				TableName:   randomSchemaObjectIdentifier(t),
				ColumnNames: []string{"COLUMN_3", "COLUMN_4"},
				Match:       Pointer(FullMatchType),
				On: &ForeignKeyOnAction{
					OnUpdate: Pointer(ForeignKeySetNullAction),
					OnDelete: Pointer(ForeignKeyRestrictAction),
				},
			},
		}
		stageFileFormat := []StageFileFormat{{
			StageFileFormatInnerValue{
			FormatType: Pointer(FileFormatTypeCSV),
			Options: &FileFormatTypeOptions{
				CSVCompression: Pointer(CSVCompressionAuto),
			},
			},
		}}
		stageCopyOptions := []StageCopyOptions{{
			StageCopyOptionsInnerValue{ 
			OnError:           StageCopyOptionsOnErrorSkipFileNum{10},
		}},
		}
		rowAccessPolicy := RowAccessPolicy{
			Name: randomSchemaObjectIdentifier(t),
			On:   []string{"COLUMN_1", "COLUMN_2"},
		}
		opts := &TableCreateOptions{
			name: id,
			Columns: []TableColumn{{
				Name:    columnName,
				Type:    columnType,
				Collate: &collation,
				Comment: &columnComment,
				DefaultValue: &ColumnDefaultValue{
					Identity: &ColumnIdentity{
						Start:     10,
						Increment: 1,
					},
				},
				NotNull:          Bool(true),
				MaskingPolicy:    &maskingPolicy,
				Tags:             columnTags,
				InlineConstraint: &inlineConstraint,
			}},
			OutOfLineConstraint:        &outOfLineConstraint,
			ClusterBy:                  []string{"COLUMN_1", "COLUMN_2"},
			EnableSchemaEvolution:      Bool(true),
			StageFileFormat:            stageFileFormat,
			StageCopyOptions:           stageCopyOptions,
			DataRetentionTimeInDays:    Int(10),
			MaxDataRetentionTimeInDays: Int(100),
			ChangeTracking:             Bool(true),
			DefaultDDLCollation:        String("en"),
			CopyGrants:                 Bool(true),
			RowAccessPolicy:            &rowAccessPolicy,
			Tags:                       tableTags,
			Comment:                    &tableComment,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(
			`CREATE TABLE %s ( %s %s COLLATE 'de' COMMENT '%s' IDENTITY START 10 INCREMENT 1 NOT NULL MASKING POLICY %s USING (FOO, BAR) TAG ("db"."schema"."column_tag1" = 'v1', "db"."schema"."column_tag2" = 'v2') CONSTRAINT INLINE_CONSTRAINT PRIMARY KEY CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT ) CLUSTER BY (COLUMN_1, COLUMN_2) ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV COMPRESSION = AUTO) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE_10) DATA_RETENTION_TIME_IN_DAYS = 10 MAX_DATA_RETENTION_TIME_IN_DAYS = 100 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en' COPY GRANTS ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) TAG ("db"."schema"."table_tag1" = 'v1', "db"."schema"."table_tag2" = 'v2') COMMENT = '%s'`,
			id.FullyQualifiedName(),
			columnName,
			columnType,
			columnComment,
			maskingPolicy.Name.FullyQualifiedName(),
			outOfLineConstraint.ForeignKey.TableName.FullyQualifiedName(),
			rowAccessPolicy.Name.FullyQualifiedName(),
			tableComment,
		)
		assert.Equal(t, expected, actual)
	})

}
