package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// there is no direct way to get the account identifier from Snowflake API, but you can get it if you know
// the account locator and by filtering the list of accounts in replication accounts by the account locator
func getAccountIdentifier(t *testing.T, client *Client) AccountIdentifier {
	t.Helper()
	ctx := context.Background()
	currentAccountLocator, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAccounts(ctx)
	require.NoError(t, err)
	for _, replicationAccount := range replicationAccounts {
		if replicationAccount.AccountLocator == currentAccountLocator {
			return AccountIdentifier{
				organizationName: replicationAccount.OrganizationName,
				accountName:      replicationAccount.AccountName,
			}
		}
	}
	return AccountIdentifier{}
}

func getSecondaryAccountIdentifier(t *testing.T) AccountIdentifier {
	t.Helper()
	client := testSecondaryClient(t)
	return getAccountIdentifier(t, client)
}

//func randomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
//	t.Helper()
//	return NewSchemaObjectIdentifier(randomStringN(t, 12), randomStringN(t, 12), randomStringN(t, 12))
//}
//
//func randomDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
//	t.Helper()
//	return NewDatabaseObjectIdentifier(randomStringN(t, 12), randomStringN(t, 12))
//}
//
//func alphanumericDatabaseObjectIdentifier(t *testing.T) DatabaseObjectIdentifier {
//	t.Helper()
//	return NewDatabaseObjectIdentifier(randomAlphanumericN(t, 12), randomAlphanumericN(t, 12))
//}
//
//func randomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
//	t.Helper()
//	return NewAccountObjectIdentifier(randomStringN(t, 12))
//}

//func useDatabase(t *testing.T, client *Client, databaseID AccountObjectIdentifier) func() {
//	t.Helper()
//	ctx := context.Background()
//	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
//	require.NoError(t, err)
//	err = client.Sessions.UseDatabase(ctx, databaseID)
//	require.NoError(t, err)
//	return func() {
//		err := client.Sessions.UseDatabase(ctx, NewAccountObjectIdentifier(orgDB))
//		require.NoError(t, err)
//	}
//}
//
//func useSchema(t *testing.T, client *Client, schemaID DatabaseObjectIdentifier) func() {
//	t.Helper()
//	ctx := context.Background()
//	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
//	require.NoError(t, err)
//	orgSchema, err := client.ContextFunctions.CurrentSchema(ctx)
//	require.NoError(t, err)
//	err = client.Sessions.UseSchema(ctx, schemaID)
//	require.NoError(t, err)
//	return func() {
//		err := client.Sessions.UseSchema(ctx, NewDatabaseObjectIdentifier(orgDB, orgSchema))
//		require.NoError(t, err)
//	}
//}

//func useWarehouse(t *testing.T, client *Client, warehouseID AccountObjectIdentifier) func() {
//	t.Helper()
//	ctx := context.Background()
//	orgWarehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
//	require.NoError(t, err)
//	err = client.Sessions.UseWarehouse(ctx, warehouseID)
//	require.NoError(t, err)
//	return func() {
//		err := client.Sessions.UseWarehouse(ctx, NewAccountObjectIdentifier(orgWarehouse))
//		require.NoError(t, err)
//	}
//}

