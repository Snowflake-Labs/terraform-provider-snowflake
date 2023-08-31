package sdk

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableCreate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *createTableOptions {
		return &createTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: no columns", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errTableNeedsAtLeastOneColumn)
	})
	t.Run("validation: both expression and identity of a column are present ", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		opts.Columns = []TableColumn{{
			Name: "",
			DefaultValue: &ColumnDefaultValue{
				Expression: String(""),
				Identity: &ColumnIdentity{
					Start:     10,
					Increment: 1,
				},
			},
		}}
		assertOptsInvalidJoinedErrors(t, opts, errColumnDefaultValueNeedsExactlyOneValue)
	})
	t.Run("validation: column masking policy incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		opts.Columns = []TableColumn{{
			Name: "",
			MaskingPolicy: &ColumnMaskingPolicy{
				Name: NewSchemaObjectIdentifier("", "", ""),
			},
		}}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: column tag association's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		opts.Columns = []TableColumn{{
			Name: "",
			Tags: []TagAssociation{{
				Name:  NewSchemaObjectIdentifier("", "", ""),
				Value: "v1",
			},
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: outOfLineConstraint's foreign key incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.OutOfLineConstraint = &CreateOutOfLineConstraint{
			ForeignKey: &OutOfLineForeignKey{
				TableName: NewSchemaObjectIdentifier("", "", ""),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: stageFileFormat's both format name and format type are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.StageFileFormat = []StageFileFormat{
			{
				InnerValue: StageFileFormatInnerValue{
					FormatName: String(""),
					FormatType: Pointer(FileFormatTypeCSV),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errStageFileFormatValueNeedsExactlyOneValue)
	})

	t.Run("validation: stageFileFormat's both format name and format type are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.StageFileFormat = []StageFileFormat{
			{
				InnerValue: StageFileFormatInnerValue{
					FormatName: String(""),
					FormatType: Pointer(FileFormatTypeCSV),
				},
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errStageFileFormatValueNeedsExactlyOneValue)
	})
	t.Run("validation: rowAccessPolicy's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.RowAccessPolicy = &RowAccessPolicy{
			Name: NewSchemaObjectIdentifier("", "", ""),
		}
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableOptions{}
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
		outOfLineConstraint := CreateOutOfLineConstraint{
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
		stageCopyOptions := []StageCopyOption{{
			StageCopyOptionsInnerValue{
				OnError: StageCopyOptionsOnErrorSkipFileNum{10},
			}},
		}
		rowAccessPolicy := RowAccessPolicy{
			Name: randomSchemaObjectIdentifier(t),
			On:   []string{"COLUMN_1", "COLUMN_2"},
		}
		opts := &createTableOptions{
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
			MaxDataExtentionTimeInDays: Int(100),
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
			`CREATE TABLE %s ( %s %s COLLATE 'de' COMMENT '%s' IDENTITY START 10 INCREMENT 1 NOT NULL MASKING POLICY %s USING (FOO, BAR) TAG ("db"."schema"."column_tag1" = 'v1', "db"."schema"."column_tag2" = 'v2') CONSTRAINT INLINE_CONSTRAINT PRIMARY KEY CONSTRAINT OUT_OF_LINE_CONSTRAINT FOREIGN KEY (COLUMN_1, COLUMN_2) REFERENCES %s (COLUMN_3, COLUMN_4) MATCH FULL ON UPDATE SET NULL ON DELETE RESTRICT ) CLUSTER BY (COLUMN_1, COLUMN_2) ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV COMPRESSION = AUTO) STAGE_COPY_OPTIONS = (ON_ERROR = SKIP_FILE_10) DATA_RETENTION_TIME_IN_DAYS = 10 MAX_DATA_EXTENSION_TIME_IN_DAYS = 100 CHANGE_TRACKING = true DEFAULT_DDL_COLLATION = 'en' COPY GRANTS ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) TAG ("db"."schema"."table_tag1" = 'v1', "db"."schema"."table_tag2" = 'v2') COMMENT = '%s'`,
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

func TestTableCreateAsSelect(t *testing.T) {

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableAsSelectOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE TABLE ( )"
		assert.Equal(t, expected, actual)
	})
	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier(t)
		columnName := "FIRST_COLUMN"
		columnType, err := ToDataType("VARCHAR")
		maskingPolicy := TableAsSelectColumnMaskingPolicy{
			With: Bool(true),
			Name: randomSchemaObjectIdentifier(t),
		}
		rowAccessPolicy := RowAccessPolicy{
			Name: randomSchemaObjectIdentifier(t),
			On:   []string{"COLUMN_1", "COLUMN_2"},
		}
		opts := &createTableAsSelectOptions{
			OrReplace: Bool(true),
			name:      id,
			Columns: []TableAsSelectColumn{
				{
					Name:          columnName,
					Type:          Pointer(columnType),
					MaskingPolicy: &maskingPolicy,
				},
			},
			ClusterBy:  []string{"COLUMN_1", "COLUMN_2"},
			CopyGrants: Bool(true),

			RowAccessPolicy: &rowAccessPolicy,
			Query:           String("* FROM ANOTHER_TABLE"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("CREATE OR REPLACE TABLE %s ( FIRST_COLUMN VARCHAR WITH MASKING POLICY %s ) CLUSTER BY (COLUMN_1, COLUMN_2) COPY GRANTS ROW ACCESS POLICY %s ON (COLUMN_1, COLUMN_2) AS SELECT * FROM ANOTHER_TABLE",
			id.FullyQualifiedName(),
			maskingPolicy.Name.FullyQualifiedName(),
			rowAccessPolicy.Name.FullyQualifiedName(),
		)
		assert.Equal(t, expected, actual)
	})
}

