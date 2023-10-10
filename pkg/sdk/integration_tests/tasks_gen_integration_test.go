package sdk_integration_tests

import (
	"context"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Tasks(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := sdk.createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := sdk.createSchema(t, client, database)
	t.Cleanup(schemaCleanup)

	sql := "SELECT CURRENT_TIMESTAMP"

	assertTask := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.name, task.Name)
		assert.NotEmpty(t, task.Id)
		assert.Equal(t, database.Name, task.DatabaseName)
		assert.Equal(t, schema.Name, task.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", task.Owner)
		assert.Equal(t, "", task.Comment)
		assert.Equal(t, "", task.Warehouse)
		assert.Equal(t, "", task.Schedule)
		assert.Equal(t, "[]", task.Predecessors)
		assert.Equal(t, "suspended", task.State)
		assert.Equal(t, sql, task.Definition)
		assert.Equal(t, "", task.Condition)
		assert.Equal(t, false, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Equal(t, "ROLE", task.OwnerRoleType)
		assert.Empty(t, task.Config)
		assert.Empty(t, task.Budget)
	}

	assertTaskWithOptions := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, comment string, warehouse string, schedule string, condition string, allowOverlappingExecution bool, config string, predecessor *sdk.SchemaObjectIdentifier) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.name, task.Name)
		assert.NotEmpty(t, task.Id)
		assert.Equal(t, database.Name, task.DatabaseName)
		assert.Equal(t, schema.Name, task.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", task.Owner)
		assert.Equal(t, comment, task.Comment)
		assert.Equal(t, warehouse, task.Warehouse)
		assert.Equal(t, schedule, task.Schedule)
		assert.Equal(t, "suspended", task.State)
		assert.Equal(t, sql, task.Definition)
		assert.Equal(t, condition, task.Condition)
		assert.Equal(t, allowOverlappingExecution, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Equal(t, "ROLE", task.OwnerRoleType)
		assert.Equal(t, config, task.Config)
		assert.Empty(t, task.Budget)
		if predecessor != nil {
			// Predecessors list is formatted, so matching it is unnecessarily complicated:
			// e.g. `[\n  \"\\\"qgb)Z1KcNWJ(\\\".\\\"glN@JtR=7dzP$7\\\".\\\"_XEL(7N_F?@frgT5>dQS>V|vSy,J\\\"\"\n]`.
			// We just match parts of the expected predecessor. Later we can parse the output while constructing Task object.
			assert.Contains(t, task.Predecessors, predecessor.DatabaseName())
			assert.Contains(t, task.Predecessors, predecessor.SchemaName())
			assert.Contains(t, task.Predecessors, predecessor.Name())
		} else {
			assert.Equal(t, "[]", task.Predecessors)
		}
	}

	assertTaskTerse := func(t *testing.T, task *sdk.Task, id sdk.SchemaObjectIdentifier, schedule string) {
		t.Helper()
		assert.Equal(t, id, task.ID())
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, id.name, task.Name)
		assert.Equal(t, database.Name, task.DatabaseName)
		assert.Equal(t, schema.Name, task.SchemaName)
		assert.Equal(t, schedule, task.Schedule)

		// all below are not contained in the terse response, that's why all of them we expect to be empty
		assert.Empty(t, task.Id)
		assert.Empty(t, task.Owner)
		assert.Empty(t, task.Comment)
		assert.Empty(t, task.Warehouse)
		assert.Empty(t, task.Predecessors)
		assert.Empty(t, task.State)
		assert.Empty(t, task.Definition)
		assert.Empty(t, task.Condition)
		assert.Empty(t, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Empty(t, task.OwnerRoleType)
		assert.Empty(t, task.Config)
		assert.Empty(t, task.Budget)
	}

	cleanupTaskProvider := func(id sdk.SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
			require.NoError(t, err)
		}
	}

	createTaskBasicRequest := func(t *testing.T) *sdk.CreateTaskRequest {
		t.Helper()
		name := sdk.randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		return sdk.NewCreateTaskRequest(id, sql)
	}

	createTaskWithRequest := func(t *testing.T, request *sdk.CreateTaskRequest) *sdk.Task {
		t.Helper()
		id := request.name

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		return task
	}

	createTask := func(t *testing.T) *sdk.Task {
		t.Helper()
		return createTaskWithRequest(t, createTaskBasicRequest(t))
	}

	t.Run("create task: no optionals", func(t *testing.T) {
		request := createTaskBasicRequest(t)

		task := createTaskWithRequest(t, request)

		assertTask(t, task, request.name)
	})

	t.Run("create task: with initial warehouse", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithWarehouse(sdk.NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(&sdk.WarehouseSizeXSmall))

		task := createTaskWithRequest(t, request)

		assertTask(t, task, request.name)
	})

	t.Run("create task: almost complete case", func(t *testing.T) {
		warehouse, warehouseCleanup := sdk.createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		request := createTaskBasicRequest(t).
			WithOrReplace(sdk.Bool(true)).
			WithWarehouse(sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(warehouse.ID()))).
			WithSchedule(sdk.String("10 MINUTE")).
			WithConfig(sdk.String(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`)).
			WithAllowOverlappingExecution(sdk.Bool(true)).
			WithSessionParameters(&sdk.SessionParameters{
				JSONIndent: sdk.Int(4),
			}).
			WithUserTaskTimeoutMs(sdk.Int(500)).
			WithSuspendTaskAfterNumFailures(sdk.Int(3)).
			WithComment(sdk.String("some comment")).
			WithWhen(sdk.String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))
		id := request.name

		task := createTaskWithRequest(t, request)

		assertTaskWithOptions(t, task, id, "some comment", warehouse.Name, "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, true, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`, nil)
	})

	t.Run("create task: with after", func(t *testing.T) {
		otherName := sdk.randomString(t)
		otherId := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, otherName)

		request := sdk.NewCreateTaskRequest(otherId, sql).WithSchedule(sdk.String("10 MINUTE"))

		createTaskWithRequest(t, request)

		request = createTaskBasicRequest(t).
			WithAfter([]sdk.SchemaObjectIdentifier{otherId})

		task := createTaskWithRequest(t, request)

		assertTaskWithOptions(t, task, request.name, "", "", "", "", false, "", &otherId)
	})

	// TODO: this fails with `syntax error line 1 at position 89 unexpected 'GRANTS'`.
	// The reason is that in the documentation there is a note: "This parameter is not supported currently.".
	// t.Run("create task: with grants", func(t *testing.T) {
	//	name := randomString(t)
	//	id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	//
	//	request := NewCreateTaskRequest(id, sql).
	//		WithOrReplace(Bool(true)).
	//		WithCopyGrants(Bool(true))
	//
	//	err := client.Tasks.Create(ctx, request)
	//	require.NoError(t, err)
	//	t.Cleanup(cleanupTaskProvider(id))
	//
	//	task, err := client.Tasks.ShowByID(ctx, id)
	//
	//	require.NoError(t, err)
	//	assertTaskWithOptions(t, task, id, name, "", "", "", "", false, "", nil)
	// })

	t.Run("create task: with tags", func(t *testing.T) {
		tag, tagCleanup := sdk.createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		request := createTaskBasicRequest(t).
			WithTag([]sdk.TagAssociation{{
				Name:  tag.ID(),
				Value: "v1",
			}})

		task := createTaskWithRequest(t, request)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, "v1", returnedTagValue)
	})

	t.Run("clone task: default", func(t *testing.T) {
		sourceTask := createTask(t)

		name := sdk.randomString(t)
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		request := sdk.NewCloneTaskRequest(id, sourceTask.ID())

		err := client.Tasks.Clone(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assertTask(t, task, request.name)
	})

	t.Run("drop task: existing", func(t *testing.T) {
		request := createTaskBasicRequest(t)
		id := request.name

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)

		err = client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
		require.NoError(t, err)

		_, err = client.Tasks.ShowByID(ctx, id)
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("drop task: non-existing", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, "does_not_exist")

		err := client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
		assert.ErrorIs(t, err, sdk.errObjectNotExistOrAuthorized)
	})

	t.Run("alter task: set value and unset value", func(t *testing.T) {
		task := createTask(t)
		id := task.ID()

		alterRequest := sdk.NewAlterTaskRequest(id).WithSet(sdk.NewTaskSetRequest().WithComment(sdk.String("new comment")))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredTask.Comment)

		alterRequest = sdk.NewAlterTaskRequest(id).WithUnset(sdk.NewTaskUnsetRequest().WithComment(sdk.Bool(true)))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredTask.Comment)
	})

	t.Run("alter task: set and unset tag", func(t *testing.T) {
		tag, tagCleanup := sdk.createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		task := createTask(t)
		id := task.ID()

		tagValue := "abc"
		tags := []sdk.TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := sdk.NewAlterTaskRequest(id).WithSetTags(tags)

		err := client.Tasks.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []sdk.ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := sdk.NewAlterTaskRequest(id).WithUnsetTags(unsetTags)

		err = client.Tasks.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, sdk.ObjectTypeTask)
		require.Error(t, err)
	})

	t.Run("alter task: resume and suspend", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(sdk.String("10 MINUTE"))

		task := createTaskWithRequest(t, request)
		id := task.ID()

		assert.Equal(t, "suspended", task.State)

		alterRequest := sdk.NewAlterTaskRequest(id).WithResume(sdk.Bool(true))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "started", alteredTask.State)

		alterRequest = sdk.NewAlterTaskRequest(id).WithSuspend(sdk.Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "suspended", alteredTask.State)
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(sdk.String("10 MINUTE"))

		otherTask := createTaskWithRequest(t, request)
		otherId := otherTask.ID()

		request = createTaskBasicRequest(t).
			WithAfter([]sdk.SchemaObjectIdentifier{otherId})

		task := createTaskWithRequest(t, request)
		id := task.ID()

		assert.Contains(t, task.Predecessors, otherId.Name())

		alterRequest := sdk.NewAlterTaskRequest(id).WithRemoveAfter([]sdk.SchemaObjectIdentifier{otherId})

		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, "[]", task.Predecessors)

		alterRequest = sdk.NewAlterTaskRequest(id).WithAddAfter([]sdk.SchemaObjectIdentifier{otherId})

		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assert.Contains(t, task.Predecessors, otherId.Name())
	})

	t.Run("alter task: modify when and as", func(t *testing.T) {
		task := createTask(t)
		id := task.ID()

		newSql := "SELECT CURRENT_DATE"
		alterRequest := sdk.NewAlterTaskRequest(id).WithModifyAs(sdk.String(newSql))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, newSql, alteredTask.Definition)

		newWhen := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
		alterRequest = sdk.NewAlterTaskRequest(id).WithModifyWhen(sdk.String(newWhen))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, newWhen, alteredTask.Condition)
	})

	t.Run("show task: default", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := sdk.NewShowTaskRequest()
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.Contains(t, returnedTasks, *task2)
	})

	t.Run("show task: terse", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(sdk.String("10 MINUTE"))

		task := createTaskWithRequest(t, request)

		showRequest := sdk.NewShowTaskRequest().WithTerse(sdk.Bool(true))
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedTasks))
		assertTaskTerse(t, &returnedTasks[0], task.ID(), "10 MINUTE")
	})

	t.Run("show task: with options", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := sdk.NewShowTaskRequest().
			WithLike(&sdk.Like{&task1.Name}).
			WithIn(&sdk.In{Schema: sdk.NewDatabaseObjectIdentifier(database.Name, schema.Name)}).
			WithLimit(&sdk.LimitFrom{Rows: sdk.Int(5)})
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)

		require.NoError(t, err)
		assert.Equal(t, 1, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.NotContains(t, returnedTasks, *task2)
	})

	t.Run("describe task: default", func(t *testing.T) {
		task := createTask(t)

		returnedTask, err := client.Tasks.Describe(ctx, task.ID())
		require.NoError(t, err)

		assertTask(t, returnedTask, task.ID())
	})

	t.Run("execute task: default", func(t *testing.T) {
		task := createTask(t)

		executeRequest := sdk.NewExecuteTaskRequest(task.ID())
		err := client.Tasks.Execute(ctx, executeRequest)
		require.NoError(t, err)
	})
}