func testClient(t *testing.T) *Client {
	t.Helper()

	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

const (
	secondaryAccountProfile = "secondary_test_account"
)

func testSecondaryClient(t *testing.T) *Client {
	t.Helper()

	client, err := testClientFromProfile(t, secondaryAccountProfile)
	if err != nil {
		t.Skipf("Snowflake secondary account not configured. Must be set in ~./snowflake/config.yml with profile name: %s", secondaryAccountProfile)
	}

	return client
}

func testClientFromProfile(t *testing.T, profile string) (*Client, error) {
	t.Helper()
	config, err := ProfileConfig(profile)
	if err != nil {
		return nil, err
	}
	return NewClient(config)
}

//func randomUUID(t *testing.T) string {
//	t.Helper()
//	v, err := uuid.GenerateUUID()
//	require.NoError(t, err)
//	return v
//}
//
//func randomComment(t *testing.T) string {
//	t.Helper()
//	return gofakeit.Sentence(10)
//}
//
//func randomBool(t *testing.T) bool {
//	t.Helper()
//	return gofakeit.Bool()
//}
//
//func randomString(t *testing.T) string {
//	t.Helper()
//	return gofakeit.Password(true, true, true, true, false, 28)
//}
//
//func randomStringN(t *testing.T, num int) string {
//	t.Helper()
//	return gofakeit.Password(true, true, true, true, false, num)
//}
//
//func randomAlphanumericN(t *testing.T, num int) string {
//	t.Helper()
//	return gofakeit.Password(true, true, true, false, false, num)
//}
//
//func randomStringRange(t *testing.T, min, max int) string {
//	t.Helper()
//	if min > max {
//		t.Errorf("min %d is greater than max %d", min, max)
//	}
//	return gofakeit.Password(true, true, true, true, false, randomIntRange(t, min, max))
//}
//
//func randomIntRange(t *testing.T, min, max int) int {
//	t.Helper()
//	if min > max {
//		t.Errorf("min %d is greater than max %d", min, max)
//	}
//	return gofakeit.IntRange(min, max)
//}

//func createSessionPolicy(t *testing.T, client *Client, database *Database, schema *Schema) (*SessionPolicy, func()) {
//	t.Helper()
//	id := NewSchemaObjectIdentifier(database.Name, schema.Name, randomStringN(t, 12))
//	return createSessionPolicyWithOptions(t, client, id, NewCreateSessionPolicyRequest(id))
//}
//
//func createSessionPolicyWithOptions(t *testing.T, client *Client, id SchemaObjectIdentifier, request *CreateSessionPolicyRequest) (*SessionPolicy, func()) {
//	t.Helper()
//	ctx := context.Background()
//	err := client.SessionPolicies.Create(ctx, request)
//	require.NoError(t, err)
//	sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
//	require.NoError(t, err)
//	return sessionPolicy, func() {
//		err := client.SessionPolicies.Drop(ctx, NewDropSessionPolicyRequest(id))
//		require.NoError(t, err)
//	}
//}

//func createResourceMonitor(t *testing.T, client *Client) (*ResourceMonitor, func()) {
//	t.Helper()
//	return createResourceMonitorWithOptions(t, client, &CreateResourceMonitorOptions{
//		With: &ResourceMonitorWith{
//			CreditQuota: Pointer(100),
//			Triggers: []TriggerDefinition{
//				{
//					Threshold:     100,
//					TriggerAction: TriggerActionSuspend,
//				},
//				{
//					Threshold:     70,
//					TriggerAction: TriggerActionSuspendImmediate,
//				},
//				{
//					Threshold:     90,
//					TriggerAction: TriggerActionNotify,
//				},
//			},
//		},
//	})
//}
//
//func createResourceMonitorWithOptions(t *testing.T, client *Client, opts *CreateResourceMonitorOptions) (*ResourceMonitor, func()) {
//	t.Helper()
//	id := randomAccountObjectIdentifier(t)
//	ctx := context.Background()
//	err := client.ResourceMonitors.Create(ctx, id, opts)
//	require.NoError(t, err)
//	resourceMonitor, err := client.ResourceMonitors.ShowByID(ctx, id)
//	require.NoError(t, err)
//	return resourceMonitor, func() {
//		err := client.ResourceMonitors.Drop(ctx, id)
//		require.NoError(t, err)
//	}
//}

func createFailoverGroup(t *testing.T, client *Client) (*FailoverGroup, func()) {
	t.Helper()
	objectTypes := []PluralObjectType{PluralObjectTypeRoles}
	ctx := context.Background()
	currentAccount, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	accountID := NewAccountIdentifierFromAccountLocator(currentAccount)
	allowedAccounts := []AccountIdentifier{accountID}
	return createFailoverGroupWithOptions(t, client, objectTypes, allowedAccounts, nil)
}

func createFailoverGroupWithOptions(t *testing.T, client *Client, objectTypes []PluralObjectType, allowedAccounts []AccountIdentifier, opts *CreateFailoverGroupOptions) (*FailoverGroup, func()) {
	t.Helper()
	id := randomAccountObjectIdentifier(t)
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

func createShare(t *testing.T, client *Client) (*Share, func()) {
	t.Helper()
	return createShareWithOptions(t, client, &CreateShareOptions{})
}

func createShareWithOptions(t *testing.T, client *Client, opts *CreateShareOptions) (*Share, func()) {
	t.Helper()
	id := randomAccountObjectIdentifier(t)
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

func createFileFormat(t *testing.T, client *Client, schema DatabaseObjectIdentifier) (*FileFormat, func()) {
	t.Helper()
	return createFileFormatWithOptions(t, client, schema, &CreateFileFormatOptions{
		Type: FileFormatTypeCSV,
	})
}

func createFileFormatWithOptions(t *testing.T, client *Client, schema DatabaseObjectIdentifier, opts *CreateFileFormatOptions) (*FileFormat, func()) {
	t.Helper()
	id := NewSchemaObjectIdentifier(schema.databaseName, schema.name, randomString(t))
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

//func createWarehouse(t *testing.T, client *Client) (*Warehouse, func()) {
//	t.Helper()
//	return createWarehouseWithOptions(t, client, &CreateWarehouseOptions{})
//}
//
//func createWarehouseWithOptions(t *testing.T, client *Client, opts *CreateWarehouseOptions) (*Warehouse, func()) {
//	t.Helper()
//	name := randomStringRange(t, 8, 28)
//	id := NewAccountObjectIdentifier(name)
//	ctx := context.Background()
//	err := client.Warehouses.Create(ctx, id, opts)
//	require.NoError(t, err)
//	return &Warehouse{
//			Name: name,
//		}, func() {
//			err := client.Warehouses.Drop(ctx, id, nil)
//			require.NoError(t, err)
//		}
//}

//func createDatabase(t *testing.T, client *Client) (*Database, func()) {
//	t.Helper()
//	return createDatabaseWithOptions(t, client, randomAccountObjectIdentifier(t), &CreateDatabaseOptions{})
//}
//
//func createDatabaseWithIdentifier(t *testing.T, client *Client, id AccountObjectIdentifier) (*Database, func()) {
//	t.Helper()
//	return createDatabaseWithOptions(t, client, id, &CreateDatabaseOptions{})
//}
//
//func createDatabaseWithOptions(t *testing.T, client *Client, id AccountObjectIdentifier, opts *CreateDatabaseOptions) (*Database, func()) {
//	t.Helper()
//	ctx := context.Background()
//	err := client.Databases.Create(ctx, id, opts)
//	require.NoError(t, err)
//	database, err := client.Databases.ShowByID(ctx, id)
//	require.NoError(t, err)
//	return database, func() {
//		err := client.Databases.Drop(ctx, id, nil)
//		if errors.Is(err, errObjectNotExistOrAuthorized) {
//			return
//		}
//		require.NoError(t, err)
//	}
//}

//func createSchema(t *testing.T, client *Client, database *Database) (*Schema, func()) {
//	t.Helper()
//	return createSchemaWithIdentifier(t, client, database, randomStringRange(t, 8, 28))
//}
//
//func createSchemaWithIdentifier(t *testing.T, client *Client, database *Database, name string) (*Schema, func()) {
//	t.Helper()
//	ctx := context.Background()
//	schemaID := NewDatabaseObjectIdentifier(database.Name, name)
//	err := client.Schemas.Create(ctx, schemaID, nil)
//	require.NoError(t, err)
//	schema, err := client.Schemas.ShowByID(ctx, NewDatabaseObjectIdentifier(database.Name, name))
//	require.NoError(t, err)
//	return schema, func() {
//		err := client.Schemas.Drop(ctx, schemaID, nil)
//		if errors.Is(err, errObjectNotExistOrAuthorized) {
//			return
//		}
//		require.NoError(t, err)
//	}
//}

//func createTable(t *testing.T, client *Client, database *Database, schema *Schema) (*Table, func()) {
//	t.Helper()
//	name := randomStringRange(t, 8, 28)
//	ctx := context.Background()
//	_, err := client.exec(ctx, fmt.Sprintf("CREATE TABLE \"%s\".\"%s\".\"%s\" (id NUMBER)", database.Name, schema.Name, name))
//	require.NoError(t, err)
//	return &Table{
//			DatabaseName: database.Name,
//			SchemaName:   schema.Name,
//			Name:         name,
//		}, func() {
//			_, err := client.exec(ctx, fmt.Sprintf("DROP TABLE \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
//			require.NoError(t, err)
//		}
//}

//func createTag(t *testing.T, client *Client, database *Database, schema *Schema) (*Tag, func()) {
//	t.Helper()
//	return createTagWithOptions(t, client, database, schema, &CreateTagOptions{})
//}
//
//func createTagWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, _ *CreateTagOptions) (*Tag, func()) {
//	t.Helper()
//	name := randomStringRange(t, 8, 28)
//	ctx := context.Background()
//	_, err := client.exec(ctx, fmt.Sprintf("CREATE TAG \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
//	require.NoError(t, err)
//	return &Tag{
//			Name:         name,
//			DatabaseName: database.Name,
//			SchemaName:   schema.Name,
//		}, func() {
//			_, err := client.exec(ctx, fmt.Sprintf("DROP TAG \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
//			require.NoError(t, err)
//		}
//}

//func createPasswordPolicyWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, options *CreatePasswordPolicyOptions) (*PasswordPolicy, func()) {
//	t.Helper()
//	var databaseCleanup func()
//	if database == nil {
//		database, databaseCleanup = createDatabase(t, client)
//	}
//	var schemaCleanup func()
//	if schema == nil {
//		schema, schemaCleanup = createSchema(t, client, database)
//	}
//	name := randomUUID(t)
//	id := NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
//	ctx := context.Background()
//	err := client.PasswordPolicies.Create(ctx, id, options)
//	require.NoError(t, err)
//
//	showOptions := &ShowPasswordPolicyOptions{
//		Like: &Like{
//			Pattern: String(name),
//		},
//		In: &In{
//			Schema: schema.ID(),
//		},
//	}
//	passwordPolicyList, err := client.PasswordPolicies.Show(ctx, showOptions)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(passwordPolicyList))
//	return &passwordPolicyList[0], func() {
//		err := client.PasswordPolicies.Drop(ctx, id, nil)
//		require.NoError(t, err)
//		if schemaCleanup != nil {
//			schemaCleanup()
//		}
//		if databaseCleanup != nil {
//			databaseCleanup()
//		}
//	}
//}
//
//func createPasswordPolicy(t *testing.T, client *Client, database *Database, schema *Schema) (*PasswordPolicy, func()) {
//	t.Helper()
//	return createPasswordPolicyWithOptions(t, client, database, schema, nil)
//}

func createMaskingPolicyWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, signature []TableColumnSignature, returns DataType, expression string, options *CreateMaskingPolicyOptions) (*MaskingPolicy, func()) {
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
	id := NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.MaskingPolicies.Create(ctx, id, signature, returns, expression, options)
	require.NoError(t, err)

	showOptions := &ShowMaskingPolicyOptions{
		Like: &Like{
			Pattern: String(name),
		},
		In: &In{
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

func createRole(t *testing.T, client *Client) (*Role, func()) {
	t.Helper()
	return createRoleWithRequest(t, client, NewCreateRoleRequest(randomAccountObjectIdentifier(t)))
}

func createRoleWithRequest(t *testing.T, client *Client, req *CreateRoleRequest) (*Role, func()) {
	t.Helper()
	require.True(t, validObjectidentifier(req.name))
	ctx := context.Background()
	err := client.Roles.Create(ctx, req)
	require.NoError(t, err)
	role, err := client.Roles.ShowByID(ctx, NewShowByIdRoleRequest(req.name))
	require.NoError(t, err)
	return role, func() {
		err = client.Roles.Drop(ctx, NewDropRoleRequest(req.name))
		require.NoError(t, err)
	}
}

func createDatabaseRole(t *testing.T, client *Client, database *Database) (*DatabaseRole, func()) {
	t.Helper()
	name := randomString(t)
	id := NewDatabaseObjectIdentifier(database.Name, name)
	ctx := context.Background()

	err := client.DatabaseRoles.Create(ctx, NewCreateDatabaseRoleRequest(id))
	require.NoError(t, err)

	databaseRole, err := client.DatabaseRoles.ShowByID(ctx, id)
	require.NoError(t, err)

	return databaseRole, cleanupDatabaseRoleProvider(t, ctx, client, id)
}

func cleanupDatabaseRoleProvider(t *testing.T, ctx context.Context, client *Client, id DatabaseObjectIdentifier) func() {
	t.Helper()
	return func() {
		err := client.DatabaseRoles.Drop(ctx, NewDropDatabaseRoleRequest(id))
		require.NoError(t, err)
	}
}

func createMaskingPolicy(t *testing.T, client *Client, database *Database, schema *Schema) (*MaskingPolicy, func()) {
	t.Helper()
	signature := []TableColumnSignature{
		{
			Name: randomString(t),
			Type: DataTypeVARCHAR,
		},
	}
	n := randomIntRange(t, 0, 5)
	for i := 0; i < n; i++ {
		signature = append(signature, TableColumnSignature{
			Name: randomString(t),
			Type: DataTypeVARCHAR,
		})
	}
	expression := "REPLACE('X', 1, 2)"
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, DataTypeVARCHAR, expression, &CreateMaskingPolicyOptions{})
}

func createAlertWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, warehouse *Warehouse, schedule string, condition string, action string, opts *CreateAlertOptions) (*Alert, func()) {
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
	id := NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.Alerts.Create(ctx, id, warehouse.ID(), schedule, condition, action, opts)
	require.NoError(t, err)

	showOptions := &ShowAlertOptions{
		Like: &Like{
			Pattern: String(name),
		},
		In: &In{
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

func createAlert(t *testing.T, client *Client, database *Database, schema *Schema, warehouse *Warehouse) (*Alert, func()) {
	t.Helper()
	schedule := "USING CRON * * * * * UTC"
	condition := "SELECT 1"
	action := "SELECT 1"
	return createAlertWithOptions(t, client, database, schema, warehouse, schedule, condition, action, &CreateAlertOptions{})
}

func ParseTimestampWithOffset(s string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02T15:04:05-07:00", s)
	if err != nil {
		return nil, err
	}
	_, offset := t.Zone()
	adjustedTime := t.Add(-time.Duration(offset) * time.Second)
	return &adjustedTime, nil
}

//func createUser(t *testing.T, client *Client) (*User, func()) {
//	t.Helper()
//	name := randomStringRange(t, 8, 28)
//	id := NewAccountObjectIdentifier(name)
//	return createUserWithOptions(t, client, id, &CreateUserOptions{})
//}
//
//func createUserWithName(t *testing.T, client *Client, name string) (*User, func()) {
//	t.Helper()
//	id := NewAccountObjectIdentifier(name)
//	return createUserWithOptions(t, client, id, &CreateUserOptions{})
//}
//
//func createUserWithOptions(t *testing.T, client *Client, id AccountObjectIdentifier, opts *CreateUserOptions) (*User, func()) {
//	t.Helper()
//	ctx := context.Background()
//	err := client.Users.Create(ctx, id, opts)
//	require.NoError(t, err)
//	user, err := client.Users.ShowByID(ctx, id)
//	require.NoError(t, err)
//	return user, func() {
//		err := client.Users.Drop(ctx, id)
//		require.NoError(t, err)
//	}
//}

//func createPipe(t *testing.T, client *Client, database *Database, schema *Schema, name string, copyStatement string) (*Pipe, func()) {
//	t.Helper()
//	require.NotNil(t, database, "database has to be created")
//	require.NotNil(t, schema, "schema has to be created")
//
//	id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
//	ctx := context.Background()
//
//	pipeCleanup := func() {
//		err := client.Pipes.Drop(ctx, id)
//		require.NoError(t, err)
//	}
//
//	err := client.Pipes.Create(ctx, id, copyStatement, &CreatePipeOptions{})
//	if err != nil {
//		return nil, pipeCleanup
//	}
//	require.NoError(t, err)
//
//	createdPipe, errDescribe := client.Pipes.Describe(ctx, id)
//	if errDescribe != nil {
//		return nil, pipeCleanup
//	}
//	require.NoError(t, errDescribe)
//
//	return createdPipe, pipeCleanup
//}

//func createStageWithName(t *testing.T, client *Client, name string) (*string, func()) {
//	t.Helper()
//	ctx := context.Background()
//	stageCleanup := func() {
//		_, err := client.exec(ctx, fmt.Sprintf("DROP STAGE %s", name))
//		require.NoError(t, err)
//	}
//	_, err := client.exec(ctx, fmt.Sprintf("CREATE STAGE %s", name))
//	if err != nil {
//		return nil, stageCleanup
//	}
//	require.NoError(t, err)
//	return &name, stageCleanup
//}
//
//func createStage(t *testing.T, client *Client, database *Database, schema *Schema, name string) (*Stage, func()) {
//	t.Helper()
//	require.NotNil(t, database, "database has to be created")
//	require.NotNil(t, schema, "schema has to be created")
//
//	id := NewSchemaObjectIdentifier(database.Name, schema.Name, name)
//	ctx := context.Background()
//
//	stageCleanup := func() {
//		_, err := client.exec(ctx, fmt.Sprintf("DROP STAGE %s", id.FullyQualifiedName()))
//		require.NoError(t, err)
//	}
//
//	_, err := client.exec(ctx, fmt.Sprintf("CREATE STAGE %s", id.FullyQualifiedName()))
//	if err != nil {
//		return nil, stageCleanup
//	}
//	require.NoError(t, err)
//
//	return &Stage{
//		DatabaseName: database.Name,
//		SchemaName:   schema.Name,
//		Name:         name,
//	}, stageCleanup
//}

func createDynamicTable(t *testing.T, client *Client) (*DynamicTable, func()) {
	t.Helper()
	return createDynamicTableWithOptions(t, client, nil, nil, nil, nil)
}

func createDynamicTableWithOptions(t *testing.T, client *Client, warehouse *Warehouse, database *Database, schema *Schema, table *Table) (*DynamicTable, func()) {
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
	name := NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, randomString(t))
	targetLag := TargetLag{
		Lagtime: String("2 minutes"),
	}
	query := "select id from " + table.ID().FullyQualifiedName()
	comment := randomComment(t)
	ctx := context.Background()
	err := client.DynamicTables.Create(ctx, NewCreateDynamicTableRequest(name, warehouse.ID(), targetLag, query).WithOrReplace(true).WithComment(&comment))
	require.NoError(t, err)
	entities, err := client.DynamicTables.Show(ctx, NewShowDynamicTableRequest().WithLike(&Like{Pattern: String(name.Name())}).WithIn(&In{Schema: schema.ID()}))
	require.NoError(t, err)
	require.Equal(t, 1, len(entities))
	return &entities[0], func() {
		require.NoError(t, client.DynamicTables.Drop(ctx, NewDropDynamicTableRequest(name)))
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

//func createStageWithURL(t *testing.T, client *Client, name AccountObjectIdentifier, url string) (*Stage, func()) {
//	t.Helper()
//	ctx := context.Background()
//	_, err := client.exec(ctx, fmt.Sprintf(`CREATE STAGE "%s" URL = '%s'`, name.Name(), url))
//	require.NoError(t, err)
//
//	return nil, func() {
//		_, err := client.exec(ctx, fmt.Sprintf(`DROP STAGE "%s"`, name.Name()))
//		require.NoError(t, err)
//	}
//}

//func createNetworkPolicy(t *testing.T, client *Client, req *CreateNetworkPolicyRequest) (error, func()) {
//	t.Helper()
//	ctx := context.Background()
//	err := client.NetworkPolicies.Create(ctx, req)
//	return err, func() {
//		err := client.NetworkPolicies.Drop(ctx, NewDropNetworkPolicyRequest(req.name))
//		require.NoError(t, err)
//	}
//}
