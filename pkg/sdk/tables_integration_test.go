package sdk

import (
	"context"
	"testing"
)

func TestInt_TableCreate(t *testing.T) {
	client := testClient(t)
	ctx := context.Background()

	schemaIdentifier := NewSchemaIdentifier("TXR@=9,TBnLj", "tcK1>AJ+")
	database, databaseCleanup := createDatabaseWithIdentifier(t, client, AccountObjectIdentifier{schemaIdentifier.databaseName})
	t.Cleanup(databaseCleanup)

	schema, schemaCleanup := createSchemaWithIdentifier(t, client, database, schemaIdentifier.schemaName)
	t.Cleanup(schemaCleanup)

	t.Run("test complete", func(t *testing.T) {
		tableName := randomString(t)
		id := NewSchemaObjectIdentifier(database.Name, schema.Name, tableName)
		emptyColumns := []TableColumn{}
		options := &CreateTableOptions{
			Scope: Pointer(GlobalTableScope),
			Kind:  Pointer(TemporaryTableKind),
			name:  id,
		}
		//przerób interfejsy też, żeby używaly tego buildera juz
		client.Tables.Create(ctx, id, emptyColumns, options)
	})
}
