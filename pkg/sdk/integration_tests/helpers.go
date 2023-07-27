package sdk_integration_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T) *sdk.Client {
	t.Helper()

	client, err := sdk.NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

func createDatabaseWithIdentifier(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return createDatabaseWithOptions(t, client, id, &sdk.CreateDatabaseOptions{})
}

func createDatabaseWithOptions(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier, _ *sdk.CreateDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Databases.Create(ctx, id, nil)
	require.NoError(t, err)
	database, err := client.Databases.ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := client.Databases.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

func createSchema(t *testing.T, client *sdk.Client, database *sdk.Database) (*sdk.Schema, func()) {
	t.Helper()
	return createSchemaWithIdentifier(t, client, database, randomStringRange(t, 8, 28))
}

func createSchemaWithIdentifier(t *testing.T, client *sdk.Client, database *sdk.Database, name string) (*sdk.Schema, func()) {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf("CREATE SCHEMA \"%s\".\"%s\"", database.Name, name))
	require.NoError(t, err)
	return &sdk.Schema{
			DatabaseName: database.Name,
			Name:         name,
		}, func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP SCHEMA \"%s\".\"%s\"", database.Name, name))
			require.NoError(t, err)
		}
}

func useDatabase(t *testing.T, client *sdk.Client, databaseID sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseDatabase(ctx, databaseID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseDatabase(ctx, sdk.NewAccountObjectIdentifier(orgDB))
		require.NoError(t, err)
	}
}

func useSchema(t *testing.T, client *sdk.Client, schemaID sdk.SchemaIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	orgSchema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseSchema(ctx, schemaID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseSchema(ctx, sdk.NewSchemaIdentifier(orgDB, orgSchema))
		require.NoError(t, err)
	}
}

func createTable(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.Table, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf("CREATE TABLE \"%s\".\"%s\".\"%s\" (id NUMBER)", database.Name, schema.Name, name))
	require.NoError(t, err)
	return &sdk.Table{
			DatabaseName: database.Name,
			SchemaName:   schema.Name,
			Name:         name,
		}, func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP TABLE \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
			require.NoError(t, err)
		}
}

func createTag(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.Tag, func()) {
	t.Helper()
	return createTagWithOptions(t, client, database, schema, &sdk.TagCreateOptions{})
}

func createTagWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, _ *sdk.TagCreateOptions) (*sdk.Tag, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf("CREATE TAG \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
	require.NoError(t, err)
	return &sdk.Tag{
			Name:         name,
			DatabaseName: database.Name,
			SchemaName:   schema.Name,
		}, func() {
			_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP TAG \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
			require.NoError(t, err)
		}
}

func createStage(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, name string) (*sdk.Stage, func()) {
	t.Helper()
	require.NotNil(t, database, "database has to be created")
	require.NotNil(t, schema, "schema has to be created")

	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	ctx := context.Background()

	stageCleanup := func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP STAGE %s", id.FullyQualifiedName()))
		require.NoError(t, err)
	}

	_, err := client.ExecForTests(ctx, fmt.Sprintf("CREATE STAGE %s", id.FullyQualifiedName()))
	if err != nil {
		return nil, stageCleanup
	}
	require.NoError(t, err)

	return &sdk.Stage{
		DatabaseName: database.Name,
		SchemaName:   schema.Name,
		Name:         name,
	}, stageCleanup
}

func createPipe(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, name string, copyStatement string) (*sdk.Pipe, func()) {
	t.Helper()
	require.NotNil(t, database, "database has to be created")
	require.NotNil(t, schema, "schema has to be created")

	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	ctx := context.Background()

	pipeCleanup := func() {
		err := client.Pipes.Drop(ctx, id)
		require.NoError(t, err)
	}

	err := client.Pipes.Create(ctx, id, copyStatement, &sdk.PipeCreateOptions{})
	if err != nil {
		return nil, pipeCleanup
	}
	require.NoError(t, err)

	createdPipe, errDescribe := client.Pipes.Describe(ctx, id)
	if errDescribe != nil {
		return nil, pipeCleanup
	}
	require.NoError(t, errDescribe)

	return createdPipe, pipeCleanup
}
