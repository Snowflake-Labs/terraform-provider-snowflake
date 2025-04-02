package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestInt_SafeShowByIdOnAccountObjectIdentifier(t *testing.T) {
	networkPolicyShowById := func(ctx context.Context, id sdk.AccountObjectIdentifier) (*sdk.NetworkPolicy, error) {
		return testClientHelper().NetworkPolicy.Show(t, id)
	}

	networkPolicy, cleanupNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(cleanupNetworkPolicy)

	value, err := sdk.SafeShowById(testClientHelper().Client(), networkPolicyShowById, context.Background(), networkPolicy.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupNetworkPolicy()

	value, err = sdk.SafeShowById(testClientHelper().Client(), networkPolicyShowById, context.Background(), networkPolicy.ID())
	assert.Nil(t, value)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
}

func TestInt_SafeShowByIdOnDatabaseObjectIdentifier(t *testing.T) {
	databaseRoleShowById := func(ctx context.Context, id sdk.DatabaseObjectIdentifier) (*sdk.DatabaseRole, error) {
		return testClientHelper().DatabaseRole.Show(t, id)
	}

	// TODO: Use common database
	database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	databaseRole, cleanupDatabaseRole := testClientHelper().DatabaseRole.CreateDatabaseRoleInDatabase(t, database.ID())

	value, err := sdk.SafeShowById(testClientHelper().Client(), databaseRoleShowById, context.Background(), databaseRole.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupDatabaseRole()

	value, err = sdk.SafeShowById(testClientHelper().Client(), databaseRoleShowById, context.Background(), databaseRole.ID())
	assert.Nil(t, value)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	cleanupDatabase()

	value, err = sdk.SafeShowById(testClientHelper().Client(), databaseRoleShowById, context.Background(), databaseRole.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_SafeShowByIdOnSchemaObjectIdentifier(t *testing.T) {
	tableShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifier) (*sdk.Table, error) {
		return testClientHelper().Table.Show(t, id)
	}

	// TODO: Use common database and schema
	database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	schema, cleanupSchema := testClientHelper().Schema.CreateSchemaInDatabase(t, database.ID())
	table, cleanupTable := testClientHelper().Table.CreateInSchema(t, schema.ID())

	value, err := sdk.SafeShowById(testClientHelper().Client(), tableShowById, context.Background(), table.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupTable()

	value, err = sdk.SafeShowById(testClientHelper().Client(), tableShowById, context.Background(), table.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	cleanupSchema()

	value, err = sdk.SafeShowById(testClientHelper().Client(), tableShowById, context.Background(), table.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)

	cleanupDatabase()

	value, err = sdk.SafeShowById(testClientHelper().Client(), tableShowById, context.Background(), table.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
}

func TestInt_SafeShowByIdOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	procedureShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, error) {
		return testClientHelper().Procedure.Show(t, id)
	}

	// TODO: Common database and schema
	database, cleanupDatabase := testClientHelper().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	schema, cleanupSchema := testClientHelper().Schema.CreateSchemaInDatabase(t, database.ID())
	procedure, cleanupProcedure := testClientHelper().Procedure.CreateInSchema(t, schema.ID(), sdk.DataTypeInt)

	value, err := sdk.SafeShowById(testClientHelper().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupProcedure()

	value, err = sdk.SafeShowById(testClientHelper().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	cleanupSchema()

	value, err = sdk.SafeShowById(testClientHelper().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

	cleanupDatabase()

	value, err = sdk.SafeShowById(testClientHelper().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}
