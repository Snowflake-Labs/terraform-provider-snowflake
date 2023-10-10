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

func createUser(t *testing.T, client *sdk.Client) (*sdk.User, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	id := sdk.NewAccountObjectIdentifier(name)
	return createUserWithOptions(t, client, id, &sdk.CreateUserOptions{})
}

func createUserWithName(t *testing.T, client *sdk.Client, name string) (*sdk.User, func()) {
	t.Helper()
	id := sdk.NewAccountObjectIdentifier(name)
	return createUserWithOptions(t, client, id, &sdk.CreateUserOptions{})
}

func createUserWithOptions(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier, opts *sdk.CreateUserOptions) (*sdk.User, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Users.Create(ctx, id, opts)
	require.NoError(t, err)
	user, err := client.Users.ShowByID(ctx, id)
	require.NoError(t, err)
	return user, func() {
		err := client.Users.Drop(ctx, id)
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

func createDynamicTable(t *testing.T, client *sdk.Client) (*sdk.DynamicTable, func()) {
	t.Helper()
	return createDynamicTableWithOptions(t, client, nil, nil, nil, nil)
}

func createDynamicTableWithOptions(t *testing.T, client *sdk.Client, warehouse *sdk.Warehouse, database *sdk.Database, schema *sdk.Schema, table *sdk.Table) (*sdk.DynamicTable, func()) {
	t.Helper()
	var warehouseCleanup func()
	if warehouse == nil {
		warehouse, warehouseCleanup = createWarehouse(t, client)
	}
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = createDatabase(t, client)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = createSchema(t, client, database)
	}
	var tableCleanup func()
	if table == nil {
		table, tableCleanup = createTable(t, client, database, schema)
	}
	name := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, randomString(t))
	targetLag := sdk.TargetLag{
		Lagtime: sdk.String("2 minutes"),
	}
	query := "select id from " + table.ID().FullyQualifiedName()
	comment := randomComment(t)
	ctx := context.Background()
	err := client.DynamicTables.Create(ctx, sdk.NewCreateDynamicTableRequest(name, warehouse.ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
	require.NoError(t, err)
	entities, err := client.DynamicTables.Show(ctx, sdk.NewShowDynamicTableRequest().WithLike(&sdk.Like{Pattern: sdk.String(name.Name())}).WithIn(&sdk.In{Schema: schema.ID()}))
	require.NoError(t, err)
	require.Equal(t, 1, len(entities))
	return &entities[0], func() {
		require.NoError(t, client.DynamicTables.Drop(ctx, sdk.NewDropDynamicTableRequest(name)))
		if tableCleanup != nil {
			tableCleanup()
		}
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
		if warehouseCleanup != nil {
			warehouseCleanup()
		}
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

func createStageWithURL(t *testing.T, client *sdk.Client, name sdk.AccountObjectIdentifier, url string) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE STAGE "%s" URL = '%s'`, name.Name(), url))
	require.NoError(t, err)

	return nil, func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP STAGE "%s"`, name.Name()))
		require.NoError(t, err)
	}
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

func createNetworkPolicy(t *testing.T, client *sdk.Client, req *sdk.CreateNetworkPolicyRequest) (error, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.NetworkPolicies.Create(ctx, req)
	return err, func() {
		err := client.NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(req.GetName()))
		require.NoError(t, err)
	}
}

func createSessionPolicy(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.SessionPolicy, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, randomStringN(t, 12))
	return createSessionPolicyWithOptions(t, client, id, sdk.NewCreateSessionPolicyRequest(id))
}

func createSessionPolicyWithOptions(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier, request *sdk.CreateSessionPolicyRequest) (*sdk.SessionPolicy, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.SessionPolicies.Create(ctx, request)
	require.NoError(t, err)
	sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
	require.NoError(t, err)
	return sessionPolicy, func() {
		err := client.SessionPolicies.Drop(ctx, sdk.NewDropSessionPolicyRequest(id))
		require.NoError(t, err)
	}
}

