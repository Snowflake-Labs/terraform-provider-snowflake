package resources_test

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	m := func(minLength int, maxLength int, minUpperCaseChars int, minLowerCaseChars int, minNumericChars int, minSpecialChars int, minAgeDays int, maxAgeDays int, maxRetries int, lockoutTimeMins int, history int, comment string) map[string]config.Variable {
		return map[string]config.Variable{
			"name":                 config.StringVariable(accName),
			"database":             config.StringVariable(acc.TestDatabaseName),
			"schema":               config.StringVariable(acc.TestSchemaName),
			"min_length":           config.IntegerVariable(minLength),
			"max_length":           config.IntegerVariable(maxLength),
			"min_upper_case_chars": config.IntegerVariable(minUpperCaseChars),
			"min_lower_case_chars": config.IntegerVariable(minLowerCaseChars),
			"min_numeric_chars":    config.IntegerVariable(minNumericChars),
			"min_special_chars":    config.IntegerVariable(minSpecialChars),
			"min_age_days":         config.IntegerVariable(minAgeDays),
			"max_age_days":         config.IntegerVariable(maxAgeDays),
			"max_retries":          config.IntegerVariable(maxRetries),
			"lockout_time_mins":    config.IntegerVariable(lockoutTimeMins),
			"history":              config.IntegerVariable(history),
			"comment":              config.StringVariable(comment),
		}
	}
	variables1 := m(10, 30, 2, 3, 4, 5, 6, 7, 8, 9, 10, "this is a test resource")
	variables2 := m(20, 50, 1, 2, 3, 4, 5, 6, 7, 8, 9, "this is a test resource")
	variables3 := m(20, 50, 1, 2, 3, 4, 5, 6, 7, 8, 9, "")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", accName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "30"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_upper_case_chars", "2"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_lower_case_chars", "3"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_numeric_chars", "4"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_special_chars", "5"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_age_days", "6"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "7"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_retries", "8"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "lockout_time_mins", "9"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "history", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", "this is a test resource"),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "50"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_upper_case_chars", "1"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_lower_case_chars", "2"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_numeric_chars", "3"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_special_chars", "4"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_age_days", "5"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "6"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_retries", "7"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "lockout_time_mins", "8"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "history", "9"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", "this is a test resource"),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", ""),
				),
			},
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   variables3,
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_PasswordPolicyMaxAgeDays(t *testing.T) {
	name := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	newName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))

	m := func(maxAgeDays int) map[string]config.Variable {
		return map[string]config.Variable{
			"name":         config.StringVariable(name),
			"database":     config.StringVariable(acc.TestDatabaseName),
			"schema":       config.StringVariable(acc.TestSchemaName),
			"max_age_days": config.IntegerVariable(maxAgeDays),
		}
	}

	configValueWithNewName := m(10)
	configValueWithNewName["name"] = config.StringVariable(newName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			// Creation sets zero properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_withMaxAgeDays"),
				ConfigVariables: m(0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_withMaxAgeDays"),
				ConfigVariables: m(10),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "10"),
				),
			},
			// Update sets zero properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_withMaxAgeDays"),
				ConfigVariables: m(0),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "0"),
				),
			},
			// Rename + Unsets properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_noOptionals"),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(newName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"schema":   config.StringVariable(acc.TestSchemaName),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("snowflake_password_policy.pa", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", newName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "90"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_noOptionals"),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(name),
					"database": config.StringVariable(acc.TestDatabaseName),
					"schema":   config.StringVariable(acc.TestSchemaName),
				},
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
