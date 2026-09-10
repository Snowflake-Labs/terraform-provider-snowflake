//go:build account_level_tests

package testacc

import (
	"regexp"
	"strings"
	"testing"

	accconfig "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/providermodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testdatatypes"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testprofiles"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// The two tests below exercise the released provider (v2.20.0) against an Iceberg table granted through
// snowflake_grant_privileges_to_database_role, with the 2026_06 bundle disabled and enabled.
// Background: https://github.com/snowflakedb/terraform-provider-snowflake/issues/5200.

func setUpIcebergTableGrantTest(t *testing.T) (icebergTableId sdk.SchemaObjectIdentifier, databaseRoleId sdk.DatabaseObjectIdentifier) {
	t.Helper()

	awsBaseUrl := testenvs.GetOrSkipTest(t, testenvs.AwsExternalBucketUrl)
	awsKeyId := testenvs.GetOrSkipTest(t, testenvs.AwsExternalKeyId)
	awsSecretKey := testenvs.GetOrSkipTest(t, testenvs.AwsExternalSecretKey)

	s3CompatBaseUrl := strings.Replace(awsBaseUrl, "s3://", "s3compat://", 1)
	s3CompatEndpoint := "s3.us-west-2.amazonaws.com"

	externalVolumeId, externalVolumeCleanup := secondaryTestClient().ExternalVolume.CreateS3Compat(t, s3CompatBaseUrl, s3CompatEndpoint, awsKeyId, awsSecretKey)
	t.Cleanup(externalVolumeCleanup)

	db, dbCleanup := secondaryTestClient().Database.CreateDatabaseWithRequest(t,
		secondaryTestClient().Database.TestParametersSet(secondaryTestClient().Ids.RandomAccountObjectIdentifier()).WithExternalVolume(externalVolumeId))
	t.Cleanup(dbCleanup)
	schemaId := sdk.NewDatabaseObjectIdentifier(db.ID().Name(), "PUBLIC")

	databaseRole, databaseRoleCleanup := secondaryTestClient().DatabaseRole.CreateDatabaseRoleInDatabase(t, db.ID())
	t.Cleanup(databaseRoleCleanup)

	icebergTableId = secondaryTestClient().Ids.RandomSchemaObjectIdentifierInSchema(schemaId)
	_, icebergTableCleanup := secondaryTestClient().IcebergTable.CreateWithRequest(t,
		sdk.NewCreateIcebergTableRequest(icebergTableId, sdk.IcebergTableColumnsAndConstraintsRequest{
			Columns: []sdk.IcebergTableColumnRequest{{Name: "ID", ColumnType: testdatatypes.DataTypeNumber}},
		}).
			WithCatalog(sdk.IcebergTableCatalogSnowflake).
			WithExternalVolume(externalVolumeId).
			WithBaseLocation("bcr_2371_grants_acceptance_test"))
	t.Cleanup(icebergTableCleanup)

	return icebergTableId, databaseRole.ID()
}

// TestAcc_GrantPrivilegesToDatabaseRole_bcr2026_06_IcebergTable_BundleDisabled_SpecificObjectType proves
// that with the bundle disabled the canonical-looking configuration is the broken one: granted_on is TABLE,
// so object_type = "ICEBERG TABLE" never matches on Read and the resource is stuck in a permadiff.
// Reference: https://github.com/snowflakedb/terraform-provider-snowflake/issues/5200
func TestAcc_GrantPrivilegesToDatabaseRole_bcr2026_06_IcebergTable_BundleDisabled_SpecificObjectType(t *testing.T) {
	icebergTableId, databaseRoleId := setUpIcebergTableGrantTest(t)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRoleId.FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeSelect).
		WithOnSchemaObjectObject(sdk.ObjectTypeIcebergTable, icebergTableId.FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.20.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.DisableBcrBundle(t, "2026_06")
				},
				Config: accconfig.FromModels(t, providerModel, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					// The permadiff: Read matched nothing, so the refresh drops the privileges and the next plan grants them again, forever.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(grantModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(grantModel.ResourceReference(), "privileges", sdk.String("[SELECT]"), sdk.String("[]")),
						planchecks.ExpectChange(grantModel.ResourceReference(), "privileges", tfjson.ActionUpdate, sdk.String("[]"), sdk.String("[SELECT]")),
						// Only privileges moves: the object type and the id are stable, so this is a perpetual re-grant rather than a replacement.
						planchecks.ExpectNoChangeOnField(grantModel.ResourceReference(), "on_schema_object.0.object_type"),
						planchecks.ExpectNoChangeOnField(grantModel.ResourceReference(), "id"),
					},
				},
				// Required, otherwise the framework fails the step on the same non-empty plan.
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(grantModel.ResourceReference(), "on_schema_object.0.object_type", string(sdk.ObjectTypeIcebergTable)),
					// Populated right after apply - Terraform keeps the planned value - and emptied by the refresh that follows.
					resource.TestCheckResourceAttr(grantModel.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(grantModel.ResourceReference(), "privileges.0", string(sdk.SchemaObjectPrivilegeSelect)),
				),
			},
		},
	})
}

