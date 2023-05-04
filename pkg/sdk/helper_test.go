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

func testBuilder(t *testing.T) *sqlBuilder {
	t.Helper()
	return &sqlBuilder{}
}

func testClient(t *testing.T) *Client {
	t.Helper()
	client, err := NewDefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	return client
}

// mock structs until we have more of the SDK implemented.
type DatabaseCreateOptions struct{}

type Database struct {
	Name string
}

func (v *Database) ID() AccountObjectIdentifier {
	return NewAccountObjectIdentifier(v.Name)
}

type Schema struct {
	DatabaseName string
	Name         string
}

func (v *Schema) ID() SchemaIdentifier {
	return NewSchemaIdentifier(v.DatabaseName, v.Name)
}

type TagCreateOptions struct{}

type Tag struct {
	DatabaseName string
	SchemaName   string
	Name         string
}

func (v *Tag) ID() SchemaObjectIdentifier {
	return NewSchemaObjectIdentifier(v.DatabaseName, v.SchemaName, v.Name)
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

/*
func randomDataType(t *testing.T) DataType {
	t.Helper()
	dataTypeList := []DataType{
		DataTypeNumber,
		DataTypeFloat,
		DataTypeVARCHAR,
		DataTypeBinary,
		DataTypeBoolean,
		DataTypeDate,
		DataTypeTime,
		DataTypeTimestampLTZ,
		DataTypeTimestampNTZ,
		DataTypeTimestampTZ,
		DataTypeVariant,
		DataTypeObject,
		DataTypeArray,
		DataTypeGeography,
		DataTypeGeometry,
	}
	return dataTypeList[randomIntN(t, 0, len(dataTypeList))]
}*/

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
