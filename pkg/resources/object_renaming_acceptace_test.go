package resources_test

import (
	"fmt"
	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"regexp"
	"testing"
)

// The following tests are showing the behavior of the provider in cases where objects higher in the hierarchy
// like database or schema are renamed when the objects lower in the hierarchy are in the Terraform configuration.
// Learn about it in TODO(SNOW-1672319): link public document.
//
// Shallow hierarchy (database + schema)
// - is in config - renamed internally - with dependency
// - is in config - renamed internally - without dependency
// TODO: More cases for is in config
//
// - is in config - renamed externally - with dependency - after rename referencing old database name
// - is in config - renamed externally - with dependency - after rename referencing new database name
// - is in config - renamed externally - without dependency - after rename referencing old database name
// - is in config - renamed externally - without dependency - after rename referencing new database name
// - is not in config - renamed externally - referencing old database name
// - is not in config - renamed externally - referencing new database name
//
// Deep hierarchy (database + schema + schema object)
// - only database is in config - renamed internally
// - only database is in config - renamed externally
// - only schema is in config - renamed internally
// - only schema is in config - renamed externally
// - both database and schema are in config - renamed internally
// - both database and schema are in config - renamed externally
// - both database and schema are not in config - renamed internally
// - both database and schema are not in config - renamed externally

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithDependency(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModelWithNewId) + configSchemaWithDatabaseReference(databaseConfigModelWithNewId.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithoutDependency(t *testing.T) {
	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				// TODO(): The error may be non-deterministic because of the Terraform resource ordering (
				ExpectError: regexp.MustCompile("does not exist or not authorized"), // (Error deleting schema)
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_AfterRenameReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
			{
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_AfterRenameReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
			{
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsNotInConfig_RenamedExternally_ReferencingOldName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					_, databaseCleanup := acc.TestClient().Database.CreateDatabaseWithIdentifier(t, databaseId)
					t.Cleanup(databaseCleanup)
				},
				Config: config.FromModel(t, schemaModelConfig),
			},
			{
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModel(t, schemaModelConfig),
				ExpectError: regexp.MustCompile("object does not exist or not authorized"),
			},
		},
	})
}

func configSchemaWithDatabaseReference(databaseReference string, schemaName string) string {
	return fmt.Sprintf(`
resource "snowflake_schema" "test" {
	database = %[1]s.name
	name = "%[2]s"
}
`, databaseReference, schemaName)
}
