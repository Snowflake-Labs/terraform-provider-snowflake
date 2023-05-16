package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

func randomSchemaObjectIdentifier(t *testing.T) SchemaObjectIdentifier {
	t.Helper()
	return NewSchemaObjectIdentifier(randomStringRange(t, 8, 12), randomStringRange(t, 8, 12), randomStringRange(t, 8, 12))
}

func randomSchemaIdentifier(t *testing.T) SchemaIdentifier {
	t.Helper()
	return NewSchemaIdentifier(randomStringRange(t, 8, 12), randomStringRange(t, 8, 12))
}

func randomAccountObjectIdentifier(t *testing.T) AccountObjectIdentifier {
	t.Helper()
	return NewAccountObjectIdentifier(randomStringRange(t, 8, 12))
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

func createWarehouse(t *testing.T, client *Client) (*Warehouse, func()) {
	t.Helper()
	return createWarehouseWithOptions(t, client, &WarehouseCreateOptions{})
}

func createWarehouseWithOptions(t *testing.T, client *Client, opts *WarehouseCreateOptions) (*Warehouse, func()) {
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
	return createDatabaseWithOptions(t, client, &DatabaseCreateOptions{})
}

func createDatabaseWithOptions(t *testing.T, client *Client, _ *DatabaseCreateOptions) (*Database, func()) {
	t.Helper()
	name := randomStringRange(t, 8, 28)
	ctx := context.Background()
	_, err := client.exec(ctx, fmt.Sprintf("CREATE DATABASE \"%s\"", name))
	require.NoError(t, err)
	return &Database{
			Name: name,
		}, func() {
			_, err := client.exec(ctx, fmt.Sprintf("DROP DATABASE \"%s\"", name))
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

func createPasswordPolicyWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, options *PasswordPolicyCreateOptions) (*PasswordPolicy, func()) {
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

func createMaskingPolicyWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, signature []TableColumnSignature, returns DataType, expression string, options *MaskingPolicyCreateOptions) (*MaskingPolicy, func()) {
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

	showOptions := &MaskingPolicyShowOptions{
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
	return createMaskingPolicyWithOptions(t, client, database, schema, signature, DataTypeVARCHAR, expression, &MaskingPolicyCreateOptions{})
}
