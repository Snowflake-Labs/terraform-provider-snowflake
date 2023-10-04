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

	t.Run("Create", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Alter", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Drop", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Show", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("ShowByID", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Describe", func(t *testing.T) {
		// TODO: fill me
	})

	t.Run("Execute", func(t *testing.T) {
		// TODO: fill me
	})
}
