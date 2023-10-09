package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_Tasks(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchema(t, client, database)
	t.Cleanup(schemaCleanup)

	sql := "SELECT CURRENT_TIMESTAMP"

	assertTask := func(t *testing.T, task *Task, id SchemaObjectIdentifier) {
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

	assertTaskWithOptions := func(t *testing.T, task *Task, id SchemaObjectIdentifier, comment string, warehouse string, schedule string, condition string, allowOverlappingExecution bool, config string, predecessor *SchemaObjectIdentifier) {
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

	assertTaskTerse := func(t *testing.T, task *Task, id SchemaObjectIdentifier, schedule string) {
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

	cleanupTaskProvider := func(id SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tasks.Drop(ctx, NewDropTaskRequest(id))
			require.NoError(t, err)
		}
	}

	createTaskBasicRequest := func(t *testing.T) *CreateTaskRequest {
		t.Helper()
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		return NewCreateTaskRequest(id, sql)
	}

	createTaskWithRequest := func(t *testing.T, request *CreateTaskRequest) *Task {
		t.Helper()
		id := request.name

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		return task
	}

	createTask := func(t *testing.T) *Task {
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
			WithWarehouse(NewCreateTaskWarehouseRequest().WithUserTaskManagedInitialWarehouseSize(&WarehouseSizeXSmall))

		task := createTaskWithRequest(t, request)

		assertTask(t, task, request.name)
	})

	t.Run("create task: almost complete case", func(t *testing.T) {
		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		request := createTaskBasicRequest(t).
			WithOrReplace(Bool(true)).
			WithWarehouse(NewCreateTaskWarehouseRequest().WithWarehouse(Pointer(warehouse.ID()))).
			WithSchedule(String("10 MINUTE")).
			WithConfig(String(`$${"output_dir": "/temp/test_directory/", "learning_rate": 0.1}$$`)).
			WithAllowOverlappingExecution(Bool(true)).
			WithSessionParameters(&SessionParameters{
				JSONIndent: Int(4),
			}).
			WithUserTaskTimeoutMs(Int(500)).
			WithSuspendTaskAfterNumFailures(Int(3)).
			WithComment(String("some comment")).
			WithWhen(String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))
		id := request.name

		task := createTaskWithRequest(t, request)

		assertTaskWithOptions(t, task, id, "some comment", warehouse.Name, "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, true, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`, nil)
	})

	t.Run("create task: with after", func(t *testing.T) {
		otherName := randomString(t)
		otherId := NewSchemaObjectIdentifier(database.Name, schema.Name, otherName)

		request := NewCreateTaskRequest(otherId, sql).WithSchedule(String("10 MINUTE"))

		createTaskWithRequest(t, request)

		request = createTaskBasicRequest(t).
			WithAfter([]SchemaObjectIdentifier{otherId})

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
		tag, tagCleanup := createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		request := createTaskBasicRequest(t).
			WithTag([]TagAssociation{{
				Name:  tag.ID(),
				Value: "v1",
			}})

		task := createTaskWithRequest(t, request)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), task.ID(), ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, "v1", returnedTagValue)
	})

	t.Run("drop task: existing", func(t *testing.T) {
		request := createTaskBasicRequest(t)
		id := request.name

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)

		err = client.Tasks.Drop(ctx, NewDropTaskRequest(id))
		require.NoError(t, err)

		_, err = client.Tasks.ShowByID(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("drop task: non-existing", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, "does_not_exist")

		err := client.Tasks.Drop(ctx, NewDropTaskRequest(id))
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("alter task: set value and unset value", func(t *testing.T) {
		task := createTask(t)
		id := task.ID()

		alterRequest := NewAlterTaskRequest(id).WithSet(NewTaskSetRequest().WithComment(String("new comment")))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredTask.Comment)

		alterRequest = NewAlterTaskRequest(id).WithUnset(NewTaskUnsetRequest().WithComment(Bool(true)))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "", alteredTask.Comment)
	})

	t.Run("alter task: set and unset tag", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		task := createTask(t)
		id := task.ID()

		tagValue := "abc"
		tags := []TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := NewAlterTaskRequest(id).WithSetTags(tags)

		err := client.Tasks.Alter(ctx, alterRequestSetTags)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), id, ObjectTypeTask)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		unsetTags := []ObjectIdentifier{
			tag.ID(),
		}
		alterRequestUnsetTags := NewAlterTaskRequest(id).WithUnsetTags(unsetTags)

		err = client.Tasks.Alter(ctx, alterRequestUnsetTags)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), id, ObjectTypeTask)
		require.Error(t, err)
	})

	t.Run("alter task: resume and suspend", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(String("10 MINUTE"))

		task := createTaskWithRequest(t, request)
		id := task.ID()

		assert.Equal(t, "suspended", task.State)

		alterRequest := NewAlterTaskRequest(id).WithResume(Bool(true))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "started", alteredTask.State)

		alterRequest = NewAlterTaskRequest(id).WithSuspend(Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "suspended", alteredTask.State)
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(String("10 MINUTE"))

		otherTask := createTaskWithRequest(t, request)
		otherId := otherTask.ID()

		request = createTaskBasicRequest(t).
			WithAfter([]SchemaObjectIdentifier{otherId})

		task := createTaskWithRequest(t, request)
		id := task.ID()

		assert.Contains(t, task.Predecessors, otherId.Name())

		alterRequest := NewAlterTaskRequest(id).WithRemoveAfter([]SchemaObjectIdentifier{otherId})

		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		task, err = client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, "[]", task.Predecessors)

		alterRequest = NewAlterTaskRequest(id).WithAddAfter([]SchemaObjectIdentifier{otherId})

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
		alterRequest := NewAlterTaskRequest(id).WithModifyAs(String(newSql))
		err := client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, newSql, alteredTask.Definition)

		newWhen := `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`
		alterRequest = NewAlterTaskRequest(id).WithModifyWhen(String(newWhen))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, newWhen, alteredTask.Condition)
	})

	t.Run("show task: default", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := NewShowTaskRequest()
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 2, len(returnedTasks))
		assert.Contains(t, returnedTasks, *task1)
		assert.Contains(t, returnedTasks, *task2)
	})

	t.Run("show task: terse", func(t *testing.T) {
		request := createTaskBasicRequest(t).
			WithSchedule(String("10 MINUTE"))

		task := createTaskWithRequest(t, request)

		showRequest := NewShowTaskRequest().WithTerse(Bool(true))
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedTasks))
		assertTaskTerse(t, &returnedTasks[0], task.ID(), "10 MINUTE")
	})

	t.Run("show task: with options", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := NewShowTaskRequest().
			WithLike(&Like{&task1.Name}).
			WithIn(&In{Schema: NewDatabaseObjectIdentifier(database.Name, schema.Name)}).
			WithLimit(&LimitFrom{Rows: Int(5)})
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

		executeRequest := NewExecuteTaskRequest(task.ID())
		err := client.Tasks.Execute(ctx, executeRequest)
		require.NoError(t, err)
	})
}
