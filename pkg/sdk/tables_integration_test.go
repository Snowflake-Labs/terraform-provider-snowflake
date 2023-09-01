package sdk

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Table(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schema, _ := createSchema(t, client, database)

	cleanupTableProvider := func(id SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tables.Drop(ctx, NewDropTableRequest(id))
			require.NoError(t, err)
		}
	}
	tag1, _ := createTag(t, client, database, schema)
	tag2, _ := createTag(t, client, database, schema)

	t.Run("create table: no optionals", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("FIRST_COLUMN", DataTypeNumber).WithDefaultValue(NewColumnDefaultValueRequest().WithIdentity(NewColumnIdentityRequest(1, 1))),
			*NewTableColumnRequest("SECOND_COLUMN", DataTypeNumber).WithDefaultValue(NewColumnDefaultValueRequest().WithIdentity(NewColumnIdentityRequest(1, 1))),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.Name, name)
	})
	t.Run("create table: complete optionals", func(t *testing.T) {
		maskingPolicy, _ := createMaskingPolicyWithOptions(t, client, database, schema, []TableColumnSignature{
			{
				Name: "col1",
				Type: DataTypeVARCHAR,
			},
			{
				Name: "col2",
				Type: DataTypeVARCHAR,
			},
		}, DataTypeVARCHAR, "REPLACE('X', 1, 2)", nil)
		table2, _ := createTable(t, client, database, schema)
		name := randomString(t)
		comment := randomString(t)

		columnTags := []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_3", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR).
				WithDefaultValue(NewColumnDefaultValueRequest().WithExpression(String("'default'"))).
				WithWith(Bool(true)).
				WithMaskingPolicy(NewColumnMaskingPolicyRequest(NewSchemaObjectIdentifier(database.Name, schema.Name, maskingPolicy.Name)).WithUsing([]string{"COLUMN_1", "COLUMN_3"})).
				WithTags(columnTags).
				WithNotNull(Bool(true)),
			*NewTableColumnRequest("COLUMN_2", DataTypeNumber).WithDefaultValue(NewColumnDefaultValueRequest().WithIdentity(NewColumnIdentityRequest(1, 1))),
		}
		outOfLineConstraint := NewOutOfLineConstraintRequest("OUT_OF_LINE_CONSTRAINT", ColumnConstraintTypeForeignKey).
			WithColumns([]string{"COLUMN_1"}).
			WithForeignKey(NewOutOfLineForeignKeyRequest(NewSchemaObjectIdentifier(database.Name, schema.Name, table2.Name), []string{"id"}).
				WithMatch(Pointer(FullMatchType)).
				WithOn(NewForeignKeyOnActionRequest().
					WithOnDelete(Pointer(ForeignKeySetNullAction)).WithOnUpdate(Pointer(ForeignKeyRestrictAction))))
		stageFileFormat := NewStageFileFormatRequest().
			WithFormatType(Pointer(FileFormatTypeCSV)).
			WithOptions(NewFileFormatTypeOptionsRequest().WithCSVCompression(Pointer(CSVCompressionAuto)))
		stageCopyOptions := NewStageCopyOptionsRequest(StageCopyOptionsOnErrorSkipFileNumPercentage{Value: 10})
		request := NewCreateTableRequest(id, columns).
			WithOutOfLineConstraint(outOfLineConstraint).
			WithStageFileFormat([]StageFileFormatRequest{*stageFileFormat}).
			WithStageCopyOptions([]StageCopyOptionsRequest{*stageCopyOptions}).
			WithComment(&comment).
			WithDataRetentionTimeInDays(Int(30)).
			WithMaxDataExtensionTimeInDays(Int(30))

		err := client.Tables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.Name, name)
		assert.Equal(t, table.Comment, comment)
		assert.Equal(t, 30, table.RetentionTime)
		param, err := client.Sessions.ShowObjectParameter(ctx, "MAX_DATA_EXTENSION_TIME_IN_DAYS", ObjectTypeTable, table.ID())
		assert.NoError(t, err)
		assert.Equal(t, "30", param.Value)
		tableColumns := tableColumns(t, ctx, client, schema.Name, table.Name)
		assert.Equal(t, 3, len(tableColumns))
	})

	t.Run("create table as select", func(t *testing.T) {
		maskingPolicy, _ := createMaskingPolicyWithOptions(t, client, database, schema, []TableColumnSignature{
			{
				Name: "col1",
				Type: DataTypeVARCHAR,
			},
		}, DataTypeVARCHAR, "REPLACE('X', 1)", nil)
		columns := []TableAsSelectColumnRequest{
			*NewTableAsSelectColumnRequest("COLUMN_3").
				WithType_(Pointer(DataTypeVARCHAR)).
				WithCopyGrants(Bool(true)).
				WithOrReplace(Bool(true)),
			*NewTableAsSelectColumnRequest("COLUMN_1").
				WithType_(Pointer(DataTypeVARCHAR)).
				WithCopyGrants(Bool(true)).
				WithOrReplace(Bool(true)),
			*NewTableAsSelectColumnRequest("COLUMN_2").
				WithType_(Pointer(DataTypeVARCHAR)).
				WithCopyGrants(Bool(true)).
				WithOrReplace(Bool(true)).WithMaskingPolicyName(Pointer(NewSchemaObjectIdentifier(database.Name, schema.Name, maskingPolicy.Name))),
		}

		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		request := NewCreateTableAsSelectRequest(id, columns)
		err := client.Tables.CreateAsSelect(ctx, request)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		tableColumns := tableColumns(t, ctx, client, schema.Name, table.Name)
		assert.Equal(t, 3, len(tableColumns))
		t.Cleanup(cleanupTableProvider(id))
	})

	t.Run("create table using template", func(t *testing.T) {
		fileFormat, _ := createFileFormat(t, client, schema.ID())
		stage, stageCleanup := createStageWithName(t, client, "new_stage")
		t.Cleanup(stageCleanup)
		warehouse, warehouseCleanup := createWarehouse(t, client)
		f, err := os.CreateTemp("/tmp", "data.csv")
		require.NoError(t, err)
		w := bufio.NewWriter(f)
		_, err = w.WriteString(` [{"name": "column1", "type" "INTEGER"},
									 {"name": "column2", "type" "INTEGER"} ]`)
		require.NoError(t, err)
		err = w.Flush()
		require.NoError(t, err)
		t.Cleanup(warehouseCleanup)
		_, err = client.exec(ctx, fmt.Sprintf("PUT file://%s @%s", f.Name(), *stage))
		require.NoError(t, err)
		err = os.Remove(f.Name())
		require.NoError(t, err)

		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		err = client.Sessions.UseWarehouse(ctx, warehouse.ID())
		require.NoError(t, err)
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '@%s', FILE_FORMAT=>'%s', ignore_case => true))`, *stage, fileFormat.ID().FullyQualifiedName())
		request := NewCreateTableUsingTemplateRequest(id, query)
		err = client.Tables.CreateUsingTemplate(ctx, request)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		tableColumns := tableColumns(t, ctx, client, schema.Name, table.Name)
		assert.Equal(t, 3, len(tableColumns))
		t.Cleanup(cleanupTableProvider(id))
	})

	t.Run("create table like", func(t *testing.T) {
		sourceTable, _ := createTable(t, client, database, schema)
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		request := NewCreateTableLikeRequest(id, sourceTable.ID()).WithCopyGrants(Bool(true))
		err := client.Tables.CreateLike(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))
		sourceTableColumns := tableColumns(t, ctx, client, schema.Name, sourceTable.Name)
		likeTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		likeTableColumns := tableColumns(t, ctx, client, schema.Name, likeTable.Name)
		assert.Equal(t, len(sourceTableColumns), len(likeTableColumns))
		for i := range sourceTableColumns {
			assert.True(t, sourceTableColumns[i] == likeTableColumns[i])
		}
	})
	t.Run("create table clone", func(t *testing.T) {
		sourceTable, _ := createTable(t, client, database, schema)
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		request := NewCreateTableCloneRequest(id, sourceTable.ID()).
			WithCopyGrants(Bool(true)).WithClonePoint(NewClonePointRequest().
			WithAt(*NewTimeTravelRequest().WithOffset(Pointer(0))).
			WithMoment(CloneMomentAt))
		err := client.Tables.CreateClone(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))
		sourceTableColumns := tableColumns(t, ctx, client, schema.Name, sourceTable.Name)
		cloneTable, err := client.Tables.ShowByID(ctx, id)
		cloneTableColumns := tableColumns(t, ctx, client, schema.Name, cloneTable.Name)
		assert.Equal(t, len(sourceTableColumns), len(cloneTableColumns))
		for i := range sourceTableColumns {
			assert.True(t, sourceTableColumns[i] == cloneTableColumns[i])
		}
	})

	t.Run("alter table: rename", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		newName := randomString(t)
		newId := NewSchemaObjectIdentifier(database.Name, schema.Name, newName)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_3", DataTypeVARCHAR),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).WithNewName(&newId)
		err = client.Tables.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupTableProvider(id))
		} else {
			t.Cleanup(cleanupTableProvider(newId))
		}
		require.NoError(t, err)
		_, err = client.Tables.ShowByID(ctx, id)
		assert.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)

		table, err := client.Tables.ShowByID(ctx, newId)
		require.NoError(t, err)
		assert.Equal(t, table.Name, newName)
	})

	t.Run("alter table: swap with", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		secondTableName := randomString(t)
		secondTableId := NewSchemaObjectIdentifier(database.Name, schema.Name, secondTableName)
		secondTableColumns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		err = client.Tables.Create(ctx, NewCreateTableRequest(secondTableId, secondTableColumns))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).WithSwapWith(&secondTableId)
		err = client.Tables.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupTableProvider(id))
		} else {
			t.Cleanup(cleanupTableProvider(secondTableId))
		}
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, secondTableId)
		assert.Equal(t, table.Name, secondTableName)
		secondTable, err := client.Tables.ShowByID(ctx, id)
		assert.Equal(t, secondTable.Name, name)
		require.NoError(t, err)
	})

	t.Run("alter table: cluster by", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		require.NoError(t, err)
		clusterByColumns := []string{"COLUMN_1", "COLUMN_2"}
		alterRequest := NewAlterTableRequest(id).WithClusteringAction(NewTableClusteringActionRequest().WithClusterBy(clusterByColumns))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.Comment, "")
		clusterByString := "LINEAR(" + strings.Join(clusterByColumns, ", ") + ")"
		assert.Equal(t, table.ClusterBy, clusterByString)
	})

	t.Run("alter table: resume recluster", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithClusteringAction(NewTableClusteringActionRequest().
				WithChangeReclusterState(Pointer(ReclusterStateSuspend)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter table: drop clustering key", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithClusteringAction(NewTableClusteringActionRequest().
				WithDropClusteringKey(Bool(true)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter table: add a column", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithColumnAction(NewTableColumnActionRequest().
				WithAdd(NewTableColumnAddActionRequest("COLUMN_3", DataTypeVARCHAR)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		currentColumns := tableColumns(t, ctx, client, schema.Name, table.Name)

		assert.Equal(t, table.Comment, "")
		assert.Equal(t, len(currentColumns), 3)
	})
	t.Run("alter table: rename column", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithColumnAction(NewTableColumnActionRequest().
				WithRename(NewTableColumnRenameActionRequest("COLUMN_1", "COLUMN_3")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		currentColumns := tableColumns(t, ctx, client, schema.Name, table.Name)

		assert.Equal(t, table.Comment, "")
		assert.Equal(t, len(currentColumns), 2)
		containsNewColumn := false
		containsOldColumn := false
		for _, column := range currentColumns {
			if column == "COLUMN_3" {
				containsNewColumn = true
			} else if column == "COLUMN_1" {
				containsOldColumn = true
			}
		}
		assert.True(t, containsNewColumn)
		assert.False(t, containsOldColumn)
	})

	t.Run("alter table: unset masking policy", func(t *testing.T) {
		maskingPolicy, _ := createMaskingPolicyWithOptions(t, client, database, schema, []TableColumnSignature{
			{
				Name: "col1",
				Type: DataTypeVARCHAR,
			},
		}, DataTypeVARCHAR, "REPLACE('X', 1)", nil)
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR).WithMaskingPolicy(NewColumnMaskingPolicyRequest(maskingPolicy.ID())),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		err := client.Tables.
			Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithColumnAction(NewTableColumnActionRequest().WithUnsetMaskingPolicy(NewTableColumnAlterUnsetMaskingPolicyActionRequest("COLUMN_1")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})
	t.Run("alter table: set tags", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		columnTags := []TagAssociationRequest{
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}

		alterRequest := NewAlterTableRequest(id).WithSetTags(columnTags)
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})
	t.Run("alter table: unset tags", func(t *testing.T) {
		columnTags := []TagAssociationRequest{
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).
			WithClusterBy(clusterBy).
			WithTags(columnTags))
		require.NoError(t, err)
		columnNames := []ObjectIdentifier{
			NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
			NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
		}

		alterRequest := NewAlterTableRequest(id).WithUnsetTags(columnNames)
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter table: drop columns", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithColumnAction(NewTableColumnActionRequest().WithDropColumns([]string{"COLUMN_1"}))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		currentColumns := tableColumns(t, ctx, client, schema.Name, table.Name)
		assert.Equal(t, len(currentColumns), 1)
		containsOldColumn := false
		for _, column := range currentColumns {
			if column == "COLUMN_1" {
				containsOldColumn = true
			}
		}
		assert.False(t, containsOldColumn)
		assert.Equal(t, table.Comment, "")
	})

	t.Run("alter constraint: add", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
		}
		secondTableName := randomString(t)
		secondTableId := NewSchemaObjectIdentifier(database.Name, schema.Name, secondTableName)
		secondTableColumns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_3", DataTypeVARCHAR).WithInlineConstraint(NewColumnInlineConstraintRequest("pkey", ColumnConstraintTypePrimaryKey)),
		}
		err := client.Tables.Create(ctx, NewCreateTableRequest(secondTableId, secondTableColumns))
		require.NoError(t, err)
		err = client.Tables.Create(ctx, NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithConstraintAction(NewTableConstraintActionRequest().
				WithAdd(NewOutOfLineConstraintRequest("OUT_OF_LINE_CONSTRAINT", ColumnConstraintTypeForeignKey).WithColumns([]string{"COLUMN_1"}).
					WithForeignKey(NewOutOfLineForeignKeyRequest(NewSchemaObjectIdentifier(database.Name, schema.Name, secondTableName), []string{"COLUMN_3"}))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter constraint: rename", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		oldConstraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(oldConstraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		newConstraintName := "NEW_OUT_OF_LINE_CONSTRAINT_NAME"
		alterRequest := NewAlterTableRequest(id).
			WithConstraintAction(NewTableConstraintActionRequest().
				WithRename(NewTableConstraintRenameActionRequest().
					WithOldName(oldConstraintName).
					WithNewName(newConstraintName)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})
	t.Run("alter constraint: alter", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithConstraintAction(NewTableConstraintActionRequest().WithAlter(NewTableConstraintAlterActionRequest().WithConstraintName(String(constraintName)).WithEnforced(Bool(true))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter constraint: drop", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithConstraintAction(NewTableConstraintActionRequest().WithDrop(NewTableConstraintDropActionRequest().WithConstraintName(String(constraintName))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("external table: add", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithExternalTableAction(NewTableExternalTableActionRequest().WithAdd(NewTableExternalTableColumnAddActionRequest().WithName("COLUMN_3").WithType(DataTypeNumber).WithExpression([]string{"1 + 1"})))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := tableColumns(t, ctx, client, schema.Name, table.Name)
		assert.Equal(t, len(currentColumns), 3)
		newColumnExists := false
		for _, column := range currentColumns {
			if column == "COLUMN_3" {
				newColumnExists = true
			}
		}
		assert.True(t, newColumnExists)
	})

	t.Run("external table: rename", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithExternalTableAction(NewTableExternalTableActionRequest().WithRename(NewTableExternalTableColumnRenameActionRequest().WithOldName("COLUMN_1").WithNewName("COLUMN_3")))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.Comment, "")
		currentColumns := tableColumns(t, ctx, client, schema.Name, table.Name)
		assert.Equal(t, len(currentColumns), 2)
		oldColumnExists := false
		newColumnExists := false
		for _, column := range currentColumns {
			if column == "COLUMN_1" {
				oldColumnExists = true
			} else if column == "COLUMN_3" {
				newColumnExists = true
			}
		}
		assert.True(t, newColumnExists)
		assert.False(t, oldColumnExists)
	})

	t.Run("external table: drop", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithExternalTableAction(NewTableExternalTableActionRequest().WithDrop(NewTableExternalTableColumnDropActionRequest().WithColumns([]string{"COLUMN_2"})))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)
		currentColumns := tableColumns(t, ctx, client, schema.Name, table.Name)

		assert.Equal(t, len(currentColumns), 1)
	})

	t.Run("add search optimiaztion", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithSearchOptimizationAction(NewTableSearchOptimizationActionRequest().WithAddSearchOptimizationOn([]string{"SUBSTRING(*)", "GEO(*)"}))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("set: with complete options", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		comment := randomString(t)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		require.NoError(t, err)
		stageFileFormats := []StageFileFormatRequest{
			{
				FormatType: Pointer(FileFormatTypeCSV),
			},
		}
		stageCopyOptions := []StageCopyOptionsRequest{
			{
				OnError: StageCopyOptionsOnErrorSkipFileNumPercentage{Value: 10},
			},
		}
		alterRequest := NewAlterTableRequest(id).
			WithSet(NewTableSetRequest().
				WithEnableSchemaEvolution(Bool(true)).
				WithStageFileFormat(stageFileFormats).
				WithStageCopyOptions(stageCopyOptions).
				WithDataRetentionTimeInDays(Int(30)).
				WithMaxDataExtensionTimeInDays(Int(90)).
				WithChangeTracking(Bool(false)).
				WithDefaultDDLCollation(String("us")).
				WithComment(&comment))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.Comment, comment)
		assert.Equal(t, table.RetentionTime, 30)
		assert.Equal(t, table.ChangeTracking, false)
		assert.Equal(t, table.EnableSchemaEvolution, true)
	})

	t.Run("set tags", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		setTags := []TagAssociationRequest{
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithSetTags(setTags)

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("alter: unset tags", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columnTags := []TagAssociation{
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}
		columns := []TableColumnRequest{
			*NewTableColumnRequest("COLUMN_1", DataTypeVARCHAR).WithTags(columnTags),
			*NewTableColumnRequest("COLUMN_2", DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := NewOutOfLineConstraintRequest(constraintName, ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})
		err := client.Tables.Create(ctx, NewCreateTableRequest(id, columns).WithOutOfLineConstraint(outOfLineConstraint))
		unsetTags := []ObjectIdentifier{
			NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
			NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
		}
		require.NoError(t, err)
		alterRequest := NewAlterTableRequest(id).
			WithUnsetTags(unsetTags)

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("drop table", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)
		err := client.Tables.Drop(ctx, NewDropTableRequest(table.ID()).WithIfExists(Bool(true)))
		require.NoError(t, err)
		_, err = client.Tables.ShowByID(ctx, table.ID())
		require.ErrorIs(t, err, ErrObjectNotExistOrAuthorized)
	})

	t.Run("show tables", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		table2, table2Cleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)
		t.Cleanup(table2Cleanup)
		tables, err := client.Tables.Show(ctx, nil)
		require.NoError(t, err)
		tableIds := make([]SchemaObjectIdentifier, len(tables))
		for i, table := range tables {
			tableIds[i] = table.ID()
		}
		assert.Contains(t, tableIds, table.ID())
		assert.Contains(t, tableIds, table2.ID())

		require.NoError(t, err)
	})

	t.Run("with terse", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)
		tables, err := client.Tables.Show(ctx, NewShowTableRequest().WithTerse(Bool(true)).WithLikePattern(table.ID().Name()))
		assert.Equal(t, 1, len(tables))
		table2 := tables[0]
		assert.Equal(t, table2.Name, table.Name)
		assert.NotEmpty(t, table2.CreatedOn)
		assert.Empty(t, table2.Owner)
		require.NoError(t, err)
	})
	t.Run("with starts with", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)
		tables, err := client.Tables.Show(ctx, NewShowTableRequest().WithStartsWith(String(table.Name)))
		assert.Equal(t, 1, len(tables))
		table2 := tables[0]
		assert.Equal(t, table2.Name, table.Name)
		assert.NotEmpty(t, table2.CreatedOn)
		assert.NotEmpty(t, table2.Owner)
		require.NoError(t, err)
	})
	t.Run("when searching a non-existent table", func(t *testing.T) {
		tables, err := client.Tables.Show(ctx, NewShowTableRequest().WithLikePattern("non-existent"))
		require.NoError(t, err)
		assert.Equal(t, 0, len(tables))
	})
}

func tableColumns(t *testing.T, ctx context.Context, client *Client, schemaName, tableName string) []string {
	t.Helper()
	warehouse, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err := client.Sessions.UseWarehouse(ctx, warehouse.ID())
	require.NoError(t, err)
	query := fmt.Sprintf("SELECT column_name\nFROM information_schema.columns\nWHERE table_schema = '%s'\n  AND table_name = '%s'\nORDER BY ordinal_position", schemaName, tableName)
	var columnNames []string
	err = client.query(ctx, &columnNames, query)
	require.NoError(t, err)
	return columnNames
}
