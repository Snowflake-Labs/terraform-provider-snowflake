package helpers

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
)

func (c *TestClient) SetUpLegacyServiceUserWithAccessToTestDatabaseAndWarehouse(t *testing.T, pass string) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	tmpUserId := c.Ids.RandomAccountObjectIdentifier()
	_, userCleanup := c.User.CreateUserWithOptions(t, tmpUserId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Type:     sdk.Pointer(sdk.UserTypeLegacyService),
		Password: sdk.String(pass),
	}})
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := c.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	c.Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, c.Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, c.Ids.SnowflakeWarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return tmpUserId, tmpRoleId
}

// TODO [this PR]: merge with the above function
func (c *TestClient) SetUpServiceUserWithAccessToTestDatabaseAndWarehouse(t *testing.T, publicKey string) (sdk.AccountObjectIdentifier, sdk.AccountObjectIdentifier) {
	tmpUserId := c.Ids.RandomAccountObjectIdentifier()
	_, userCleanup := c.User.CreateUserWithOptions(t, tmpUserId, &sdk.CreateUserOptions{ObjectProperties: &sdk.UserObjectProperties{
		Type:         sdk.Pointer(sdk.UserTypeService),
		RSAPublicKey: sdk.String(publicKey),
	}})
	t.Cleanup(userCleanup)

	tmpRole, roleCleanup := c.Role.CreateRole(t)
	t.Cleanup(roleCleanup)

	tmpRoleId := tmpRole.ID()

	c.Grant.GrantPrivilegesOnDatabaseToAccountRole(t, tmpRoleId, c.Ids.DatabaseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Grant.GrantPrivilegesOnWarehouseToAccountRole(t, tmpRoleId, c.Ids.SnowflakeWarehouseId(), []sdk.AccountObjectPrivilege{sdk.AccountObjectPrivilegeUsage}, false)
	c.Role.GrantRoleToUser(t, tmpRoleId, tmpUserId)

	return tmpUserId, tmpRoleId
}
