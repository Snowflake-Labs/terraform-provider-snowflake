package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAcc_ObjectParameter(t *testing.T) {
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
	userName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acc.TestAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.RequireAbove(tfversion.Version1_5_0),
		},
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: userParameterConfigBasic(userName, "ENABLE_UNREDACTED_QUERY_SYNTAX_ERROR", "true"),
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

func userParameterConfigBasic(userName string, key string, value string) string {
	s := `
resource "snowflake_user" "user" {
	name = "%s"
}
resource "snowflake_object_parameter" "p" {
	key = "%s"
	value = "%s"
	object_type = "USER"
	object_identifier {
		name = snowflake_user.user.name
	}
}
`
	return fmt.Sprintf(s, userName, key, value)
}