func createResourceMonitor(t *testing.T, client *sdk.Client) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	return createResourceMonitorWithOptions(t, client, &sdk.CreateResourceMonitorOptions{
		With: &sdk.ResourceMonitorWith{
			CreditQuota: sdk.Pointer(100),
			Triggers: []sdk.TriggerDefinition{
				{
					Threshold:     100,
					TriggerAction: sdk.TriggerActionSuspend,
				},
				{
					Threshold:     70,
					TriggerAction: sdk.TriggerActionSuspendImmediate,
				},
				{
					Threshold:     90,
					TriggerAction: sdk.TriggerActionNotify,
				},
			},
		},
	})
}

func createResourceMonitorWithOptions(t *testing.T, client *sdk.Client, opts *sdk.CreateResourceMonitorOptions) (*sdk.ResourceMonitor, func()) {
	t.Helper()
	id := randomAccountObjectIdentifier(t)
	ctx := context.Background()
	err := client.ResourceMonitors.Create(ctx, id, opts)
	require.NoError(t, err)
	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
	require.NoError(t, err)
	return resourceMonitor, func() {
		err := client.ResourceMonitors.Drop(ctx, id)
		require.NoError(t, err)
	}
}

func createMaskingPolicy(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	signature := []sdk.TableColumnSignature{
		{
			Name: randomString(t),
			Type: sdk.DataTypeVARCHAR,
		},
	}
	n := randomIntRange(t, 0, 5)
	for i := 0; i < n; i++ {
		signature = append(signature, sdk.TableColumnSignature{
			Name: randomString(t),
			Type: sdk.DataTypeVARCHAR,
		})
	}
	expression := "REPLACE('X', 1, 2)"
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, sdk.DataTypeVARCHAR, expression, &sdk.CreateMaskingPolicyOptions{})
}

func createMaskingPolicyWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, signature []sdk.TableColumnSignature, returns sdk.DataType, expression string, options *sdk.CreateMaskingPolicyOptions) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = createDatabase(t, client)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = createSchema(t, client, database)
	}
	name := randomString(t)
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.MaskingPolicies.Create(ctx, id, signature, returns, expression, options)
	require.NoError(t, err)

	showOptions := &sdk.ShowMaskingPolicyOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(name),
		},
		In: &sdk.In{
			Schema: schema.ID(),
		},
	}
	maskingPolicyList, err := client.MaskingPolicies.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(maskingPolicyList))
	return &maskingPolicyList[0], func() {
		err := client.MaskingPolicies.Drop(ctx, id)
		require.NoError(t, err)
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
	}
}

func createAlert(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, warehouse *sdk.Warehouse) (*sdk.Alert, func()) {
	t.Helper()
	schedule := "USING CRON * * * * * UTC"
	condition := "SELECT 1"
	action := "SELECT 1"
	return createAlertWithOptions(t, client, database, schema, warehouse, schedule, condition, action, &sdk.CreateAlertOptions{})
}

func createAlertWithOptions(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, warehouse *sdk.Warehouse, schedule string, condition string, action string, opts *sdk.CreateAlertOptions) (*sdk.Alert, func()) {
	t.Helper()
	var databaseCleanup func()
	if database == nil {
		database, databaseCleanup = createDatabase(t, client)
	}
	var schemaCleanup func()
	if schema == nil {
		schema, schemaCleanup = createSchema(t, client, database)
	}
	var warehouseCleanup func()
	if warehouse == nil {
		warehouse, warehouseCleanup = createWarehouse(t, client)
	}

	name := randomString(t)
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.Alerts.Create(ctx, id, warehouse.ID(), schedule, condition, action, opts)
	require.NoError(t, err)

	showOptions := &sdk.ShowAlertOptions{
		Like: &sdk.Like{
			Pattern: sdk.String(name),
		},
		In: &sdk.In{
			Schema: schema.ID(),
		},
	}
	alertList, err := client.Alerts.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(alertList))
	return &alertList[0], func() {
		err := client.Alerts.Drop(ctx, id)
		require.NoError(t, err)
		if schemaCleanup != nil {
			schemaCleanup()
		}
		if databaseCleanup != nil {
			databaseCleanup()
		}
		if warehouseCleanup != nil {
			warehouseCleanup()
		}
	}
}
