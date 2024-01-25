package testint

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPipeCopyStatement(t *testing.T, table *sdk.Table, stage *sdk.Stage) string {
	t.Helper()
	require.NotNil(t, table, "table has to be created")
	require.NotNil(t, stage, "stage has to be created")
	return fmt.Sprintf("COPY INTO %s\nFROM @%s", table.ID().FullyQualifiedName(), stage.ID().FullyQualifiedName())
}

// TestInt_CreatePipeWithStrangeSchemaName documented previous bad behavior. It changed with Snowflake 8.3.1 release.
// We leave the test for future reference.
func TestInt_CreatePipeWithStrangeSchemaName(t *testing.T) {
	schemaIdentifier := sdk.NewDatabaseObjectIdentifier(testDb(t).Name, "tcK1>AJ+")

	// creating a new schema on purpose
	schema, schemaCleanup := createSchemaWithIdentifier(t, itc.client, testDb(t), schemaIdentifier.Name())
	t.Cleanup(schemaCleanup)

	table, tableCleanup := createTable(t, itc.client, testDb(t), schema)
	t.Cleanup(tableCleanup)

	stageName := random.AlphanumericN(20)
	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, stageName))
	t.Cleanup(stageCleanup)

	t.Run("if we have special characters in db or schema name, create pipe succeeds", func(t *testing.T) {
		err := itc.client.Pipes.Create(
			itc.ctx,
			sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphanumericN(20)),
			createPipeCopyStatement(t, table, stage),
			&sdk.CreatePipeOptions{},
		)

		require.NoError(t, err)
	})

	t.Run("the same does not work when using non fully qualified name for table", func(t *testing.T) {
		createCopyStatementWithoutQualifiersForStage := func(t *testing.T, table *sdk.Table, stage *sdk.Stage) string {
			t.Helper()
			require.NotNil(t, table, "table has to be created")
			require.NotNil(t, stage, "stage has to be created")
			return fmt.Sprintf("COPY INTO %s\nFROM @\"%s\"", table.ID().FullyQualifiedName(), stage.Name)
		}

		err := itc.client.Pipes.Create(
			itc.ctx,
			sdk.NewSchemaObjectIdentifier(testDb(t).Name, schema.Name, random.AlphanumericN(20)),
			createCopyStatementWithoutQualifiersForStage(t, table, stage),
			&sdk.CreatePipeOptions{},
		)

		require.Error(t, err)
		require.ErrorContains(t, err, "object does not exist or not authorized")
	})
}

