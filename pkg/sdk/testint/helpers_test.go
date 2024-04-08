package testint

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"github.com/stretchr/testify/require"
)

const (
	nycWeatherDataURL = "s3://snowflake-workshop-lab/weather-nyc"
)

// there is no direct way to get the account identifier from Snowflake API, but you can get it if you know
// the account locator and by filtering the list of accounts in replication accounts by the account locator
func getAccountIdentifier(t *testing.T, client *sdk.Client) sdk.AccountIdentifier {
	t.Helper()
	ctx := context.Background()
	currentAccountLocator, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	require.NoError(t, err)
	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == currentAccountLocator {
			return sdk.NewAccountIdentifier(replicationAccount.OrganizationName, replicationAccount.AccountName)
		}
	}
	return sdk.AccountIdentifier{}
}

func useWarehouse(t *testing.T, client *sdk.Client, warehouseID sdk.AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	err := client.Sessions.UseWarehouse(ctx, warehouseID)
	require.NoError(t, err)
	return func() {
		err = client.Sessions.UseWarehouse(ctx, testWarehouse(t).ID())
		require.NoError(t, err)
	}
}

func createDatabase(t *testing.T, client *sdk.Client) (*sdk.Database, func()) {
	t.Helper()
	return createDatabaseWithOptions(t, client, sdk.RandomAccountObjectIdentifier(), &sdk.CreateDatabaseOptions{})
}

func createDatabaseWithOptions(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier, opts *sdk.CreateDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Databases.Create(ctx, id, opts)
	require.NoError(t, err)
	database, err := client.Databases.ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := client.Databases.Drop(ctx, id, nil)
		require.NoError(t, err)
		err = client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(TestDatabaseName, TestSchemaName))
		require.NoError(t, err)
	}
}

func createSecondaryDatabase(t *testing.T, client *sdk.Client, externalId sdk.ExternalObjectIdentifier) (*sdk.Database, func()) {
	t.Helper()
	return createSecondaryDatabaseWithOptions(t, client, sdk.RandomAccountObjectIdentifier(), externalId, &sdk.CreateSecondaryDatabaseOptions{})
}

func createSecondaryDatabaseWithOptions(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier, externalId sdk.ExternalObjectIdentifier, opts *sdk.CreateSecondaryDatabaseOptions) (*sdk.Database, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Databases.CreateSecondary(ctx, id, externalId, opts)
	require.NoError(t, err)
	database, err := client.Databases.ShowByID(ctx, id)
	require.NoError(t, err)
	return database, func() {
		err := client.Databases.Drop(ctx, id, nil)
		require.NoError(t, err)

		// TODO [926148]: make this wait better with tests stabilization
		// waiting because sometimes dropping primary db right after dropping the secondary resulted in error
		time.Sleep(1 * time.Second)
		err = client.Sessions.UseSchema(ctx, sdk.NewDatabaseObjectIdentifier(TestDatabaseName, TestSchemaName))
		require.NoError(t, err)
	}
}

func createSchema(t *testing.T, client *sdk.Client, database *sdk.Database) (*sdk.Schema, func()) {
	t.Helper()
	return createSchemaWithIdentifier(t, client, database, random.StringRange(8, 28))
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
		err = client.Sessions.UseSchema(ctx, testSchema(t).ID())
		require.NoError(t, err)
	}
}

func createWarehouse(t *testing.T, client *sdk.Client) (*sdk.Warehouse, func()) {
	t.Helper()
	return createWarehouseWithOptions(t, client, &sdk.CreateWarehouseOptions{})
}

func createWarehouseWithOptions(t *testing.T, client *sdk.Client, opts *sdk.CreateWarehouseOptions) (*sdk.Warehouse, func()) {
	t.Helper()
	name := random.StringRange(8, 28)
	id := sdk.NewAccountObjectIdentifier(name)
	ctx := context.Background()
	err := client.Warehouses.Create(ctx, id, opts)
	require.NoError(t, err)
	return &sdk.Warehouse{
			Name: name,
		}, func() {
			err := client.Warehouses.Drop(ctx, id, nil)
			require.NoError(t, err)
			err = client.Sessions.UseWarehouse(ctx, testWarehouse(t).ID())
			require.NoError(t, err)
		}
}

