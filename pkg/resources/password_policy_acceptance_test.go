package resources_test

import (
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_PasswordPolicy(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	m := func(minLength int, maxLength int, comment string) map[string]config.Variable {
		return map[string]config.Variable{
			"name":       config.StringVariable(accName),
			"database":   config.StringVariable(acc.TestDatabaseName),
			"schema":     config.StringVariable(acc.TestSchemaName),
			"min_length": config.IntegerVariable(minLength),
			"max_length": config.IntegerVariable(maxLength),
			"comment":    config.StringVariable(comment),
		}
	}

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
				ConfigVariables: m(10, 30, "this is a test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", accName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "10"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "30"),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: m(20, 50, "this is a test resource"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "min_length", "20"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_length", "50"),
				),
			},
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: m(20, 50, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", ""),
				),
			},
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   m(20, 50, ""),
				ResourceName:      "snowflake_password_policy.pa",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAcc_PasswordPolicyMaxAgeDays(t *testing.T) {
	accName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	m := func(maxAgeDays int) map[string]config.Variable {
		return map[string]config.Variable{
			"name":         config.StringVariable(accName),
			"database":     config.StringVariable(acc.TestDatabaseName),
			"schema":       config.StringVariable(acc.TestSchemaName),
			"max_age_days": config.IntegerVariable(maxAgeDays),
		}
	}

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
			// Unsets properly
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_noOptionals"),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(accName),
					"database": config.StringVariable(acc.TestDatabaseName),
					"schema":   config.StringVariable(acc.TestSchemaName),
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "max_age_days", "90"),
				),
			},
			{
				ConfigDirectory: acc.ConfigurationDirectory("TestAcc_PasswordPolicy_noOptionals"),
				ConfigVariables: map[string]config.Variable{
					"name":     config.StringVariable(accName),
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
