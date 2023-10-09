package resources_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ObjectParameter(t *testing.T) {
	prefix := "tst-terraform" + strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: objectParameterConfigBasic(prefix, "USER_TASK_TIMEOUT_MS", "1000"),
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
	resource.ParallelTest(t, resource.TestCase{
		Providers:    acc.TestAccProviders(),
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

func objectParameterConfigBasic(prefix, key, value string) string {
	s := `
resource "snowflake_database" "d" {
	name = "%s"
}	
resource "snowflake_object_parameter" "p" {
	key = "%s"
	value = "%s"
	object_type = "DATABASE"
	object_identifier {
		name = snowflake_database.d.name
	}
}
`
	return fmt.Sprintf(s, prefix, key, value)
}