func TestTableCreateUsingTemplate(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *createTableAsSelectOptions {
		return &createTableAsSelectOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableAsSelectOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: no columns", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsInvalidJoinedErrors(t, opts, errTableNeedsAtLeastOneColumn)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableUsingTemplateOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE TABLE USING TEMPLATE"
		assert.Equal(t, expected, actual)
	})
	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier(t)
		opts := &createTableUsingTemplateOptions{
			OrReplace:  Bool(true),
			name:       id,
			CopyGrants: Bool(true),
			Query:      []string{"sample_data"},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("CREATE OR REPLACE TABLE %s COPY GRANTS USING TEMPLATE sample_data",
			id.FullyQualifiedName(),
		)
		assert.Equal(t, expected, actual)
	})
}

func TestTableCreateLike(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *createTableLikeOptions {
		return &createTableLikeOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableLikeOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: source table's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SourceTable = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableLikeOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE TABLE LIKE"
		assert.Equal(t, expected, actual)
	})
	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier(t)
		sourceTable := randomSchemaObjectIdentifier(t)
		opts := &createTableLikeOptions{
			OrReplace:   Bool(true),
			name:        id,
			SourceTable: sourceTable,
			ClusterBy:   []string{"date", "id"},
			CopyGrants:  Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("CREATE OR REPLACE TABLE %s LIKE %s CLUSTER BY (date, id) COPY GRANTS",
			id.FullyQualifiedName(),
			sourceTable.FullyQualifiedName(),
		)
		assert.Equal(t, expected, actual)
	})
}