func TestInt_PipesShowAndDescribe(t *testing.T) {
	table1, table1Cleanup := createTable(t, itc.client, testDb(t), testSchema(t))
	t.Cleanup(table1Cleanup)

	table2, table2Cleanup := createTable(t, itc.client, testDb(t), testSchema(t))
	t.Cleanup(table2Cleanup)

	stageName := random.AlphanumericN(20)
	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, stageName))
	t.Cleanup(stageCleanup)

	pipe1Name := random.AlphanumericN(20)
	pipe1CopyStatement := createPipeCopyStatement(t, table1, stage)
	pipe1, pipe1Cleanup := createPipe(t, itc.client, testDb(t), testSchema(t), pipe1Name, pipe1CopyStatement)
	t.Cleanup(pipe1Cleanup)

	pipe2Name := random.AlphanumericN(20)
	pipe2CopyStatement := createPipeCopyStatement(t, table2, stage)
	pipe2, pipe2Cleanup := createPipe(t, itc.client, testDb(t), testSchema(t), pipe2Name, pipe2CopyStatement)
	t.Cleanup(pipe2Cleanup)

	t.Run("show: without options", func(t *testing.T) {
		pipes, err := itc.client.Pipes.Show(itc.ctx, &sdk.ShowPipeOptions{})

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, *pipe1)
		assert.Contains(t, pipes, *pipe2)
	})

	t.Run("show: in schema", func(t *testing.T) {
		showOptions := &sdk.ShowPipeOptions{
			In: &sdk.In{
				Schema: testSchema(t).ID(),
			},
		}
		pipes, err := itc.client.Pipes.Show(itc.ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 2, len(pipes))
		assert.Contains(t, pipes, *pipe1)
		assert.Contains(t, pipes, *pipe2)
	})

	t.Run("show: like", func(t *testing.T) {
		showOptions := &sdk.ShowPipeOptions{
			Like: &sdk.Like{
				Pattern: sdk.String(pipe1Name),
			},
		}
		pipes, err := itc.client.Pipes.Show(itc.ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 1, len(pipes))
		assert.Contains(t, pipes, *pipe1)
	})

	t.Run("show: non-existent pipe", func(t *testing.T) {
		showOptions := &sdk.ShowPipeOptions{
			Like: &sdk.Like{
				Pattern: sdk.String("non-existent"),
			},
		}
		pipes, err := itc.client.Pipes.Show(itc.ctx, showOptions)

		require.NoError(t, err)
		assert.Equal(t, 0, len(pipes))
	})

	t.Run("describe: existing pipe", func(t *testing.T) {
		pipe, err := itc.client.Pipes.Describe(itc.ctx, pipe1.ID())

		require.NoError(t, err)
		assert.Equal(t, pipe1.Name, pipe.Name)
	})

	t.Run("describe: non-existing pipe", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testDb(t).Name, "does_not_exist")

		_, err := itc.client.Pipes.Describe(itc.ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_PipeCreate(t *testing.T) {
	table, tableCleanup := createTable(t, itc.client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	stageName := random.AlphanumericN(20)
	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, stageName))
	t.Cleanup(stageCleanup)

	copyStatement := createPipeCopyStatement(t, table, stage)

	assertPipe := func(t *testing.T, pipeDetails *sdk.Pipe, expectedName string, expectedComment string) {
		t.Helper()
		assert.NotEmpty(t, pipeDetails.CreatedOn)
		assert.Equal(t, expectedName, pipeDetails.Name)
		assert.Equal(t, testDb(t).Name, pipeDetails.DatabaseName)
		assert.Equal(t, testSchema(t).Name, pipeDetails.SchemaName)
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
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)
		comment := random.Comment()

		err := itc.client.Pipes.Create(itc.ctx, id, copyStatement, &sdk.CreatePipeOptions{
			OrReplace:   sdk.Bool(false),
			IfNotExists: sdk.Bool(true),
			AutoIngest:  sdk.Bool(false),
			Comment:     sdk.String(comment),
		})
		require.NoError(t, err)

		pipe, err := itc.client.Pipes.Describe(itc.ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, name, comment)
	})

	t.Run("test if not exists and or replace are incompatible", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		err := itc.client.Pipes.Create(itc.ctx, id, copyStatement, &sdk.CreatePipeOptions{
			OrReplace:   sdk.Bool(true),
			IfNotExists: sdk.Bool(true),
		})
		require.ErrorContains(t, err, "(0A000): SQL compilation error:\noptions IF NOT EXISTS and OR REPLACE are incompatible")
	})

	t.Run("test no options", func(t *testing.T) {
		name := random.String()
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, name)

		err := itc.client.Pipes.Create(itc.ctx, id, copyStatement, nil)
		require.NoError(t, err)

		pipe, err := itc.client.Pipes.Describe(itc.ctx, id)

		require.NoError(t, err)
		assertPipe(t, pipe, name, "")
	})
}

func TestInt_PipeDrop(t *testing.T) {
	table, tableCleanup := createTable(t, itc.client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	stageName := random.AlphanumericN(20)
	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, stageName))
	t.Cleanup(stageCleanup)

	t.Run("pipe exists", func(t *testing.T) {
		pipeName := random.AlphanumericN(20)
		pipeCopyStatement := createPipeCopyStatement(t, table, stage)
		pipe, _ := createPipe(t, itc.client, testDb(t), testSchema(t), pipeName, pipeCopyStatement)

		err := itc.client.Pipes.Drop(itc.ctx, pipe.ID())

		require.NoError(t, err)
		_, err = itc.client.Pipes.Describe(itc.ctx, pipe.ID())
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})

	t.Run("pipe does not exist", func(t *testing.T) {
		id := sdk.NewSchemaObjectIdentifier(testDb(t).Name, testDb(t).Name, "does_not_exist")

		err := itc.client.Alerts.Drop(itc.ctx, id)
		assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	})
}

