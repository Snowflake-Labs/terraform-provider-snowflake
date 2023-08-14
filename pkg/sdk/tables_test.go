package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {

	t.Run("empty options", func(t *testing.T) {
		opts := &CreateTableOptions{}
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
				OnError: StageCopyOptionsOnErrorSkipFileNum{10},
			}},
		}
		rowAccessPolicy := RowAccessPolicy{
			Name: randomSchemaObjectIdentifier(t),
			On:   []string{"COLUMN_1", "COLUMN_2"},
		}
		opts := &CreateTableOptions{
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

func TestTableAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	t.Run("empty options", func(t *testing.T) {
		opts := &AlterTableOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "ALTER TABLE"
		assert.Equal(t, expected, actual)
	})

	t.Run("rename table", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		newID := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, randomUUID(t))
		opts := &AlterTableOptions{
			name:    id,
			NewName: newID,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("swap with", func(t *testing.T) {
		targetTableId := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, randomUUID(t))
		opts := &AlterTableOptions{
			name:     id,
			SwapWith: targetTableId,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s SWAP WITH %s", id.FullyQualifiedName(), targetTableId.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("cluster by", func(t *testing.T) {
		clusterByColumns := []string{"date", "id"}
		opts := &AlterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				ClusterBy: clusterByColumns,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s CLUSTER BY (date, id)", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("recluster", func(t *testing.T) {
		condition := "name = 'John'"
		opts := &AlterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				Recluster: &TableReclusterAction{
					MaxSize:   Int(1024),
					Condition: &condition,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s RECLUSTER MAX_SIZE = 1024 WHERE name = 'John'", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("suspend recluster", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				ChangeReclusterState: &TableReclusterChangeState{
					State: ReclusterStateSuspend,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s SUSPEND RECLUSTER", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("drop clustering key", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				DropClusteringKey: Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s DROP CLUSTERING KEY", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("add new column", func(t *testing.T) {
		columnName := "NEXT_COLUMN"
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				Add: &TableColumnAddAction{
					Column: Bool(true),
					Name:   columnName,
					Type:   DataTypeBoolean,
					DefaultValue: &ColumnDefaultValue{
						Identity: &ColumnIdentity{
							Start:     10,
							Increment: 1,
						},
					},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ADD COLUMN NEXT_COLUMN BOOLEAN IDENTITY START 10 INCREMENT 1", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("rename column", func(t *testing.T) {
		oldColumn := "OLD_NAME"
		newColumnName := "NEW_NAME"
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				Rename: &TableColumnRenameAction{
					OldName: oldColumn,
					NewName: newColumnName,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN OLD_NAME TO NEW_NAME", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("alter column", func(t *testing.T) {
		//column_1
		columnOneName := "COLUMN_1"
		alterActionsForColumnOne := []TableColumnAlterAction{
			{
				Name:        columnOneName,
				DropDefault: Bool(true),
			},
			{
				Name:       columnOneName,
				SetDefault: Pointer(SequenceName("SEQUENCE_1")),
			},
			{
				Name:         columnOneName,
				UnsetComment: Bool(true),
			},
		}

		columnTwoName := "COLUMN_2"
		alterActionsForColumnTwo := []TableColumnAlterAction{
			{
				Name:        columnTwoName,
				DropDefault: Bool(true),
			},
			{
				Name:       columnTwoName,
				SetDefault: Pointer(SequenceName("SEQUENCE_2")),
			},
			{
				Name:    columnTwoName,
				Comment: String("comment"),
			},
			{
				Name: columnTwoName,
				Type: Pointer(DataTypeBoolean),
			},
			{
				Name:              columnTwoName,
				NotNullConstraint: &TableColumnNotNullConstraint{Drop: Bool(true)},
			},
		}
		var actions []TableColumnAlterAction
		actions = append(actions, alterActionsForColumnOne...)
		actions = append(actions, alterActionsForColumnTwo...)

		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				Alter: actions,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN_1 DROP DEFAULT, COLUMN_1 SET DEFAULT SEQUENCE_1.NEXTVAL, COLUMN_1 UNSET COMMENT, COLUMN_2 DROP DEFAULT, COLUMN_2 SET DEFAULT SEQUENCE_2.NEXTVAL, COLUMN_2 COMMENT 'comment', COLUMN_2 SET DATA TYPE BOOLEAN, COLUMN_2 DROP NOT NULL", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("alter: unset masking policy", func(t *testing.T) {
		maskingPolicyName := randomSchemaObjectIdentifier(t)
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				SetMaskingPolicy: &TableColumnAlterSetMaskingPolicyAction{
					ColumnName:        "COLUMN_1",
					MaskingPolicyName: maskingPolicyName,
					Using:             []string{"FOO", "BAR"},
					Force:             Bool(true),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN COLUMN_1 SET MASKING POLICY %s USING (FOO, BAR) FORCE", id.FullyQualifiedName(), maskingPolicyName.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("alter: unset masking policy", func(t *testing.T) {
		maskingPolicyName := randomSchemaObjectIdentifier(t)
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				UnsetMaskingPolicy: &TableColumnAlterUnsetMaskingPolicyAction{
					ColumnName:        "COLUMN_1",
					MaskingPolicyName: maskingPolicyName,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ALTER COLUMN COLUMN_1 UNSET MASKING POLICY %s", id.FullyQualifiedName(), maskingPolicyName.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("alter: set tags", func(t *testing.T) {
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
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				SetTags: &TableColumnAlterSetTagsAction{
					ColumnName: "COLUMN_1",
					Tags:       columnTags,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN COLUMN_1 SET TAG "db"."schema"."column_tag1" = 'v1', "db"."schema"."column_tag2" = 'v2'`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("alter: unset tags", func(t *testing.T) {
		columnTags := []ObjectIdentifier{
			NewSchemaObjectIdentifier("db", "schema", "column_tag1"),
			NewSchemaObjectIdentifier("db", "schema", "column_tag2"),
		}
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				UnsetTags: &TableColumnAlterUnsetTagsAction{
					ColumnName: "COLUMN_1",
					Tags:       columnTags,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN COLUMN_1 UNSET TAG "db"."schema"."column_tag1", "db"."schema"."column_tag2"`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("alter: drop columns", func(t *testing.T) {
		columns := []string{"COLUMN_1", "COLUMN_2"}
		opts := &AlterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				DropColumns: &TableColumnAlterDropColumns{
					Columns: columns,
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s DROP COLUMN COLUMN_1, COLUMN_2", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("alter constraint: add", func(t *testing.T) {
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
		opts := &AlterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Add: &outOfLineConstraint,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT", id.FullyQualifiedName(), outOfLineConstraint.ForeignKey.TableName.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("alter constraint: rename", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Rename: &TableConstraintRenameAction{
					OldName: "OLD_NAME_CONSTRAINT",
					NewName: "NEW_NAME_CONSTRAINT",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s RENAME CONSTRAINT OLD_NAME_CONSTRAINT TO NEW_NAME_CONSTRAINT", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("alter constraint: alter", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Alter: &TableConstraintAlterAction{
					ConstraintName: String("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					NotEnforced:    Bool(true),
					Valiate:        Bool(true),
					Rely:           Bool(true),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ALTER CONSTRAINT OUT_OF_LINE_CONSTRAINT (COLUMN_3, COLUMN_4) NOT ENFORCED VALIDATE RELY", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("alter constraint: drop", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Drop: &TableConstraintDropAction{
					ConstraintName: String("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					Cascade:        Bool(true),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s DROP CONSTRAINT OUT_OF_LINE_CONSTRAINT (COLUMN_3, COLUMN_4) CASCADE", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("external table: add", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ExternalTableAction: &TableExternalTableAction{
				Add: &TableExternalTableColumnAddAction{
					Name:       "COLUMN_1",
					Type:       DataTypeBoolean,
					Expression: []string{"SELECT 1"},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ADD COLUMN COLUMN_1 BOOLEAN AS (SELECT 1)", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("external table: rename", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ExternalTableAction: &TableExternalTableAction{
				Rename: &TableExternalTableColumnRenameAction{
					OldName: "OLD_NAME_COLUMN",
					NewName: "NEW_NAME_COLUMN",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s RENAME COLUMN OLD_NAME_COLUMN TO NEW_NAME_COLUMN", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("external table: drop", func(t *testing.T) {
		opts := &AlterTableOptions{
			name: id,
			ExternalTableAction: &TableExternalTableAction{
				Drop: &TableExternalTableColumnDropAction{
					Columns: []string{"COLUMN_3", "COLUMN_4"},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s DROP COLUMN COLUMN_3, COLUMN_4", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
	//TODO remove
}