func createUser(t *testing.T, client *sdk.Client) (*sdk.User, func()) {
	t.Helper()
	name := random.StringRange(8, 28)
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
	columns := []sdk.TableColumnRequest{
		*sdk.NewTableColumnRequest("id", sdk.DataTypeNumber),
	}
	return createTableWithColumns(t, client, database, schema, columns)
}

func createTableWithColumns(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, columns []sdk.TableColumnRequest) (*sdk.Table, func()) {
	t.Helper()
	name := random.StringRange(8, 28)
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	ctx := context.Background()

	dbCreateRequest := sdk.NewCreateTableRequest(id, columns)
	err := client.Tables.Create(ctx, dbCreateRequest)
	require.NoError(t, err)

	table, err := client.Tables.ShowByID(ctx, id)
	require.NoError(t, err)

	return table, func() {
		dropErr := client.Tables.Drop(ctx, sdk.NewDropTableRequest(id))
		require.NoError(t, dropErr)
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
	name := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, random.String())
	targetLag := sdk.TargetLag{
		MaximumDuration: sdk.String("2 minutes"),
	}
	query := "select id from " + table.ID().FullyQualifiedName()
	comment := random.Comment()
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
	name := random.StringRange(8, 28)
	ctx := context.Background()
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	err := client.Tags.Create(ctx, sdk.NewCreateTagRequest(id))
	require.NoError(t, err)
	tag, err := client.Tags.ShowByID(ctx, id)
	require.NoError(t, err)
	return tag, func() {
		err := client.Tags.Drop(ctx, sdk.NewDropTagRequest(id))
		require.NoError(t, err)
	}
}

func createStageWithDirectory(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, name string) (*sdk.Stage, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, name)
	return createStageWithOptions(t, client, id, func(request *sdk.CreateInternalStageRequest) *sdk.CreateInternalStageRequest {
		return request.WithDirectoryTableOptions(sdk.NewInternalDirectoryTableOptionsRequest().WithEnable(sdk.Bool(true)))
	})
}

func createStage(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier) (*sdk.Stage, func()) {
	t.Helper()
	return createStageWithOptions(t, client, id, func(request *sdk.CreateInternalStageRequest) *sdk.CreateInternalStageRequest { return request })
}

func createStageWithURL(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier, url string) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Stages.CreateOnS3(ctx, sdk.NewCreateOnS3StageRequest(id).
		WithExternalStageParams(sdk.NewExternalS3StageParamsRequest(url)))
	require.NoError(t, err)

	stage, err := client.Stages.ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, func() {
		err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
		require.NoError(t, err)
	}
}

func createStageWithOptions(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier, reqMapping func(*sdk.CreateInternalStageRequest) *sdk.CreateInternalStageRequest) (*sdk.Stage, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Stages.CreateInternal(ctx, reqMapping(sdk.NewCreateInternalStageRequest(id)))
	require.NoError(t, err)

	stage, err := client.Stages.ShowByID(ctx, id)
	require.NoError(t, err)

	return stage, func() {
		err := client.Stages.Drop(ctx, sdk.NewDropStageRequest(id))
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
	name := random.UUID()
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
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, random.StringN(12))
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
	id := sdk.RandomAccountObjectIdentifier()
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
			Name: random.String(),
			Type: sdk.DataTypeVARCHAR,
		},
	}
	n := random.IntRange(0, 5)
	for i := 0; i < n; i++ {
		signature = append(signature, sdk.TableColumnSignature{
			Name: random.String(),
			Type: sdk.DataTypeVARCHAR,
		})
	}
	expression := "REPLACE('X', 1, 2)"
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, sdk.DataTypeVARCHAR, expression, &sdk.CreateMaskingPolicyOptions{})
}

func createMaskingPolicyIdentity(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema, columnType sdk.DataType) (*sdk.MaskingPolicy, func()) {
	t.Helper()
	name := "a"
	signature := []sdk.TableColumnSignature{
		{
			Name: name,
			Type: columnType,
		},
	}
	expression := "a"
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, columnType, expression, &sdk.CreateMaskingPolicyOptions{})
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
	name := random.String()
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

	name := random.String()
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