// TestAcc_GrantPrivilegesToDatabaseRole_bcr2026_06_IcebergTable_BundleEnabled_GenericObjectType walks the
// whole reported scenario with the bundle enabled - the state that becomes permanent once the bundle can no
// longer be disabled.
// Reference: https://github.com/snowflakedb/terraform-provider-snowflake/issues/5200
func TestAcc_GrantPrivilegesToDatabaseRole_bcr2026_06_IcebergTable_BundleEnabled_GenericObjectType(t *testing.T) {
	icebergTableId, databaseRoleId := setUpIcebergTableGrantTest(t)

	providerModel := providermodel.SnowflakeProvider().WithProfile(testprofiles.Secondary)
	grantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRoleId.FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeSelect).
		WithOnSchemaObjectObject(sdk.ObjectTypeTable, icebergTableId.FullyQualifiedName())
	correctedGrantModel := model.GrantPrivilegesToDatabaseRole("test", databaseRoleId.FullyQualifiedName()).
		WithSchemaObjectPrivileges(sdk.SchemaObjectPrivilegeSelect).
		WithOnSchemaObjectObject(sdk.ObjectTypeIcebergTable, icebergTableId.FullyQualifiedName())

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.20.0"),
				PreConfig: func() {
					secondaryTestClient().BcrBundles.EnableBcrBundle(t, "2026_06")
				},
				Config: accconfig.FromModels(t, providerModel, grantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(grantModel.ResourceReference(), plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(grantModel.ResourceReference(), "privileges", sdk.String("[SELECT]"), sdk.String("[]")),
						planchecks.ExpectChange(grantModel.ResourceReference(), "privileges", tfjson.ActionUpdate, sdk.String("[]"), sdk.String("[SELECT]")),
						planchecks.ExpectNoChangeOnField(grantModel.ResourceReference(), "on_schema_object.0.object_type"),
						planchecks.ExpectNoChangeOnField(grantModel.ResourceReference(), "id"),
					},
				},
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(grantModel.ResourceReference(), "on_schema_object.0.object_type", string(sdk.ObjectTypeTable)),
					resource.TestCheckResourceAttr(grantModel.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(grantModel.ResourceReference(), "privileges.0", string(sdk.SchemaObjectPrivilegeSelect)),
				),
			},
			// Correcting object_type on the released provider: the ForceNew destroy cannot build its revoke
			// from the emptied state, so the replacement fails before any SQL is sent.
			{
				ExternalProviders: ExternalProviderWithExactVersion("2.20.0"),
				Config:            accconfig.FromModels(t, providerModel, correctedGrantModel),
				ExpectError:       regexp.MustCompile(`exactly one of\s+DatabaseRoleGrantPrivileges`),
			},
			// The same configuration on the current provider: the revoke is built from the resource id, the
			// replacement succeeds, and the corrected object type matches what Snowflake reports.
			{
				ProtoV6ProviderFactories: secondaryAccountProviderFactory,
				Config:                   accconfig.FromModels(t, providerModel, correctedGrantModel),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PostApplyPostRefresh: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(correctedGrantModel.ResourceReference(), "on_schema_object.0.object_type", string(sdk.ObjectTypeIcebergTable)),
					resource.TestCheckResourceAttr(correctedGrantModel.ResourceReference(), "privileges.#", "1"),
					resource.TestCheckResourceAttr(correctedGrantModel.ResourceReference(), "privileges.0", string(sdk.SchemaObjectPrivilegeSelect)),
				),
			},
		},
	})
}
