package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

// there is no direct way to get the account identifier from Snowflake API, but you can get it if you know
// the account locator and by filtering the list of accounts in replication accounts by the account locator
func getAccountIdentifier(t *testing.T, client *Client) AccountIdentifier {
	t.Helper()
	ctx := context.Background()
	currentAccountLocator, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	replicationAccounts, err := client.ReplicationFunctions.ShowReplicationAcccounts(ctx)
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

func getPrimaryAccountIdentifier(t *testing.T) AccountIdentifier {
	t.Helper()
	client := testClient(t)
	return getAccountIdentifier(t, client)
}

func getSecondaryAccountIdentifier(t *testing.T) AccountIdentifier {
	t.Helper()
	client := testSecondaryClient(t)
	return getAccountIdentifier(t, client)
}

func randomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	return NewSchemaObjectIdentifier(randomStringN(t, 12), randomStringN(t, 12), randomStringN(t, 12))
}

func randomSchemaIdentifier(t *testing.T) SchemaIdentifier {
	t.Helper()
	return NewSchemaIdentifier(randomStringN(t, 12), randomStringN(t, 12))
}

func randomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	return NewAccountObjectIdentifier(randomStringN(t, 12))
}

func useDatabase(t *testing.T, client *Client, databaseID AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseDatabase(ctx, databaseID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseDatabase(ctx, NewAccountObjectIdentifier(orgDB))
		require.NoError(t, err)
	}
}

func useSchema(t *testing.T, client *Client, schemaID SchemaIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgDB, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	orgSchema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseSchema(ctx, schemaID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseSchema(ctx, NewSchemaIdentifier(orgDB, orgSchema))
		require.NoError(t, err)
	}
}

func useWarehouse(t *testing.T, client *Client, warehouseID AccountObjectIdentifier) func() {
	t.Helper()
	ctx := context.Background()
	orgWarehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	err = client.Sessions.UseWarehouse(ctx, warehouseID)
	require.NoError(t, err)
	return func() {
		err := client.Sessions.UseWarehouse(ctx, NewAccountObjectIdentifier(orgWarehouse))
		require.NoError(t, err)
	}
}

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

func randomUUID(t *testing.T) string {
	t.Helper()
	v, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return v
}

func randomComment(t *testing.T) string {
	t.Helper()
	return gofakeit.Sentence(10)
}

func randomBool(t *testing.T) bool {
	t.Helper()
	return gofakeit.Bool()
}

func randomString(t *testing.T) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, 28)
}

func randomStringN(t *testing.T, num int) string {
	t.Helper()
	return gofakeit.Password(true, true, true, true, false, num)
}

func randomStringRange(t *testing.T, min, max int) string {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.Password(true, true, true, true, false, randomIntRange(t, min, max))
}

func randomIntRange(t *testing.T, min, max int) int {
	t.Helper()
	if min > max {
		t.Errorf("min %d is greater than max %d", min, max)
	}
	return gofakeit.IntRange(min, max)
}

func createSessionPolicy(t *testing.T, client *Client, database *Database, schema *Schema) (*SessionPolicy, func()) {
	t.Helper()
	id := NewSchemaObjectIdentifier(database.Name, schema.Name, randomStringN(t, 12))
	return createSessionPolicyWithOptions(t, client, id, &CreateSessionPolicyOptions{})
}

func createSessionPolicyWithOptions(t *testing.T, client *Client, id SchemaObjectIdentifier, opts *CreateSessionPolicyOptions) (*SessionPolicy, func()) {
	t.Helper()
	ctx := context.Background()
	err := client.SessionPolicies.Create(ctx, id, opts)
	require.NoError(t, err)
	sessionPolicy, err := client.SessionPolicies.ShowByID(ctx, id)
	require.NoError(t, err)
	return sessionPolicy, func() {
		err := client.SessionPolicies.Drop(ctx, id, nil)
		require.NoError(t, err)
	}
}

func createResourceMonitor(t *testing.T, client *Client) (*ResourceMonitor, func()) {
	t.Helper()
	return createResourceMonitorWithOptions(t, client, &CreateResourceMonitorOptions{})
}

func createResourceMonitorWithOptions(t *testing.T, client *Client, opts *CreateResourceMonitorOptions) (*ResourceMonitor, func()) {
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

func createWarehouse(t *testing.T, client *Client) (*Warehouse, func()) {
	t.Helper()
	return createWarehouseWithOptions(t, client, &CreateWarehouseOptions{})
}

func createWarehouseWithOptions(t *testing.T, client *Client, opts *CreateWarehouseOptions) (*Warehouse, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	id := NewAccountObjectIdentifier(name)
	ctx := context.Background()
	err := client.Warehouses.Create(ctx, id, opts)
	require.NoError(t, err)
	return &Warehouse{
			Name: name,
		}, func() {
			err := client.Warehouses.Drop(ctx, id, nil)
			require.NoError(t, err)
		}
}

func createDatabase(t *testing.T, client *Client) (*Database, func()) {
	t.Helper()
	return createDatabaseWithOptions(t, client, &CreateDatabaseOptions{})
}

func createDatabaseWithOptions(t *testing.T, client *Client, _ *CreateDatabaseOptions) (*Database, func()) {
	t.Helper()
	id := randomAccountObjectIdentifier(t)
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

func createSchema(t *testing.T, client *Client, database *Database) (*Schema, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf("CREATE SCHEMA \"%s\".\"%s\"", database.Name, name))
	require.NoError(t, err)
	return &Schema{
			DatabaseName: database.Name,
			Name:         name,
		}, func() {
			_, err := client.exec(ctx, fmt.Sprintf("DROP SCHEMA \"%s\".\"%s\"", database.Name, name))
			require.NoError(t, err)
		}
}

func createTag(t *testing.T, client *Client, database *Database, schema *Schema) (*Tag, func()) {
	t.Helper()
	return createTagWithOptions(t, client, database, schema, &TagCreateOptions{})
}

func createTagWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, _ *TagCreateOptions) (*Tag, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf("CREATE TAG \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
	require.NoError(t, err)
	return &Tag{
			Name:         name,
			DatabaseName: database.Name,
			SchemaName:   schema.Name,
		}, func() {
			_, err := client.exec(ctx, fmt.Sprintf("DROP TAG \"%s\".\"%s\".\"%s\"", database.Name, schema.Name, name))
			require.NoError(t, err)
		}
}

func createPasswordPolicyWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, options *CreatePasswordPolicyOptions) (*PasswordPolicy, func()) {
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
	id := NewSchemaObjectIdentifier(schema.DatabaseName, schema.Name, name)
	ctx := context.Background()
	err := client.PasswordPolicies.Create(ctx, id, options)
	require.NoError(t, err)

	showOptions := &PasswordPolicyShowOptions{
		Like: &Like{
			Pattern: String(name),
		},
		In: &In{
			Schema: schema.ID(),
		},
	}
	passwordPolicyList, err := client.PasswordPolicies.Show(ctx, showOptions)
	require.NoError(t, err)
	require.Equal(t, 1, len(passwordPolicyList))
	return passwordPolicyList[0], func() {
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

func createPasswordPolicy(t *testing.T, client *Client, database *Database, schema *Schema) (*PasswordPolicy, func()) {
	t.Helper()
	return createPasswordPolicyWithOptions(t, client, database, schema, nil)
}

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
	return maskingPolicyList[0], func() {
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