func useRole(t *testing.T, client *sdk.Client, roleName string) func() {
	t.Helper()
	ctx := context.Background()

	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, sdk.NewAccountObjectIdentifier(roleName))
	require.NoError(t, err)

	return func() {
		err = client.Sessions.UseRole(ctx, sdk.NewAccountObjectIdentifier(currentRole))
		require.NoError(t, err)
	}
}

func createRole(t *testing.T, client *sdk.Client) (*sdk.Role, func()) {
	t.Helper()
	return createRoleWithRequest(t, client, sdk.NewCreateRoleRequest(sdk.RandomAccountObjectIdentifier()))
}

func createRoleGrantedToCurrentUser(t *testing.T, client *sdk.Client) (*sdk.Role, func()) {
	t.Helper()
	role, roleCleanup := createRoleWithRequest(t, client, sdk.NewCreateRoleRequest(sdk.RandomAccountObjectIdentifier()))

	ctx := context.Background()
	currentUser, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{
		User: sdk.Pointer(sdk.NewAccountObjectIdentifier(currentUser)),
	}))
	require.NoError(t, err)

	return role, roleCleanup
}

func createRoleWithRequest(t *testing.T, client *sdk.Client, req *sdk.CreateRoleRequest) (*sdk.Role, func()) {
	t.Helper()
	require.True(t, sdk.ValidObjectIdentifier(req.GetName()))
	ctx := context.Background()
	err := client.Roles.Create(ctx, req)
	require.NoError(t, err)
	role, err := client.Roles.ShowByID(ctx, sdk.NewShowByIdRoleRequest(req.GetName()))
	require.NoError(t, err)
	return role, func() {
		err = client.Roles.Drop(ctx, sdk.NewDropRoleRequest(req.GetName()))
		require.NoError(t, err)
	}
}

func createDatabaseRole(t *testing.T, client *sdk.Client, database *sdk.Database) (*sdk.DatabaseRole, func()) {
	t.Helper()
	name := random.String()
	id := sdk.NewDatabaseObjectIdentifier(database.Name, name)
	ctx := context.Background()

	err := client.DatabaseRoles.Create(ctx, sdk.NewCreateDatabaseRoleRequest(id))
	require.NoError(t, err)

	databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
	require.NoError(t, err)

	return databaseRole, cleanupDatabaseRoleProvider(t, ctx, client, id)
}

func cleanupDatabaseRoleProvider(t *testing.T, ctx context.Context, client *sdk.Client, id sdk.DatabaseObjectIdentifier) func() {
	t.Helper()
	return func() {
		err := client.DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id))
		require.NoError(t, err)
	}
}

func createFailoverGroup(t *testing.T, client *sdk.Client) (*sdk.FailoverGroup, func()) {
	t.Helper()
	objectTypes := []sdk.PluralObjectType{sdk.PluralObjectTypeRoles}
	ctx := context.Background()
	currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	accountID := sdk.NewAccountIdentifierFromAccountLocator(currentAccount)
	allowedAccounts := []sdk.AccountIdentifier{accountID}
	return createFailoverGroupWithOptions(t, client, objectTypes, allowedAccounts, nil)
}

