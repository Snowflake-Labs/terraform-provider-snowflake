package sdk

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-uuid"
	"github.com/stretchr/testify/require"
)

func testClient(t *testing.T) *Client {
	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

// mock structs until we have more of the SDK implemented
type DatabaseCreateOptions struct{}

type Database struct {
	Name string
}

func (v *Database) Identifier() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

type Schema struct {
	DatabaseName string
	Name         string
}

func (v *Schema) Identifier() SchemaIdentifier {
	return NewSchemaIdentifier(v.DatabaseName, v.Name)
}

func randomString(t *testing.T) string {
	v, err := uuid.GenerateUUID()
	require.NoError(t, err)
	return v
}

func createDatabase(t *testing.T, client *Client) (*Database, func()) {
	return createDatabaseWithOptions(t, client, &DatabaseCreateOptions{})
}

func createDatabaseWithOptions(t *testing.T, client *Client, opts *DatabaseCreateOptions) (*Database, func()) {
	name := randomString(t)
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
	name := randomString(t)
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

func createPasswordPolicyWithOptions(t *testing.T, client *Client, database *Database, schema *Schema, options *PasswordPolicyCreateOptions) (*PasswordPolicy, func()) {
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
	err := client.PasswordPolicies.Create(ctx, id, options)
	require.NoError(t, err)

	showOptions := &PasswordPolicyShowOptions{
		Like: &Like{
			Pattern: String(name),
		},
		In: &In{
			Schema: schema.Identifier(),
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
	return createPasswordPolicyWithOptions(t, client, database, schema, nil)
}
