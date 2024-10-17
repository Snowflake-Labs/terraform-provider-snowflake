package resources_test

import (
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/provider/resources"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_AuthenticationPolicy(t *testing.T) {
	accName := acc.TestClient().Ids.Alpha()
	comment := "This is a test resource"
	m := func(authenticationMethods []string, mfaAuthenticationMethods []string, mfaEnrollment string, clientTypes []string, securityIntegrations []string) map[string]config.Variable {
		authenticationMethodsStringVariables := make([]config.Variable, len(authenticationMethods))
		for i, v := range authenticationMethods {
			authenticationMethodsStringVariables[i] = config.StringVariable(v)
		}
		mfaAuthenticationMethodsStringVariables := make([]config.Variable, len(mfaAuthenticationMethods))
		for i, v := range mfaAuthenticationMethods {
			mfaAuthenticationMethodsStringVariables[i] = config.StringVariable(v)
		}
		clientTypesStringVariables := make([]config.Variable, len(clientTypes))
		for i, v := range clientTypes {
			clientTypesStringVariables[i] = config.StringVariable(v)
		}
		securityIntegrationsStringVariables := make([]config.Variable, len(securityIntegrations))
		for i, v := range securityIntegrations {
			securityIntegrationsStringVariables[i] = config.StringVariable(v)
		}

		return map[string]config.Variable{
			"name":                       config.StringVariable(accName),
			"database":                   config.StringVariable(acc.TestDatabaseName),
			"schema":                     config.StringVariable(acc.TestSchemaName),
			"authentication_methods":     config.SetVariable(authenticationMethodsStringVariables...),
			"mfa_authentication_methods": config.SetVariable(mfaAuthenticationMethodsStringVariables...),
			"mfa_enrollment":             config.StringVariable(mfaEnrollment),
			"client_types":               config.SetVariable(clientTypesStringVariables...),
			"security_integrations":      config.SetVariable(securityIntegrationsStringVariables...),
			"comment":                    config.StringVariable(comment),
		}
	}
	variables1 := m([]string{"PASSWORD"}, []string{"PASSWORD"}, "REQUIRED", []string{"SNOWFLAKE_UI"}, []string{"ALL"})

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: acc.CheckDestroy(t, resources.AuthenticationPolicy),
		Steps: []resource.TestStep{
			{
				ConfigDirectory: config.TestNameDirectory(),
				ConfigVariables: variables1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "name", accName),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "authentication_methods.0", "PASSWORD"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "mfa_authentication_methods.0", "PASSWORD"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "mfa_enrollment", "REQUIRED"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "client_types.0", "SNOWFLAKE_UI"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "security_integrations.0", "ALL"),
					resource.TestCheckResourceAttr("snowflake_password_policy.pa", "comment", comment),
				),
			},
			{
				ConfigDirectory:   config.TestNameDirectory(),
				ConfigVariables:   variables1,
				ResourceName:      "snowflake_authentication_policy.authentication_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
