package resources_test

import (
	"context"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	resourcenames "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_SafeShowByIdOnAccountObjectIdentifier(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	networkPolicyShowById := func(ctx context.Context, id sdk.AccountObjectIdentifier) (*sdk.NetworkPolicy, error) {
		return acc.TestClient().NetworkPolicy.Show(t, id)
	}

	networkPolicy, cleanupNetworkPolicy := acc.TestClient().NetworkPolicy.CreateNetworkPolicy(t)
	t.Cleanup(cleanupNetworkPolicy)

	value, shouldRemoveFromState, diags := resources.SafeShowById(resourcenames.NetworkPolicy, acc.TestClient().Client(), networkPolicyShowById, context.Background(), networkPolicy.ID())
	assert.NotNil(t, value)
	assert.Equal(t, false, shouldRemoveFromState)
	assert.Len(t, diags, 0)

	cleanupNetworkPolicy()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.NetworkPolicy, acc.TestClient().Client(), networkPolicyShowById, context.Background(), networkPolicy.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "Failed to query account object.")
}

func Test_SafeShowByIdOnDatabaseObjectIdentifier(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseRoleShowById := func(ctx context.Context, id sdk.DatabaseObjectIdentifier) (*sdk.DatabaseRole, error) {
		return acc.TestClient().DatabaseRole.Show(t, id)
	}

	database, cleanupDatabase := acc.TestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	databaseRole, cleanupDatabaseRole := acc.TestClient().DatabaseRole.CreateDatabaseRoleInDatabase(t, database.ID())

	value, shouldRemoveFromState, diags := resources.SafeShowById(resourcenames.DatabaseRole, acc.TestClient().Client(), databaseRoleShowById, context.Background(), databaseRole.ID())
	assert.NotNil(t, value)
	assert.Equal(t, false, shouldRemoveFromState)
	assert.Len(t, diags, 0)

	cleanupDatabaseRole()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.DatabaseRole, acc.TestClient().Client(), databaseRoleShowById, context.Background(), databaseRole.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	assert.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "Failed to query database object.")

	cleanupDatabase()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.DatabaseRole, acc.TestClient().Client(), databaseRoleShowById, context.Background(), databaseRole.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	assert.Len(t, diags, 2)
	assert.Contains(t, diags[0].Summary, "Failed to query database object.")
	assert.Contains(t, diags[1].Summary, "Failed to query database for snowflake_database_role.")
}

func Test_SafeShowByIdOnSchemaObjectIdentifier(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	tableShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifier) (*sdk.Table, error) {
		return acc.TestClient().Table.Show(t, id)
	}

	database, cleanupDatabase := acc.TestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	schema, cleanupSchema := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	table, cleanupTable := acc.TestClient().Table.CreateInSchema(t, schema.ID())

	value, shouldRemoveFromState, diags := resources.SafeShowById(resourcenames.Table, acc.TestClient().Client(), tableShowById, context.Background(), table.ID())
	assert.NotNil(t, value)
	assert.Equal(t, false, shouldRemoveFromState)
	assert.Len(t, diags, 0)

	cleanupTable()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.Table, acc.TestClient().Client(), tableShowById, context.Background(), table.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "Failed to query schema object.")

	cleanupSchema()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.Table, acc.TestClient().Client(), tableShowById, context.Background(), table.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 2)
	assert.Contains(t, diags[0].Summary, "Failed to query schema object.")
	assert.Contains(t, diags[1].Summary, "Failed to query schema for snowflake_table.")

	cleanupDatabase()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.Table, acc.TestClient().Client(), tableShowById, context.Background(), table.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 2)
	assert.Contains(t, diags[0].Summary, "Failed to query schema object.")
	assert.Contains(t, diags[1].Summary, "Failed to query schema for snowflake_table.")
}

func Test_SafeShowByIdOnSchemaObjectIdentifierWithArguments(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	procedureShowById := func(ctx context.Context, id sdk.SchemaObjectIdentifierWithArguments) (*sdk.Procedure, error) {
		return acc.TestClient().Procedure.Show(t, id)
	}

	database, cleanupDatabase := acc.TestClient().Database.CreateDatabase(t)
	t.Cleanup(cleanupDatabase)
	schema, cleanupSchema := acc.TestClient().Schema.CreateSchemaInDatabase(t, database.ID())
	procedure := acc.TestClient().Procedure.CreateInSchema(t, schema.ID(), sdk.DataTypeInt)

	value, shouldRemoveFromState, diags := resources.SafeShowById(resourcenames.ProcedureSql, acc.TestClient().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.NotNil(t, value)
	assert.Equal(t, false, shouldRemoveFromState)
	assert.Len(t, diags, 0)

	acc.TestClient().Procedure.DropProcedureFunc(t, procedure.ID())()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.ProcedureSql, acc.TestClient().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 1)
	assert.Contains(t, diags[0].Summary, "Failed to query schema object.")

	cleanupSchema()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.ProcedureSql, acc.TestClient().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 2)
	assert.Contains(t, diags[0].Summary, "Failed to query schema object.")
	assert.Contains(t, diags[1].Summary, "Failed to query schema for snowflake_procedure_sql.")

	cleanupDatabase()

	value, shouldRemoveFromState, diags = resources.SafeShowById(resourcenames.ProcedureSql, acc.TestClient().Client(), procedureShowById, context.Background(), procedure.ID())
	assert.Nil(t, value)
	assert.Equal(t, true, shouldRemoveFromState)
	require.Len(t, diags, 2)
	assert.Contains(t, diags[0].Summary, "Failed to query schema object.")
	assert.Contains(t, diags[1].Summary, "Failed to query schema for snowflake_procedure_sql.")
}
