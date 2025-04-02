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

	value, err := sdk.SafeShowById(testClient(t), networkPolicyShowById, testContext(t), networkPolicy.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupNetworkPolicy()

	value, err = sdk.SafeShowById(testClient(t), networkPolicyShowById, testContext(t), networkPolicy.ID())
	assert.Nil(t, value)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
}

func TestInt_SafeShowByIdOnDatabaseObjectIdentifier(t *testing.T) {
	databaseRoleShowById := func(ctx context.Context, id sdk.DatabaseObjectIdentifier) (*sdk.DatabaseRole, error) {
		return testClientHelper().DatabaseRole.Show(t, id)
	}

	databaseRole, cleanupDatabaseRole := testClientHelper().DatabaseRole.CreateDatabaseRole(t)

	value, err := sdk.SafeShowById(testClient(t), databaseRoleShowById, testContext(t), databaseRole.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupDatabaseRole()

	value, err = sdk.SafeShowById(testClient(t), databaseRoleShowById, testContext(t), databaseRole.ID())
	assert.Nil(t, value)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidDatabaseRoleId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)

	value, err = sdk.SafeShowById(testClient(t), databaseRoleShowById, testContext(t), invalidDatabaseRoleId)
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}

func TestInt_SafeShowByIdOnSchemaObjectIdentifier(t *testing.T) {
	tableShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifier) (*sdk.Table, error) {
		return testClientHelper().Table.Show(t, id)
	}

	table, cleanupTable := testClientHelper().Table.Create(t)

	value, err := sdk.SafeShowById(testClient(t), tableShowById, testContext(t), table.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupTable()

	value, err = sdk.SafeShowById(testClient(t), tableShowById, testContext(t), table.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	invalidSchemaIdOnValidDatabase := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
	invalidTableIdOnValidDatabase := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(invalidSchemaIdOnValidDatabase)

	value, err = sdk.SafeShowById(testClient(t), tableShowById, testContext(t), invalidTableIdOnValidDatabase)
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)
	invalidTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(invalidSchemaId)

	value, err = sdk.SafeShowById(testClient(t), tableShowById, testContext(t), invalidTableId)
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
}

func TestInt_SafeShowByIdOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	procedureShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, error) {
		return testClientHelper().Procedure.Show(t, id)
	}

	procedure, cleanupProcedure := testClientHelper().Procedure.Create(t, sdk.DataTypeInt)

	value, err := sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), procedure.ID())
	assert.NotNil(t, value)
	assert.NoError(t, err)

	cleanupProcedure()

	value, err = sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), procedure.ID())
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)

	invalidSchemaIdOnValidDatabase := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
	invalidProcedureIdOnValidDatabase := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(invalidSchemaIdOnValidDatabase)

	value, err = sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), invalidProcedureIdOnValidDatabase)
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)
	invalidProcedureId := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(invalidSchemaId)

	value, err = sdk.SafeShowById(testClient(t), procedureShowById, testContext(t), invalidProcedureId)
	assert.Nil(t, value)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
}
