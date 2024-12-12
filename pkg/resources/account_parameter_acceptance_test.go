package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/assert/resourceassert"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/bettertestspoc/config/model"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AccountParameter(t *testing.T) {
	model := model.AccountParameter("test", "ALLOW_ID_TOKEN", "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.AccountParameterResource(t, model.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterAllowIDToken)).
					HasValueString("true"),
				),
			},
		},
	})
}

func accountParameterBasic(key, value string) string {
	s := `
resource "snowflake_account_parameter" "p" {
	key = "%s"
	value = "%s"
}
`
	return fmt.Sprintf(s, key, value)
}

func TestAcc_AccountParameter_PREVENT_LOAD_FROM_INLINE_URL(t *testing.T) {
	model := model.AccountParameter("test", string(sdk.AccountParameterPreventLoadFromInlineURL), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.AccountParameterResource(t, model.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterPreventLoadFromInlineURL)).
					HasValueString("true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_REQUIRE_STORAGE_INTEGRATION_FOR_STAGE_CREATION(t *testing.T) {
	model := model.AccountParameter("test", string(sdk.AccountParameterRequireStorageIntegrationForStageCreation), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.AccountParameterResource(t, model.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterRequireStorageIntegrationForStageCreation)).
					HasValueString("true"),
				),
			},
		},
	})
}

func TestAcc_AccountParameter_Issue2573(t *testing.T) {
	model := model.AccountParameter("test", string(sdk.AccountParameterTimezone), "Etc/UTC")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.AccountParameterResource(t, model.ResourceReference()).
					HasKeyString(string(sdk.AccountParameterTimezone)).
					HasValueString("Etc/UTC"),
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
	model := model.AccountParameter("test", string(sdk.AccountParameterOAuthAddPrivilegedRolesToBlockedList), "true")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: config.FromModel(t, model),
				Check: assert.AssertThat(t, resourceassert.AccountParameterResource(t, model.ResourceReference()).
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

// TODO(next pr): add more acc tests for the remaining parameters
// TODO(next pr): check unsetting in CheckDestroy
