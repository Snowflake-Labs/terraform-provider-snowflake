//go:build non_account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_GrantOwnership_OnDatabaseRole proves that SHOW GRANTS ON DATABASE ROLE reports
// granted_on = DATABASE_ROLE.
func TestAcc_GrantOwnership_OnDatabaseRole(t *testing.T) {
	database, databaseCleanup := testClient().Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseRoleId := testClient().Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()
	accountRoleId := testClient().Ids.RandomAccountObjectIdentifier()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	accountRoleModel := model.AccountRole("test", accountRoleId.Name())
	dbRoleModel := model.DatabaseRole("test", databaseId.Name(), databaseRoleId.Name())
	grantModel := model.GrantOwnershipWithRawOn("test").
		WithAccountRoleName(accountRoleId.Name()).
		WithOnObject(sdk.ObjectTypeDatabaseRole, databaseRoleFullyQualifiedName).
		WithDependsOn(accountRoleModel.ResourceReference(), dbRoleModel.ResourceReference())

	resourceName := "snowflake_grant_ownership.test"
	checkResourceAttrs := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(resourceName, "account_role_name", accountRoleId.Name()),
		resource.TestCheckResourceAttr(resourceName, "on.0.object_type", string(sdk.ObjectTypeDatabaseRole)),
		resource.TestCheckResourceAttr(resourceName, "on.0.object_name", databaseRoleFullyQualifiedName),
		resource.TestCheckResourceAttr(resourceName, "id", fmt.Sprintf("ToAccountRole|%s||OnObject|DATABASE ROLE|%s", accountRoleFullyQualifiedName, databaseRoleFullyQualifiedName)),
	)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: accconfig.FromModels(t, accountRoleModel, dbRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					checkResourceAttrs,
					checkDatabaseRoleOwnershipGrantedOn(t, databaseRoleId, sdk.ObjectTypeDatabaseRole, accountRoleId.Name()),
				),
			},
		},
	})
}

func checkDatabaseRoleOwnershipGrantedOn(t *testing.T, databaseRoleId sdk.DatabaseObjectIdentifier, grantOn sdk.ObjectType, roleName string) resource.TestCheckFunc {
	t.Helper()
	objectName := databaseRoleId.FullyQualifiedName()
	return func(_ *terraform.State) error {
		grants, err := testClient().Grant.ShowGrantsOnObject(t, sdk.ObjectTypeDatabaseRole, databaseRoleId)
		if err != nil {
			return err
		}
		var found bool
		for _, grant := range grants {
			if grant.Privilege == "OWNERSHIP" &&
				(grant.GrantedOn == grantOn || grant.GrantOn == grantOn) &&
				grant.GranteeName.Name() == roleName && grant.Name.FullyQualifiedName() == objectName {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unable to find ownership privilege on %s granted to %s, expected name: %s", grantOn, roleName, objectName)
		}
		return nil
	}
}
