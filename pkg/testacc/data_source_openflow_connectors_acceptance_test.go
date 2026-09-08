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

// Connectors are schema-level, so SHOW accepts an `in` filter; the account, database and schema variants are
// covered alongside the other three filters.
func TestAcc_OpenflowConnectors_BasicUseCase_DifferentFiltering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	currentRole := testClient().Context.CurrentRole(t)
	runtimeId := testClient().OpenflowRuntime.ActiveRuntime(t)
	definition := testClient().OpenflowConnectorDefinition.ForTesting(t)

	id := testClient().Ids.RandomSchemaObjectIdentifierInSchema(runtimeId.SchemaId())
	comment := random.Comment()
	displayName := random.AlphaN(12)

	connectorModel := model.OpenflowConnector("t", id.DatabaseName(), id.SchemaName(), id.Name(), definition.Name, runtimeId.FullyQualifiedName()).
		WithDisplayName(displayName).
		WithComment(comment)

	dsLikeExact := datasourcemodel.OpenflowConnectors("test").
		WithLike(id.Name()).
		WithDependsOn(connectorModel.ResourceReference())
	dsLikeMatchingNothing := datasourcemodel.OpenflowConnectors("test").
		WithLike("non_existing_connector").
		WithDependsOn(connectorModel.ResourceReference())
	dsWithoutDescribe := datasourcemodel.OpenflowConnectors("test").
		WithLike(id.Name()).
		WithWithDescribe(false).
		WithDependsOn(connectorModel.ResourceReference())
	dsStartsWith := datasourcemodel.OpenflowConnectors("test").
		WithStartsWith(id.Name()).
		WithDependsOn(connectorModel.ResourceReference())
	dsInSchema := datasourcemodel.OpenflowConnectors("test").
		WithInSchema(id.SchemaId()).
		WithLike(id.Name()).
		WithDependsOn(connectorModel.ResourceReference())
	dsInDatabase := datasourcemodel.OpenflowConnectors("test").
		WithInDatabase(id.DatabaseId()).
		WithLike(id.Name()).
		WithDependsOn(connectorModel.ResourceReference())
	dsInAccount := datasourcemodel.OpenflowConnectors("test").
		WithInAccount().
		WithLike(id.Name()).
		WithDependsOn(connectorModel.ResourceReference())
	// The account holds other connectors, so LIMIT 1 is asserted as returning exactly one row rather than as
	// returning this connector.
	dsLimit := datasourcemodel.OpenflowConnectors("test").
		WithLimit(1).
		WithDependsOn(connectorModel.ResourceReference())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// The connector here is only a fixture for the data source to read; its own resource tests cover destroy.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// 1. like on an exact name, describe on by default - assert every exposed field
			{
				Config: config.FromModels(t, connectorModel, dsLikeExact),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "openflow_connectors.#", "1")),
					resourceshowoutputassert.OpenflowConnectorsDatasourceShowOutput(t, dsLikeExact.DatasourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.OpenflowConnectorStatusStopped).
						HasRuntime(runtimeId.Name()).
						HasConnectorDefinition(definition.Name).
						HasDatabaseName(id.DatabaseName()).
						HasSchemaName(id.SchemaName()).
						HasDisplayName(displayName).
						HasComment(comment).
						HasOwner(currentRole.Name()).
						HasDefaultVersion("LAST").
						HasDefaultVersionName("").
						HasDefaultVersionAlias("").
						HasDefaultVersionLocationUri("").
						HasDefaultVersionSourceLocationUri("").
						HasCreatedOnNotEmpty().
						HasUpdatedOnNotEmpty().
						HasConnectorUrlNotEmpty().
						HasLiveVersionLocationUriNotEmpty(),
					resourceshowoutputassert.OpenflowConnectorsDatasourceDescribeOutput(t, dsLikeExact.DatasourceReference()).
						HasName(id.Name()).
						HasStatus(sdk.OpenflowConnectorStatusStopped).
						HasRuntime(runtimeId.Name()).
						HasConnectorDefinition(definition.Name).
						HasDisplayName(displayName).
						HasComment(comment).
						HasOwner(currentRole.Name()).
						HasDefaultVersion("LAST").
						HasDefaultVersionName("").
						HasDefaultVersionAlias("").
						HasDefaultVersionLocationUri("").
						HasDefaultVersionSourceLocationUri("").
						HasDefaultVersionGitCommitHash("").
						HasLastVersionName("").
						HasLastVersionAlias("").
						HasLastVersionLocationUri("").
						HasLastVersionSourceLocationUri("").
						HasLastVersionGitCommitHash("").
						HasConnectorUrlNotEmpty().
						HasLiveVersionLocationUriNotEmpty(),
				),
			},
			// 2. like matching nothing
			{
				Config: config.FromModels(t, connectorModel, dsLikeMatchingNothing),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeMatchingNothing.DatasourceReference(), "openflow_connectors.#", "0")),
				),
			},
			// 3. with_describe = false suppresses the DESCRIBE call, so the block must be empty
			{
				Config: config.FromModels(t, connectorModel, dsWithoutDescribe),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "openflow_connectors.#", "1")),
					resourceshowoutputassert.OpenflowConnectorsDatasourceShowOutput(t, dsWithoutDescribe.DatasourceReference()).
						HasName(id.Name()),
					assert.Check(resource.TestCheckResourceAttr(dsWithoutDescribe.DatasourceReference(), "openflow_connectors.0.describe_output.#", "0")),
				),
			},
			// 4. starts_with
			{
				Config: config.FromModels(t, connectorModel, dsStartsWith),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsStartsWith.DatasourceReference(), "openflow_connectors.#", "1")),
					resourceshowoutputassert.OpenflowConnectorsDatasourceShowOutput(t, dsStartsWith.DatasourceReference()).
						HasName(id.Name()),
				),
			},
			// 5. in schema
			{
				Config: config.FromModels(t, connectorModel, dsInSchema),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsInSchema.DatasourceReference(), "openflow_connectors.#", "1")),
					resourceshowoutputassert.OpenflowConnectorsDatasourceShowOutput(t, dsInSchema.DatasourceReference()).
						HasName(id.Name()).
						HasSchemaName(id.SchemaName()),
				),
			},
			// 6. in database
			{
				Config: config.FromModels(t, connectorModel, dsInDatabase),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsInDatabase.DatasourceReference(), "openflow_connectors.#", "1")),
					resourceshowoutputassert.OpenflowConnectorsDatasourceShowOutput(t, dsInDatabase.DatasourceReference()).
						HasName(id.Name()).
						HasDatabaseName(id.DatabaseName()),
				),
			},
			// 7. in account
			{
				Config: config.FromModels(t, connectorModel, dsInAccount),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsInAccount.DatasourceReference(), "openflow_connectors.#", "1")),
					resourceshowoutputassert.OpenflowConnectorsDatasourceShowOutput(t, dsInAccount.DatasourceReference()).
						HasName(id.Name()),
				),
			},
			// 8. limit
			{
				Config: config.FromModels(t, connectorModel, dsLimit),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "openflow_connectors.#", "1")),
				),
			},
		},
	})
}
