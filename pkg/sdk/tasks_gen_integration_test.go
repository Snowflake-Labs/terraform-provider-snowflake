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
		assert.NotEmpty(t, task.CreatedOn)
		// TODO: fill out
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
		assertTask(t, task, id)
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
		assertTask(t, task, id)
	})

	t.Run("drop task: existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("drop task: non-existing", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter task: set value and unset value", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter task: set and unset tag", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter task: suspend and resume", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter task: remove after and add after", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("alter task: modify when and as", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show task: default", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show task: terse", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("show task: with options", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("describe task: default", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("execute task: default", func(t *testing.T) {
		// TODO: fill me
	})
}
