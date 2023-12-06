package testint

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/collections"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type expectedColumn struct {
	Name string
	Type sdk.DataType
}

func TestInt_Table(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	database := testDb(t)
	schema := testSchema(t)

	cleanupTableProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(id))
			require.NoError(t, err)
		}
	}
	tag1, _ := createTag(t, client, database, schema)
	tag2, _ := createTag(t, client, database, schema)

	assertColumns := func(t *testing.T, expectedColumns []expectedColumn, createdColumns []informationSchemaColumns) {
		t.Helper()

		require.Len(t, createdColumns, len(expectedColumns))
		for i, expectedColumn := range expectedColumns {
			assert.Equal(t, strings.ToUpper(expectedColumn.Name), createdColumns[i].ColumnName)
			createdColumnDataType, err := sdk.ToDataType(createdColumns[i].DataType)
			assert.NoError(t, err)
			assert.Equal(t, expectedColumn.Type, createdColumnDataType)
		}
	}

	assertTable := func(t *testing.T, table *sdk.Table, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, table.ID())
		assert.NotEmpty(t, table.CreatedOn)
		assert.Equal(t, id.Name(), table.Name)
		assert.Equal(t, testDb(t).Name, table.DatabaseName)
		assert.Equal(t, testSchema(t).Name, table.SchemaName)
		assert.Equal(t, "TABLE", table.Kind)
		assert.Equal(t, 0, table.Rows)
		assert.Equal(t, "ACCOUNTADMIN", table.Owner)
	}

	assertTableTerse := func(t *testing.T, table *sdk.Table, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, table.ID())
		assert.NotEmpty(t, table.CreatedOn)
		assert.Equal(t, id.Name(), table.Name)
		assert.Equal(t, testDb(t).Name, table.DatabaseName)
		assert.Equal(t, testSchema(t).Name, table.SchemaName)
		assert.Equal(t, "TABLE", table.Kind)
		assert.Empty(t, table.Rows)
		assert.Empty(t, table.Owner)
	}

	t.Run("create table: no optionals", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("FIRST_COLUMN", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
			*sdk.NewTableColumnRequest("SECOND_COLUMN", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
	})

	t.Run("create table: complete optionals", func(t *testing.T) {
		maskingPolicy, _ := createMaskingPolicyWithOptions(t, client, database, schema, []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: sdk.DataTypeVARCHAR,
			},
			{
				Name: "col2",
				Type: sdk.DataTypeVARCHAR,
			},
		}, sdk.DataTypeVARCHAR, "REPLACE('X', 1, 2)", nil)
		table2, _ := createTable(t, client, database, schema)
		name := random.String()
		comment := random.String()

		columnTags := []sdk.TagAssociation{
			{
				Name:  sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_3", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR).
				WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithExpression(sdk.String("'default'"))).
				WithMaskingPolicy(sdk.NewColumnMaskingPolicyRequest(sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, maskingPolicy.Name)).WithUsing([]string{"COLUMN_1", "COLUMN_3"})).
				WithTags(columnTags).
				WithNotNull(sdk.Bool(true)),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeNumber).WithDefaultValue(sdk.NewColumnDefaultValueRequest().WithIdentity(sdk.NewColumnIdentityRequest(1, 1))),
		}
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest("OUT_OF_LINE_CONSTRAINT", sdk.ColumnConstraintTypeForeignKey).
			WithColumns([]string{"COLUMN_1"}).
			WithForeignKey(sdk.NewOutOfLineForeignKeyRequest(sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, table2.Name), []string{"id"}).
				WithMatch(sdk.Pointer(sdk.FullMatchType)).
				WithOn(sdk.NewForeignKeyOnAction().
					WithOnDelete(sdk.Pointer(sdk.ForeignKeySetNullAction)).WithOnUpdate(sdk.Pointer(sdk.ForeignKeyRestrictAction))))
		stageFileFormat := sdk.NewStageFileFormatRequest().
			WithType(sdk.Pointer(sdk.FileFormatTypeCSV)).
			WithOptions(sdk.NewFileFormatTypeOptionsRequest().WithCSVCompression(sdk.Pointer(sdk.CSVCompressionAuto)))
		stageCopyOptions := sdk.NewStageCopyOptionsRequest().WithOnError(sdk.NewStageCopyOnErrorOptionsRequest().WithSkipFile())
		request := sdk.NewCreateTableRequest(id, columns).
			WithOutOfLineConstraint(*outOfLineConstraint).
			WithStageFileFormat(*stageFileFormat).
			WithStageCopyOptions(*stageCopyOptions).
			WithComment(&comment).
			WithDataRetentionTimeInDays(sdk.Int(30)).
			WithMaxDataExtensionTimeInDays(sdk.Int(30))

		err := client.Tables.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
		assert.Equal(t, table.Comment, comment)
		assert.Equal(t, 30, table.RetentionTime)

		param, err := client.Parameters.ShowObjectParameter(ctx, sdk.ObjectParameterMaxDataExtensionTimeInDays, sdk.Object{ObjectType: sdk.ObjectTypeTable, Name: table.ID()})
		assert.NoError(t, err)
		assert.Equal(t, "30", param.Value)

		tableColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeNumber},
		}
		assertColumns(t, expectedColumns, tableColumns)
	})

	t.Run("create table as select", func(t *testing.T) {
		maskingPolicy, _ := createMaskingPolicyWithOptions(t, client, database, schema, []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: sdk.DataTypeVARCHAR,
			},
		}, sdk.DataTypeVARCHAR, "REPLACE('X', 1)", nil)
		columns := []sdk.TableAsSelectColumnRequest{
			*sdk.NewTableAsSelectColumnRequest("COLUMN_3").
				WithType_(sdk.Pointer(sdk.DataTypeVARCHAR)).
				WithCopyGrants(sdk.Bool(true)).
				WithOrReplace(sdk.Bool(true)),
			*sdk.NewTableAsSelectColumnRequest("COLUMN_1").
				WithType_(sdk.Pointer(sdk.DataTypeVARCHAR)).
				WithCopyGrants(sdk.Bool(true)).
				WithOrReplace(sdk.Bool(true)),
			*sdk.NewTableAsSelectColumnRequest("COLUMN_2").
				WithType_(sdk.Pointer(sdk.DataTypeVARCHAR)).
				WithCopyGrants(sdk.Bool(true)).
				WithOrReplace(sdk.Bool(true)).WithMaskingPolicyName(sdk.Pointer(sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, maskingPolicy.Name))),
		}

		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		request := sdk.NewCreateTableAsSelectRequest(id, columns)

		err := client.Tables.CreateAsSelect(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		tableColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, tableColumns)
	})

	// TODO: fix this test, it should create two integer column but is creating 3 text ones instead
	t.Run("create table using template", func(t *testing.T) {
		fileFormat, fileFormatCleanup := createFileFormat(t, client, schema.ID())
		t.Cleanup(fileFormatCleanup)
		stage, stageCleanup := createStageWithName(t, client, "new_stage")
		t.Cleanup(stageCleanup)

		f, err := os.CreateTemp("/tmp", "data.csv")
		require.NoError(t, err)
		w := bufio.NewWriter(f)
		_, err = w.WriteString(` [{"name": "column1", "type" "INTEGER"},
									 {"name": "column2", "type" "INTEGER"} ]`)
		require.NoError(t, err)
		err = w.Flush()
		require.NoError(t, err)
		_, err = client.ExecForTests(ctx, fmt.Sprintf("PUT file://%s @%s", f.Name(), *stage))
		require.NoError(t, err)
		err = os.Remove(f.Name())
		require.NoError(t, err)

		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		query := fmt.Sprintf(`SELECT ARRAY_AGG(OBJECT_CONSTRUCT(*)) WITHIN GROUP (ORDER BY order_id) FROM TABLE (INFER_SCHEMA(location => '@%s', FILE_FORMAT=>'%s', ignore_case => true))`, *stage, fileFormat.ID().FullyQualifiedName())
		request := sdk.NewCreateTableUsingTemplateRequest(id, query)

		err = client.Tables.CreateUsingTemplate(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		returnedTableColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"C1", sdk.DataTypeVARCHAR},
			{"C2", sdk.DataTypeVARCHAR},
			{"C3", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, returnedTableColumns)
	})

	t.Run("create table like", func(t *testing.T) {
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", "NUMBER"),
			*sdk.NewTableColumnRequest("col2", "VARCHAR"),
			*sdk.NewTableColumnRequest("col3", "BOOLEAN"),
		}
		sourceTable, sourceTableCleanup := createTableWithColumns(t, client, database, schema, columns)
		t.Cleanup(sourceTableCleanup)

		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		request := sdk.NewCreateTableLikeRequest(id, sourceTable.ID()).WithCopyGrants(sdk.Bool(true))

		err := client.Tables.CreateLike(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		sourceTableColumns := getTableColumnsFor(t, client, sourceTable.ID())
		expectedColumns := []expectedColumn{
			{"id", sdk.DataTypeNumber},
			{"col2", sdk.DataTypeVARCHAR},
			{"col3", sdk.DataTypeBoolean},
		}
		assertColumns(t, expectedColumns, sourceTableColumns)

		likeTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		likeTableColumns := getTableColumnsFor(t, client, likeTable.ID())
		assertColumns(t, expectedColumns, likeTableColumns)
	})

	t.Run("create table clone", func(t *testing.T) {
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("id", "NUMBER"),
			*sdk.NewTableColumnRequest("col2", "VARCHAR"),
			*sdk.NewTableColumnRequest("col3", "BOOLEAN"),
		}
		sourceTable, sourceTableCleanup := createTableWithColumns(t, client, database, schema, columns)
		t.Cleanup(sourceTableCleanup)

		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		request := sdk.NewCreateTableCloneRequest(id, sourceTable.ID()).
			WithCopyGrants(sdk.Bool(true)).WithClonePoint(sdk.NewClonePointRequest().
			WithAt(*sdk.NewTimeTravelRequest().WithOffset(sdk.Pointer(0))).
			WithMoment(sdk.CloneMomentAt))

		err := client.Tables.CreateClone(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		sourceTableColumns := getTableColumnsFor(t, client, sourceTable.ID())
		expectedColumns := []expectedColumn{
			{"id", sdk.DataTypeNumber},
			{"col2", sdk.DataTypeVARCHAR},
			{"col3", sdk.DataTypeBoolean},
		}
		assertColumns(t, expectedColumns, sourceTableColumns)

		cloneTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		cloneTableColumns := getTableColumnsFor(t, client, cloneTable.ID())
		assertColumns(t, expectedColumns, cloneTableColumns)
	})

	t.Run("alter table: rename", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		newName := random.String()
		newId := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, newName)

		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_3", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)

		alterRequest := sdk.NewAlterTableRequest(id).WithNewName(&newId)
		err = client.Tables.Alter(ctx, alterRequest)
		if err != nil {
			t.Cleanup(cleanupTableProvider(id))
		} else {
			t.Cleanup(cleanupTableProvider(newId))
		}
		require.NoError(t, err)

		_, err = client.Tables.ShowByID(ctx, id)
		assert.ErrorIs(t, err, collections.ErrObjectNotFound)

		table, err := client.Tables.ShowByID(ctx, newId)
		require.NoError(t, err)
		assertTable(t, table, newId)
	})

	t.Run("alter table: swap with", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		secondTableName := random.String()
		secondTableId := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, secondTableName)
		secondTableColumns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		err = client.Tables.Create(ctx, sdk.NewCreateTableRequest(secondTableId, secondTableColumns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(secondTableId))

		alterRequest := sdk.NewAlterTableRequest(id).WithSwapWith(&secondTableId)
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, secondTableId)
		require.NoError(t, err)

		assertTable(t, table, secondTableId)

		secondTable, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, secondTable, id)
	})

	t.Run("alter table: cluster by", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		clusterByColumns := []string{"COLUMN_1", "COLUMN_2"}
		alterRequest := sdk.NewAlterTableRequest(id).WithClusteringAction(sdk.NewTableClusteringActionRequest().WithClusterBy(clusterByColumns))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTable(t, table, id)
		assert.Equal(t, "", table.Comment)
		clusterByString := "LINEAR(" + strings.Join(clusterByColumns, ", ") + ")"
		assert.Equal(t, clusterByString, table.ClusterBy)
	})

	t.Run("alter table: resume recluster", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithClusteringAction(sdk.NewTableClusteringActionRequest().
				WithChangeReclusterState(sdk.Pointer(sdk.ReclusterStateSuspend)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		clusterByString := "LINEAR(" + strings.Join(clusterBy, ", ") + ")"
		assert.Equal(t, clusterByString, table.ClusterBy)
	})

	t.Run("alter table: drop clustering key", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithClusteringAction(sdk.NewTableClusteringActionRequest().
				WithDropClusteringKey(sdk.Bool(true)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", table.ClusterBy)
	})

	t.Run("alter table: add a column", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		clusterBy := []string{"COLUMN_1", "COLUMN_2"}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithClusterBy(clusterBy))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().
				WithAdd(sdk.NewTableColumnAddActionRequest("COLUMN_3", sdk.DataTypeVARCHAR)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
			{"COLUMN_3", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)

		assert.Equal(t, "", table.Comment)
	})

	t.Run("alter table: rename column", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().
				WithRename(sdk.NewTableColumnRenameActionRequest("COLUMN_1", "COLUMN_3")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)

		assert.Equal(t, "", table.Comment)
	})

	t.Run("alter table: unset masking policy", func(t *testing.T) {
		maskingPolicy, maskingPolicyCleanup := createMaskingPolicyWithOptions(t, client, database, schema, []sdk.TableColumnSignature{
			{
				Name: "col1",
				Type: sdk.DataTypeVARCHAR,
			},
		}, sdk.DataTypeVARCHAR, "REPLACE('X', 1)", nil)
		t.Cleanup(maskingPolicyCleanup)

		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR).WithMaskingPolicy(sdk.NewColumnMaskingPolicyRequest(maskingPolicy.ID())),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		tableDetails, err := client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
		require.NoError(t, err)

		require.Equal(t, 2, len(tableDetails))
		assert.Equal(t, maskingPolicy.ID().FullyQualifiedName(), *tableDetails[0].PolicyName)

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().WithUnsetMaskingPolicy(sdk.NewTableColumnAlterUnsetMaskingPolicyActionRequest("COLUMN_1")))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		tableDetails, err = client.Tables.DescribeColumns(ctx, sdk.NewDescribeTableColumnsRequest(id))
		require.NoError(t, err)

		require.Equal(t, 2, len(tableDetails))
		assert.Empty(t, tableDetails[0].PolicyName)
	})

	t.Run("alter table: set and unset tags", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		columnTags := []sdk.TagAssociationRequest{
			{
				Name:  sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, tag1.Name),
				Value: "v1",
			},
			{
				Name:  sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, tag2.Name),
				Value: "v2",
			},
		}

		alterRequest := sdk.NewAlterTableRequest(id).WithSetTags(columnTags)
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag1.ID(), id, sdk.ObjectTypeTable)
		require.NoError(t, err)

		assert.Equal(t, "v1", returnedTagValue)

		returnedTagValue, err = client.SystemFunctions.GetTag(ctx, tag2.ID(), id, sdk.ObjectTypeTable)
		require.NoError(t, err)

		assert.Equal(t, "v2", returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag1.ID(),
			tag2.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterTableRequest(id).WithUnsetTags(unsetTags)

		err = client.Tables.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag1.ID(), id, sdk.ObjectTypeTable)
		require.Error(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag2.ID(), id, sdk.ObjectTypeTable)
		require.Error(t, err)
	})

	t.Run("alter table: drop columns", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithColumnAction(sdk.NewTableColumnActionRequest().WithDropColumns([]string{"COLUMN_1"}))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)

		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)

		assert.Equal(t, table.Comment, "")
	})

	// TODO: check added constraints
	// Add method similar to getTableColumnsFor based on https://docs.snowflake.com/en/sql-reference/info-schema/table_constraints.
	t.Run("alter constraint: add", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
		}

		secondTableName := random.String()
		secondTableId := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, secondTableName)
		secondTableColumns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_3", sdk.DataTypeVARCHAR).WithInlineConstraint(sdk.NewColumnInlineConstraintRequest("pkey", sdk.ColumnConstraintTypePrimaryKey)),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		err = client.Tables.Create(ctx, sdk.NewCreateTableRequest(secondTableId, secondTableColumns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(secondTableId))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().
				WithAdd(sdk.NewOutOfLineConstraintRequest("OUT_OF_LINE_CONSTRAINT", sdk.ColumnConstraintTypeForeignKey).WithColumns([]string{"COLUMN_1"}).
					WithForeignKey(sdk.NewOutOfLineForeignKeyRequest(sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, secondTableName), []string{"COLUMN_3"}))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO: check renamed constraint
	t.Run("alter constraint: rename", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		oldConstraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(oldConstraintName, sdk.ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		newConstraintName := "NEW_OUT_OF_LINE_CONSTRAINT_NAME"
		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().
				WithRename(sdk.NewTableConstraintRenameActionRequest().
					WithOldName(oldConstraintName).
					WithNewName(newConstraintName)))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO: check altered constraint
	t.Run("alter constraint: alter", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(constraintName, sdk.ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().WithAlter(sdk.NewTableConstraintAlterActionRequest().WithConstraintName(sdk.String(constraintName)).WithEnforced(sdk.Bool(true))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO: check dropped constraint
	t.Run("alter constraint: drop", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}
		constraintName := "OUT_OF_LINE_CONSTRAINT"
		outOfLineConstraint := sdk.NewOutOfLineConstraintRequest(constraintName, sdk.ColumnConstraintTypePrimaryKey).WithColumns([]string{"COLUMN_1"})

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns).WithOutOfLineConstraint(*outOfLineConstraint))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithConstraintAction(sdk.NewTableConstraintActionRequest().WithDrop(sdk.NewTableConstraintDropActionRequest().WithConstraintName(sdk.String(constraintName))))
		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	t.Run("external table: add", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(sdk.NewTableExternalTableActionRequest().WithAdd(sdk.NewTableExternalTableColumnAddActionRequest().WithName("COLUMN_3").WithType(sdk.DataTypeNumber).WithExpression("1 + 1")))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
			{"COLUMN_3", sdk.DataTypeNumber},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	t.Run("external table: rename", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(sdk.NewTableExternalTableActionRequest().WithRename(sdk.NewTableExternalTableColumnRenameActionRequest().WithOldName("COLUMN_1").WithNewName("COLUMN_3")))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, table.Comment, "")
		currentColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_3", sdk.DataTypeVARCHAR},
			{"COLUMN_2", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	t.Run("external table: drop", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithExternalTableAction(sdk.NewTableExternalTableActionRequest().WithDrop(sdk.NewTableExternalTableColumnDropActionRequest().WithColumns([]string{"COLUMN_2"})))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
		table, err := client.Tables.ShowByID(ctx, id)
		require.NoError(t, err)

		currentColumns := getTableColumnsFor(t, client, table.ID())
		expectedColumns := []expectedColumn{
			{"COLUMN_1", sdk.DataTypeVARCHAR},
		}
		assertColumns(t, expectedColumns, currentColumns)
	})

	// TODO: check search optimization - after adding https://docs.snowflake.com/en/sql-reference/sql/desc-search-optimization
	t.Run("add search optimization", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		alterRequest := sdk.NewAlterTableRequest(id).
			WithSearchOptimizationAction(sdk.NewTableSearchOptimizationActionRequest().WithAddSearchOptimizationOn([]string{"SUBSTRING(*)", "GEO(*)"}))

		err = client.Tables.Alter(ctx, alterRequest)
		require.NoError(t, err)
	})

	// TODO: try to check more sets (ddl collation, max data extension time in days, etc.)
	t.Run("set: with complete options", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		comment := random.String()
		columns := []sdk.TableColumnRequest{
			*sdk.NewTableColumnRequest("COLUMN_1", sdk.DataTypeVARCHAR),
			*sdk.NewTableColumnRequest("COLUMN_2", sdk.DataTypeVARCHAR),
		}

		err := client.Tables.Create(ctx, sdk.NewCreateTableRequest(id, columns))
		require.NoError(t, err)
		t.Cleanup(cleanupTableProvider(id))

		stageFileFormats := sdk.StageFileFormatRequest{
			Type: sdk.Pointer(sdk.FileFormatTypeCSV),
		}
		stageCopyOptions := sdk.StageCopyOptionsRequest{
			OnError: sdk.NewStageCopyOnErrorOptionsRequest().WithSkipFile(),
		}
		alterRequest := sdk.NewAlterTableRequest(id).
			WithSet(sdk.NewTableSetRequest().
				WithEnableSchemaEvolution(sdk.Bool(true)).
				WithStageFileFormat(stageFileFormats).
				WithStageCopyOptions(stageCopyOptions).
				WithDataRetentionTimeInDays(sdk.Int(30)).
				WithMaxDataExtensionTimeInDays(sdk.Int(90)).
				WithChangeTracking(sdk.Bool(false)).
				WithDefaultDDLCollation(sdk.String("us")).
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

	t.Run("drop table", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		err := client.Tables.Drop(ctx, sdk.NewDropTableRequest(table.ID()).WithIfExists(sdk.Bool(true)))
		if err != nil {
			t.Cleanup(tableCleanup)
		}
		require.NoError(t, err)

		_, err = client.Tables.ShowByID(ctx, table.ID())
		require.ErrorIs(t, err, collections.ErrObjectNotFound)
	})

	t.Run("show tables", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)
		table2, table2Cleanup := createTable(t, client, database, schema)
		t.Cleanup(table2Cleanup)

		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest())
		require.NoError(t, err)

		t1, err := collections.FindOne(tables, func(t sdk.Table) bool { return t.ID().FullyQualifiedName() == table.ID().FullyQualifiedName() })
		require.NoError(t, err)
		t2, err := collections.FindOne(tables, func(t sdk.Table) bool { return t.ID().FullyQualifiedName() == table2.ID().FullyQualifiedName() })
		require.NoError(t, err)

		assertTable(t, t1, table.ID())
		assertTable(t, t2, table2.ID())
	})

	t.Run("with terse", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)

		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithTerse(sdk.Bool(true)).WithLikePattern(table.ID().Name()))
		require.NoError(t, err)
		assert.Equal(t, 1, len(tables))

		assertTableTerse(t, &tables[0], table.ID())
	})

	t.Run("with starts with", func(t *testing.T) {
		table, tableCleanup := createTable(t, client, database, schema)
		t.Cleanup(tableCleanup)

		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithStartsWith(sdk.String(table.Name)))
		require.NoError(t, err)
		assert.Equal(t, 1, len(tables))

		assertTable(t, &tables[0], table.ID())
	})

	t.Run("when searching a non-existent table", func(t *testing.T) {
		tables, err := client.Tables.Show(ctx, sdk.NewShowTableRequest().WithLikePattern("non-existent"))
		require.NoError(t, err)
		assert.Equal(t, 0, len(tables))
	})
}
