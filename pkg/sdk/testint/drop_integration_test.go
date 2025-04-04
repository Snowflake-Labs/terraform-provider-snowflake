package testint

import (
	"context"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
)

func TestInt_SafeDropOnAccountObjectIdentifier(t *testing.T) {
	networkPolicy, cleanupNetworkPolicy := testClientHelper().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(cleanupNetworkPolicy)

	ctx := context.Background()
	networkPolicyDrop := func() error {
		return testClient(t).NetworkPolicies.Drop(ctx, sdk.NewDropNetworkPolicyRequest(networkPolicy.ID()).WithIfExists(true))
	}

	err := sdk.SafeDrop(testClient(t), networkPolicyDrop, ctx, networkPolicy.ID())
	assert.NoError(t, err)

	err = sdk.SafeDrop(testClient(t), networkPolicyDrop, ctx, networkPolicy.ID())
	assert.NoError(t, err)
}

func TestInt_SafeDropOnDatabaseObjectIdentifier(t *testing.T) {
	databaseRole, cleanupDatabaseRole := testClientHelper().DatabaseRole.CreateDatabaseRole(t)
	t.Cleanup(cleanupDatabaseRole)

	ctx := context.Background()
	databaseRoleDrop := func(id sdk.DatabaseObjectIdentifier) func() error {
		return func() error {
			return testClient(t).DatabaseRoles.Drop(ctx, sdk.NewDropDatabaseRoleRequest(id).WithIfExists(true))
		}
	}

	err := sdk.SafeDrop(testClient(t), databaseRoleDrop(databaseRole.ID()), testContext(t), databaseRole.ID())
	assert.NoError(t, err)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidDatabaseRoleId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)

	err = sdk.SafeDrop(testClient(t), databaseRoleDrop(invalidDatabaseRoleId), testContext(t), invalidDatabaseRoleId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrSkippable)
}

func TestInt_SafeDropOnSchemaObjectIdentifier(t *testing.T) {
	table, cleanupTable := testClientHelper().Table.Create(t)
	t.Cleanup(cleanupTable)

	ctx := context.Background()
	tableDrop := func(id sdk.SchemaObjectIdentifier) func() error {
		return func() error {
			return testClient(t).Tables.Drop(ctx, sdk.NewDropTableRequest(id).WithIfExists(sdk.Bool(true)))
		}
	}

	err := sdk.SafeDrop(testClient(t), tableDrop(table.ID()), ctx, table.ID())
	assert.NoError(t, err)

	invalidSchemaIdOnValidDatabase := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
	invalidTableIdOnValidDatabase := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(invalidSchemaIdOnValidDatabase)

	err = sdk.SafeDrop(testClient(t), tableDrop(invalidTableIdOnValidDatabase), ctx, invalidTableIdOnValidDatabase)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrSkippable)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)
	invalidTableId := testClientHelper().Ids.RandomSchemaObjectIdentifierInSchema(invalidSchemaId)

	err = sdk.SafeDrop(testClient(t), tableDrop(invalidTableId), ctx, invalidTableId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrSkippable)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
}

func TestInt_SafeDropOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	procedure := testClientHelper().Procedure.Create(t, sdk.DataTypeInt)

	ctx := context.Background()
	procedureDrop := func(id sdk.SchemaObjectIdentifierWithArguments) func() error {
		return func() error {
			return testClient(t).Procedures.Drop(ctx, sdk.NewDropProcedureRequest(id).WithIfExists(true))
		}
	}

	err := sdk.SafeDrop(testClient(t), procedureDrop(procedure.ID()), ctx, procedure.ID())
	assert.NoError(t, err)

	invalidSchemaIdOnValidDatabase := testClientHelper().Ids.RandomDatabaseObjectIdentifier()
	invalidProcedureIdOnValidDatabase := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(invalidSchemaIdOnValidDatabase, sdk.DataTypeInt)

	err = sdk.SafeDrop(testClient(t), procedureDrop(invalidProcedureIdOnValidDatabase), ctx, invalidProcedureIdOnValidDatabase)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrSkippable)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)

	invalidDatabaseId := testClientHelper().Ids.RandomAccountObjectIdentifier()
	invalidSchemaId := testClientHelper().Ids.RandomDatabaseObjectIdentifierInDatabase(invalidDatabaseId)
	invalidProcedureId := testClientHelper().Ids.RandomSchemaObjectIdentifierWithArgumentsInSchema(invalidSchemaId, sdk.DataTypeInt)

	err = sdk.SafeDrop(testClient(t), procedureDrop(invalidProcedureId), ctx, invalidProcedureId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, sdk.ErrSkippable)
	assert.ErrorIs(t, err, sdk.ErrObjectNotFound)
	assert.ErrorIs(t, err, sdk.ErrObjectNotExistOrAuthorized)
	assert.ErrorIs(t, err, sdk.ErrDoesNotExistOrOperationCannotBePerformed)
}
