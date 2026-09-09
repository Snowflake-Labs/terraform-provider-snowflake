//go:build account_level_tests

package testacc

import (
	"fmt"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// TestAcc_GrantOwnership_OnDatabaseRole proves that SHOW GRANTS ON DATABASE ROLE reports
// granted_on = ROLE when the 2026_06 bundle (BCR-2371) is disabled and DATABASE_ROLE when it is
// enabled. Toggling a behavior change bundle is account-wide, so this runs on the secondary account.
// TODO(SNOW-4075337): skip or remove the disabled-bundle step once the 2026_06 bundle is enforced
// and can no longer be disabled.
func TestAcc_GrantOwnership_OnDatabaseRole(t *testing.T) {
	tc := secondaryTestClient()
	database, databaseCleanup := tc.Database.CreateDatabaseWithParametersSet(t)
	t.Cleanup(databaseCleanup)

	databaseId := database.ID()
	databaseRoleId := tc.Ids.RandomDatabaseObjectIdentifierInDatabase(databaseId)
	databaseRoleFullyQualifiedName := databaseRoleId.FullyQualifiedName()
	accountRoleId := tc.Ids.RandomAccountObjectIdentifier()
	accountRoleFullyQualifiedName := accountRoleId.FullyQualifiedName()

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
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
		ProtoV6ProviderFactories: secondaryAccountProviderFactory,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					tc.BcrBundles.DisableBcrBundle(t, "2026_06")
				},
				Config: accconfig.FromModels(t, providerModel, accountRoleModel, dbRoleModel, grantModel),
				Check: resource.ComposeTestCheckFunc(
					checkResourceAttrs,
					checkDatabaseRoleOwnershipGrantedOn(t, databaseRoleId, sdk.ObjectTypeRole, accountRoleId.Name()),
				),
			},
			{
				PreConfig: func() {
					tc.BcrBundles.EnableBcrBundle(t, "2026_06")
				},
				Config: accconfig.FromModels(t, providerModel, accountRoleModel, dbRoleModel, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
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
		grants, err := secondaryTestClient().Grant.ShowGrantsOnObject(t, sdk.ObjectTypeDatabaseRole, databaseRoleId)
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
