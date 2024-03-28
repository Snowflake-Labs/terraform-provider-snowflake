package testint

import (
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk/internal/random"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInt_CurrentAccount(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	account, err := client.ContextFunctions.CurrentAccount(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, account)
}

func TestInt_CurrentRole(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	role, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, role)
}

func TestInt_CurrentRegion(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	region, err := client.ContextFunctions.CurrentRegion(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, region)
}

func TestInt_CurrentSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	session, err := client.ContextFunctions.CurrentSession(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, session)
}

func TestInt_CurrentUser(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, user)
}

func TestInt_CurrentSessionDetails(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	account, err := client.ContextFunctions.CurrentSessionDetails(ctx)
	require.NoError(t, err)
	assert.NotNil(t, account)
	assert.NotEmpty(t, account.Account)
	assert.NotEmpty(t, account.Role)
	assert.NotEmpty(t, account.Region)
	assert.NotEmpty(t, account.Session)
	assert.NotEmpty(t, account.User)
}

func TestInt_CurrentDatabase(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	err := client.Sessions.UseDatabase(ctx, databaseTest.ID())
	require.NoError(t, err)
	db, err := client.ContextFunctions.CurrentDatabase(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, db)
}

func TestInt_CurrentSchema(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// new database and schema created on purpose
	databaseTest, databaseCleanup := createDatabase(t, client)
	t.Cleanup(databaseCleanup)
	schemaTest, schemaCleanup := createSchema(t, client, databaseTest)
	t.Cleanup(schemaCleanup)
	err := client.Sessions.UseSchema(ctx, schemaTest.ID())
	require.NoError(t, err)
	schema, err := client.ContextFunctions.CurrentSchema(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, schema)
}

func TestInt_CurrentWarehouse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	// new warehouse created on purpose
	warehouseTest, warehouseCleanup := createWarehouse(t, client)
	t.Cleanup(warehouseCleanup)
	err := client.Sessions.UseWarehouse(ctx, warehouseTest.ID())
	require.NoError(t, err)
	warehouse, err := client.ContextFunctions.CurrentWarehouse(ctx)
	require.NoError(t, err)
	assert.NotEmpty(t, warehouse)
}

func TestInt_IsRoleInSession(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)
	role, err := client.ContextFunctions.IsRoleInSession(ctx, sdk.NewAccountObjectIdentifier(currentRole))
	require.NoError(t, err)
	assert.True(t, role)
}

func TestInt_RolesUse(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	currentRoleID := sdk.NewAccountObjectIdentifier(currentRole)
	require.NoError(t, err)

	role, cleanup := createRole(t, client)
	t.Cleanup(cleanup)
	require.NotEqual(t, currentRole, role.Name)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{Role: &currentRoleID}))
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, role.ID())
	require.NoError(t, err)

	activeRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	assert.Equal(t, activeRole, role.Name)

	err = client.Sessions.UseRole(ctx, currentRoleID)
	require.NoError(t, err)
}

func TestInt_RolesUseSecondaryRoles(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)
	currentRole, err := client.ContextFunctions.CurrentRole(ctx)
	require.NoError(t, err)

	role, cleanup := createRole(t, client)
	t.Cleanup(cleanup)
	require.NotEqual(t, currentRole, role.Name)

	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	userID := sdk.NewAccountObjectIdentifier(user)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{User: &userID}))
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, role.ID())
	require.NoError(t, err)

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.SecondaryRolesAll)
	require.NoError(t, err)

	r, err := client.ContextFunctions.CurrentSecondaryRoles(ctx)
	require.NoError(t, err)

	names := make([]string, len(r.Roles))
	for i, v := range r.Roles {
		names[i] = v.Name()
	}
	assert.Equal(t, sdk.SecondaryRolesAll, r.Value)
	assert.Contains(t, names, currentRole)

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.SecondaryRolesNone)
	require.NoError(t, err)

	secondaryRolesAfter, err := client.ContextFunctions.CurrentSecondaryRoles(ctx)
	require.NoError(t, err)

	assert.Equal(t, sdk.SecondaryRolesNone, secondaryRolesAfter.Value)
	assert.Equal(t, 0, len(secondaryRolesAfter.Roles))

	t.Cleanup(func() {
		err = client.Sessions.UseRole(ctx, sdk.NewAccountObjectIdentifier(currentRole))
		require.NoError(t, err)
	})
}

func TestInt_RolesUseSecondaryRolesUsage(t *testing.T) {
	client := testClient(t)
	ctx := testContext(t)

	role, cleanupRole := createRole(t, client)
	t.Cleanup(cleanupRole)

	privilegedRole, cleanupPrivilegedRole := createRole(t, client)
	t.Cleanup(cleanupPrivilegedRole)

	user, err := client.ContextFunctions.CurrentUser(ctx)
	require.NoError(t, err)
	userID := sdk.NewAccountObjectIdentifier(user)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(role.ID(), sdk.GrantRole{User: &userID}))
	require.NoError(t, err)

	err = client.Roles.Grant(ctx, sdk.NewGrantRoleRequest(privilegedRole.ID(), sdk.GrantRole{User: &userID}))
	require.NoError(t, err)

	err = client.Grants.GrantPrivilegesToAccountRole(ctx, &sdk.AccountRoleGrantPrivileges{
		GlobalPrivileges: []sdk.GlobalPrivilege{
			sdk.GlobalPrivilegeCreateDatabase,
		},
	},
		&sdk.AccountRoleGrantOn{
			Account: sdk.Bool(true),
		},
		privilegedRole.ID(),
		new(sdk.GrantPrivilegesToAccountRoleOptions),
	)
	require.NoError(t, err)

	err = client.Sessions.UseRole(ctx, role.ID())
	require.NoError(t, err)

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.SecondaryRolesNone)
	require.NoError(t, err)

	err = client.Databases.Create(ctx, sdk.NewAccountObjectIdentifier(random.AlphaN(20)), new(sdk.CreateDatabaseOptions))
	require.ErrorContains(t, err, "Insufficient privileges to operate on account 'IYA62698'")

	err = client.Sessions.UseSecondaryRoles(ctx, sdk.SecondaryRolesAll)
	require.NoError(t, err)

	dbId := sdk.NewAccountObjectIdentifier(random.AlphaN(20))
	err = client.Databases.Create(ctx, dbId, new(sdk.CreateDatabaseOptions))
	require.NoError(t, err)

	//err = client.Sessions.UseRole(ctx, privilegedRole.ID())
	//require.NoError(t, err)
	//
	//err = client.Databases.Create(ctx, dbId, new(sdk.CreateDatabaseOptions))
	//require.NoError(t, err)
	//
	//// Only needed for cleanup
	//err = client.Sessions.UseRole(ctx, sdk.NewAccountObjectIdentifier("ACCOUNTADMIN"))
	//require.NoError(t, err)

	t.Cleanup(func() {
		err = client.Databases.Drop(ctx, dbId, new(sdk.DropDatabaseOptions))
		require.NoError(t, err)
	})
}
