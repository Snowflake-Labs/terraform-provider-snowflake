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

	assertTask := func(t *testing.T, task *Task, id SchemaObjectIdentifier, name string) {
		t.Helper()
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, name, task.Name)
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

	assertTaskWithOptions := func(t *testing.T, task *Task, id SchemaObjectIdentifier, name string, comment string, warehouse string, schedule string, condition string, config string) {
		t.Helper()
		assert.NotEmpty(t, task.CreatedOn)
		assert.Equal(t, name, task.Name)
		assert.NotEmpty(t, task.Id)
		assert.Equal(t, database.Name, task.DatabaseName)
		assert.Equal(t, schema.Name, task.SchemaName)
		assert.Equal(t, "ACCOUNTADMIN", task.Owner)
		assert.Equal(t, comment, task.Comment)
		assert.Equal(t, warehouse, task.Warehouse)
		assert.Equal(t, schedule, task.Schedule)
		assert.Equal(t, "[]", task.Predecessors)
		assert.Equal(t, "suspended", task.State)
		assert.Equal(t, sql, task.Definition)
		assert.Equal(t, condition, task.Condition)
		assert.Equal(t, true, task.AllowOverlappingExecution)
		assert.Empty(t, task.ErrorIntegration)
		assert.Empty(t, task.LastCommittedOn)
		assert.Empty(t, task.LastSuspendedOn)
		assert.Equal(t, "ROLE", task.OwnerRoleType)
		assert.Equal(t, config, task.Config)
		assert.Empty(t, task.Budget)
	}

	assertTaskTerse := func(t *testing.T, task *Task, id SchemaObjectIdentifier) {
		t.Helper()
		assert.NotEmpty(t, task.CreatedOn)
		// TODO: fill out
	}

	cleanupTaskProvider := func(id SchemaObjectIdentifier) func() {
		return func() {
			err := client.Tasks.Drop(ctx, NewDropTaskRequest(id))
			require.NoError(t, err)
		}
	}

	createTask := func(t *testing.T) *Task {
		t.Helper()
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Tasks.Create(ctx, NewCreateTaskRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		return task
	}

	_, _ = assertTaskTerse, createTask

	t.Run("create task: no optionals", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		request := NewCreateTaskRequest(id, sql)

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assertTask(t, task, id, name)
	})

	t.Run("create task: almost complete case", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		warehouse, warehouseCleanup := createWarehouse(t, client)
		t.Cleanup(warehouseCleanup)

		//tag, tagCleanup := createTag(t, client, database, schema)
		//t.Cleanup(tagCleanup)

		//otherTask := createTask(t)

		request := NewCreateTaskRequest(id, sql).
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
			//WithCopyGrants(Bool(true)).
			//WithAfter([]SchemaObjectIdentifier{otherTask.ID()}).
			WithComment(String("some comment")).
			//WithTag([]TagAssociation{{
			//	Name:  tag.ID(),
			//	Value: "v1",
			//}}).
			WithWhen(String(`SYSTEM$STREAM_HAS_DATA('MYSTREAM')`))

		err := client.Tasks.Create(ctx, request)
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		task, err := client.Tasks.ShowByID(ctx, id)

		require.NoError(t, err)
		assertTaskWithOptions(t, task, id, name, "some comment", warehouse.Name, "10 MINUTE", `SYSTEM$STREAM_HAS_DATA('MYSTREAM')`, `{"output_dir": "/temp/test_directory/", "learning_rate": 0.1}`)
	})

	t.Run("drop task: existing", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Tasks.Create(ctx, NewCreateTaskRequest(id, sql))
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
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Tasks.Create(ctx, NewCreateTaskRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		alterRequest := NewAlterTaskRequest(id).WithSet(NewTaskSetRequest().WithComment(String("new comment")))
		err = client.Tasks.Alter(ctx, alterRequest)
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

		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Tasks.Create(ctx, NewCreateTaskRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		tagValue := "abc"
		tags := []TagAssociation{
			{
				Name:  tag.ID(),
				Value: tagValue,
			},
		}
		alterRequestSetTags := NewAlterTaskRequest(id).WithSetTags(tags)

		err = client.Tasks.Alter(ctx, alterRequestSetTags)
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

	t.Run("alter task: suspend and resume", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Tasks.Create(ctx, NewCreateTaskRequest(id, sql).WithSchedule(String("10 MINUTE")))
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		alterRequest := NewAlterTaskRequest(id).WithSuspend(Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err := client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "suspended", alteredTask.State)

		alterRequest = NewAlterTaskRequest(id).WithResume(Bool(true))
		err = client.Tasks.Alter(ctx, alterRequest)
		require.NoError(t, err)

		alteredTask, err = client.Tasks.ShowByID(ctx, id)
		require.NoError(t, err)

		assert.Equal(t, "started", alteredTask.State)
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter task: modify when and as", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Tasks.Create(ctx, NewCreateTaskRequest(id, sql))
		require.NoError(t, err)
		t.Cleanup(cleanupTaskProvider(id))

		newSql := "SELECT CURRENT_DATE"
		alterRequest := NewAlterTaskRequest(id).WithModifyAs(String(newSql))
		err = client.Tasks.Alter(ctx, alterRequest)
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
		task1 := createTask(t)

		showRequest := NewShowTaskRequest().WithTerse(Bool(true))
		returnedTasks, err := client.Tasks.Show(ctx, showRequest)
		require.NoError(t, err)

		assert.Equal(t, 1, len(returnedTasks))
		assertTaskTerse(t, &returnedTasks[0], task1.ID())
	})

	t.Run("show task: with options", func(t *testing.T) {
		task1 := createTask(t)
		task2 := createTask(t)

		showRequest := NewShowTaskRequest().
			WithLike(&Like{&task1.Name}).
			WithIn(&In{Schema: NewDatabaseObjectIdentifier(database.Name, schema.Name)}).
			WithLimit(Int(5))
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

		assertTask(t, returnedTask, task.ID(), task.Name)
	})

	t.Run("execute task: default", func(t *testing.T) {
		task := createTask(t)

		executeRequest := NewExecuteTaskRequest(task.ID())
		err := client.Tasks.Execute(ctx, executeRequest)
		require.NoError(t, err)
	})
}