func createFailoverGroupWithOptions(t *testing.T, client *sdk.Client, objectTypes []sdk.PluralObjectType, allowedAccounts []sdk.AccountIdentifier, opts *sdk.CreateFailoverGroupOptions) (*sdk.FailoverGroup, func()) {
	t.Helper()
	id := sdk.RandomAlphanumericAccountObjectIdentifier()
	ctx := context.Background()
	err := client.FailoverGroups.Create(ctx, id, objectTypes, allowedAccounts, opts)
	require.NoError(t, err)
	failoverGroup, err := client.FailoverGroups.ShowByID(ctx, id)
	require.NoError(t, err)
	return failoverGroup, func() {
		err := client.FailoverGroups.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

func createShare(t *testing.T, client *sdk.Client) (*sdk.Share, func()) {
	t.Helper()
	// TODO(SNOW-1058419): Try with identifier containing dot during identifiers rework
	id := sdk.RandomAlphanumericAccountObjectIdentifier()
	return createShareWithOptions(t, client, id, &sdk.CreateShareOptions{})
}

func createShareWithOptions(t *testing.T, client *sdk.Client, id sdk.AccountObjectIdentifier, opts *sdk.CreateShareOptions) (*sdk.Share, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.Shares.Create(ctx, id, opts)
	require.NoError(t, err)
	share, err := client.Shares.ShowByID(ctx, id)
	require.NoError(t, err)
	return share, func() {
		err := client.Shares.Drop(ctx, id)
		require.NoError(t, err)
	}
}

func createFileFormat(t *testing.T, client *sdk.Client, schema sdk.DatabaseObjectIdentifier) (*sdk.FileFormat, func()) {
	t.Helper()
	return createFileFormatWithOptions(t, client, schema, &sdk.CreateFileFormatOptions{
		Type: sdk.FileFormatTypeCSV,
	})
}

func createFileFormatWithOptions(t *testing.T, client *sdk.Client, schema sdk.DatabaseObjectIdentifier, opts *sdk.CreateFileFormatOptions) (*sdk.FileFormat, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName(), schema.Name(), random.String())
	ctx := context.Background()
	err := client.FileFormats.Create(ctx, id, opts)
	require.NoError(t, err)
	fileFormat, err := client.FileFormats.ShowByID(ctx, id)
	require.NoError(t, err)
	return fileFormat, func() {
		err := client.FileFormats.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

func createView(t *testing.T, client *sdk.Client, viewId sdk.SchemaObjectIdentifier, asQuery string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE VIEW %s AS %s`, viewId.FullyQualifiedName(), asQuery))
	require.NoError(t, err)

	return func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP VIEW %s`, viewId.FullyQualifiedName()))
		require.NoError(t, err)
	}
}

func putOnStage(t *testing.T, client *sdk.Client, stage *sdk.Stage, filename string) {
	t.Helper()
	ctx := context.Background()

	path, err := filepath.Abs("./testdata/" + filename)
	require.NoError(t, err)
	absPath := "file://" + path

	_, err = client.ExecForTests(ctx, fmt.Sprintf(`PUT '%s' @%s AUTO_COMPRESS = FALSE`, absPath, stage.ID().FullyQualifiedName()))
	require.NoError(t, err)
}

func putOnStageWithContent(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier, filename string, content string) {
	t.Helper()
	ctx := context.Background()

	tf := fmt.Sprintf("/tmp/%s", filename)
	f, err := os.Create(tf)
	require.NoError(t, err)
	if content != "" {
		_, err = f.Write([]byte(content))
		require.NoError(t, err)
	}
	f.Close()
	defer os.Remove(f.Name())

	_, err = client.ExecForTests(ctx, fmt.Sprintf(`PUT file://%s @%s AUTO_COMPRESS = FALSE OVERWRITE = TRUE`, f.Name(), id.FullyQualifiedName()))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err = client.ExecForTests(ctx, fmt.Sprintf(`REMOVE @%s/%s`, id.FullyQualifiedName(), filename))
		require.NoError(t, err)
	})
}

func createApplicationPackage(t *testing.T, client *sdk.Client, name string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE APPLICATION PACKAGE "%s"`, name))
	require.NoError(t, err)
	return func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP APPLICATION PACKAGE "%s"`, name))
		require.NoError(t, err)
	}
}

func addApplicationPackageVersion(t *testing.T, client *sdk.Client, stage *sdk.Stage, appPackageName string, versionName string) {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`ALTER APPLICATION PACKAGE "%s" ADD VERSION %v USING '@%s'`, appPackageName, versionName, stage.ID().FullyQualifiedName()))
	require.NoError(t, err)
}

func createApplication(t *testing.T, client *sdk.Client, name string, packageName string, version string) func() {
	t.Helper()
	ctx := context.Background()
	_, err := client.ExecForTests(ctx, fmt.Sprintf(`CREATE APPLICATION "%s" FROM APPLICATION PACKAGE "%s" USING VERSION %s`, name, packageName, version))
	require.NoError(t, err)
	return func() {
		_, err := client.ExecForTests(ctx, fmt.Sprintf(`DROP APPLICATION "%s"`, name))
		require.NoError(t, err)
	}
}

