//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_TestResource_DataTypeDiffHandling(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)
	acc.TestAccPreCheck(t)

	envName := fmt.Sprintf("%s_%s", testenvs.TestResourceDataTypeDiffHandlingEnv, strings.ToUpper(random.AlphaN(10)))

	testConfig := func(configValue string) string {
		return fmt.Sprintf(`
resource "snowflake_test_resource_data_type_diff_handling" "test" {
	env_name = "%[2]s"
	return_data_type = "%[1]s"
}
`, configValue, envName)
	}

	testCases := []struct {
		ConfigValue    string
		NewConfigValue string
		ExternalValue  string
		ExpectChanges  bool
	}{
		{"NUMBER(20, 4)", "VARCHAR(20)", "", true},
		// TODO: add the rest of the use cases
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(fmt.Sprintf("TestAcc_TestResource_DataTypeDiffHandling config value: %s, new config value: %s, external value: %s, expecitng changes: %t", tc.ConfigValue, tc.NewConfigValue, tc.ExternalValue, tc.ExpectChanges), func(t *testing.T) {
			var planchecks []plancheck.PlanCheck
			if tc.ExpectChanges {
				planchecks = []plancheck.PlanCheck{
					plancheck.ExpectResourceAction("snowflake_test_resource_data_type_diff_handling.test", plancheck.ResourceActionCreate),
				}
			} else {
				planchecks = []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan("snowflake_test_resource_data_type_diff_handling.test"),
				}
			}

			var newValue string
			if tc.NewConfigValue != "" {
				newValue = tc.NewConfigValue
			} else {
				newValue = tc.ConfigValue
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Schema),
				Steps: []resource.TestStep{
					{
						Config: testConfig(tc.ConfigValue),
					},
					{
						PreConfig: func() {
							t.Setenv(envName, tc.ExternalValue)
						},
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: planchecks,
						},
						Config: testConfig(newValue),
					},
				},
			})
		})
	}
}