func TestTableCreateClone(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *createTableCloneOptions {
		return &createTableCloneOptions{
			name: id,
		}
	}
	t.Run("validation: nil options", func(t *testing.T) {
		var opts *createTableCloneOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &createTableCloneOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "CREATE TABLE CLONE"
		assert.Equal(t, expected, actual)
	})
	t.Run("with complete options", func(t *testing.T) {
		id := randomSchemaObjectIdentifier(t)
		sourceTable := randomSchemaObjectIdentifier(t)
		clonePoint := ClonePoint{
			Moment: CloneMomentAt,
			At: TimeTravel{
				Offset: Int(0),
			},
		}
		opts := &createTableCloneOptions{
			OrReplace:   Bool(true),
			name:        id,
			SourceTable: sourceTable,
			ClonePoint:  &clonePoint,
			CopyGrants:  Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("CREATE OR REPLACE TABLE %s CLONE %s AT (OFFSET => 0) COPY GRANTS",
			id.FullyQualifiedName(),
			sourceTable.FullyQualifiedName(),
		)
		assert.Equal(t, expected, actual)
	})
}

func TestTableAlter(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *alterTableOptions {
		return &alterTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *alterTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		opts.NewName = Pointer(NewSchemaObjectIdentifier("test", "test", "test"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("validation: both NewName and SwapWith are present ", func(t *testing.T) {
		opts := defaultOpts()
		opts.NewName = Pointer(NewSchemaObjectIdentifier("test", "test", "test"))
		opts.SwapWith = Pointer(NewSchemaObjectIdentifier("test", "test", "test"))

		assertOptsInvalidJoinedErrors(t, opts, errAlterTableNeedsExactlyOneAction)
	})
	t.Run("validation: NewName's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.NewName = Pointer(NewSchemaObjectIdentifier("", "", ""))
		opts.SwapWith = Pointer(NewSchemaObjectIdentifier("test", "test", "test"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: SwapWith's incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.SwapWith = Pointer(NewSchemaObjectIdentifier("", "", ""))
		opts.NewName = Pointer(NewSchemaObjectIdentifier("test", "test", "test"))
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})
	t.Run("validation: Clustering action's ClusterBy and Recluster are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ClusteringAction = &TableClusteringAction{
			ClusterBy: []string{"date"},
			Recluster: &TableReclusterAction{
				MaxSize:   Int(10),
				Condition: String("true"),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableClusteringActionNeedsExactlyOneAction)
	})
	t.Run("validation: Column action's Add and Rename are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ColumnAction = &TableColumnAction{
			Add: &TableColumnAddAction{},
			Rename: &TableColumnRenameAction{
				NewName: "new",
				OldName: "old",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableColumnActionNeedsExactlyOneAction)
	})
	t.Run("validation: Column alter action's DropDefault and SetDefault are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ColumnAction = &TableColumnAction{
			Alter: []TableColumnAlterAction{{
				DropDefault: Bool(true),
				SetDefault:  Pointer(SequenceName("sequence")),
			}},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableColumnAlterActionNeedsExactlyOneAction)
	})
	t.Run("validation: Constraint alter action's ConstraintName and Unique are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ConstraintAction = &TableConstraintAction{
			Alter: &TableConstraintAlterAction{
				ConstraintName: String("constraint"),
				Unique:         Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableConstraintAlterActionNeedsExactlyOneAction)
	})
	t.Run("validation: Constraint drop action's ConstraintName and Unique are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ConstraintAction = &TableConstraintAction{
			Drop: &TableConstraintDropAction{
				ConstraintName: String("constraint"),
				Unique:         Bool(true),
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableConstraintDropActionNeedsExactlyOneAction)
	})
	t.Run("validation: External action's Add and Rename are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.ExternalTableAction = &TableExternalTableAction{
			Add: &TableExternalTableColumnAddAction{},
			Rename: &TableExternalTableColumnRenameAction{
				OldName: "old_name",
				NewName: "new_name",
			},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableExternalActionNeedsExactlyOneAction)
	})
	t.Run("validation: SearchOptimization action's Add and Drop are present", func(t *testing.T) {
		opts := defaultOpts()
		opts.SearchOptimizationAction = &TableSearchOptimizationAction{
			Add:  &AddSearchOptimaztion{},
			Drop: &DropSearchOptimaztion{},
		}
		assertOptsInvalidJoinedErrors(t, opts, errTableSearchOptimizationActionNeedsExactlyOneAction)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := &alterTableOptions{}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := "ALTER TABLE"
		assert.Equal(t, expected, actual)
	})

	t.Run("table with name", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("rename", func(t *testing.T) {
		newID := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, randomUUID(t))
		opts := &alterTableOptions{
			name:    id,
			NewName: &newID,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s RENAME TO %s", id.FullyQualifiedName(), newID.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("swap with", func(t *testing.T) {
		targetTableId := NewSchemaObjectIdentifier(id.databaseName, id.schemaName, randomUUID(t))
		opts := &alterTableOptions{
			name:     id,
			SwapWith: &targetTableId,
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s SWAP WITH %s", id.FullyQualifiedName(), targetTableId.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("cluster by", func(t *testing.T) {
		clusterByColumns := []string{"date", "id"}
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
			name: id,
			ClusteringAction: &TableClusteringAction{
				ChangeReclusterState: &TableReclusterChangeState{
					State: Pointer(ReclusterStateSuspend),
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s SUSPEND RECLUSTER", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})

	t.Run("drop clustering key", func(t *testing.T) {
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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

		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
			name: id,
			ColumnAction: &TableColumnAction{
				UnsetMaskingPolicy: &TableColumnAlterUnsetMaskingPolicyAction{
					ColumnName: "COLUMN_1",
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		outOfLineConstraint := AlterOutOfLineConstraint{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
			name: id,
			ConstraintAction: &TableConstraintAction{
				Alter: &TableConstraintAlterAction{
					ConstraintName: String("OUT_OF_LINE_CONSTRAINT"),
					Columns:        []string{"COLUMN_3", "COLUMN_4"},
					NotEnforced:    Bool(true),
					Validate:       Bool(true),
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
		opts := &alterTableOptions{
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
	t.Run("add search optimiaztion", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SearchOptimizationAction: &TableSearchOptimizationAction{
				Add: &AddSearchOptimaztion{
					On: []string{"SUBSTRING(*)", "GEO(*)"},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s ADD SEARCH OPTIMIZATION ON SUBSTRING(*), GEO(*)", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)

	})

	t.Run("drop search optimiaztion", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SearchOptimizationAction: &TableSearchOptimizationAction{
				Drop: &DropSearchOptimaztion{
					On: []string{"SUBSTRING(*)", "FOO"},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s DROP SEARCH OPTIMIZATION ON SUBSTRING(*), FOO", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)

	})

	t.Run("drop search optimiaztion", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SearchOptimizationAction: &TableSearchOptimizationAction{
				Drop: &DropSearchOptimaztion{
					On: []string{"SUBSTRING(*)", "FOO"},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf("ALTER TABLE %s DROP SEARCH OPTIMIZATION ON SUBSTRING(*), FOO", id.FullyQualifiedName())
		assert.Equal(t, expected, actual)

	})

	t.Run("set: with complete options", func(t *testing.T) {
		comment := randomString(t)
		opts := &alterTableOptions{
			name: id,
			Set: &TableSet{
				EnableSchemaEvolution: Bool(true),
				StageFileFormat: []StageFileFormat{
					{
						InnerValue: StageFileFormatInnerValue{
							FormatType: Pointer(FileFormatTypeCSV),
						},
					},
				},
				StageCopyOptions: []StageCopyOption{
					{
						InnerValue: StageCopyOptionsInnerValue{
							OnError: StageCopyOptionsOnErrorSkipFileNumPercentage{Value: 10},
						},
					},
				},
				DataRetentionTimeInDays:    Int(30),
				MaxDataExtensionTimeInDays: Int(90),
				ChangeTracking:             Bool(false),
				DefaultDDLCollation:        String("us"),
				Comment:                    &comment,
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s SET ENABLE_SCHEMA_EVOLUTION = true STAGE_FILE_FORMAT = (TYPE = CSV) STAGE_COPY_OPTIONS = (ON_ERROR = 'SKIP_FILE_10%%') DATA_RETENTION_TIME_IN_DAYS = 30 MAX_DATA_EXTENSION_TIME_IN_DAYS = 90 CHANGE_TRACKING = false DEFAULT_DDL_COLLATION = 'us' COMMENT = '%s'`, id.FullyQualifiedName(), comment)
		assert.Equal(t, expected, actual)
	})
	t.Run("set tags", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			SetTags: []TagAssociation{
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "table_tag1"),
					Value: "v1",
				},
				{
					Name:  NewSchemaObjectIdentifier("db", "schema", "table_tag2"),
					Value: "v2",
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s SET TAG "db"."schema"."table_tag1" = 'v1', "db"."schema"."table_tag2" = 'v2'`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("unset tags", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			UnsetTags: []ObjectIdentifier{
				NewSchemaObjectIdentifier("db", "schema", "table_tag1"),
				NewSchemaObjectIdentifier("db", "schema", "table_tag2"),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s UNSET TAG "db"."schema"."table_tag1", "db"."schema"."table_tag2"`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("unset: complete options", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			Unset: &TableUnset{
				DataRetentionTimeInDays:    Bool(true),
				MaxDataExtensionTimeInDays: Bool(true),
				ChangeTracking:             Bool(true),
				DefaultDDLCollation:        Bool(true),
				EnableSchemaEvolution:      Bool(true),
				Comment:                    Bool(true),
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s UNSET DATA_RETENTION_TIME_IN_DAYS MAX_DATA_EXTENSION_TIME_IN_DAYS CHANGE_TRACKING DEFAULT_DDL_COLLATION ENABLE_SCHEMA_EVOLUTION COMMENT`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("add row access policy", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			AddRowAccessPolicy: &AddRowAccessPolicy{
				PolicyName:  "ROW_ACCESS_POLICY_1",
				ColumnNames: []string{"FIRST_COLUMN"},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s ADD ROW ACCESS POLICY ROW_ACCESS_POLICY_1 ON (FIRST_COLUMN)`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("drop row access policy", func(t *testing.T) {
		opts := &alterTableOptions{
			name:                id,
			DropRowAccessPolicy: String("ROW_ACCESS_POLICY_1"),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s DROP ROW ACCESS POLICY ROW_ACCESS_POLICY_1`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("drop and add row access policy", func(t *testing.T) {
		opts := &alterTableOptions{
			name: id,
			DropAndAddRowAccessPolicy: &DropAndAddRowAccessPolicy{
				DroppedPolicyName: "ROW_ACCESS_POLICY_1",
				AddedPolicy: &AddRowAccessPolicy{
					PolicyName:  "ROW_ACCESS_POLICY_2",
					ColumnNames: []string{"FIRST_COLUMN"},
				},
			},
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s DROP ROW ACCESS POLICY ROW_ACCESS_POLICY_1 , ADD ROW ACCESS POLICY ROW_ACCESS_POLICY_2 ON (FIRST_COLUMN)`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
	t.Run("drop all row access policies", func(t *testing.T) {
		opts := &alterTableOptions{
			name:                     id,
			DropAllAccessRowPolicies: Bool(true),
		}
		actual, err := structToSQL(opts)
		require.NoError(t, err)
		expected := fmt.Sprintf(`ALTER TABLE %s DROP ALL ROW ACCESS POLICIES`, id.FullyQualifiedName())
		assert.Equal(t, expected, actual)
	})
}

func TestTableDrop(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *dropTableOptions {
		return &dropTableOptions{
			name: id,
		}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *dropTableOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: incorrect identifier", func(t *testing.T) {
		opts := defaultOpts()
		opts.name = NewSchemaObjectIdentifier("", "", "")
		assertOptsInvalidJoinedErrors(t, opts, ErrInvalidObjectIdentifier)
	})

	t.Run("empty options", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `DROP TABLE %s`, id.FullyQualifiedName())
	})

	t.Run("with if exists", func(t *testing.T) {
		opts := defaultOpts()
		opts.IfExists = Bool(true)
		assertOptsValidAndSQLEquals(t, opts, `DROP TABLE IF EXISTS %s`, id.FullyQualifiedName())
	})
}

func TestTableShow(t *testing.T) {
	id := randomSchemaObjectIdentifier(t)

	defaultOpts := func() *showTableOptions {
		return &showTableOptions{}
	}

	t.Run("validation: nil options", func(t *testing.T) {
		var opts *PipeShowOptions = nil
		assertOptsInvalidJoinedErrors(t, opts, errNilOptions)
	})

	t.Run("validation: empty like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{}
		assertOptsInvalidJoinedErrors(t, opts, errPatternRequiredForLikeKeyword)
	})

	t.Run("show", func(t *testing.T) {
		opts := defaultOpts()
		assertOptsValidAndSQLEquals(t, opts, `SHOW TABLES`)
	})

	t.Run("show with like", func(t *testing.T) {
		opts := defaultOpts()
		opts.Like = &Like{
			Pattern: String(id.Name()),
		}
		assertOptsValidAndSQLEquals(t, opts, `SHOW TABLES LIKE '%s'`, id.Name())
	})
}
