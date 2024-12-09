package resources_test

import (
	"fmt"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance/testenvs"
	"github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/sdk"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ObjectParameter(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	// TODO(SNOW-1528546): Remove after parameter-setting resources are using UNSET in the delete operation.
	t.Cleanup(func() {
		acc.TestClient().Database.Alter(t, acc.TestClient().Ids.DatabaseId(), &sdk.AlterDatabaseOptions{
			Unset: &sdk.DatabaseUnset{
				UserTaskTimeoutMs: sdk.Bool(true),
			},
		})
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigBasic("USER_TASK_TIMEOUT_MS", "1000", acc.TestDatabaseName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", "USER_TASK_TIMEOUT_MS"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "1000"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func TestAcc_ObjectParameterAccount(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	// TODO(SNOW-1528546): Remove after parameter-setting resources are using UNSET in the delete operation.
	t.Cleanup(func() {
		acc.TestClient().Parameter.UnsetAccountParameter(t, sdk.AccountParameterDataRetentionTimeInDays)
	})
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigOnAccount("DATA_RETENTION_TIME_IN_DAYS", "0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", "DATA_RETENTION_TIME_IN_DAYS"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "0"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "true"),
				),
			},
		},
	})
}

func TestAcc_UserParameter(t *testing.T) {
	_ = testenvs.GetOrSkipTest(t, testenvs.EnableAcceptance)

	user, userCleanup := acc.TestClient().User.CreateUser(t)
	t.Cleanup(userCleanup)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userParameterConfigBasic(user.ID(), "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR", "true"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "key", "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "value", "true"),
					resource.TestCheckResourceAttr("snowflake_object_parameter.p", "on_account", "false"),
				),
			},
		},
	})
}

func objectParameterConfigOnAccount(key, value string) string {
	s := `
resource "snowflake_object_parameter" "p" {
	key = "%s"
	value = "%s"
	on_account = true
}
`
	return fmt.Sprintf(s, key, value)
}

func objectParameterConfigBasic(key, value, databaseName string) string {
	s := `
resource "snowflake_object_parameter" "p" {
	key = "%s"
	value = "%s"
	object_type = "DATABASE"
	object_identifier {
		name = "%s"
	}
}
`
	return fmt.Sprintf(s, key, value, databaseName)
}

func userParameterConfigBasic(userId sdk.AccountObjectIdentifier, key string, value string) string {
	return fmt.Sprintf(`
resource "snowflake_object_parameter" "p" {
	key = "%[2]s"
	value = "%[3]s"
	object_type = "USER"
	object_identifier {
		name = %[1]s
	}
}
`, userId.FullyQualifiedName(), key, value)
}
