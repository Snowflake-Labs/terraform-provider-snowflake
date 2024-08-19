package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRole(t *testing.T) {
	resourceName := "snowflake_database_role.test_db_role"
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()
	comment2 := random.Comment()

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.DatabaseRole),
		Steps: []resource.TestStep{
			{
				Config: databaseRoleConfig(id.Name(), acc.TestDatabaseName, comment),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment),
				),
			},
			{
				Config: databaseRoleConfig(id.Name(), acc.TestDatabaseName, comment2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", id.Name()),
					resource.TestCheckResourceAttr(resourceName, "fully_qualified_name", id.FullyQualifiedName()),
					resource.TestCheckResourceAttr(resourceName, "database", acc.TestDatabaseName),
					resource.TestCheckResourceAttr(resourceName, "comment", comment2),
				),
			},
		},
	})
}

func databaseRoleConfig(dbRoleName string, databaseName string, comment string) string {
	s := `
resource "snowflake_database_role" "test_db_role" {
	name     	  = "%s"
	database  	  = "%s"
	comment       = "%s"
}
	`
	return fmt.Sprintf(s, dbRoleName, databaseName, comment)
}

func TestAcc_DatabaseRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				Config: databaseRoleConfig(id.Name(), id.DatabaseName(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "id", fmt.Sprintf(`%s|%s`, id.DatabaseName(), id.Name())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseRoleConfig(id.Name(), id.DatabaseName(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_DatabaseRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	quotedDatabaseRoleId := fmt.Sprintf(`\"%s\"`, id.Name())
	comment := random.Comment()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"snowflake": {
						VersionConstraint: "=0.94.1",
						Source:            "Snowflake-Labs/snowflake",
					},
				},
				ExpectNonEmptyPlan: true,
				Config:             databaseRoleConfig(quotedDatabaseRoleId, id.DatabaseName(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "id", fmt.Sprintf(`%s|%s`, id.DatabaseName(), id.Name())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   databaseRoleConfig(quotedDatabaseRoleId, id.DatabaseName(), comment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database_role.test_db_role", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}