func createRowAccessPolicy(t *testing.T, client *sdk.Client, schema *sdk.Schema) (sdk.SchemaObjectIdentifier, func()) {
	t.Helper()
	ctx := context.Background()
	id := sdk.NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, random.String())

	arg := sdk.NewCreateRowAccessPolicyArgsRequest("A", sdk.DataTypeNumber)
	body := "true"
	createRequest := sdk.NewCreateRowAccessPolicyRequest(id, []sdk.CreateRowAccessPolicyArgsRequest{*arg}, body)
	err := client.RowAccessPolicies.Create(ctx, createRequest)
	require.NoError(t, err)

	return id, func() {
		err := client.RowAccessPolicies.Drop(ctx, sdk.NewDropRowAccessPolicyRequest(id))
		require.NoError(t, err)
	}
}

// TODO: extract getting row access policies as resource (like getting tag in system functions)
// getRowAccessPolicyFor is based on https://docs.snowflake.com/en/user-guide/security-row-intro#obtain-database-objects-with-a-row-access-policy.
func getRowAccessPolicyFor(t *testing.T, client *sdk.Client, id sdk.SchemaObjectIdentifier, objectType sdk.ObjectType) (*policyReference, error) {
	t.Helper()
	ctx := context.Background()

	s := &policyReference{}
	policyReferencesId := sdk.NewSchemaObjectIdentifier(id.DatabaseName(), "INFORMATION_SCHEMA", "POLICY_REFERENCES")
	err := client.QueryOneForTests(ctx, s, fmt.Sprintf(`SELECT * FROM TABLE(%s(REF_ENTITY_NAME => '%s', REF_ENTITY_DOMAIN => '%v'))`, policyReferencesId.FullyQualifiedName(), id.FullyQualifiedName(), objectType))

	return s, err
}

type policyReference struct {
	PolicyDb          string         `db:"POLICY_DB"`
	PolicySchema      string         `db:"POLICY_SCHEMA"`
	PolicyName        string         `db:"POLICY_NAME"`
	PolicyKind        string         `db:"POLICY_KIND"`
	RefDatabaseName   string         `db:"REF_DATABASE_NAME"`
	RefSchemaName     string         `db:"REF_SCHEMA_NAME"`
	RefEntityName     string         `db:"REF_ENTITY_NAME"`
	RefEntityDomain   string         `db:"REF_ENTITY_DOMAIN"`
	RefColumnName     sql.NullString `db:"REF_COLUMN_NAME"`
	RefArgColumnNames string         `db:"REF_ARG_COLUMN_NAMES"`
	TagDatabase       sql.NullString `db:"TAG_DATABASE"`
	TagSchema         sql.NullString `db:"TAG_SCHEMA"`
	TagName           sql.NullString `db:"TAG_NAME"`
	PolicyStatus      string         `db:"POLICY_STATUS"`
}

// TODO: extract getting table columns as resource (like getting tag in system functions)
// getTableColumnsFor is based on https://docs.snowflake.com/en/sql-reference/info-schema/columns.
func getTableColumnsFor(t *testing.T, client *sdk.Client, tableId sdk.SchemaObjectIdentifier) []informationSchemaColumns {
	t.Helper()
	ctx := context.Background()

	var columns []informationSchemaColumns
	query := fmt.Sprintf("SELECT * FROM information_schema.columns WHERE table_schema = '%s'  AND table_name = '%s' ORDER BY ordinal_position", tableId.SchemaName(), tableId.Name())
	err := client.QueryForTests(ctx, &columns, query)
	require.NoError(t, err)

	return columns
}

