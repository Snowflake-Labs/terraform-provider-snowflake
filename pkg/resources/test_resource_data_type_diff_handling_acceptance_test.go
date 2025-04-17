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
	resourceType := "snowflake_test_resource_data_type_diff_handling"
	resourceName := "test"
	resourceReference := fmt.Sprintf("%s.%s", resourceType, resourceName)

	testConfig := func(configValue string) string {
		return fmt.Sprintf(`
resource "%[3]s" "%[4]s" {
	env_name = "%[2]s"
	return_data_type = "%[1]s"
}
`, configValue, envName, resourceType, resourceName)
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
					// TODO: add more checks (change from-to)
					plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionCreate),
				}
			} else {
				planchecks = []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(resourceReference),
				}
			}

			var newConfigValue string
			if tc.NewConfigValue != "" {
				newConfigValue = tc.NewConfigValue
			} else {
				newConfigValue = tc.ConfigValue
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				CheckDestroy: acc.CheckDestroy(t, resources.Schema),
				Steps: []resource.TestStep{
					{
						// our test resource does not set the env, so we set it proactively
						PreConfig: func() {
							t.Setenv(envName, tc.ConfigValue)
						},
						Config: testConfig(tc.ConfigValue),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, "return_data_type", tc.ConfigValue),
						),
					},
					{
						PreConfig: func() {
							t.Setenv(envName, tc.ExternalValue)
						},
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: planchecks,
						},
						Config: testConfig(newConfigValue),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, "return_data_type", newConfigValue),
						),
					},
				},
			})
		})
	}
}
