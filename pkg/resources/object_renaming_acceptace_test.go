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

/*
The following tests are showing the behavior of the provider in cases where objects higher in the hierarchy
like database or schema are renamed when the objects lower in the hierarchy are in the Terraform configuration.
Learn about it in TODO(SNOW-1672319): link public document.

Shallow hierarchy (database + schema)
- is in config - renamed internally - with implicit dependency

- is in config - renamed internally - without dependency - after rename schema referencing old database name
- is in config - renamed internally - without dependency - after rename schema referencing new database name

- is in config - renamed internally - with depends_on - after rename schema referencing old database name
- is in config - renamed internally - with depends_on - after rename schema referencing new database name

- is in config - renamed externally - with implicit dependency - database holding the same name in config
- is in config - renamed externally - with implicit dependency - database holding the new name in config

- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing old database name
- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing old database name - check impact of the resource order on the plan
- is in config - renamed externally - without dependency - after rename database referencing old name and schema referencing new database name
- is in config - renamed externally - without dependency - after rename database referencing new name and schema referencing old database name
- is in config - renamed externally - without dependency - after rename database referencing new name and schema referencing new database name

- is in config - renamed externally - with depends_on - after rename database referencing old name and schema referencing old database name
- is in config - renamed externally - with depends_on - after rename database referencing old name and schema referencing new database name
- is in config - renamed externally - with depends_on - after rename database referencing new name and schema referencing old database name
- is in config - renamed externally - with depends_on - after rename database referencing new name and schema referencing new database name

- is not in config - renamed externally - referencing old database name
- is not in config - renamed externally - referencing new database name

Deep hierarchy (database + schema + schema object)
- only database is in config - renamed internally
- only database is in config - renamed externally
- only schema is in config - renamed internally
- only schema is in config - renamed externally
- both database and schema are in config - renamed internally
- both database and schema are in config - renamed externally
- both database and schema are not in config - renamed internally
- both database and schema are not in config - renamed externally

// TODO: Add Ticket number to TODOs next to Skips
*/

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithImplicitDependency(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModelWithNewId) + configSchemaWithDatabaseReference(databaseConfigModelWithNewId.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithoutDependency_AfterRenameSchemaReferencingOldDatabaseName(t *testing.T) {
	// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
	t.Skip("Not able to handle the error produced by Delete operation that results in test always failing")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithoutDependency_AfterRenameSchemaReferencingNewDatabaseName(t *testing.T) {
	// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
	t.Skip("Not able to handle the error produced by Delete operation that results in test always failing")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithDependsOn_AfterRenameSchemaReferencingOldDatabaseName(t *testing.T) {
	// Error happens during schema's Read operation and then Delete operation (schema cannot be removed).
	t.Skip("Not able to handle the error produced by Delete operation that results in test always failing")
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: resource.ComposeAggregateTestCheckFunc(),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
			{
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionNoop),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedInternally_WithDependsOn_AfterRenameSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())

	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId.SetDependsOn(databaseConfigModelWithNewId.ResourceReference())

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
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionUpdate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionDestroyBeforeCreate),
					},
				},
				Config: config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithImplicitDependency_DatabaseHoldingTheOldNameInConfig(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
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
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
						NewName: &newDatabaseId,
					})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModel(t, databaseConfigModel) + configSchemaWithDatabaseReference(databaseConfigModel.ResourceReference(), schemaName),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithImplicitDependency_DatabaseHoldingTheNewNameInConfig(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{
						NewName: &newDatabaseId,
					})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// TODO: Why Create? This case could be handled in a way that it shouldn't require any action (implicit import)
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModel(t, databaseConfigModelWithNewId) + configSchemaWithDatabaseReference(databaseConfigModelWithNewId.ResourceReference(), schemaName),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	t.Skip("Test results are inconsistent because Terraform execution order is non-deterministic")
	// Although the above applies, it seems to be consistently failing on delete operation after the test (because the database is dropped before schema).
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

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
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				// TODO: This test may have inconsistent result (depending on Terraform execution order which is not deterministic).
				Config:      config.FromModels(t, databaseConfigModel, schemaModelConfig),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

// This test checks if the order of the configuration resources has any impact on the order of resource execution (it seems to have no impact).
func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName_ConfigOrderSwap(t *testing.T) {
	t.Skip("Test results are inconsistent because Terraform execution order is non-deterministic")
	// Although the above applies, it seems to be consistently failing on delete operation after the test (because the database is dropped before schema).
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.Database),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, schemaModelConfig, databaseConfigModel),
			},
			{
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				// TODO: This test may have inconsistent result (depending on Terraform execution order which is not deterministic).
				Config:      config.FromModels(t, schemaModelConfig, databaseConfigModel),
				ExpectError: regexp.MustCompile("does not exist or not authorized"),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModel, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("Failed to create schema"), // already exists
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)

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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithoutDependency_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())

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
				PreConfig: func() {
					newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config: config.FromModels(t, databaseConfigModel, schemaModelConfig),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingOldNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId.SetDependsOn(databaseConfigModel.ResourceReference())

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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModel, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("Failed to create schema"), // already exists
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingOldDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())

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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfig),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
			},
		},
	})
}

func TestAcc_ShallowHierarchy_IsInConfig_RenamedExternally_WithDependsOn_AfterRenameDatabaseReferencingNewNameAndSchemaReferencingNewDatabaseName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()

	databaseConfigModel := model.Database("test", databaseId.Name())
	databaseConfigModelWithNewId := model.Database("test", newDatabaseId.Name())
	schemaModelConfig := model.Schema("test", databaseId.Name(), schemaName)
	schemaModelConfig.SetDependsOn(databaseConfigModel.ResourceReference())
	schemaModelConfigWithNewDatabaseId := model.Schema("test", newDatabaseId.Name(), schemaName)
	schemaModelConfigWithNewDatabaseId.SetDependsOn(databaseConfigModelWithNewId.ResourceReference())

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
				PreConfig: func() {
					acc.TestClient().Database.Alter(t, databaseId, &sdk.AlterDatabaseOptions{NewName: &newDatabaseId})
					t.Cleanup(acc.TestClient().Database.DropDatabaseFunc(t, newDatabaseId))
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database.test", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("snowflake_schema.test", plancheck.ResourceActionCreate),
					},
				},
				Config:      config.FromModels(t, databaseConfigModelWithNewId, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile(fmt.Sprintf(`Object '%s' already exists`, newDatabaseId.Name())),
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

func TestAcc_ShallowHierarchy_IsNotInConfig_RenamedExternally_ReferencingNewName(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	databaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	newDatabaseId := acc.TestClient().Ids.RandomAccountObjectIdentifier()
	schemaName := acc.TestClient().Ids.Alpha()
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
				Config:      config.FromModel(t, schemaModelConfigWithNewDatabaseId),
				ExpectError: regexp.MustCompile("Failed to create schema"), // already exists
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
