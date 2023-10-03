package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPipeCopyStatement(t *testing.T, table *Table, stage *Stage) string {
	t.Helper()
	require.NotNil(t, table, "table has to be created")
	require.NotNil(t, stage, "stage has to be created")
	return fmt.Sprintf("COPY INTO %s\nFROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
}

func TestInt_IncorrectCreatePipeBehaviour(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := NewDatabaseObjectIdentifier("TXR@=9,TBnLj", "tcK1>AJ+")
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, AccountObjectIdentifier{schemaIdentifier.databaseName})
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.name)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, client, database, schema)
	t.Cleanup(tableCleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	t.Run("if we have special characters in db or schema name, create pipe returns error in copy <> from <> section", func(t *testing.T) {
		err := client.Pipes.Create(
			ctx,
			NewSchemaObjectIdentifier(database.Name, schema.Name, randomAlphanumericN(t, 20)),
			createPipeCopyStatement(t, table, stage),
			&CreatePipeOptions{},
		)

		require.ErrorContains(t, err, "(42000): SQL compilation error:\nsyntax error line")
		require.ErrorContains(t, err, "at position")
		require.ErrorContains(t, err, "unexpected ','")
	})

	t.Run("the same works with using db and schema statements", func(t *testing.T) {
		useDatabaseCleanup := useDatabase(t, client, database.ID())
		t.Cleanup(useDatabaseCleanup)
		useSchemaCleanup := useSchema(t, client, schema.ID())
		t.Cleanup(useSchemaCleanup)

		createCopyStatementWithoutQualifiersForStage := func(t *testing.T, table *Table, stage *Stage) string {
			t.Helper()
			require.NotNil(t, table, "table has to be created")
			require.NotNil(t, stage, "stage has to be created")
			return fmt.Sprintf("COPY INTO %s\nFROM @\"%s\"", table.ID().FullyQualifiedName(), stage.Name)
		}

		err := client.Pipes.Create(
			ctx,
			NewSchemaObjectIdentifier(database.Name, schema.Name, randomAlphanumericN(t, 20)),
			createCopyStatementWithoutQualifiersForStage(t, table, stage),
			&CreatePipeOptions{},
		)

		require.NoError(t, err)
	})
}

func TestInt_PipesShowAndDescribe(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := alphanumericDatabaseObjectIdentifier(t)
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, AccountObjectIdentifier{schemaIdentifier.databaseName})
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.name)
	t.Cleanup(schemaCleanup)

	table1, table1Cleanup := createTable(t, client, database, schema)
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := createTable(t, client, database, schema)
	t.Cleanup(table2Cleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	pipe1Name := randomAlphanumericN(t, 20)
	pipe1CopyStatement := createPipeCopyStatement(t, table1, stage)
	pipe1, pipe1Cleanup := createPipe(t, client, database, schema, pipe1Name, pipe1CopyStatement)
	t.Cleanup(pipe1Cleanup)

	pipe2Name := randomAlphanumericN(t, 20)
	pipe2CopyStatement := createPipeCopyStatement(t, table2, stage)
	pipe2, pipe2Cleanup := createPipe(t, client, database, schema, pipe2Name, pipe2CopyStatement)
	t.Cleanup(pipe2Cleanup)

	t.Run("show: without options", func(t *testing.T) {
		pipes, err := client.Pipes.Show(ctx, &ShowPipeOptions{})

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, *pipe1)
		assert.Contains(t, pipes, *pipe2)
	})

	t.Run("show: in schema", func(t *testing.T) {
		showOptions := &ShowPipeOptions{
			In: &In{
				Schema: schema.ID(),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, *pipe1)
		assert.Contains(t, pipes, *pipe2)
	})

	t.Run("show: like", func(t *testing.T) {
		showOptions := &ShowPipeOptions{
			Like: &Like{
				Pattern: String(pipe1Name),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 1, len(pipes))
		assert.Contains(t, pipes, *pipe1)
	})

	t.Run("show: non-existent pipe", func(t *testing.T) {
		showOptions := &ShowPipeOptions{
			Like: &Like{
				Pattern: String("non-existent"),
			},
		}
		pipes, err := client.Pipes.Show(ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 0, len(pipes))
	})

	t.Run("describe: existing pipe", func(t *testing.T) {
		pipe, err := client.Pipes.Describe(ctx, pipe1.ID())

		require.NoError(t, err)
		assert.Equal(t, pipe1.Name, pipe.Name)
	})

	t.Run("describe: non-existing pipe", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(database.Name, database.Name, "does_not_exist")

		_, err := client.Pipes.Describe(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})
}

func TestInt_PipeCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := alphanumericDatabaseObjectIdentifier(t)
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, AccountObjectIdentifier{schemaIdentifier.databaseName})
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.name)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, client, database, schema)
	t.Cleanup(tableCleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)

	assertPipe := func(t *testing.T, pipeDetails *Pipe, expectedName string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, pipeDetails.CreatedOn)
		assert.Equal(t, expectedName, pipeDetails.Name)
		assert.Equal(t, database.Name, pipeDetails.DatabaseName)
		assert.Equal(t, schema.Name, pipeDetails.SchemaName)
		assert.Equal(t, copyStatement, pipeDetails.Definition)
		assert.Equal(t, "ACCOUNTADMIN", pipeDetails.Owner)
		assert.Empty(t, pipeDetails.NotificationChannel)
		assert.Equal(t, expectedComment, pipeDetails.Comment)
		assert.Empty(t, pipeDetails.Integration)
		assert.Empty(t, pipeDetails.Pattern)
		assert.Empty(t, pipeDetails.ErrorIntegration)
		assert.Equal(t, "ROLE", pipeDetails.OwnerRoleType)
		assert.Empty(t, pipeDetails.InvalidReason)
	}

	// TODO: test error integration, aws sns topic and integration when we have them in project
	t.Run("test complete case", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
		comment := randomComment(t)

		err := client.Pipes.Create(ctx, id, copyStatement, &CreatePipeOptions{
			OrReplace:   Bool(false),
			IfNotExists: Bool(true),
			AutoIngest:  Bool(false),
			Comment:     String(comment),
		})
		require.NoError(t, err)

		pipe, err := client.Pipes.Describe(ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, name, comment)
	})

	t.Run("test if not exists and or replace are incompatible", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Pipes.Create(ctx, id, copyStatement, &CreatePipeOptions{
			OrReplace:   Bool(true),
			IfNotExists: Bool(true),
		})
		require.ErrorContains(t, err, "(0A000): SQL compilation error:\noptions IF NOT EXISTS and OR REPLACE are incompatible")
	})

	t.Run("test no options", func(t *testing.T) {
		name := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)

		err := client.Pipes.Create(ctx, id, copyStatement, nil)
		require.NoError(t, err)

		pipe, err := client.Pipes.Describe(ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, name, "")
	})
}

