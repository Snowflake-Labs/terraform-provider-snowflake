package resources_test

import (
	"fmt"
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/objectassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/helpers"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_DatabaseRole(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()
	databaseRoleModel := model.DatabaseRole("test", id.DatabaseName(), id.Name())
	databaseRoleModelWithComment := model.DatabaseRole("test", id.DatabaseName(), id.Name()).WithComment(comment)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.DatabaseRole),
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, databaseRoleModel),
				Check: assert.AssertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(id.Name()).
						HasComment(""),
					objectassert.DatabaseRole(t, id).
						HasName(id.Name()).
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_database_role.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedDatabaseRoleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(""),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(""),
				),
			},
			// set comment
			{
				Config: config.FromModel(t, databaseRoleModelWithComment),
				Check: assert.AssertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasCommentString(comment).
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(id.Name()).
						HasComment(comment),
					objectassert.DatabaseRole(t, id).
						HasName(id.Name()).
						HasComment(comment),
				),
			},
			{
				ResourceName: "snowflake_database_role.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedDatabaseRoleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(comment),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(comment),
				),
			},
			// unset comment
			{
				Config: config.FromModel(t, databaseRoleModel),
				Check: assert.AssertThat(t,
					resourceassert.DatabaseRoleResource(t, "snowflake_database_role.test").
						HasNameString(id.Name()).
						HasDatabaseString(id.DatabaseName()).
						HasCommentString("").
						HasFullyQualifiedNameString(id.FullyQualifiedName()),
					resourceshowoutputassert.DatabaseRoleShowOutput(t, "snowflake_database_role.test").
						HasName(id.Name()).
						HasComment(""),
					objectassert.DatabaseRole(t, id).
						HasName(id.Name()).
						HasComment(""),
				),
			},
			{
				ResourceName: "snowflake_database_role.test",
				ImportState:  true,
				ImportStateCheck: assert.AssertThatImport(t,
					resourceassert.ImportedDatabaseRoleResource(t, helpers.EncodeResourceIdentifier(id)).
						HasNameString(id.Name()).
						HasCommentString(""),
					resourceshowoutputassert.ImportedWarehouseShowOutput(t, helpers.EncodeResourceIdentifier(id)).
						HasName(id.Name()).
						HasComment(""),
				),
			},
		},
	})
}

func TestAcc_DatabaseRole_migrateFromV0941_ensureSmoothUpgradeWithNewResourceId(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	comment := random.Comment()
	databaseRoleModelWithComment := model.DatabaseRole("test", id.DatabaseName(), id.Name()).WithComment(comment)

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
				Config: config.FromModel(t, databaseRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", fmt.Sprintf(`%s|%s`, id.DatabaseName(), id.Name())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModel(t, databaseRoleModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}

func TestAcc_DatabaseRole_IdentifierQuotingDiffSuppression(t *testing.T) {
	id := acc.TestClient().Ids.RandomDatabaseObjectIdentifier()
	quotedDatabaseRoleId := fmt.Sprintf(`"%s"`, id.Name())
	comment := random.Comment()
	databaseRoleModelWithComment := model.DatabaseRole("test", id.DatabaseName(), quotedDatabaseRoleId).WithComment(comment)

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
				Config:             config.FromModel(t, databaseRoleModelWithComment),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", fmt.Sprintf(`%s|%s`, id.DatabaseName(), id.Name())),
				),
			},
			{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				Config:                   config.FromModel(t, databaseRoleModelWithComment),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_database_role.test", plancheck.ResourceActionNoop),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_database_role.test", "database", id.DatabaseName()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "name", id.Name()),
					resource.TestCheckResourceAttr("snowflake_database_role.test", "id", id.FullyQualifiedName()),
				),
			},
		},
	})
}
