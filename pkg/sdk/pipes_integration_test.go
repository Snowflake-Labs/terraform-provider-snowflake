package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_PipesShow(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	database, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchema(t, client, database)
	t.Cleanup(schemaCleanup)

	table1, table1Cleanup := createTable(t, client, database, schema)
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := createTable(t, client, database, schema)
	t.Cleanup(table2Cleanup)

	stageName := randomString(t)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	createCopyStatement := func(table *Table, stage *Stage) string {
		require.NotNil(t, table, "table has to be created")
		require.NotNil(t, stage, "stage has to be created")
		return fmt.Sprintf("COPY INTO %s FROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
	}

	pipe1Name := randomString(t)
	pipe1CopyStatement := createCopyStatement(table1, stage)
	pipe1, pipe1Cleanup := createPipe(t, client, database, schema, pipe1Name, pipe1CopyStatement)
	t.Cleanup(pipe1Cleanup)

	pipe2Name := randomString(t)
	pipe2CopyStatement := createCopyStatement(table2, stage)
	pipe2, pipe2Cleanup := createPipe(t, client, database, schema, pipe2Name, pipe2CopyStatement)
	t.Cleanup(pipe2Cleanup)

	t.Run("without show options", func(t *testing.T) {
		pipes, err := client.Pipes.Show(ctx, &PipeShowOptions{})

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, pipe1)
		assert.Contains(t, pipes, pipe2)
	})

	t.Run("show in schema", func(t *testing.T) {
		showOptions := &PipeShowOptions{
			In: &In{
				Schema: schema.ID(),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, pipe1)
		assert.Contains(t, pipes, pipe2)
	})

	t.Run("show like", func(t *testing.T) {
		showOptions := &PipeShowOptions{
			Like: &Like{
				Pattern: String(pipe1Name),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 1, len(pipes))
		assert.Contains(t, pipes, pipe1)
	})

	t.Run("search for non-existent pipe", func(t *testing.T) {
		showOptions := &PipeShowOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 0, len(pipes))
	})
}
