package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

// Use only parameters that can be set only on the account level for the time-being.
// TODO [SNOW-1866453]: add more acc tests for the remaining parameters

func TestAcc_AccountParameter(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterAllowIDToken), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterAllowIDToken),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterAllowIDToken)).
					HasValueString("true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_PREVENT_LOAD_FROM_INLINE_URL(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterPreventLoadFromInlineURL), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterPreventLoadFromInlineURL),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterPreventLoadFromInlineURL)).
					HasValueString("true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterRequireStorageIntegrationForStageCreation), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterRequireStorageIntegrationForStageCreation),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterRequireStorageIntegrationForStageCreation)).
					HasValueString("true"),
				),
			},
		},
	})
}

// Proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2573 is solved.
// Instead of TIMEZONE, we used INITIAL_REPLICATION_SIZE_LIMIT_IN_TB which is only settable on account so it does not mess with other tests.
func TestAcc_AccountParameter_Issue2573(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterInitialReplicationSizeLimitInTB), "3.0")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterTimezone),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterInitialReplicationSizeLimitInTB)).
					HasValueString("3.0"),
				),
			},
			{
				ResourceName:            "snowflake_account_parameter.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAcc_AccountParameter_Issue3025(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList)).
					HasValueString("true"),
				),
			},
			{
				ResourceName:            "snowflake_account_parameter.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func TestAcc_AccountParameter_ENFORCE_NETWORK_RULES_FOR_INTERNAL_STAGES(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterRequireStorageIntegrationForStageCreation), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterRequireStorageIntegrationForStageCreation),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterRequireStorageIntegrationForStageCreation)).
					HasValueString("true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_INITIAL_REPLICATION_SIZE_LIMIT_IN_TB(t *testing.T) {
	accountParameterModel := model.AccountParameter("test", string(sdk.AccountParameterInitialReplicationSizeLimitInTB), "3.0")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckAccountParameterUnset(t, sdk.AccountParameterInitialReplicationSizeLimitInTB),
		Steps: []resource.TestStep{
			{
				Config: config.FromModels(t, accountParameterModel),
				Check: assertThat(t, resourceassert.AccountParameterResource(t, accountParameterModel.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterInitialReplicationSizeLimitInTB)).
					HasValueString("3.0"),
				),
			},
		},
	})
}
