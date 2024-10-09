package datasources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ParametersOnAccount(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnAccount(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "pattern", "AUTOCOMMIT"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "parameters.#", "1"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "parameters.0.key", "AUTOCOMMIT"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
				),
			},
		},
	})
}

func TestAcc_ParametersOnSession(t *testing.T) {
	userName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnSession(userName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.value"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "user", userName),
				),
			},
		},
	})
}

func TestAcc_ParametersOnObject(t *testing.T) {
	dbName := acc.TestClient().Ids.Alpha()
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: parametersConfigOnObject(dbName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.#"),
					resource.TestCheckResourceAttrSet("data.snowflake_parameters.p", "parameters.0.key"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "object_type", "DATABASE"),
					resource.TestCheckResourceAttr("data.snowflake_parameters.p", "object_name", dbName),
				),
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2353 is fixed
func TestAcc_Parameters_TransactionAbortOnErrorCanBeSet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: `resource "snowflake_account_parameter" "test" {
				   key   = "TRANSACTION_ABORT_ON_ERROR"
				   value = "true"
				}`,
			},
		},
	})
}

// proves https://github.com/Snowflake-Labs/terraform-provider-snowflake/issues/2353 is fixed
// done on user, to not interfere with other parallel tests on the same account
func TestAcc_Parameters_QuotedIdentifiersIgnoreCaseCanBeSet(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	user, userCleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		PreCheck:                 func() { acc.TestAccPreCheck(t) },
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		Steps: []resource.TestStep{
			{
				Config: sessionParameterOnUser(user.ID()),
			},
		},
	})
}

func sessionParameterOnUser(userId sdk.AccountObjectIdentifier) string {
	return fmt.Sprintf(
		`
    resource "snowflake_session_parameter" "test" {
		key   = "QUOTED_IDENTIFIERS_IGNORE_CASE"
		value = "true"
		user  = %[1]s
	}`, userId.FullyQualifiedName())
}

func parametersConfigOnAccount() string {
	return `data "snowflake_parameters" "p" {
		parameter_type = "ACCOUNT"
		pattern = "AUTOCOMMIT"
	}`
}

func parametersConfigOnSession(user string) string {
	s := `
	resource "snowflake_user" "u" {
		name = "%s"
	}

	data "snowflake_parameters" "p" {
		parameter_type = "SESSION"
		user = snowflake_user.u.name
	}`
	return fmt.Sprintf(s, user)
}

func parametersConfigOnObject(name string) string {
	stmt := `
	resource "snowflake_database" "d" {
		name = "%s"
	}
	data "snowflake_parameters" "p" {
		parameter_type = "OBJECT"
		object_type = "DATABASE"
		object_name = snowflake_database.d.name
	}`
	return fmt.Sprintf(stmt, name)
}
