package sdk_integration_tests

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/require"
)

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

func useSchema(t *testing.T, client *sdk.Client, schemaID sdk.DatabaseObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	orgSchema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseSchema(ctx, schemaID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(orgDB, orgSchema))
		require.NoError(t, err)
	}
}

func useWarehouse(t *testing.T, client *sdk.Client, warehouseID sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgWarehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseWarehouse(ctx, warehouseID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseWarehouse(ctx, sdk.NewAccountObjectIdentifier(orgWarehouse))
		require.NoError(t, err)
	}
}

func createDatabase(t *testing.T, client *sdk.Client) (*sdk.Database, func()) {
	t.Helper()
	return createDatabaseWithOptions(t, client, randomAccountObjectIdentifier(t), &sdk.CreateDatabaseOptions{})
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
	schemaID := sdk.NewDatabaseObjectIdentifier(database.Name, name)
	err := client.Schemas.Create(ctx, schemaID, nil)
	require.NoError(t, err)
	schema, err := client.Schemas.ShowByID(ctx, sdk.NewDatabaseObjectIdentifier(database.Name, name))
	require.NoError(t, err)
	return schema, func() {
		err := client.Schemas.Drop(ctx, schemaID, nil)
		if errors.Is(err, sdk.ErrObjectNotExistOrAuthorized) {
			return
		}
		require.NoError(t, err)
	}
}

func createWarehouse(t *testing.T, client *sdk.Client) (*sdk.Warehouse, func()) {
	t.Helper()
	return createWarehouseWithOptions(t, client, &sdk.CreateWarehouseOptions{})
}

func createWarehouseWithOptions(t *testing.T, client *sdk.Client, opts *sdk.CreateWarehouseOptions) (*sdk.Warehouse, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	id := sdk.NewAccountObjectIdentifier(name)
	ctx := context.Background()
	err := client.Warehouses.Create(ctx, id, opts)
	require.NoError(t, err)
	return &sdk.Warehouse{
			Name: name,
		}, func() {
			err := client.Warehouses.Drop(ctx, id, nil)
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
	return createTagWithOptions(t, client, database, schema, &sdk.CreateTagOptions{})
}

func createTagWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, _ *sdk.CreateTagOptions) (*sdk.Tag, func()) {
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

func createStageWithName(t *testing.T, client *sdk.Client, name string) (*string, func()) {
	t.Helper()
	ctx := context.Background()
	stageCleanup := func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf("DROP STAGE %s", name))
		require.NoError(t, err)
	}
	_, err := client.ExecForTests(ctx, fmt.Sprintf("CREATE STAGE %s", name))
	if err != nil {
		return nil, stageCleanup
	}
	require.NoError(t, err)
	return &name, stageCleanup
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

	err := client.Pipes.Create(ctx, id, copyStatement, &sdk.CreatePipeOptions{})
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

func createPasswordPolicy(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	return createPasswordPolicyWithOptions(t, client, database, schema, nil)
}

func createPasswordPolicyWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, options *sdk.CreatePasswordPolicyOptions) (*sdk.PasswordPolicy, func()) {
	t.Helper()
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = createDatabase(t, client)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = createSchema(t, client, database)
	}
	name := randomUUID(t)
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.PasswordPolicies.Create(ctx, id, options)
	require.NoError(t, err)

	showOptions := &sdk.ShowPasswordPolicyOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(name),
		},
		In: &sdk.In{
			Schema: schema.ID(),
		},
	}
	passwordPolicyList, err := client.PasswordPolicies.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(passwordPolicyList))
	return &passwordPolicyList[0], func() {
		err := client.PasswordPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
	}
}
