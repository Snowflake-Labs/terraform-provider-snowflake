//go:build !account_level_tests

package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	tfjson "github.com/hashicorp/terraform-json"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/helpers/random"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/planchecks"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
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
	propertyName := "top_level_datatype"

	testConfig := func(configValue string) string {
		return fmt.Sprintf(`
resource "%[3]s" "%[4]s" {
	env_name = "%[2]s"
	%[5]s = "%[1]s"
}
`, configValue, envName, resourceType, resourceName, propertyName)
	}

	type DataTypeDiffHandlingTestCase struct {
		ConfigValue    string
		NewConfigValue string
		ExternalValue  string
		ExpectChanges  bool
	}

	changeInConfig := func(configValue string, newConfigValue string, expectChanges bool) DataTypeDiffHandlingTestCase {
		return DataTypeDiffHandlingTestCase{
			ConfigValue:    configValue,
			NewConfigValue: newConfigValue,
			ExpectChanges:  expectChanges,
		}
	}

	externalChange := func(configValue string, externalValue string, expectChanges bool) DataTypeDiffHandlingTestCase {
		return DataTypeDiffHandlingTestCase{
			ConfigValue:   configValue,
			ExternalValue: externalValue,
			ExpectChanges: expectChanges,
		}
	}

	testCases := []DataTypeDiffHandlingTestCase{
		// different data type
		changeInConfig("NUMBER(20, 4)", "VARCHAR(20)", true),
		changeInConfig("NUMBER(20, 4)", "VARCHAR", true),
		changeInConfig("NUMBER(20)", "VARCHAR(20)", true),
		changeInConfig("NUMBER", "VARCHAR(20)", true),
		changeInConfig("NUMBER", "VARCHAR", true),

		// same data type - no attributes before
		changeInConfig("NUMBER", "NUMBER", false),
		changeInConfig("NUMBER", "NUMBER(20)", true),
		changeInConfig("NUMBER", "NUMBER(20, 4)", true),
		changeInConfig("NUMBER", "NUMBER(38)", false),
		changeInConfig("NUMBER", "NUMBER(38, 0)", false),

		// same data type - one attribute before
		changeInConfig("NUMBER(20)", "NUMBER(20)", false),
		changeInConfig("NUMBER(20)", "NUMBER", true),
		changeInConfig("NUMBER(20)", "NUMBER(21)", true),
		changeInConfig("NUMBER(20)", "NUMBER(20, 0)", false),
		changeInConfig("NUMBER(20)", "NUMBER(20, 4)", true),
		changeInConfig("NUMBER(20)", "NUMBER(21, 4)", true),

		// same data type - two attributes before
		changeInConfig("NUMBER(20, 3)", "NUMBER(20, 3)", false),
		changeInConfig("NUMBER(20, 3)", "NUMBER", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(20)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(20, 4)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(21)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(21, 3)", true),
		changeInConfig("NUMBER(20, 3)", "NUMBER(21, 4)", true),

		// same data type - one attribute but default before
		changeInConfig("NUMBER(38)", "NUMBER(38)", false),
		changeInConfig("NUMBER(38)", "NUMBER", false),
		changeInConfig("NUMBER(38)", "NUMBER(20)", true),
		changeInConfig("NUMBER(38)", "NUMBER(20, 3)", true),
		changeInConfig("NUMBER(38)", "NUMBER(38, 2)", true),
		changeInConfig("NUMBER(38)", "NUMBER(38, 0)", false),

		// same data type - two attributes but default before
		changeInConfig("NUMBER(38, 0)", "NUMBER(38, 0)", false),
		changeInConfig("NUMBER(38, 0)", "NUMBER", false),
		changeInConfig("NUMBER(38, 0)", "NUMBER(38)", false),
		changeInConfig("NUMBER(38, 0)", "NUMBER(20)", true),
		changeInConfig("NUMBER(38, 0)", "NUMBER(20, 3)", true),
		changeInConfig("NUMBER(38, 0)", "NUMBER(38, 2)", true),

		// different data type
		externalChange("NUMBER(20, 4)", "VARCHAR(20)", true),
		externalChange("NUMBER(20, 4)", "VARCHAR", true),
		externalChange("NUMBER(20)", "VARCHAR(20)", true),
		externalChange("NUMBER", "VARCHAR(20)", true),
		externalChange("NUMBER", "VARCHAR", true),

		// same data type - no attributes before
		externalChange("NUMBER", "NUMBER", false),
		externalChange("NUMBER", "NUMBER(20)", true),
		externalChange("NUMBER", "NUMBER(20, 4)", true),
		externalChange("NUMBER", "NUMBER(38)", false),
		externalChange("NUMBER", "NUMBER(38, 0)", false),

		// same data type - one attribute before
		externalChange("NUMBER(20)", "NUMBER(20)", false),
		externalChange("NUMBER(20)", "NUMBER", false),
		externalChange("NUMBER(20)", "NUMBER(21)", true),
		externalChange("NUMBER(20)", "NUMBER(20, 0)", false),
		externalChange("NUMBER(20)", "NUMBER(20, 4)", true),
		externalChange("NUMBER(20)", "NUMBER(21, 4)", true),

		// same data type - two attributes before
		externalChange("NUMBER(20, 3)", "NUMBER(20, 3)", false),
		externalChange("NUMBER(20, 3)", "NUMBER", false),
		externalChange("NUMBER(20, 3)", "NUMBER(20)", false),
		externalChange("NUMBER(20, 3)", "NUMBER(20, 4)", true),
		externalChange("NUMBER(20, 3)", "NUMBER(21)", true),
		externalChange("NUMBER(20, 3)", "NUMBER(21, 3)", true),
		externalChange("NUMBER(20, 3)", "NUMBER(21, 4)", true),

		// same data type - one attribute but default before
		externalChange("NUMBER(38)", "NUMBER(38)", false),
		externalChange("NUMBER(38)", "NUMBER", false),
		externalChange("NUMBER(38)", "NUMBER(20)", true),
		externalChange("NUMBER(38)", "NUMBER(20, 3)", true),
		externalChange("NUMBER(38)", "NUMBER(38, 2)", true),
		externalChange("NUMBER(38)", "NUMBER(38, 0)", false),

		// same data type - two attributes but default before
		externalChange("NUMBER(38, 0)", "NUMBER(38, 0)", false),
		externalChange("NUMBER(38, 0)", "NUMBER", false),
		externalChange("NUMBER(38, 0)", "NUMBER(38)", false),
		externalChange("NUMBER(38, 0)", "NUMBER(20)", true),
		externalChange("NUMBER(38, 0)", "NUMBER(20, 3)", true),
		externalChange("NUMBER(38, 0)", "NUMBER(38, 2)", true),
	}

	for _, testCase := range testCases {
		tc := testCase
		t.Run(fmt.Sprintf("TestAcc_TestResource_DataTypeDiffHandling config value: %s, new config value: %s, external value: %s, expecitng changes: %t", tc.ConfigValue, tc.NewConfigValue, tc.ExternalValue, tc.ExpectChanges), func(t *testing.T) {
			var checks []plancheck.PlanCheck
			if tc.ExpectChanges {
				if tc.ExternalValue != "" {
					checks = []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectDrift(resourceReference, propertyName, sdk.String(tc.ConfigValue), sdk.String(tc.ExternalValue)),
						planchecks.ExpectChange(resourceReference, propertyName, tfjson.ActionUpdate, sdk.String(tc.ExternalValue), sdk.String(tc.ConfigValue)),
					}
				} else {
					checks = []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceReference, plancheck.ResourceActionUpdate),
						planchecks.ExpectChange(resourceReference, propertyName, tfjson.ActionUpdate, sdk.String(tc.ConfigValue), sdk.String(tc.NewConfigValue)),
					}
				}
			} else {
				checks = []plancheck.PlanCheck{
					plancheck.ExpectEmptyPlan(),
				}
			}

			var newConfigValue string
			if tc.NewConfigValue != "" {
				newConfigValue = tc.NewConfigValue
			} else {
				newConfigValue = tc.ConfigValue
			}

			var expectedValue string
			if tc.ExpectChanges {
				expectedValue = newConfigValue
			} else {
				expectedValue = tc.ConfigValue
			}

			resource.Test(t, resource.TestCase{
				ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
				TerraformVersionChecks: []tfversion.TerraformVersionCheck{
					tfversion.RequireAbove(tfversion.Version1_5_0),
				},
				Steps: []resource.TestStep{
					{
						// our test resource manages this env, so we remove it before the test start
						PreConfig: func() {
							t.Setenv(envName, "")
						},
						Config: testConfig(tc.ConfigValue),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, propertyName, tc.ConfigValue),
						),
					},
					{
						PreConfig: func() {
							t.Setenv(envName, tc.ExternalValue)
						},
						ConfigPlanChecks: resource.ConfigPlanChecks{
							PreApply: checks,
						},
						Config: testConfig(newConfigValue),
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckResourceAttr(resourceReference, propertyName, expectedValue),
						),
					},
				},
			})
		})
	}
}
