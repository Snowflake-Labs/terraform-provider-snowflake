//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Unlike SHOW OPENFLOW DEPLOYMENTS, SHOW OPENFLOW RUNTIMES is schema-scoped and accepts an `in` filter, so
// the account, database and schema variants are covered alongside the other three filters.
func TestAcc_OpenflowRuntimes_BasicUseCase_DifferentFiltering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)
	deploymentId := testClient().OpenflowDeployment.ActiveDeploymentForRuntimes(t)

	id := testClient().Ids.RandomSchemaObjectIdentifier()
	comment := random.Comment()
	displayName := random.AlphaN(12)

	runtimeModel := model.OpenflowRuntime("t", id.DatabaseName(), id.SchemaName(), id.Name(),
		deploymentId.Name(), currentRole.Name(), 1, 1, string(sdk.OpenflowRuntimeNodeTypeSmall)).
		WithDisplayName(displayName).
		WithComment(comment)

	dsLikeExact := datasourcemodel.OpenflowRuntimes("test").
		WithLike(id.Name()).
		WithDependsOn(runtimeModel.ResourceReference())
	dsLikeMatchingNothing := datasourcemodel.OpenflowRuntimes("test").
		WithLike("non_existing_runtime").
		WithDependsOn(runtimeModel.ResourceReference())
	dsWithoutDescribe := datasourcemodel.OpenflowRuntimes("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(runtimeModel.ResourceReference())
	dsStartsWith := datasourcemodel.OpenflowRuntimes("test").
		WithStartsWith(id.Name()).
		WithDependsOn(runtimeModel.ResourceReference())
	dsInSchema := datasourcemodel.OpenflowRuntimes("test").
		WithInSchema(id.SchemaId()).
		WithLike(id.Name()).
		WithDependsOn(runtimeModel.ResourceReference())
	dsInDatabase := datasourcemodel.OpenflowRuntimes("test").
		WithInDatabase(id.DatabaseId()).
		WithLike(id.Name()).
		WithDependsOn(runtimeModel.ResourceReference())
	dsInAccount := datasourcemodel.OpenflowRuntimes("test").
		WithInAccount().
		WithLike(id.Name()).
		WithDependsOn(runtimeModel.ResourceReference())
	// The account holds other runtimes, so LIMIT 1 is asserted as returning exactly one row rather than as
	// returning this runtime.
	dsLimit := datasourcemodel.OpenflowRuntimes("test").
		WithLimit(1).
		WithDependsOn(runtimeModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// The runtime here is only a fixture for the data source to read; its own resource tests cover destroy.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// 1. like on an exact name, describe on by default - assert every exposed field
			{
				Config: config.FromModels(t, runtimeModel, dsLikeExact),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "openflow_runtimes.#", "1")),
					resourceshowoutputassert.OpenflowRuntimesDatasourceShowOutput(t, dsLikeExact.DatasourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.OpenflowRuntimeStatusActive).
						HasDeployment(deploymentId.Name()).
						HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
						HasMinNodes(1).
						HasMaxNodes(1).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasExecuteAsRole(currentRole.Name()).
						HasDisplayName(displayName).
						HasComment(comment).
						HasOwner(currentRole.Name()).
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty(),
					resourceshowoutputassert.OpenflowRuntimesDatasourceDescribeOutput(t, dsLikeExact.DatasourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.OpenflowRuntimeStatusActive).
						HasDeployment(deploymentId.Name()).
						HasNodeType(sdk.OpenflowRuntimeNodeTypeSmall).
						HasMinNodes(1).
						HasMaxNodes(1).
						HasExecuteAsRole(currentRole.Name()).
						HasDisplayName(displayName).
						HasComment(comment).
						HasOwner(currentRole.Name()).
						HasKeyNotEmpty(),
				),
			},
			// 2. like matching nothing
			{
				Config: config.FromModels(t, runtimeModel, dsLikeMatchingNothing),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeMatchingNothing.DatasourceReference(), "openflow_runtimes.#", "0")),
				),
			},
			// 3. describe off - show_output is still populated, describe_output is empty
			{
				Config: config.FromModels(t, runtimeModel, dsWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "openflow_runtimes.#", "1")),
					resourceshowoutputassert.OpenflowRuntimesDatasourceShowOutput(t, dsWithoutDescribe.DatasourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "openflow_runtimes.0.describe_output.#", "0")),
				),
			},
			// 4. starts_with
			{
				Config: config.FromModels(t, runtimeModel, dsStartsWith),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "openflow_runtimes.#", "1")),
					resourceshowoutputassert.OpenflowRuntimesDatasourceShowOutput(t, dsStartsWith.DatasourceReference()).
						HasName(id.Name()),
				),
			},
			// 5. in schema
			{
				Config: config.FromModels(t, runtimeModel, dsInSchema),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "openflow_runtimes.#", "1")),
					resourceshowoutputassert.OpenflowRuntimesDatasourceShowOutput(t, dsInSchema.DatasourceReference()).
						HasName(id.Name()).
						HasSchemaName(id.SchemaName()),
				),
			},
			// 6. in database
			{
				Config: config.FromModels(t, runtimeModel, dsInDatabase),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsInDatabase.DatasourceReference(), "openflow_runtimes.#", "1")),
					resourceshowoutputassert.OpenflowRuntimesDatasourceShowOutput(t, dsInDatabase.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()),
				),
			},
			// 7. in account
			{
				Config: config.FromModels(t, runtimeModel, dsInAccount),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsInAccount.DatasourceReference(), "openflow_runtimes.#", "1")),
					resourceshowoutputassert.OpenflowRuntimesDatasourceShowOutput(t, dsInAccount.DatasourceReference()).
						HasName(id.Name()),
				),
			},
			// 8. limit
			{
				Config: config.FromModels(t, runtimeModel, dsLimit),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "openflow_runtimes.#", "1")),
				),
			},
		},
	})
}