func TestInt_PipeDrop(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := alphanumericDatabaseObjectIdentifier(t)
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, AccountObjectIdentifier{schemaIdentifier.databaseName})
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.name)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, client, database, schema)
	t.Cleanup(tableCleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	t.Run("pipe exists", func(t *testing.T) {
		pipeName := randomAlphanumericN(t, 20)
		pipeCopyStatement := createPipeCopyStatement(t, table, stage)
		pipe, _ := createPipe(t, client, database, schema, pipeName, pipeCopyStatement)

		err := client.Pipes.Drop(ctx, pipe.ID())

		require.NoError(t, err)
		_, err = client.Pipes.Describe(ctx, pipe.ID())
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})

	t.Run("pipe does not exist", func(t *testing.T) {
		id := NewSchemaObjectIdentifier(database.Name, database.Name, "does_not_exist")

		err := client.Alerts.Drop(ctx, id)
		assert.ErrorIs(t, err, errObjectNotExistOrAuthorized)
	})
}

func TestInt_PipeAlter(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := alphanumericDatabaseObjectIdentifier(t)
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, AccountObjectIdentifier{schemaIdentifier.databaseName})
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.name)
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, client, database, schema)
	t.Cleanup(tableCleanup)

	stageName := randomAlphanumericN(t, 20)
	stage, stageCleanup := createStage(t, client, database, schema, stageName)
	t.Cleanup(stageCleanup)

	pipeCopyStatement := createPipeCopyStatement(t, table, stage)

	// TODO: test error integration when we have them in project
	t.Run("set value and unset value", func(t *testing.T) {
		pipeName := randomAlphanumericN(t, 20)
		pipe, pipeCleanup := createPipe(t, client, database, schema, pipeName, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		alterOptions := &AlterPipeOptions{
			Set: &PipeSet{
				Comment:             String("new comment"),
				PipeExecutionPaused: Bool(true),
			},
		}

		err := client.Pipes.Alter(ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		alteredPipe, err := client.Pipes.ShowByID(ctx, pipe.ID())
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredPipe.Comment)

		alterOptions = &AlterPipeOptions{
			Unset: &PipeUnset{
				Comment:             Bool(true),
				PipeExecutionPaused: Bool(true),
			},
		}

		err = client.Pipes.Alter(ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		alteredPipe, err = client.Pipes.ShowByID(ctx, pipe.ID())
		require.NoError(t, err)

		assert.Equal(t, "", alteredPipe.Comment)
	})

	t.Run("set and unset tag", func(t *testing.T) {
		tag, tagCleanup := createTag(t, client, database, schema)
		t.Cleanup(tagCleanup)

		pipeName := randomAlphanumericN(t, 20)
		pipe, pipeCleanup := createPipe(t, client, database, schema, pipeName, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		tagValue := "abc"
		alterOptions := &AlterPipeOptions{
			SetTags: &PipeSetTags{
				Tag: []TagAssociation{
					{
						Name:  tag.ID(),
						Value: tagValue,
					},
				},
			},
		}

		err := client.Pipes.Alter(ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		returnedTagValue, err := client.SystemFunctions.GetTag(ctx, tag.ID(), pipe.ID(), ObjectTypePipe)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		alterOptions = &AlterPipeOptions{
			UnsetTags: &PipeUnsetTags{
				Tag: []ObjectIdentifier{
					tag.ID(),
				},
			},
		}

		err = client.Pipes.Alter(ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		_, err = client.SystemFunctions.GetTag(ctx, tag.ID(), pipe.ID(), ObjectTypePipe)
		assert.Error(t, err)
	})

	t.Run("refresh with all", func(t *testing.T) {
		pipeName := randomAlphanumericN(t, 20)
		pipe, pipeCleanup := createPipe(t, client, database, schema, pipeName, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		alterOptions := &AlterPipeOptions{
			Refresh: &PipeRefresh{
				Prefix:        String("/d1"),
				ModifiedAfter: String("2018-07-30T13:56:46-07:00"),
			},
		}

		err := client.Pipes.Alter(ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)
	})
}
