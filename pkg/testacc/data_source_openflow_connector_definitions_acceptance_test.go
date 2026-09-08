//go:build non_account_level_tests

package testacc

import (
	"testing"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceshowoutputassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/datasourcemodel"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Connector definitions are Snowflake-managed and read-only, so this needs no fixture of its own and there is
// no describe output to cover. SHOW is account-scoped, so there is no `in` filter either.
func TestAcc_OpenflowConnectorDefinitions_BasicUseCase_DifferentFiltering(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.TestOpenflow)

	definition := testClient().OpenflowConnectorDefinition.ForTesting(t)

	dsLikeExact := datasourcemodel.OpenflowConnectorDefinitions("test").WithLike(definition.Name)
	dsLikeMatchingNothing := datasourcemodel.OpenflowConnectorDefinitions("test").WithLike("non_existing_definition")
	dsLikePattern := datasourcemodel.OpenflowConnectorDefinitions("test").WithLike(definition.Name[:len(definition.Name)-1] + "%")
	// The account holds other definitions, so LIMIT 1 is asserted as returning exactly one row rather than as
	// returning this definition.
	dsLimit := datasourcemodel.OpenflowConnectorDefinitions("test").WithLimit(1)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		// Nothing is created, so there is nothing to destroy.
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// 1. like on an exact name - assert every exposed field against what SHOW reports
			{
				Config: config.FromModels(t, dsLikeExact),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeExact.DatasourceReference(), "openflow_connector_definitions.#", "1")),
					resourceshowoutputassert.OpenflowConnectorDefinitionsDatasourceShowOutput(t, dsLikeExact.DatasourceReference()).
						HasName(definition.Name).
						HasProvider(definition.Provider).
						HasVersion(definition.Version).
						HasDisplayName(definition.DisplayName).
						HasDescription(definition.Description).
						HasCategories(definition.Categories...).
						// This definition reports no minimum node type; the maximum is set.
						HasMinRuntimeNodeType("").
						HasMaxNodeCount(*definition.MaxNodeCount),
				),
			},
			// 2. like matching nothing
			{
				Config: config.FromModels(t, dsLikeMatchingNothing),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikeMatchingNothing.DatasourceReference(), "openflow_connector_definitions.#", "0")),
				),
			},
			// 3. like with a pattern rather than an exact name
			{
				Config: config.FromModels(t, dsLikePattern),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLikePattern.DatasourceReference(), "openflow_connector_definitions.#", "1")),
					resourceshowoutputassert.OpenflowConnectorDefinitionsDatasourceShowOutput(t, dsLikePattern.DatasourceReference()).
						HasName(definition.Name),
				),
			},
			// 4. limit
			{
				Config: config.FromModels(t, dsLimit),
				Check: assertThat(
					t,
					assert.Check(resource.TestCheckResourceAttr(dsLimit.DatasourceReference(), "openflow_connector_definitions.#", "1")),
				),
			},
		},
	})
}