type informationSchemaColumns struct {
	TableCatalog           string         `db:"TABLE_CATALOG"`
	TableSchema            string         `db:"TABLE_SCHEMA"`
	TableName              string         `db:"TABLE_NAME"`
	ColumnName             string         `db:"COLUMN_NAME"`
	OrdinalPosition        string         `db:"ORDINAL_POSITION"`
	ColumnDefault          sql.NullString `db:"COLUMN_DEFAULT"`
	IsNullable             string         `db:"IS_NULLABLE"`
	DataType               string         `db:"DATA_TYPE"`
	CharacterMaximumLength sql.NullString `db:"CHARACTER_MAXIMUM_LENGTH"`
	CharacterOctetLength   sql.NullString `db:"CHARACTER_OCTET_LENGTH"`
	NumericPrecision       sql.NullString `db:"NUMERIC_PRECISION"`
	NumericPrecisionRadix  sql.NullString `db:"NUMERIC_PRECISION_RADIX"`
	NumericScale           sql.NullString `db:"NUMERIC_SCALE"`
	DatetimePrecision      sql.NullString `db:"DATETIME_PRECISION"`
	IntervalType           sql.NullString `db:"INTERVAL_TYPE"`
	IntervalPrecision      sql.NullString `db:"INTERVAL_PRECISION"`
	CharacterSetCatalog    sql.NullString `db:"CHARACTER_SET_CATALOG"`
	CharacterSetSchema     sql.NullString `db:"CHARACTER_SET_SCHEMA"`
	CharacterSetName       sql.NullString `db:"CHARACTER_SET_NAME"`
	CollationCatalog       sql.NullString `db:"COLLATION_CATALOG"`
	CollationSchema        sql.NullString `db:"COLLATION_SCHEMA"`
	CollationName          sql.NullString `db:"COLLATION_NAME"`
	DomainCatalog          sql.NullString `db:"DOMAIN_CATALOG"`
	DomainSchema           sql.NullString `db:"DOMAIN_SCHEMA"`
	DomainName             sql.NullString `db:"DOMAIN_NAME"`
	UdtCatalog             sql.NullString `db:"UDT_CATALOG"`
	UdtSchema              sql.NullString `db:"UDT_SCHEMA"`
	UdtName                sql.NullString `db:"UDT_NAME"`
	ScopeCatalog           sql.NullString `db:"SCOPE_CATALOG"`
	ScopeSchema            sql.NullString `db:"SCOPE_SCHEMA"`
	ScopeName              sql.NullString `db:"SCOPE_NAME"`
	MaximumCardinality     sql.NullString `db:"MAXIMUM_CARDINALITY"`
	DtdIdentifier          sql.NullString `db:"DTD_IDENTIFIER"`
	IsSelfReferencing      string         `db:"IS_SELF_REFERENCING"`
	IsIdentity             string         `db:"IS_IDENTITY"`
	IdentityGeneration     sql.NullString `db:"IDENTITY_GENERATION"`
	IdentityStart          sql.NullString `db:"IDENTITY_START"`
	IdentityIncrement      sql.NullString `db:"IDENTITY_INCREMENT"`
	IdentityMaximum        sql.NullString `db:"IDENTITY_MAXIMUM"`
	IdentityMinimum        sql.NullString `db:"IDENTITY_MINIMUM"`
	IdentityCycle          sql.NullString `db:"IDENTITY_CYCLE"`
	IdentityOrdered        sql.NullString `db:"IDENTITY_ORDERED"`
	Comment                sql.NullString `db:"COMMENT"`
}

func updateAccountParameterTemporarily(t *testing.T, client *sdk.Client, parameter sdk.AccountParameter, newValue string) func() {
	t.Helper()
	ctx := context.Background()

	param, err := client.Parameters.ShowAccountParameter(ctx, parameter)
	require.NoError(t, err)
	oldValue := param.Value

	err = client.Parameters.SetAccountParameter(ctx, parameter, newValue)
	require.NoError(t, err)

	return func() {
		err = client.Parameters.SetAccountParameter(ctx, parameter, oldValue)
		require.NoError(t, err)
	}
}

func createTaskWithRequest(t *testing.T, client *sdk.Client, request *sdk.CreateTaskRequest) (*sdk.Task, func()) {
	t.Helper()
	ctx := context.Background()

	id := request.GetName()

	err := client.Tasks.Create(ctx, request)
	require.NoError(t, err)

	task, err := client.Tasks.ShowByID(ctx, id)
	require.NoError(t, err)

	return task, func() {
		err = client.Tasks.Drop(ctx, sdk.NewDropTaskRequest(id))
		require.NoError(t, err)
	}
}

func createTask(t *testing.T, client *sdk.Client, database *sdk.Database, schema *sdk.Schema) (*sdk.Task, func()) {
	t.Helper()
	id := sdk.NewSchemaObjectIdentifier(database.Name, schema.Name, random.AlphaN(20))
	warehouseReq := sdk.NewCreateTaskWarehouseRequest().WithWarehouse(sdk.Pointer(testWarehouse(t).ID()))
	return createTaskWithRequest(t, client, sdk.NewCreateTaskRequest(id, "SELECT CURRENT_TIMESTAMP").WithSchedule(sdk.String("60 minutes")).WithWarehouse(warehouseReq))
}