func TestInt_PipeAlter(t *testing.T) {
	table, tableCleanup := createTable(t, itc.client, testDb(t), testSchema(t))
	t.Cleanup(tableCleanup)

	stageName := random.AlphanumericN(20)
	stage, stageCleanup := createStage(t, itc.client, sdk.NewSchemaObjectIdentifier(testDb(t).Name, testSchema(t).Name, stageName))
	t.Cleanup(stageCleanup)

	pipeCopyStatement := createPipeCopyStatement(t, table, stage)

	// TODO: test error integration when we have them in project
	t.Run("set value and unset value", func(t *testing.T) {
		pipeName := random.AlphanumericN(20)
		pipe, pipeCleanup := createPipe(t, itc.client, testDb(t), testSchema(t), pipeName, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		alterOptions := &sdk.AlterPipeOptions{
			Set: &sdk.PipeSet{
				Comment:             sdk.String("new comment"),
				PipeExecutionPaused: sdk.Bool(true),
			},
		}

		err := itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		alteredPipe, err := itc.client.Pipes.ShowByID(itc.ctx, pipe.ID())
		require.NoError(t, err)

		assert.Equal(t, "new comment", alteredPipe.Comment)

		alterOptions = &sdk.AlterPipeOptions{
			Unset: &sdk.PipeUnset{
				Comment:             sdk.Bool(true),
				PipeExecutionPaused: sdk.Bool(true),
			},
		}

		err = itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		alteredPipe, err = itc.client.Pipes.ShowByID(itc.ctx, pipe.ID())
		require.NoError(t, err)

		assert.Equal(t, "", alteredPipe.Comment)
	})

	t.Run("set and unset tag", func(t *testing.T) {
		tag, tagCleanup := createTag(t, itc.client, testDb(t), testSchema(t))
		t.Cleanup(tagCleanup)

		pipeName := random.AlphanumericN(20)
		pipe, pipeCleanup := createPipe(t, itc.client, testDb(t), testSchema(t), pipeName, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		tagValue := "abc"
		alterOptions := &sdk.AlterPipeOptions{
			SetTag: []sdk.TagAssociation{
				{
					Name:  tag.ID(),
					Value: tagValue,
				},
			},
		}

		err := itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		returnedTagValue, err := itc.client.SystemFunctions.GetTag(itc.ctx, tag.ID(), pipe.ID(), sdk.ObjectTypePipe)
		require.NoError(t, err)

		assert.Equal(t, tagValue, returnedTagValue)

		alterOptions = &sdk.AlterPipeOptions{
			UnsetTag: []sdk.ObjectIdentifier{
				tag.ID(),
			},
		}

		err = itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)

		_, err = itc.client.SystemFunctions.GetTag(itc.ctx, tag.ID(), pipe.ID(), sdk.ObjectTypePipe)
		assert.Error(t, err)
	})

	t.Run("refresh with all", func(t *testing.T) {
		pipeName := random.AlphanumericN(20)
		pipe, pipeCleanup := createPipe(t, itc.client, testDb(t), testSchema(t), pipeName, pipeCopyStatement)
		t.Cleanup(pipeCleanup)

		alterOptions := &sdk.AlterPipeOptions{
			Refresh: &sdk.PipeRefresh{
				Prefix:        sdk.String("/d1"),
				ModifiedAfter: sdk.String("2018-07-30T13:56:46-07:00"),
			},
		}

		err := itc.client.Pipes.Alter(itc.ctx, pipe.ID(), alterOptions)
		require.NoError(t, err)
	})
}
